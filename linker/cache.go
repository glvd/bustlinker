package linker

import (
	"encoding/json"
	"github.com/ipfs/go-ipfs/linker/config"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	cacheName = ".data"
)

type data struct {
	db  *gorm.DB
	cfg config.CacheConfig
}

type Cache struct {
	data
}

// Cacher ...
type Cacher interface {
	Load(hash string, data json.Unmarshaler) error
	Store(hash string, data json.Marshaler) error
	Update(hash string, up CacheUpdater) error
	Close() error
	Range(f func(hash string, value string) bool)
}

type CacheUpdater interface {
	json.Unmarshaler
	Do()
	json.Marshaler
}

// DataHashInfo ...
type DataHashInfo struct {
	DataHash string `json:"data_hash"`
}

// HashCache ...
func (v DataHashInfo) Hash() string {
	return v.DataHash
}

// Marshal ...
func (v DataHashInfo) Marshal() ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal ...
func (v *DataHashInfo) Unmarshal(b []byte) error {
	return json.Unmarshal(b, v)
}

// NewCache ...
func NewCache(cfg config.CacheConfig, path, name string) Cacher {

	path = filepath.Join(path, name)
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			panic(err)
		}
	}
	db, err := gorm.Open(sqlite.Open(filepath.Join(path, "linker.db")), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	c := data{
		db:  db,
		cfg: config.CacheConfig{},
	}

	return &c
}

func (c *data) UpdateBytes(hash string, b []byte) error {
	return c.db.Update(
		func(txn *badger.Txn) error {
			item, err := txn.Get([]byte(hash))
			if err != nil {
				return err
			}
			return item.Value(func(val []byte) error {
				return txn.Set([]byte(hash), b)
			})
		})
}

// Update ...
func (c *data) Update(hash string, up CacheUpdater) error {
	return c.db.Update(
		func(txn *badger.Txn) error {
			item, err := txn.Get([]byte(hash))
			if err != nil {
				return err
			}
			if up != nil {
				return item.Value(func(val []byte) error {
					err := up.UnmarshalJSON(val)
					if err != nil {
						//do nothing when have err
						return err
					}
					up.Do()
					encode, err := up.MarshalJSON()
					if err != nil {
						return err
					}
					return txn.Set([]byte(hash), encode)
				})
			}
			return nil
		})
}

// SaveNode ...
func (c *data) Store(hash string, data json.Marshaler) error {
	return c.db.Update(
		func(txn *badger.Txn) error {
			encode, err := data.MarshalJSON()
			if err != nil {
				return err
			}
			return txn.Set([]byte(hash), encode)
		})
}

// LoadNode ...
func (c *data) Load(hash string, data json.Unmarshaler) error {
	return c.db.View(
		func(txn *badger.Txn) error {
			item, err := txn.Get([]byte(hash))
			if err != nil {
				return err
			}
			return item.Value(func(val []byte) error {
				return data.UnmarshalJSON(val)
			})
		})
}

// Range ...
func (c *data) Range(f func(key, value string) bool) {
	err := c.db.View(func(txn *badger.Txn) error {
		iter := txn.NewIterator(c.iteratorOpts)
		defer iter.Close()
		var item *badger.Item
		continueFlag := true
		for iter.Rewind(); iter.Valid(); iter.Next() {
			if !continueFlag {
				return nil
			}
			item = iter.Item()
			err := iter.Item().Value(func(v []byte) error {
				key := item.Key()
				val, err := item.ValueCopy(v)
				if err != nil {
					return err
				}
				continueFlag = f(string(key), string(val))
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Errorw("range data failed", "err", err)
	}
}

// Close ...
func (c *data) Close() error {
	if c.db != nil {
		defer func() {
			c.db = nil
		}()
		return c.db.Close()
	}
	return nil
}
