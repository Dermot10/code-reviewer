package dto

import (
	"encoding/json"
	"time"
)

type WSEventType string

const (
	EventReviewStarted   WSEventType = "review.started"
	EventReviewCompleted WSEventType = "review.completed"
	EventReviewFailed    WSEventType = "review.failed"

	EventFileUpload  WSEventType = "file.uploaded"
	EventFileUpdated WSEventType = "file.updated"

	EventConversationCreate   WSEventType = "conversation.create"
	EventConversationCreated  WSEventType = "conversation.created"
	EventConversationArchive  WSEventType = "conversation.archive"
	EventConversationArchived WSEventType = "conversation.archived"
	EventConversationRename   WSEventType = "conversation.rename"
	EventConversationRenamed  WSEventType = "conversation.renamed"
	EventConvrsationDelete    WSEventType = "conversation.delete"
	EventConversationDeleted  WSEventType = "conversation.deleted"

	EventMessageSend    WSEventType = "message.send"
	EventMessageCreated WSEventType = "message.created"
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

// inbound
type ConversationCreatePayload struct {
	Title string `json:"title"`
}

type ConversationArhivePayload struct {
	ConversationID uint `json:"conversation_id"`
}

type ConversationRenamePayload struct {
	ConversationID uint   `json:"conversation_id"`
	Title          string `json:"title"`
}

type ConversationDeletePayload struct {
	ConversationID uint `json:"conversation_id"`
}

type MessageSendPayload struct {
	ConversationID uint   `json:"conversation_id"`
	Content        string `json:"content"`
}

// outbound
type ConversationCreatedPayload struct {
	ConversationID uint   `json:"id"`
	Title          string `json:"title"`
}

type ConversationArchivedPayload struct {
	ConversationID uint `json:"conversation_id"`
	Archived       bool `json:"archived"`
}

type ConversationRenamedPayload struct {
	ConversationID uint   `json:"conversation_id"`
	Title          string `json:"title"`
}

type ConversationDeletedPayload struct {
	ConversationID uint `json:"conversation_id"`
	Deleted        bool `json:"deleted"`
}

type MessageCreatedPayload struct {
	ID             uint      `json:"id"`
	ConversationID uint      `json:"conversation_id"`
	Role           string    `json:"role"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"created_at"`
}
