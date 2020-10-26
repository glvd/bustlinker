package linker

import "gorm.io/gorm"

type Channel struct {
	gorm.Model
	ID     string
	Hash   string
	IsFree bool
	IsJoin bool
}
