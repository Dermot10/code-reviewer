package dto

type ReviewTask struct {
	UserID   uint   `json:"user_id"`
	ReviewID uint   `json:"review_id"`
	Action   string `json:"action"`
}

type EnhanceTask struct {
	UserID        uint   `json:"user_id"`
	EnhancementID uint   `json:"enhancement_id"`
	Action        string `json:"action"`
}
