package service

import (
	"github.com/gorilla/websocket"
	"log/slog"
	"time"
	"net/http"
	"bytes"
)

const (
	writeWait = 10 * time.Second
	maxMessageSize = 512
	pongWait = time.Second * 60
	pingPeriod = (pongWait * 9) / 10
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	hub * Hub
	conn *websocket.Conn
	send chan []byte
}

type Hub struct {
	subscribers map[*Client] bool
	subscribe chan *Client
	unsubscribe chan *Client
	stream chan []byte
}

func NewHub(stream chan []byte) Hub {
	hub := Hub {
		subscribers: make(map[*Client] bool),
		subscribe: make(chan *Client),
		unsubscribe: make(chan *Client),
		stream: stream,
	}
	go hub.run()
	return hub
}

func (hub * Hub) run() {
	for {
		select {
		case client := <- hub.subscribe:
			hub.subscribers[client] = true
			slog.Debug("client subscribed", "count", len(hub.subscribers))
		case client := <- hub.unsubscribe:
			if _, ok := hub.subscribers[client]; ok {
				delete(hub.subscribers, client)
				close(client.send)
			}
			slog.Debug("client unsubscribed", "count", len(hub.subscribers))
		case message := <- hub.stream:
			//slog.Debug("broadcasting mesasge", "msg", message)
			for client := range hub.subscribers {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(hub.subscribers, client)
				}
			}
		}
	}
}

func (hub * Hub) Subscribe(w http.ResponseWriter, r * http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Warn(err.Error())
		return
	}
	client := NewClient(hub, conn)
	hub.subscribe <- &client
	go client.read()
	go client.write()
}

func NewClient (hub * Hub, conn * websocket.Conn) Client {
	return Client {
		hub: hub,
		conn: conn,
		send: make(chan []byte),
	}
}

func (c * Client) read() {
	defer func() {
		c.hub.unsubscribe <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error(err.Error())
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		slog.Warn("User sent unexpected message", "msg", message)
	}
}

func (c * Client) write() {
	ticker := time.NewTicker(pingPeriod)
	defer func () {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <- c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				slog.Warn("websocket failed to create writer")
				return
			}
			w.Write(message)
			if err := w.Close(); err != nil {
				slog.Warn("websocket failed to write")
				return
			}
		case <- ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
