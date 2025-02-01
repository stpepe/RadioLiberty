package audio_processor

import (
	"log/slog"
	"sync"

	"github.com/gorilla/websocket"
)

type AudioProcessor struct {
	audioReadChannel chan []byte
	currentAudioPart []byte
	audioSyncer      *sync.Cond
}

func NewAudioProcessor(audioReadChannel chan []byte) *AudioProcessor {
	return &AudioProcessor{
		audioReadChannel: audioReadChannel,
		audioSyncer:      sync.NewCond(&sync.Mutex{}),
	}
}

func (a *AudioProcessor) Run() {
	defer slog.Info("AudioProcessor stopped")

	for audioPart := range a.audioReadChannel {
		a.currentAudioPart = audioPart
		a.audioSyncer.Broadcast()
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
