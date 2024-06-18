package models

import (
	"time"
)

type Token struct {
	BaseModel
	UserID    uint      `gorm:"<-:create;not null"`
	Value     string    `gorm:"<-:create;unique;index;not null"`
	ExpiredAt time.Time `gorm:"<-:create;not null"`
}

type SessionToken struct {
	Token
}

type RefreshToken struct {
	Token
}
