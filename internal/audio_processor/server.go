package audio_processor

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	audioProcessor AudioProcessorService
}

type AudioProcessorService interface {
	AddSubscriber(subscriber *websocket.Conn)
}

func NewServer(audioProcessor AudioProcessorService) *Server {
	return &Server{
		audioProcessor: audioProcessor,
	}
}
func (s *Server) Run(port string) error {
	const errorMsg = "from server Run error: %w"

	http.HandleFunc("/stream", s.streamAudioHandler)

	slog.Info("starting server", "port", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		return fmt.Errorf(errorMsg, err)
	}

	return nil
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *Server) streamAudioHandler(w http.ResponseWriter, r *http.Request) {
	const errorMsg = "from streamAudioHandler error"
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error(errorMsg, "error", err)
		return
	}
	s.audioProcessor.AddSubscriber(conn)
}
