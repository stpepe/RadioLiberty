package queue

import (
	"RadioLiberty/pkg/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"mime/multipart"
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
	queueStorage      QueueStorage
	audioFilesStorage AudioFilesStorage
}

func (q *Queue) GetNextAudioInfo(ctx context.Context) (*models.AudioInfo, error) {
	const errorMsg = "from queue GetNextAudioInfo error: %w"
	audioInfo, err := q.queueStorage.Next()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf(errorMsg, err)
	}
	return audioInfo, nil
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

func NewQueue(queueStorage QueueStorage, audioFilesStorage AudioFilesStorage) *Queue {
	return &Queue{
		queueStorage:      queueStorage,
		audioFilesStorage: audioFilesStorage,
	}
}
