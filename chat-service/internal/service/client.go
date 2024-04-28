package service

import (
	"bufio"
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
}

func NewClient(u url.URL, nickname string) (*Client, error) {
	q := u.Query()
	q.Set("nickname", nickname)
	u.RawQuery = q.Encode()
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}

	c := Client{conn: conn}

	return &c, nil
}

func (c *Client) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	reader := bufio.NewReader(os.Stdin)
	go func() {
		defer cancel()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				read, err := reader.ReadString('\n')
				fmt.Println()
				read = strings.TrimSpace(read)
				if err != nil {
					return
				}
				if len(read) > 0 {
					err = c.conn.WriteMessage(websocket.TextMessage, []byte(read))
					if err != nil {
						return
					}
				}
			}
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				cancel()
				return
			}
			fmt.Println(string(message))
			fmt.Println()
		}
	}
}

func (c *Client) Stop() error {
	return c.conn.Close()
}
