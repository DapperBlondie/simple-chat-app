package handlers

import (
	"fmt"
	"github.com/gorilla/websocket"
	zerolog "github.com/rs/zerolog/log"
	"sort"
)

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

	go func() {
		for {
			err := conn.MyConn.ReadJSON(&payload)
			if err != nil {
				err = conn.MyConn.WriteMessage(websocket.BinaryMessage, []byte(err.Error()))
				if err != nil {
					return
				}
			} else {
				payload.UserConn = conn
				Clients[conn] = payload.Username
				WsChan <- payload
			}
		}
	}()

	<-conn.CloseChan
	err := conn.MyConn.Close()
	if err != nil {
		return
	}
	return
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

		switch e.Action {
		case "usernames":
			users := getAllClients()
			resp.Action = "UsersList"
			resp.MessageType = "JSON"
			resp.Message = "Get all usernames"
			resp.UsersList = users
			go broadCastToAll(resp)
		case "left":
			delete(Clients, e.UserConn)
			users := getAllClients()
			resp.Action = "LeftUser"
			resp.MessageType = "JSON"
			resp.Message = e.Username + " Left !"
			resp.UsersList = users
			close(e.UserConn.CloseChan)
			go broadCastToAll(resp)
		case "broadcast":
			resp.Action = "BroadCast"
			resp.MessageType = "JSON"
			resp.Message = e.Username + " : " + e.Message
			go broadCastToAll(resp)
		default:
			resp.Action = e.Action + "; Action"
			resp.Message = fmt.Sprintf("Some message you sent : %v", e.Username)
			resp.MessageType = "JSON"
			err := e.UserConn.MyConn.WriteJSON(resp)
			if err != nil {
				return
			}
		}
	}
}

// broadCastToAll use for broadCasting to all
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

// getAllClients use for getting all available clients
func getAllClients() []string {
	var users []string = []string{}
	for _, name := range Clients {
		users = append(users, name)
	}
	sort.Strings(users)

	return users
}
