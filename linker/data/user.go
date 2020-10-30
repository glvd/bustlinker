package data

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Hash        string
	IsPinned    bool
	IsSubscribe bool
}
