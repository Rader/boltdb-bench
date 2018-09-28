package main

import (
	"encoding/binary"
	"math/rand"
	"testing"
	"time"

	"github.com/icstglobal/plasma/core/types"
	bolt "go.etcd.io/bbolt"
)

func BenchmarkRandomRead(b *testing.B) {
	path := "bolt.db"
	db, _ := bolt.Open(path, 0666, nil)
	defer db.Close()

	const limit = 600000 * 1000
	id := types.UTXOID{BlockNum: uint64(rand.Int63n(limit)), TxIndex: 1, OutIndex: 0}
	key := utxoKey(id.BlockNum, id.TxIndex, id.OutIndex)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// b.StopTimer()
		blockNum := uint64(rand.Int63n(limit))
		binary.BigEndian.PutUint64(key[1:], blockNum)
		// b.StartTimer()
		if randomRead(db, key) != nil {
			b.Fatal("read db failed")
		}
	}

	b.StopTimer()
	b.Log(b.N)
}

func TestRandomRead(t *testing.T) {
	path := "bolt.db"
	db, _ := bolt.Open(path, 0666, nil)
	defer db.Close()

	const limit = 600000 * 1000
	id := types.UTXOID{BlockNum: uint64(rand.Int63n(limit)), TxIndex: 1, OutIndex: 0}

	id.BlockNum = uint64(rand.Int63n(limit))
	key := utxoKey(id.BlockNum, id.TxIndex, id.OutIndex)
	start := time.Now()
	if err := randomRead(db, key); err != nil {
		t.Fatal("read db failed,", err)
	}
	du := time.Since(start)
	t.Log("time nano seconds:", du.Nanoseconds())

}

func BenchmarkRandomGen(b *testing.B) {
	const limit = 600000 * 1000
	var n uint64
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n = uint64(rand.Int63n(limit))
	}
	b.StopTimer()
	b.Log(n)
}
