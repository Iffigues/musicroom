package ws

import (
	"bytes"
	"log"
	"net/http"
	"errors"
	"time"
	"fmt"
	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second
	pongWait = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.broadcast <- message
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func NewHub()  (c *Hub) {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func NewChannel() (c *Channel) {
	c = new(Channel)
	c.Chan = make(map[string]*Mess)
	return
}

func (c *Channel)NewChan(a string) (err error) {
	c.RLock()
	if _, ok := c.Chan[a]; ok {
		c.RUnlock()
		return errors.New("chan already exists")
	}
	c.RUnlock()
	e := &Mess{}
	fmt.Println(e)
	e.Hub = NewHub()
	fmt.Println(e)
	e.Close = make(chan bool)
	go e.Hub.run(e)
	c.Lock()
	defer c.Unlock()
	c.Chan[a] = e;
	return
}


func (c *Channel) Close(a string) {
	c.RLock()
	defer c.RUnlock()
	if val, ok := c.Chan[a]; ok {
		val.Close <- true
		delete(c.Chan, a)
	}
}

func (c *Channel) AddClient(a string, w http.ResponseWriter, r *http.Request) (err error){
	if a, ok := c.Chan[a]; ok {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return err
		}

		client := &Client{
			conn: conn,
			send: make(chan []byte, 256),
		}
		a.Hub.register <- client
		go client.writePump()
		go client.readPump()
		return nil
	}
	return errors.New("don't exists")
}
