package model

import (
	"time"
)

type User struct {
    ID           uint   `json:"id" gorm:"primaryKey"`
    Username     string `json:"username" gorm:"unique;not null"`
    Email        string `json:"email" gorm:"unique;not null"`
    Password     string `json:"-" gorm:"not null"` // 不返回密碼
    RefreshToken string `json:"-" gorm:"type:text"` // 存 refresh token
    CreatedAt    time.Time
    UpdatedAt    time.Time
}


type RefreshToken struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
    UserID    uint      `json:"user_id"`
    Token     string    `json:"token" gorm:"type:text;not null"`
    ExpiresAt time.Time `json:"expires_at"`
    CreatedAt time.Time
    User      User      `json:"user" gorm:"foreignKey:UserID"`
}