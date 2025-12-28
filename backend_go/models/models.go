package models

import "time"

type User struct {
	ID             int64     `db:"id" json:"id"`
	Username       string    `db:"username" json:"username"`
	Email          string    `db:"email" json:"email"`
	Fullname       string    `db:"fullname" json:"fullname,omitempty"`
	HashedPassword string    `db:"hashed_password" json:"-"`
	IsActive       bool      `db:"is_active" json:"is_active"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`

	// Relationships (optional, for eager loading)
	Organisations []Organisation `db:"-" json:"organisations,omitempty"`
	Projects      []Project      `db:"-" json:"projects,omitempty"`
}

type Organisation struct {
	ID        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	OwnerID   int64     `db:"owner_id" json:"owner_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`

	Owner    *User     `db:"-" json:"owner,omitempty"`
	Projects []Project `db:"-" json:"projects,omitempty"`
}

type Project struct {
	ID             int64     `db:"id" json:"id"`
	Name           string    `db:"name" json:"name"`
	Description    string    `db:"description,omitempty" json:"description,omitempty"`
	OwnerID        int64     `db:"owner_id" json:"owner_id"`
	OrganisationID *int64    `db:"organisation_id,omitempty" json:"organisation_id,omitempty"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`

	Owner        *User         `db:"-" json:"owner,omitempty"`
	Organisation *Organisation `db:"-" json:"organisation,omitempty"`
}

type Review struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ProjectID   uint      `gorm:"index;not null" json:"project_id"`
	ReviewerID  uint      `gorm:"index;not null" json:"reviewer_id"`
	Feedback    string    `gorm:"type:text;not null" json:"feedback"`
	IssuesCount int       `gorm:"default:0" json:"issues_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Project   Project    `gorm:"foreignKey:ProjectID"`
	Reviewer  User       `gorm:"foreignKey:ReviewerID"`
	AiResults []AiResult `gorm:"constraint:OnDelete:CASCADE"`
}

type AiResult struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ReviewID     uint      `gorm:"index;not null" json:"review_id"`
	Output       string    `gorm:"type:text;not null" json:"output"`
	ModelVersion string    `gorm:"size:50" json:"model_version,omitempty"`
	CreatedAt    time.Time `json:"created_at"`

	Review Review `gorm:"foreignKey:ReviewID"`
}
