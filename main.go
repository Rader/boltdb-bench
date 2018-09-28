package main

import (
	"encoding/binary"
	"errors"
	"math/big"
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/rlp"

	"github.com/icstglobal/plasma/core/types"
	log "github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

var (
	emptyAddress   [20]byte
	utxoPrefix     = []byte{'c'}
	utxoBucketName = []byte("utxoset")
)

func main() {
	path := "bolt.db"
	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		log.WithError(err).WithField("path", path).Error("can not open db")
	}
	defer db.Close()

	log.SetLevel(log.DebugLevel)
	log.WithField("path", path).Info("bolt db opened")

	// trigger benchmark manually
	// benchmarkBody(db)

	// generate test db data
	// genData(db)

	// timing random read
	// const limit = 600000 * 1000
	// id := types.UTXOID{BlockNum: uint64(rand.Int63n(limit)), TxIndex: 1, OutIndex: 0}
	// key := utxoKey(id.BlockNum, id.TxIndex, id.OutIndex)
	// randomRead(db, key)
}

func utxoKey(blockNum uint64, txIdx uint32, outIdx byte) []byte {
	key := append(utxoPrefix, make([]byte, 13)...)
	binary.BigEndian.PutUint64(key[1:], blockNum)
	binary.BigEndian.PutUint32(key[9:], txIdx)
	key[13] = outIdx

	return key
}

func randomRead(db *bolt.DB, key []byte) error {
	return db.View(func(tx *bolt.Tx) error {
		buf := tx.Bucket(utxoBucketName).Get(key)
		if len(buf) == 0 {
			err := errors.New("utxo not found")
			log.Error(err)
			return err
		}
		// ingnore desirialization for benchmark test;
		// uncomment this for function test

		// utxo := types.UTXO{}
		// err := rlp.DecodeBytes(buf, &utxo)
		// if err != nil {
		// 	log.Error(err)
		// 	return err
		// }
		// log.Printf("utxo from db:%+v", utxo)
		return nil
	})
}

func genData(db *bolt.DB) {
	const loops = 600000
	// const loops = 1
	for n := 0; n <= loops; n++ {
		if err := db.Update(func(tx *bolt.Tx) error {

			b, err := tx.CreateBucketIfNotExists(utxoBucketName)
			if err != nil {
				return err
			}
			const batch = 1000
			for i := 1; i < batch+1; i++ {
				id := types.UTXOID{BlockNum: uint64(batch*n + i), TxIndex: 1, OutIndex: 0}
				utxo := types.UTXO{UTXOID: id, Owner: emptyAddress, Amount: big.NewInt(10000)}
				if buf, err := rlp.EncodeToBytes(&utxo); err != nil {
					log.Error(err)
					return err
				} else {
					key := utxoKey(id.BlockNum, id.TxIndex, id.OutIndex)
					if err := b.Put(key, buf); err != nil {
						log.Error(err)
						return err
					}
				}
			}
			return nil
		}); err != nil {
			log.Fatal(err)
		}
	}
}

func benchmarkBody(db *bolt.DB) {
	const limit = 600000 * 1000
	id := types.UTXOID{BlockNum: uint64(rand.Int63n(limit)), TxIndex: 1, OutIndex: 0}
	// key := utxoKey(id.BlockNum, id.TxIndex, id.OutIndex)

	id.BlockNum = uint64(rand.Int63n(limit))
	key := utxoKey(id.BlockNum, id.TxIndex, id.OutIndex)
	log.Debug("start to run benchmark")
	result := testing.Benchmark(func(b *testing.B) {
		if randomRead(db, key) != nil {
			log.Fatal("read db failed")
		}
		// log.Debug("quit from benchmark body", b.N)
	})
	log.Debug("end run benchmark")
	log.Printf("result.string:%+v\n result.memstring:%v", result.String(), result.MemString())
}
