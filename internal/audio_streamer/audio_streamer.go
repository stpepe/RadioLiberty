package audio_streamer

import (
	"RadioLiberty/pkg/models"
	"context"
	"log/slog"
	"time"
)

const secondsToStream = 1

type AudioStreamer struct {
	audioQueue chan *models.AudioFile

	streamChannel chan []byte
}

func (a *AudioStreamer) stream(ctx context.Context) {
	defer close(a.streamChannel)
	defer slog.Info("AudioStream stopping")

	defaultAudioFile := <-a.audioQueue
	var audioFile *models.AudioFile
	for {
		select {
		case <-ctx.Done():
			return
		case audioFile = <-a.audioQueue:
		default:
			audioFile = defaultAudioFile
		}
		bufferSize := 65_536
		for offset := 0; offset <= len(audioFile.File); offset += bufferSize {
			rightBorder := offset + bufferSize
			if rightBorder > len(audioFile.File) {
				rightBorder = len(audioFile.File)
			}
			a.streamChannel <- audioFile.File[offset:rightBorder]
			time.Sleep(secondsToStream * time.Second)
		}
	}
}

func NewAudioStreamer(audioQueue chan *models.AudioFile) *AudioStreamer {
	return &AudioStreamer{
		audioQueue:    audioQueue,
		streamChannel: make(chan []byte),
	}
}

func (a *AudioStreamer) Run(ctx context.Context) chan []byte {
	go a.stream(ctx)

	return a.streamChannel
}
