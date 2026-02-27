package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Username       string    `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Email          string    `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Fullname       string    `gorm:"size:100" json:"fullname,omitempty"`
	HashedPassword string    `gorm:"size:256;not null" json:"-"`
	IsActive       bool      `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Review struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	Code      string    `gorm:"type:text;not null" json:"code"`
	Status    string    `gorm:"type:varchar(20);default:'pending'" json:"status"`
	Result    string    `gorm:"type:text;not null" json:"result,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type Enhancement struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"index;not null" json:"user_id"`
	Original     string    `gorm:"type:text;not null" json:"original"`
	Enhanced     string    `gorm:"type:text;not null" json:"enhanced"`
	Status       string    `gorm:"type:varchar(20);default:'pending'" json:"status"`
	ModelVersion string    `gorm:"size:50" json:"model_version,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type File struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"index; not null"`
	Name      string `gorm:"not null"`
	Path      string `gorm:"not null"`
	Content   string `gorm:"type:text"`
	ProjectID *uint  `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ChatMessage struct {
	ID             uint           `gorm:"primaryKey"`
	ConversationID uint           `gorm:"index:idx_conv_created"`
	Role           string         // for system level roles (user or assistant)
	Content        string         `gorm:"type:text"`
	TokenCount     int            // for billing/tracking
	CreatedAt      time.Time      `gorm:"index:idx_conv_created"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type Conversation struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"index:idx_user_updated"`
	Title     string
	Archived  bool `gorm:"default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time      `gorm:"index:idx_user_updated"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
