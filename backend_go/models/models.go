package models

import "time"

// -------------------- USERS --------------------

type User struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Username       string    `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Email          string    `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Fullname       string    `gorm:"size:100" json:"fullname,omitempty"`
	HashedPassword string    `gorm:"size:256;not null" json:"-"`
	IsActive       bool      `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	Organisations []Organisation `gorm:"foreignKey:OwnerID" json:"organisations,omitempty"`
	Projects      []Project      `gorm:"foreignKey:OwnerID" json:"projects,omitempty"`
}

// -------------------- ORGANISATIONS --------------------

type Organisation struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;uniqueIndex;not null" json:"name"`
	OwnerID   uint      `gorm:"index;not null" json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Owner    User      `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Projects []Project `gorm:"foreignKey:OrganisationID" json:"projects,omitempty"`
}

// -------------------- PROJECTS --------------------

type Project struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Name           string    `gorm:"size:100;not null;index:idx_org_project,unique" json:"name"`
	Description    string    `gorm:"type:text" json:"description,omitempty"`
	OwnerID        uint      `gorm:"index;not null" json:"owner_id"`
	OrganisationID *uint     `gorm:"index;index:idx_org_project,unique" json:"organisation_id,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	Owner        User          `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Organisation *Organisation `gorm:"foreignKey:OrganisationID" json:"organisation,omitempty"`
	Reviews      []Review      `gorm:"constraint:OnDelete:CASCADE" json:"reviews,omitempty"`
}

// -------------------- REVIEWS --------------------

type Review struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ProjectID   uint      `gorm:"index;not null" json:"project_id"`
	ReviewerID  uint      `gorm:"index;not null" json:"reviewer_id"`
	Feedback    string    `gorm:"type:text;not null" json:"feedback"`
	IssuesCount int       `gorm:"default:0" json:"issues_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Project   Project    `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	Reviewer  User       `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
	AiResults []AiResult `gorm:"constraint:OnDelete:CASCADE" json:"ai_results,omitempty"`
}

// -------------------- AI RESULTS --------------------

type AiResult struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ReviewID     uint      `gorm:"index;not null" json:"review_id"`
	Output       string    `gorm:"type:text;not null" json:"output"`
	ModelVersion string    `gorm:"size:50" json:"model_version,omitempty"`
	CreatedAt    time.Time `json:"created_at"`

	Review Review `gorm:"foreignKey:ReviewID" json:"review,omitempty"`
}
