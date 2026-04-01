package services

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChatService_CreateAndListConversations(t *testing.T) {
	setup := SetUpTest(t)
	db := setup.DB

	logger := newTestLogger()
	service := NewChatService(db, logger)

	conv, err := service.CreateConversation(1, "My First Chat")
	require.NoError(t, err)
	require.Equal(t, uint(1), conv.UserID)
	require.Equal(t, "My First Chat", conv.Title)

	convs, err := service.ListConversations(1)
	require.NoError(t, err)
	require.Len(t, convs, 1)
	require.Equal(t, "My First Chat", convs[0].Title)
}

func TestChatService_CreateAndListMessages(t *testing.T) {
	setup := SetUpTest(t)
	db := setup.DB

	logger := newTestLogger()
	service := NewChatService(db, logger)

	conv, _ := service.CreateConversation(1, "Test Chat")

	msg, err := service.CreateMessage(1, conv.ID, "user", "Hello!")
	require.NoError(t, err)
	require.Equal(t, "Hello!", msg.Content)

	msgs, err := service.ListMessages(1, conv.ID, 10, 0)
	require.NoError(t, err)
	require.Len(t, msgs, 1)
	require.Equal(t, "Hello!", msgs[0].Content)
}

func TestChatService_ArchiveRenameDeleteConversation(t *testing.T) {
	setup := SetUpTest(t)
	db := setup.DB

	logger := newTestLogger()
	service := NewChatService(db, logger)

	conv, _ := service.CreateConversation(1, "Chat To Modify")

	// Archive
	require.NoError(t, service.ArchiveConversation(1, conv.ID))
	archived, _ := service.ListConversations(1)
	require.Len(t, archived, 0) // archived conversations excluded

	// Rename
	conv2, _ := service.CreateConversation(1, "Chat To Rename")
	require.NoError(t, service.RenameConversation(1, conv2.ID, "New Title"))
	renamed, _ := service.ListConversations(1)
	require.Equal(t, "New Title", renamed[0].Title)

	// Delete
	require.NoError(t, service.DeleteConversation(1, conv2.ID))
	final, _ := service.ListConversations(1)
	require.Len(t, final, 0)
}
