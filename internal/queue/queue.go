package queue

import (
	"RadioLiberty/pkg/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"mime/multipart"
	"time"
)

type QueueStorage interface {
	AddToQueue(audio *models.AudioInfo) error
	Next() (*models.AudioInfo, error)
}

type AudioFilesStorage interface {
	GetAudioFile(ctx context.Context, audioInfo *models.AudioInfo, removeFlag bool) (*models.AudioFile, error)
	PutObject(ctx context.Context, file multipart.File, header *multipart.FileHeader) error
}

type Queue struct {
	QueueChannel      chan *models.AudioFile
	queueStorage      QueueStorage
	audioFilesStorage AudioFilesStorage
}

func (q *Queue) Run(ctx context.Context) {
	const errorMsg = "from queue Run"

	defer close(q.QueueChannel)
	defer log.Println("Queue service stopping")
	for {
		select {
		case <-ctx.Done():
			return
		default:
			audioInfo, err := q.queueStorage.Next()
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					time.Sleep(5 * time.Second)
					continue
				} else {
					slog.Error(errorMsg, "error", err)
					continue
				}
			}

			audioFile, err := q.audioFilesStorage.GetAudioFile(ctx, audioInfo, true)
			if err != nil {
				slog.Error(errorMsg, "error", err)
				continue
			}

			q.QueueChannel <- audioFile
		}
	}
}

func (q *Queue) PushAudio(ctx context.Context, file multipart.File, header *multipart.FileHeader) error {
	const errorMsg = "from queue PushAudio error: %w"

	audioInfo := &models.AudioInfo{
		FileName:  header.Filename,
		AudioName: header.Filename,
		Artist:    header.Filename,
	}

	err := q.audioFilesStorage.PutObject(ctx, file, header)
	if err != nil {
		return fmt.Errorf(errorMsg, err)
	}
	err = q.queueStorage.AddToQueue(audioInfo)
	if err != nil {
		return fmt.Errorf(errorMsg, err)
	}

	return nil
}

func NewQueue(queueStorage QueueStorage, audioFilesStorage AudioFilesStorage) (*Queue, error) {
	const errorMsg = "from NewQueue error: %w"

	defaultAudioInfo := &models.AudioInfo{
		FileName:  "default.mp3",
		AudioName: "default",
		Artist:    "default",
	}
	audioFile, err := audioFilesStorage.GetAudioFile(context.Background(), defaultAudioInfo, false)
	if err != nil {
		return nil, fmt.Errorf(errorMsg, err)
	}
	queue := make(chan *models.AudioFile, 1)
	queue <- audioFile
	return &Queue{
		QueueChannel:      queue,
		queueStorage:      queueStorage,
		audioFilesStorage: audioFilesStorage,
	}, nil
}
