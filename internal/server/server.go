package server

import (
	"context"
	"fmt"
	"log/slog"
	"mime/multipart"
	"net/http"

	"github.com/gorilla/websocket"
)

type AudioQueue interface {
	PushAudio(ctx context.Context, file multipart.File, header *multipart.FileHeader) error
}

type AudioProcessor interface {
	Run()
	AddSubscriber(subscriber *websocket.Conn)
}

type Server struct {
	audioQueue     AudioQueue
	AudioProcessor AudioProcessor
}

func NewServer(queue AudioQueue, audioProcessor AudioProcessor) *Server {
	return &Server{
		audioQueue:     queue,
		AudioProcessor: audioProcessor,
	}
}
func (s *Server) Run(port string) error {
	const errorMsg = "from server Run error: %w"

	http.HandleFunc("/", s.Index)
	http.HandleFunc("/upload", s.uploadAudioHandler)
	http.HandleFunc("/stream", s.streamAudioHandler)

	go s.AudioProcessor.Run()

	slog.Info("starting server", "port", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		return fmt.Errorf(errorMsg, err)
	}

	return nil
}
