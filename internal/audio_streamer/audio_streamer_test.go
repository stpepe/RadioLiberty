package audio_streamer

import (
	"RadioLiberty/pkg/models"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewAudioStreamer(t *testing.T) {
	testChan := make(chan *models.AudioFile)
	expectedResult := &AudioStreamer{
		audioQueue:    testChan,
		streamChannel: make(chan []byte),
	}

	realResult := NewAudioStreamer(testChan)

	assert.Equal(t, expectedResult.audioQueue, realResult.audioQueue)
	assert.Equal(t, len(expectedResult.streamChannel), len(realResult.streamChannel))
}

func TestAudioStreamer_Stream(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	audioQueue := make(chan *models.AudioFile, 1)
	streamChannel := make(chan []byte, 10)

	streamer := &AudioStreamer{
		audioQueue:    audioQueue,
		streamChannel: streamChannel,
	}

	testData := make([]byte, 200_000)
	for i := range testData {
		testData[i] = byte(i % 256) // Заполняем тестовыми байтами
	}

	audioFile := &models.AudioFile{File: testData}
	audioQueue <- audioFile

	go streamer.stream(ctx)

	var receivedData []byte
	timeout := time.After(5 * time.Second)
	for {
		select {
		case chunk, ok := <-streamChannel:
			if !ok {
				t.Fatal("streamChannel closed unexpectedly")
			}
			receivedData = append(receivedData, chunk...)
			if len(receivedData) >= len(testData) {
				cancel()
				return
			}
		case <-timeout:
			t.Fatal("Test timeout: did not receive expected data")
		}
	}
}

func TestRun(t *testing.T) {
	ctx := context.Background()
	audioStreamer := &AudioStreamer{
		audioQueue:    make(chan *models.AudioFile),
		streamChannel: make(chan []byte),
	}

	result := audioStreamer.Run(ctx)

	assert.Equal(t, result, audioStreamer.streamChannel)
}
