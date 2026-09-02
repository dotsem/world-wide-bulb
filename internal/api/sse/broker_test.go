package sse_test

import (
	"testing"
	"world-wide-bulb/internal/api/sse"

	"github.com/stretchr/testify/assert"
)

func TestBroker_SubscribeAndBroadcast(t *testing.T) {
	broker := sse.NewBroker()
	ch := broker.Subscribe()
	defer broker.Unsubscribe(ch)

	msg := []byte(`{"event":"state_changed","state":true}`)
	broker.Broadcast(msg)

	received := <-ch
	assert.Equal(t, msg, received)
}

func TestBroker_UnsubscribeClosesChannel(t *testing.T) {
	broker := sse.NewBroker()
	ch := broker.Subscribe()

	broker.Unsubscribe(ch)

	_, ok := <-ch
	assert.False(t, ok, "channel should be closed after unsubscribe")
}

func TestBroker_SlowClientDrop(t *testing.T) {
	broker := sse.NewBroker()
	ch := broker.Subscribe()
	defer broker.Unsubscribe(ch)

	// Fill the 16-element buffer + 1 overflow
	for range 20 {
		broker.Broadcast([]byte("ping"))
	}

	assert.Len(t, ch, 16)
}
