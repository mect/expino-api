package main

import (
	"encoding/binary"
	"encoding/json"

	bolt "github.com/coreos/bbolt"
)

func getNewsItems() ([]NewsItem, error) {
	items := []NewsItem{}

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("news"))

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			item := NewsItem{}
			json.Unmarshal(v, &item)
			items = append(items, item)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return items, nil
}

func getNewsItem(id int) (NewsItem, error) {
	item := NewsItem{}

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("news"))
		return json.Unmarshal(b.Get(itob(id)), &item)
	})
	if err != nil {
		return item, err
	}

	return item, nil
}

func addNewsItem(item NewsItem) error {
	return db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucket([]byte("news")) // to be sure
		b := tx.Bucket([]byte("news"))

		id, _ := b.NextSequence()
		item.ID = int(id)

		buf, err := json.Marshal(item)
		if err != nil {
			return err
		}

		// Persist bytes to users bucket.
		return b.Put(itob(item.ID), buf)
	})
}

func editNewsItem(item NewsItem) error {
	return db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucket([]byte("news")) // to be sure
		b := tx.Bucket([]byte("news"))

		buf, err := json.Marshal(item)
		if err != nil {
			return err
		}

		// Persist bytes to users bucket.
		return b.Put(itob(item.ID), buf)
	})
}

func deleteNewsItem(id int) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("news"))
		return b.Delete(itob(id))
	})
}

// itob returns an 8-byte big endian representation of v.
// credit to https://github.com/boltdb/bolt
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
