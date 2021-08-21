package handlers

import (
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type WSConnection struct {
	MyConn *websocket.Conn
}

type ApplicationConfig struct {
}

// WsJsonResponse use for send response to user
type WsJsonResponse struct {
	Action      string   `json:"action"`
	Message     string   `json:"message"`
	MessageType string   `json:"message_type"`
	UsersList   []string `json:"users_list,omitempty"`
}

// WsPayload use for handling user's payload
type WsPayload struct {
	Action   string        `json:"action"`
	Username string        `json:"username"`
	Message  string        `json:"message"`
	UserConn *WSConnection `json:"-"`
}

var AppConf *ApplicationConfig

var WsChan = make(chan *WsPayload)
var Clients = make(map[*WSConnection]string)

// TcpUpgrade use for upgrading HTTP request to TCP connection
var TcpUpgrade = websocket.Upgrader{
	HandshakeTimeout: time.Second * 10,
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	EnableCompression: true,
}
