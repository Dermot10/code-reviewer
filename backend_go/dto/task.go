package dto

// task dto to manage lifecycle of the request
// Redis bound
// db bound

type ReviewTask struct {
	Type     string `json:"type"`
	UserID   uint   `json:"user_id"`
	ReviewID uint   `json:"review_id"`
	Code     string `json:"code"`
	Action   string `json:"action"`
}

type EnhanceTask struct {
	Type          string `json:"type"`
	UserID        uint   `json:"user_id"`
	EnhancementID uint   `json:"enhancement_id"`
	Code          string `json:"code"`
	Action        string `json:"action"`
}
