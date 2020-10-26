package linker

import (
	"encoding/json"
	"fmt"
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

type Cache interface {
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
func NewCache(cfg config.CacheConfig, path, name string) Cache {
	path = filepath.Join(path, name)
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			panic(err)
		}
	}
	dsn := fmt.Sprintf("file:%s?_auth&_auth_user=amin&_auth_pass=admin", filepath.Join(path, "linker.db"))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	c := data{
		db:  db,
		cfg: config.CacheConfig{},
	}

	return &c
}
