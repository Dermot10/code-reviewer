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

	EventAssistantPrompt WSEventType = "assistant.prompt"
	EventAssistantStream WSEventType = "assistant.stream"

	EventConversationCreate  WSEventType = "conversation.create"
	EventConversationCreated WSEventType = "conversation.created"

	EventConversationArchive  WSEventType = "conversation.archive"
	EventConversationArchived WSEventType = "conversation.archived"

	EventConversationRename  WSEventType = "conversation.rename"
	EventConversationRenamed WSEventType = "conversation.renamed"

	EventConvrsationDelete   WSEventType = "conversation.delete"
	EventConversationDeleted WSEventType = "conversation.deleted"

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

// // inbound
type FileUpdatedPayload struct {
	FileID  uint   `json:"file_id"`
	Content string `json:"content"`
}

type PromptPayload struct {
	ConversationID uint   `json:"conversation_id"`
	Prompt         string `json:"prompt"`
}

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
type AssistantStreamPayload struct {
	ConversationID uint   `json:"conversation_id"`
	Chunk          string `json:"chunk"`
	Done           bool   `json:"done"` //end of stream flag
}

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
