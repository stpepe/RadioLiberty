package audio_processor

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAudioProcessor(t *testing.T) {
	testChan := make(chan []byte)
	expectedResult := &AudioProcessor{
		audioReadChannel: testChan,
		audioSyncer:      sync.NewCond(&sync.Mutex{}),
	}

	realResult := NewAudioProcessor(testChan)

	assert.Equal(t, expectedResult, realResult)
}
