package server

import (
	"net/http"

	"log/slog"

	"github.com/gorilla/websocket"
)

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
	s.AudioProcessor.AddSubscriber(conn)
}
