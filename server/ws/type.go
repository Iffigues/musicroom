package ws

import (
	"net/http"
	"github.com/gorilla/websocket"
	"sync"
)

type Channel struct {
	sync.RWMutex
	Chan map[string]*Mess
}


type Mess struct {
	Hub *Hub
	Client []*Client
	Close chan bool
}

type Hub struct {
	clients map[*Client]bool
	broadcast chan []byte
	register chan *Client
	unregister chan *Client
}

type Client struct {
	hub	*Hub
	upgrader  websocket.Upgrader
	conn	*websocket.Conn
	Hh	http.Header
	send chan []byte
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run(c *Mess) {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		case closes := <-c.Close:
			if closes == true {
				return
			}
		}
	}
}
