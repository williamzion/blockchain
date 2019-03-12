package main

import (
	"encoding/hex"
	"log"

	bolt "go.etcd.io/bbolt"
)

const utxoBucket = "chainstate"

// UTXOSet represents UTXO set.
type UTXOSet struct {
	Blockchain *Blockchain
}

// Reindex rebuilds the UTXO set.
func (u UTXOSet) Reindex() {
	db := u.Blockchain.db
	// bucketName := []byte(utxoBucket)

	err := db.Update(func(tx *bolt.Tx) error {
		// Removes the bucket if it exists.
		err := tx.DeleteBucket([]byte(utxoBucket))
		if err != nil && err != bolt.ErrBucketNotFound {
			log.Panic(err)
		}

		_, err = tx.CreateBucket([]byte(utxoBucket))
		if err != nil {
			log.Panic(err)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	UTXO := u.Blockchain.FindUTXO()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))

		for txID, outs := range UTXO {
			key, err := hex.DecodeString(txID)
			if err != nil {
				log.Panic(err)
			}

			err = b.Put(key, outs.Serialize())
			if err != nil {
				log.Panic(err)
			}
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}
