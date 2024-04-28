package service

import (
	"encoding/json"
	"sync"
	"time"

	"chat-service/internal/adapter"
	"chat-service/internal/domain"
	"chat-service/pkg/validate"

	"github.com/IBM/sarama"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type Chat struct {
	producer      *adapter.Producer
	storageClient *adapter.StorageClient

	activeUsers      map[*websocket.Conn]string // connection -> nickname
	activeUsersMutex sync.RWMutex

	messagesToSend chan domain.Message
}

func NewChat(kafkaBrokers []string, kafkaTopic string, storageUrl string) *Chat {
	return &Chat{
		producer:         adapter.NewKafkaProducer(kafkaBrokers, kafkaTopic),
		storageClient:    adapter.NewStorageClient(storageUrl),
		activeUsers:      make(map[*websocket.Conn]string),
		activeUsersMutex: sync.RWMutex{},
		messagesToSend:   make(chan domain.Message),
	}
}

func (c *Chat) DistributeMessages() {
	for msg := range c.messagesToSend {
		c.activeUsersMutex.Lock()
		for conn := range c.activeUsers {
			_ = conn.WriteMessage(websocket.TextMessage, []byte(msg.String()))
		}
		c.activeUsersMutex.Unlock()
	}
}

func (c *Chat) AddUser(conn *websocket.Conn, nickname string) error {
	if err := validate.Name(nickname); err != nil {
		return err
	}
	log.Infof("User %s connected", nickname)

	c.addToActiveUsersPool(conn, nickname)

	return nil
}

func (c *Chat) RemoveUser(conn *websocket.Conn) {
	c.activeUsersMutex.RLock()
	delete(c.activeUsers, conn)
	c.activeUsersMutex.RUnlock()

	_ = conn.Close()
}

func (c *Chat) SendMessage(m domain.Message) {
	c.messagesToSend <- m

	messageBytes, err := json.Marshal(m)
	if err != nil {
		log.Errorln("Failed to marshal message: ", err)
		return
	}

	c.producer.Input() <- &sarama.ProducerMessage{
		Topic: c.producer.TopicName,
		Value: sarama.ByteEncoder(messageBytes),
	}
}

func (c *Chat) ReadAndStoreMessage(conn *websocket.Conn) {
	for {
		mt, messageText, err := conn.ReadMessage()
		if err != nil || mt == websocket.CloseMessage {
			break
		}
		if err := validate.Message(string(messageText)); err != nil {
			_ = conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			continue
		}

		nickname := c.getNickname(conn)
		message := domain.Message{
			Sender: nickname,
			Text:   string(messageText),
			SendAt: time.Now().UTC(),
		}

		c.SendMessage(message)
	}
}

func (c *Chat) ShowLastMessages(conn *websocket.Conn) error {
	messages, err := c.storageClient.FetchMessageHistory()
	if err != nil {
		return err
	}

	for _, m := range messages {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(m.String())); err != nil {
			return err
		}
	}

	return nil
}

func (c *Chat) addToActiveUsersPool(conn *websocket.Conn, nickname string) {
	c.activeUsersMutex.RLock()
	c.activeUsers[conn] = nickname
	c.activeUsersMutex.RUnlock()
}

func (c *Chat) getNickname(conn *websocket.Conn) string {
	c.activeUsersMutex.RLock()
	n := c.activeUsers[conn]
	c.activeUsersMutex.RUnlock()

	return n
}
