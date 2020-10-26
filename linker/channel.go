package linker

import "gorm.io/gorm"

type Channel struct {
	gorm.Model
	ID     string
	IsFree bool
	Hash   string
	IsJoin bool
}
