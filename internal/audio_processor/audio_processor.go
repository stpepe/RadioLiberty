package audio_processor

import (
	"RadioLiberty/pkg/models"
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const secondsToStream = 1

type QueueClient interface {
	GetNextAudioInfo(ctx context.Context) (*models.AudioInfo, error)
}

type AudioFilesStorage interface {
	GetAudioFile(ctx context.Context, audioInfo *models.AudioInfo, removeFlag bool) (*models.AudioFile, error)
}

type AudioProcessor struct {
	queueClient       QueueClient
	audioFilesStorage AudioFilesStorage
	defaultAudioName  string
	currentAudioPart  []byte
	audioSyncer       *sync.Cond
}

func NewAudioProcessor(queueClient QueueClient, audioFilesStorage AudioFilesStorage, defaultAudioName string) *AudioProcessor {
	return &AudioProcessor{
		queueClient:       queueClient,
		defaultAudioName:  defaultAudioName,
		audioFilesStorage: audioFilesStorage,
		audioSyncer:       sync.NewCond(&sync.Mutex{}),
	}
}

func (a *AudioProcessor) GetNextAudio(ctx context.Context) (*models.AudioFile, error) {
	audioFileInfo, err := a.queueClient.GetNextAudioInfo(ctx)
	if err != nil {
		return nil, err
	}
	var removeFlag bool = true
	if audioFileInfo == nil {
		audioFileInfo = &models.AudioInfo{
			FileName:  a.defaultAudioName,
			AudioName: a.defaultAudioName,
			Artist:    a.defaultAudioName,
		}
		removeFlag = false
	}
	audioFile, err := a.audioFilesStorage.GetAudioFile(ctx, audioFileInfo, removeFlag)
	if err != nil {
		return nil, err
	}
	return audioFile, nil
}

func (a *AudioProcessor) Stream(ctx context.Context) {
	defer slog.Info("AudioStream stopping")
	for {
		select {
		case <-ctx.Done():
			return
		default:
			audioFile, err := a.GetNextAudio(ctx)
			if err != nil {
				slog.Error("from GetNextAudio error", "error", err)
				time.Sleep(10 * time.Second)
				continue
			}
			bufferSize := 65_536
			for offset := 0; offset <= len(audioFile.File); offset += bufferSize {
				rightBorder := offset + bufferSize
				if rightBorder > len(audioFile.File) {
					rightBorder = len(audioFile.File)
				}
				a.currentAudioPart = audioFile.File[offset:rightBorder]
				a.audioSyncer.Broadcast()
				time.Sleep(secondsToStream * time.Second)
			}
		}
	}
}

func (a *AudioProcessor) AddSubscriber(subscriber *websocket.Conn) {
	defer subscriber.Close()
	func() {
		for {
			a.audioSyncer.L.Lock()
			a.audioSyncer.Wait()

			err := subscriber.WriteMessage(websocket.BinaryMessage, a.currentAudioPart)
			a.audioSyncer.L.Unlock()
			if err != nil {
				slog.Error("send audio part to user through websocket", "error", err)
				return
			}
		}
	}()
}
