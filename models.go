package main

import (
	"time"
)

type Model struct {
	ID        uint `gorm:"type:bigserial;primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	Model        `json:"-"`
	Email        string `json:"email"`
	Username     string `json:"username"`
	Fullname     string `json:"fullname"`
	PasswordHash string `json:"-"`
	PasswordSalt string `json:"-"`
	Role         string `json:"role"`
	IsDisabled   bool   `json:"isdisabled"`
}
