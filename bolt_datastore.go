package main

import bolt "go.etcd.io/bbolt"

type BoltDatastore struct {
	db *bolt.DB
}

func NewBoltDatastore(file string) *BoltDatastore {
	db, err := bolt.Open(file, 0600, nil)
	if err != nil {
		panic(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("data"))
		return err
	})
	return &BoltDatastore{db}
}

func (datastore *BoltDatastore) Close() {
	datastore.db.Close()
}

func (datastore *BoltDatastore) Set(key, value string) error {
	return datastore.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("data")).Put([]byte(key), []byte(value))
	})
}

func (datastore *BoltDatastore) Delete(key string) error {
	return datastore.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("data")).Delete([]byte(key))
	})
}

func (datastore *BoltDatastore) Get(key string) (error, string) {
	var result string
	if err := datastore.db.View(func(tx *bolt.Tx) error {
		value := tx.Bucket([]byte("data")).Get([]byte(key))
		result = string(value)
		return nil
	}); err != nil {
		return err, ""
	}
	return nil, result
}
