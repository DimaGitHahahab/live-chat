package handler

import (
	"net/http"

	"chat-service/internal/service"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type ChatHandler struct {
	chatService *service.Chat
	mux         *http.ServeMux
}

func NewChatHandler(chat *service.Chat) *ChatHandler {
	mux := http.NewServeMux()
	handler := &ChatHandler{
		chatService: chat,
		mux:         mux,
	}
	mux.HandleFunc("/", handler.manageConnection)
	return handler
}

func (h *ChatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Infoln("Starting to distribute incoming messages")
	go h.chatService.DistributeMessages()

	h.mux.ServeHTTP(w, r)
}

func getNickname(r *http.Request) string {
	return r.URL.Query().Get("nickname")
}

func (h *ChatHandler) manageConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgradeConnection(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer h.closeConnection(conn)

	if err := h.chatService.AddUser(conn, getNickname(r)); err != nil {
		return
	}

	if err := h.chatService.ShowLastMessages(conn); err != nil {
		log.Errorln("Unable to show last messages: ", err)
		return
	}

	h.chatService.ReadAndStoreMessage(conn)
}

func (h *ChatHandler) closeConnection(conn *websocket.Conn) {
	h.chatService.RemoveUser(conn)
}
