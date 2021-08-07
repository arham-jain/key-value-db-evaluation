package main

import (
	"github.com/dgraph-io/badger/v3"
)

type BadgerDBWrapper struct {
	db *badger.DB
}

func NewBadgerDBWrapper(db *badger.DB) *BadgerDBWrapper {
	return &BadgerDBWrapper{
		db: db,
	}
}

func (b BadgerDBWrapper) BatchWriteToDB(records []KeyValueMarshal) error {
	batchWriter := b.db.NewWriteBatch()
	writeBatch := func(records []KeyValueMarshal) error {
		for _, item := range records {
			err := batchWriter.Set(item.Key, item.Value)
			if err != nil {
				return err
			}
		}
		return nil
	}
	err := writeBatch(records)
	if err != nil {
		return err
	}
	return batchWriter.Flush()
}

func (b BadgerDBWrapper) PrefixSearchCountFromDB(keyPrefix []byte) int {
	txn := b.db.NewTransaction(false)
	opts := badger.DefaultIteratorOptions
	opts.Prefix = keyPrefix
	it := txn.NewIterator(opts)
	return b.count(it)
}

func (b BadgerDBWrapper) count(iter *badger.Iterator) int {
	result := 0
	for iter.Rewind(); iter.Valid(); iter.Next() {
		result++
	}
	return result
}
