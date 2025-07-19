package model

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `gorm:"created_at" json:"-"`
	UpdatedAt time.Time      `gorm:"updated_at" json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
