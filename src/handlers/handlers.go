package handlers

import (
	"github.com/DapperBlondie/simple-chat-app/src/render"
	zerolog "github.com/rs/zerolog/log"
	"net/http"
)

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
