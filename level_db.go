package main

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type LevelDBWrapper struct {
	db *leveldb.DB
}

func NewLevelDBWrapper(db *leveldb.DB) *LevelDBWrapper {
	return &LevelDBWrapper{db: db}
}

func (l LevelDBWrapper) BatchWriteToDB(records []KeyValueMarshal) error {
	batch := new(leveldb.Batch)
	for _, item := range records {
		batch.Put(item.Key, item.Value)
	}
	err := l.db.Write(batch, nil)
	return err
}

func (l LevelDBWrapper) PrefixSearchCountFromDB(keyPrefix []byte) int {
	iter := l.db.NewIterator(util.BytesPrefix(keyPrefix), nil)
	defer iter.Release()
	return l.count(iter)
}

func (l LevelDBWrapper) count(iter iterator.Iterator) int {
	var result int
	for iter.Next() {
		result++
	}
	return result
}
