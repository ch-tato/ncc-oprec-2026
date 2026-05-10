package models

import "time"

// Note represents a single note entry in the database.
type Note struct {
	ID        uint      `gorm:"primaryKey"`
	Title     string    `gorm:"size:200;not null"`
	Content   string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
