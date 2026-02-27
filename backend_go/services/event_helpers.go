package services

import (
	"encoding/json"

	"github.com/dermot10/code-reviewer/backend_go/dto"
	"github.com/dermot10/code-reviewer/backend_go/websocket"
)

func (s *ReviewService) emitReviewStarted(userID, reviewID uint) {
	payload := dto.ReviewStartedPayload{
		ReviewID: reviewID,
		Status:   "pending",
	}

	payloadJSON, _ := json.Marshal(payload)

	event := dto.WSEvent{
		Type:    dto.EventReviewStarted,
		Payload: payloadJSON,
	}

	eventData, _ := json.Marshal(event)

	s.wsHub.Broadcast(websocket.Message{
		UserID: userID,
		Data:   eventData,
	})
}

func (s *ReviewService) emitReviewCompleted(userID, reviewID uint, result string) {
	payload := dto.ReviewCompletedPayload{
		ReviewID: reviewID,
		Status:   "completed",
		Result:   result,
	}

	payloadJSON, _ := json.Marshal(payload)

	event := dto.WSEvent{
		Type:    dto.EventReviewCompleted,
		Payload: payloadJSON,
	}

	eventData, _ := json.Marshal(event)

	s.wsHub.Broadcast(websocket.Message{
		UserID: userID,
		Data:   eventData,
	})
}

func (a *AssistantService) emitAssistantChunk(userID, conversationID uint, chunk string) {}
