package queue

import (
	"RadioLiberty/pkg/models"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	queue QueueService
}

type QueueService interface {
	PushAudio(ctx context.Context, file multipart.File, header *multipart.FileHeader) error
	GetNextAudioInfo(ctx context.Context) (*models.AudioInfo, error)
}

func NewServer(queue QueueService) *Server {
	return &Server{
		queue: queue,
	}
}
func (s *Server) Run(port string) error {
	const errorMsg = "from server Run error: %w"

	http.HandleFunc("/upload", s.uploadAudioHandler)
	http.HandleFunc("/next", s.getNextAudioInfoHandler)

	slog.Info("starting server", "port", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		return fmt.Errorf(errorMsg, err)
	}

	return nil
}

func (s *Server) uploadAudioHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to parse form file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Добавляется время в формате Unix для уникальности имени
	header.Filename = strconv.FormatInt(time.Now().Unix(), 10) + "_" + header.Filename

	err = s.queue.PushAudio(r.Context(), file, header)
	if err != nil {
		http.Error(w, "Failed to upload file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) getNextAudioInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	audioInfo, err := s.queue.GetNextAudioInfo(r.Context())
	if err != nil {
		http.Error(w, "Failed to get next audio info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if audioInfo == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	audioInfoJson, err := json.Marshal(audioInfo)
	if err != nil {
		http.Error(w, "Failed convert to json audio info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(audioInfoJson)
}
