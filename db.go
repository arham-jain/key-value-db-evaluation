package main

type DB interface {
	BatchWriteToDB(records []KeyValueMarshal) error
	PrefixSearchCountFromDB(keyPrefix []byte) int
}
