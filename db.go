package main

import (
	"encoding/binary"
	"encoding/json"

	bolt "github.com/coreos/bbolt"
)

func initDB() {
	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("news"))
		tx.CreateBucketIfNotExists([]byte("settings"))
		tx.CreateBucketIfNotExists([]byte("files"))
		tx.CreateBucketIfNotExists([]byte("ticker"))
		tx.CreateBucketIfNotExists([]byte("graphs"))

		return nil
	})
}

func getNewsItems() ([]NewsItem, error) {
	items := []NewsItem{}

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("news"))

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			item := NewNewsItem()
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
	item := NewNewsItem()

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

func editSettings(key string, value interface{}) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("settings"))
		bytes, err := json.Marshal(value)
		if err != nil {
			return err
		}
		return b.Put([]byte(key), bytes)
	})
}

func getSetting(key string, item interface{}) error {
	return db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("settings"))
		bytes := b.Get([]byte(key))
		if len(bytes) == 0 {
			return nil
		}
		return json.Unmarshal(bytes, &item)
	})
}

func saveFile(key string, content []byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("files"))
		return b.Put([]byte(key), content)
	})
}

func getFile(key string) (out []byte, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("files"))
		out = b.Get([]byte(key))

		return nil
	})

	return
}

func getTickerItems() ([]TickerItem, error) {
	items := []TickerItem{}

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("ticker"))

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			item := TickerItem{}
			json.Unmarshal(v, &item)
			items = append(items, item)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return items, err
}

func addTickerItem(item TickerItem) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("ticker"))

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

func deleteTickerItem(id int) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("ticker"))
		return b.Delete(itob(id))
	})
}

func getKeukenDienstItems() ([]KeukendienstItem, error) {
	items := []KeukendienstItem{}

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("settings"))

		data := b.Get([]byte("keukendienst"))
		json.Unmarshal(data, &items)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return items, err
}

func setKeukenDienstItems(items []KeukendienstItem) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("settings"))

		buf, err := json.Marshal(items)
		if err != nil {
			return err
		}
		return b.Put([]byte("keukendienst"), buf)
	})
}

func getGraphItems() ([]GraphItem, error) {
	items := []GraphItem{}

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("graphs"))

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			item := GraphItem{}
			json.Unmarshal(v, &item)
			items = append(items, item)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return items, err
}

func addGraphItem(item GraphItem) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("graphs"))

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

func deleteGraphItem(id int) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("graphs"))
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
