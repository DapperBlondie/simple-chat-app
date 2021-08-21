package handlers

import (
	"fmt"
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

	conn := &WSConnection{MyConn: wsConn}
	Clients[conn] = ""

	err = wsConn.WriteJSON(resp)
	if err != nil {
		zerolog.Error().Msg(err.Error())
		http.Error(w, "Error in sending the response to user over TCP connection", http.StatusInternalServerError)
		return
	}

	go ListenForWS(conn)
}

// ListenForWS listening to every request and send them to
func ListenForWS(conn *WSConnection) {
	defer func() {
		if r := recover(); r != nil {
			zerolog.Error().Msg("error ; " + fmt.Sprintf("%v", r))
		}
	}()

	payload := &WsPayload{
		Action:   "",
		Username: "",
		Message:  "",
		UserConn: nil,
	}
	for {
		err := conn.MyConn.ReadJSON(&payload)
		if err != nil {
			err = conn.MyConn.WriteMessage(websocket.BinaryMessage, []byte(err.Error()))
			if err != nil {
				return
			}
		} else {
			payload.UserConn = conn
			WsChan <- payload
		}
	}
}

// ListenToWsChannel use for listening to our websocket channel and receiving *WsPayload
func ListenToWsChannel() {
	resp := &WsJsonResponse{
		Action:      "",
		Message:     "",
		MessageType: "",
	}
	for {
		e := <-WsChan

		resp.Action = e.Action + "; Action"
		resp.Message = fmt.Sprintf("Some message you sent : %v", e.Username)
		broadCastToAll(resp)
	}
}

// broadCastToAll use for broadCasting to all users
func broadCastToAll(resp *WsJsonResponse) {
	for client := range Clients {
		err := client.MyConn.WriteJSON(resp)
		if err != nil {
			zerolog.Error().Msg(err.Error() + "; occurred in broadcasting")
			err = client.MyConn.Close()
			if err != nil {

			}
			delete(Clients, client)
		}
	}
}
