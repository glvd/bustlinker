package data

import (
	"fmt"
	"github.com/ipfs/go-ipfs/linker/config"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	cacheName = "linker.s3db"
)

type data struct {
	db  *gorm.DB
	cfg config.CacheConfig
}

type Cache interface {
}

//New ...
func New(cfg config.CacheConfig, path, name string) Cache {
	path = filepath.Join(path, name)
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			panic(err)
		}
	}
	dsn := fmt.Sprintf("file:%s?_auth&_auth_user=amin&_auth_pass=admin", filepath.Join(path, cacheName))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	c := data{
		db:  db,
		cfg: cfg,
	}

	return &c
}
