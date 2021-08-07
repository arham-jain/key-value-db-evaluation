package main

import (
	cryptoRand "crypto/rand"
	"errors"
	"github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/rcrowley/go-metrics"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	status           bool
	batches          int
	dbVolumePath     string
	metricsRegistry  metrics.Registry
	timeBatchWriter  metrics.Timer
	timePrefixReader metrics.Timer
)

func init() {
	timeBatchWriter = metrics.NewRegisteredTimer("BatchWriterTimer", metricsRegistry)
	timePrefixReader = metrics.NewRegisteredTimer("PrefixReaderTimer", metricsRegistry)
	metricsRegistry = metrics.DefaultRegistry
	batches, _ = strconv.Atoi(os.Getenv("BATCHES"))
	dbVolumePath = os.Getenv("VOLUME_PATH")
}

func server() {
	router := gin.Default()
	router.GET("/metrics", func(c *gin.Context) {
		c.JSON(http.StatusOK, metricsRegistry.GetAll())
	})
	router.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]bool{"done": status})
	})
	err := router.Run(":9092")
	if err == nil {
		log.Print("started server successfully")
	} else {
		panic(errors.New("cannot start http server. error: " + err.Error()))
	}
}

////Main for leveldb
//func main() {
//	// Initialise database
//	DBInTest, err := leveldb.OpenFile(dbVolumePath, nil)
//	if err != nil {
//		panic("can not connect to LevelDB. Error: " + err.Error())
//	}
//	kvDatabase := NewLevelDBWrapper(DBInTest)
//	go evaluate(kvDatabase, err)
//	server()
//}

// Main for badger db
func main() {
	// Initialise database
	DBInTest, err := badger.Open(badger.DefaultOptions(dbVolumePath))
	if err != nil {
		panic("can not connect to BadgerDB. Error: " + err.Error())
	}
	kvDatabase := NewBadgerDBWrapper(DBInTest)
	go evaluate(kvDatabase, err)
	server()
}

func evaluate(kvDatabase DB, err error) {
	status = false
	log.Print("Starting evaluation")
	addressCountMap := make(map[string]int, 0)
	defer timeBatchWriter.Stop()
	for i := 0; i < batches; i++ {
		noOfTxs := int(math.Pow(2, float64(i)))
		log.Printf("writing batch %d, transaction count %d", i, noOfTxs)
		address, records := createRandomTransactions(noOfTxs)
		addressCountMap[address.String()] = noOfTxs
		var err error
		ts := time.Now()
		err = kvDatabase.BatchWriteToDB(records)
		timeBatchWriter.Update(time.Since(ts))
		if err != nil {
			log.Print("error writing to database")
			panic(err)
		}
	}
	totalCount := kvDatabase.PrefixSearchCountFromDB([]byte{})
	log.Printf("total transactions written to db, count: %d", totalCount)
	defer timePrefixReader.Stop()
	for k, v := range addressCountMap {
		var count int
		ts := time.Now()
		count = kvDatabase.PrefixSearchCountFromDB(KeyPrefix(k))
		timePrefixReader.Update(time.Since(ts))
		if count != v {
			log.Print("transaction count for address does not match")
			panic(err)
		}
		log.Printf("address %s validated, expected: %d, actual: %d", k, v, count)
	}
	log.Print("evaluation completed")
	status = true
}

func createRandomTransactions(noOfTxs int) (common.Address, []KeyValueMarshal) {
	records := make([]KeyValueMarshal, 0)
	address := randomAddress()
	for i := 0; i < noOfTxs; i++ {
		// create random tx details
		randomBlockNumber, _ := cryptoRand.Prime(cryptoRand.Reader, 128)
		randomTransactionSequence := rand.Intn(256)
		randomTransactionHash := randomHash()
		key := Key{
			Address:     address.String(),
			BlockNumber: randomBlockNumber,
			TxSequence:  uint8(randomTransactionSequence),
		}
		val := Value{
			TxHash: randomTransactionHash.String(),
		}
		records = append(records, KeyValueMarshal{
			Key:   key.MarshaledKey(),
			Value: val.MarshaledValue(),
		})
	}
	return address, records
}

func randomHash() common.Hash {
	var hash common.Hash
	n, err := rand.Read(hash[:])
	if n != common.HashLength || err != nil {
		log.Print("error creating random hash")
		panic(err)
	}
	return hash
}

func randomAddress() common.Address {
	var address common.Address
	n, err := rand.Read(address[:])
	if n != common.AddressLength || err != nil {
		log.Print("error creating random address")
		panic(err)
	}
	return address
}
