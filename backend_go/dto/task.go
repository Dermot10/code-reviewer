package dto

// task dto to manage lifecycle of the request
// Redis bound

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

// internal - worker input
type AssistantTask struct {
	Type           string `json:"type"`
	UserID         uint   `json:"user_id"`
	ConversationID uint   `json:"conversation_id"`
	Prompt         string `json:"prompt"`
}

// worker output
type AssistantTaskEvent struct {
	Type           string `json:"type"` //assistant.chunk | assistant.completed
	UserID         uint   `json:"user_id"`
	ConversationID uint   `json:"conversation_id"`
	Chunk          string `json:"chunk,omitempty"`
	Content        string `json:"content,omitempty"`
}
