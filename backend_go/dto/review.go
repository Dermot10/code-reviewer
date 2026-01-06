package dto

type ReviewRequest struct {
	UserID uint   `json:"user_id"`
	Code   string `json:"code"`
}
