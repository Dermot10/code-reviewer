package dto

// Frontend DTO for the actual code request
type ReviewRequest struct {
	UserID uint   `json:"user_id"`
	Code   string `json:"code"`
}
