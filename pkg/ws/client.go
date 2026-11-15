package ws

import "github.com/gorilla/websocket"

type Client struct {
	Conn *websocket.Conn
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{Conn: conn}
}
