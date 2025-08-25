package model

import (
	"time"
)

type Vocabulary struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Vocabulary    string    `json:"vocabulary" gorm:"not null"`
	Mean    string    `json:"mean" gorm:"not null"`
	Category   string    `json:"category"`
	Difficulty int       `json:"difficulty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
