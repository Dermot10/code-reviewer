package dto

import "encoding/json"

type WSEventType string

const (
	EventReviewStarted   WSEventType = "review.started"
	EventReviewCompleted WSEventType = "review.completed"
	EventReviewFailed    WSEventType = "review.failed"
	EventFileUpload      WSEventType = "file.uploaded"
	EventFileUpdated     WSEventType = "file.updated"
	EventChatMessage     WSEventType = "chat.message"
)

// top level ws event
type WSEvent struct {
	Type    WSEventType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type ReviewStartedPayload struct {
	ReviewID uint   `json:"review_id"`
	Status   string `json:"status"`
}

type ReviewCompletedPayload struct {
	ReviewID uint   `json:"review_id"`
	Status   string `json:"status"`
	Result   string `json:"result"`
}

type ReviewFailedPayload struct {
	ReviewID uint   `json:"review_id"`
	Status   string `json:"status"`
	Error    string `json:"error"`
}

// may extend for more granular events corresponding to specific events
// e.g cursor move, selection, patch edits
type FileUpdatedPayload struct {
	FileID  uint   `json:"file_id"`
	Content string `json:"content"`
}

type ChatMessagePayload struct {
	From    string `json:"from"`
	Message string `json:"message"`
	Time    string `json:"time"`
}
