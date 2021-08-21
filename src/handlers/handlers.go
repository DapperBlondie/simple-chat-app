package handlers

import (
	"github.com/DapperBlondie/simple-chat-app/src/render"
	"github.com/gorilla/websocket"
	zerolog "github.com/rs/zerolog/log"
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
	Action      string `json:"action"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
}

// WsPayload use for handling user's payload
type WsPayload struct {
	Action   string        `json:"action"`
	Username string        `json:"username"`
	Message  string        `json:"message"`
	UserConn *WSConnection `json:"-"`
}

var AppConf *ApplicationConfig

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

// Home a first page that we send it to user
func (ac *ApplicationConfig) Home(w http.ResponseWriter, r *http.Request) {
	if http.MethodGet != r.Method {
		http.Error(w, "Error in method usage.", http.StatusMethodNotAllowed)
		return
	}

	err := render.RendererPage(w, "home.jet", nil)
	if err != nil {
		http.Error(w, "Error in rendering page", http.StatusInternalServerError)
		return
	}
	return
}

// WsEndpointHandler my first endpoint for chat application
func (ac *ApplicationConfig) WsEndpointHandler(w http.ResponseWriter, r *http.Request) {
	wsConn, err := TcpUpgrade.Upgrade(w, r, nil)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		http.Error(w, "Error in upgrade to TCP connection"+"; "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := &WsJsonResponse{
		Action:      "Check Connection",
		Message:     "Upgraded to TCP",
		MessageType: "Status",
	}

	err = wsConn.WriteJSON(resp)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		http.Error(w, "Error in sending the response to user over TCP connection", http.StatusInternalServerError)
		return
	}
}
