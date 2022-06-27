package model

import (
	"time"

	"gorm.io/gorm"
)

type Common struct {
	ID         int64     `gorm:"primaryKey;autoIncrement:false"`
	CreateTime time.Time `gorm:"index;autoCreateTime:nano"`
	UpdateTime time.Time `gorm:"index;autoUpdateTime:nano"`
}

type SoftDelete struct {
	DeleteTime gorm.DeletedAt `gorm:"index"`
}

type Version struct {
	Version uint64 `gorm:"not null;default:1"`
}

func (Common) IsFieldOutputOnly(field string) bool {
	list := [...]string{
		"id",
		"create_time",
		"update_time",
		"delete_time",
		"version",
	}

	for _, curr := range list {
		if curr == field {
			return true
		}
	}

	return false
}

func (v Version) ResourceVersion() uint64 {
	return v.Version
}
