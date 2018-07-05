package main

import (
	"fmt"
	"log"
	"os"
  "strconv"

	"github.com/boltdb/bolt"
)

const dbFile = "data_set_%s.db"
const bucket = "d"

// Blockchain implements interactions with a DB
type DataSet struct {
  height int
	db     *bolt.DB
}

// CreateBlockchain creates a new blockchain DB
func createDataSet(data, nodeID string) *DataSet {
	dbFile := fmt.Sprintf(dbFile, nodeID)
	if dbExists(dbFile) {
		fmt.Println("DB already exists.")
		os.Exit(1)
	}

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(bucket))
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("0"), []byte(data))
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("h"), []byte("1"))
		if err != nil {
			log.Panic(err)
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	dataSet := DataSet{1, db}

	return &dataSet
}

func (dataSet *DataSet) printDataSet(nodeID string) {
  dataSet.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
    h := b.Get([]byte("h"))
    height, _ := strconv.Atoi(string(h))

    for i := 0; i < height; i++ {
      data := b.Get([]byte(fmt.Sprintf("%v", i)))
      fmt.Println(string(data))
    }

		return nil
	})
}

// NewBlockchain creates a new Blockchain with genesis Block
func newDataSet(nodeID string) *DataSet {
	dbFile := fmt.Sprintf(dbFile, nodeID)
	if dbExists(dbFile) == false {
		fmt.Println("No existing db found. Create one first.")
		os.Exit(1)
	}

	var height int
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
    h := string(b.Get([]byte("h")))
		height, _ = strconv.Atoi(h)

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	dataSet := DataSet{height, db}

	return &dataSet
}

func dbExists(dbFile string) bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

func main() {
  ds := newDataSet("3001")
  ds.printDataSet("3001")
}
