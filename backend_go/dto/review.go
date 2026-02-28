package dto

// Innbound DTO for code request
type ReviewRequest struct {
	UserID uint   `json:"user_id"`
	Code   string `json:"code"`
}
