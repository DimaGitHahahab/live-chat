package handler

import (
	"encoding/json"
	"net/http"

	"storage-service/internal/service"

	log "github.com/sirupsen/logrus"
)

type StorageHandler struct {
	storageService *service.Storage
	mux            *http.ServeMux
}

func NewChatHandler(chat *service.Storage) *StorageHandler {
	mux := http.NewServeMux()
	handler := &StorageHandler{
		storageService: chat,
		mux:            mux,
	}
	mux.HandleFunc("/", handler.getLastMessages)
	return handler
}

func (h *StorageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *StorageHandler) getLastMessages(w http.ResponseWriter, _ *http.Request) {
	messages := h.storageService.GetLastMessages()
	if len(messages) > 0 {
		w.Header().Set("Content-Type", "application/json")

		err := json.NewEncoder(w).Encode(messages)
		if err != nil {
			log.Error("Failed to encode message into json: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}
