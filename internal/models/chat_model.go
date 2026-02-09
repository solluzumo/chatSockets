package models

import "time"

type Chat struct {
	ID        int    `gorm:"primaryKey;autoIncrement"`
	Title     string `gorm:"size:200;not null"`
	CreatedAt time.Time
}
