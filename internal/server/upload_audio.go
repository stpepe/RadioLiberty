package server

import (
	"net/http"
	"strconv"
	"time"
)

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

	err = s.audioQueue.PushAudio(r.Context(), file, header)
	if err != nil {
		http.Error(w, "Failed to upload file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
