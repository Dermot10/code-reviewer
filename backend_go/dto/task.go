package dto

type ReviewTask struct {
	Type     string `json:"type"`
	UserID   uint   `json:"user_id"`
	ReviewID uint   `json:"review_id"`
	Action   string `json:"action"`
}

type EnhanceTask struct {
	Type          string `json:"type"`
	UserID        uint   `json:"user_id"`
	EnhancementID uint   `json:"enhancement_id"`
	Action        string `json:"action"`
}
