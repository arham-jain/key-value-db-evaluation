package main

import (
	"bytes"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
)

type KeyValueMarshal struct {
	Key   []byte
	Value []byte
}

type Key struct {
	Address     string
	BlockNumber *big.Int
	TxSequence  uint8
}

type Value struct {
	TxHash string
}

func (k *Key) MarshaledKey() []byte {
	buf := &bytes.Buffer{}
	// 20 bytes
	resultByteArr, _ := hexutil.Decode(k.Address)
	buf.Write(resultByteArr)
	// 3 byte
	blockNumberArr := k.BlockNumber.Bytes()
	buf.Write(blockNumberArr)
	// 1 byte for sequence
	buf.WriteByte(k.TxSequence)
	return buf.Bytes()
}

func KeyPrefix(address string) []byte {
	// 20 bytes
	resultByteArr, _ := hexutil.Decode(address)
	return resultByteArr
}

func (v *Value) MarshaledValue() []byte {
	buf := &bytes.Buffer{}
	txHashByteArr, _ := hexutil.Decode(v.TxHash)
	buf.Write(txHashByteArr)
	return buf.Bytes()
}
