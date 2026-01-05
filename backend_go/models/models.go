package models

import "time"

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

// type Organisation struct {
// 	ID        uint      `gorm:"primaryKey" json:"id"`
// 	Name      string    `gorm:"size:100;uniqueIndex;not null" json:"name"`
// 	OwnerID   uint      `gorm:"index;not null" json:"owner_id"`
// 	CreatedAt time.Time `json:"created_at"`
// 	UpdatedAt time.Time `json:"updated_at"`

// 	Owner    User      `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
// 	Projects []Project `gorm:"foreignKey:OrganisationID" json:"projects,omitempty"`
// }

// type Project struct {
// 	ID             uint      `gorm:"primaryKey" json:"id"`
// 	Name           string    `gorm:"size:100;not null;index:idx_org_project,unique" json:"name"`
// 	Description    string    `gorm:"type:text" json:"description,omitempty"`
// 	OwnerID        uint      `gorm:"index;not null" json:"owner_id"`
// 	OrganisationID *uint     `gorm:"index;index:idx_org_project,unique" json:"organisation_id,omitempty"`
// 	CreatedAt      time.Time `json:"created_at"`
// 	UpdatedAt      time.Time `json:"updated_at"`

// 	Owner        User          `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
// 	Organisation *Organisation `gorm:"foreignKey:OrganisationID" json:"organisation,omitempty"`
// 	Reviews      []Review      `gorm:"constraint:OnDelete:CASCADE" json:"reviews,omitempty"`
// }

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
