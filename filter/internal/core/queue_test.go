package core

import (
	"testing"
)

func TestAddSpamMessage(t *testing.T) {
	spamBufferLock.Lock()
	spamBuffer = spamBuffer[:0]
	spamBufferLock.Unlock()

	AddSpamMessage([]byte("test spam"))

	spamBufferLock.Lock()
	defer spamBufferLock.Unlock()

	if len(spamBuffer) != 1 {
		t.Errorf("Expected 1 message in spamBuffer, got %d", len(spamBuffer))
	}
}