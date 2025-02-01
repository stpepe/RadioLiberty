package queue

import (
	"RadioLiberty/internal/queue/mock_queue"
	"RadioLiberty/pkg/models"
	"context"
	"errors"
	"mime/multipart"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueueStorage := mock_queue.NewMockQueuStorage(ctrl)
	mockAudioStorage := mock_queue.NewMockAudioFilesStorage(ctrl)

	audioInfo := &models.AudioInfo{ID: 1, FileName: "test.mp3"}
	audioFile := &models.AudioFile{File: []byte("mock audio data")}

	mockQueueStorage.EXPECT().Next().Return(audioInfo, nil).AnyTimes() // Вернет аудио
	mockAudioStorage.EXPECT().GetAudioFile(gomock.Any(), audioInfo, true).Return(audioFile, nil).AnyTimes()

	q := &Queue{
		QueueChannel:      make(chan *models.AudioFile, 1),
		queueStorage:      mockQueueStorage,
		audioFilesStorage: mockAudioStorage,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	go q.Run(ctx)

	select {
	case result := <-q.QueueChannel:
		assert.Equal(t, audioFile, result)
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for queue output")
	}
}

func TestPushAudio(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAudioStorage := mock_queue.NewMockAudioFilesStorage(ctrl)
	mockQueueStorage := mock_queue.NewMockQueuStorage(ctrl)

	ctx := context.Background()
	var testFile multipart.File
	testHeader := &multipart.FileHeader{
		Filename: "test.mp3",
		Size:     64,
	}

	q := &Queue{
		audioFilesStorage: mockAudioStorage,
		queueStorage:      mockQueueStorage,
	}

	t.Run("Success case", func(t *testing.T) {
		mockAudioStorage.EXPECT().PutObject(ctx, testFile, testHeader).Return(nil).Times(1)
		mockQueueStorage.EXPECT().AddToQueue(gomock.Any()).Return(nil).Times(1)

		err := q.PushAudio(ctx, testFile, testHeader)
		assert.NoError(t, err)
	})

	t.Run("Failure in PutObject", func(t *testing.T) {
		mockAudioStorage.EXPECT().PutObject(ctx, testFile, testHeader).Return(errors.New("upload error")).Times(1)

		err := q.PushAudio(ctx, testFile, testHeader)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "upload error")
	})

	t.Run("Failure in AddToQueue", func(t *testing.T) {
		mockAudioStorage.EXPECT().PutObject(ctx, testFile, testHeader).Return(nil).Times(1)
		mockQueueStorage.EXPECT().AddToQueue(gomock.Any()).Return(errors.New("queue error")).Times(1)

		err := q.PushAudio(ctx, testFile, testHeader)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "queue error")
	})
}

func TestNewQueue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAudioStorage := mock_queue.NewMockAudioFilesStorage(ctrl)
	mockQueueStorage := mock_queue.NewMockQueuStorage(ctrl)

	ctx := context.Background()

	t.Run("Success case", func(t *testing.T) {
		expectedAudioFile := &models.AudioFile{
			Info: models.AudioInfo{
				FileName:  "default.mp3",
				AudioName: "default",
				Artist:    "default",
			},
			File: nil,
		}
		queue := make(chan *models.AudioFile, 1)
		queue <- expectedAudioFile

		mockAudioStorage.EXPECT().GetAudioFile(ctx, gomock.Any(), gomock.Any()).Return(expectedAudioFile, nil).Times(1)

		expectedResult := &Queue{
			QueueChannel:      queue,
			queueStorage:      mockQueueStorage,
			audioFilesStorage: mockAudioStorage,
		}

		realResult, err := NewQueue(mockQueueStorage, mockAudioStorage)
		realAudioFile := <-realResult.QueueChannel

		assert.NoError(t, err)
		assert.EqualValues(t, expectedResult.audioFilesStorage, realResult.audioFilesStorage)
		assert.EqualValues(t, expectedResult.queueStorage, realResult.queueStorage)
		assert.EqualValues(t, expectedResult.audioFilesStorage, realResult.audioFilesStorage)
		assert.EqualValues(t, expectedAudioFile, realAudioFile)
	})

	t.Run("Failure case in GetAudioFile", func(t *testing.T) {
		mockAudioStorage.EXPECT().GetAudioFile(ctx, gomock.Any(), gomock.Any()).Return(nil, errors.New("get audio file error")).Times(1)

		_, err := NewQueue(mockQueueStorage, mockAudioStorage)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "get audio file error")
	})
}
