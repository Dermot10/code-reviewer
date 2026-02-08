package websocket

type MessageType string

const (
	TypeReviewCompleted MessageType = "review_completed"
	TypeReviewFailed    MessageType = "review_failed"
	TypeFileUpdated     MessageType = "file_updated"
)

// Client Msg
type Message struct {
	Type    MessageType `json:"type"`
	UserID  uint        `json:"-"` // for routing info only
	Payload interface{} `json:"payload"`
}

// Payload
type ReviewCompletedPayload struct {
	ReviewID uint   `json:"review_id"`
	Status   string `json:"status"`
	Result   string `json:"result"`
}

type FileUpdatedPayload struct {
	FileID  uint   `json:"file_id"`
	Content string `json:"content"`
}
