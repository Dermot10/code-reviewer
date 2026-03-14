package websocket

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func newTestClient(userID uint) *Client {
	return &Client{
		UserID: userID,
		Send:   make(chan []byte, 1),
	}
}

func TestHub_RegisterClient(t *testing.T) {
	// client is correctly added to the userMap

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hub := NewHub()
	go hub.Run(ctx)

	client := newTestClient(1)

	hub.Register <- client

	time.Sleep(50 * time.Millisecond)

	hub.mu.RLock()
	defer hub.mu.RUnlock()

	require.Len(t, hub.clients[1], 1)
	require.True(t, hub.clients[1][client])
}

func TestHub_UnregisterClient(t *testing.T) {
	// client is removed and clean up

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hub := NewHub()
	go hub.Run(ctx)

	client := newTestClient(1)

	hub.Register <- client
	time.Sleep(20 * time.Millisecond)

	hub.Unregister <- client
	time.Sleep(20 * time.Millisecond)

	hub.mu.RLock()
	defer hub.mu.RUnlock()

	_, exists := hub.clients[1]
	require.False(t, exists)
}

func TestHub_BroadcastToUser(t *testing.T) {
	// broadcast is delivered

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hub := NewHub()
	go hub.Run(ctx)

	client := newTestClient(42)
	hub.Register <- client

	time.Sleep(20 * time.Millisecond)

	msg := []byte(`{"type":""test}`)
	hub.Broadcast(Message{
		UserID: 42,
		Data:   msg,
	})

	select {
	case received := <-client.Send:
		require.Equal(t, msg, received)
	case <-time.After(time.Second):
		t.Fatal("message not received")
	}
}

func TestHub_BroadcastOnlyToCorrectUser(t *testing.T) {
	//  routing works

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hub := NewHub()
	go hub.Run(ctx)

	client1 := newTestClient(1)
	client2 := newTestClient(2)

	hub.Register <- client1
	hub.Register <- client2

	time.Sleep(20 * time.Millisecond)

	msg := []byte(`hello`)

	hub.Broadcast(Message{
		UserID: 1,
		Data:   msg,
	})

	select {
	case <-client1.Send:
		// expected
	case <-time.After(time.Second):
		t.Fatal("client1 did not receive message")
	}

	select {
	case <-client2.Send:
		t.Fatalf("client2 should not receive message")
	case <-time.After(200 * time.Millisecond):
		// expected
	}

}
