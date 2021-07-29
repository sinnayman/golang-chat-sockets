package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
)

var wsChan = make(chan WsPayload)
var clients = make(map[WebSocketConnection]string)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode(),
)

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func Home(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "home.jet", nil)
	if err != nil {
		log.Println(err)
	}
}

type WebSocketConnection struct {
	*websocket.Conn
}

type WsResponse struct {
	Action         string   `json:"action"`
	Message        string   `json:"message"`
	MessageType    string   `json:"message_type"`
	ConnectedUsers []string `json:"connected_users"`
}

type WsPayload struct {
	Action   string              `json:"action"`
	Message  string              `json:"message"`
	Username string              `json:"username"`
	Conn     WebSocketConnection `json:"-"`
}

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	conn := WebSocketConnection{
		ws,
	}
	clients[conn] = ""

	log.Println("client connected to endpoint")
	var response WsResponse = WsResponse{}
	response.Message = `<em><small>Connected to server</small></em>`

	err = ws.WriteJSON(response)
	if err != nil {
		log.Println(err)
	}
	go listenForWs(&conn)
}

func listenForWs(connection *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error", fmt.Sprintf("%v", r))
		}
	}()

	var payload WsPayload
	for {
		err := connection.ReadJSON(&payload)
		if err != nil {
			fmt.Println(err)
		} else {
			payload.Conn = *connection
			wsChan <- payload
		}
	}

}

func Listen() {
	var response WsResponse

	for {
		m := <-wsChan

		switch m.Action {
		case "user_joined":
			clients[m.Conn] = m.Username
			response.Action = "update_users"
			response.ConnectedUsers = getUsersList()
			fmt.Println(response)
		}
		broadcast(response)
	}
}

func getUsersList() []string {
	var usersList []string
	for _, client := range clients {
		usersList = append(usersList, client)
	}
	sort.Strings(usersList)
	return usersList
}

func broadcast(response WsResponse) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			log.Println("Websocket error")
			_ = client.Close()
			delete(clients, client)
		}
	}
}

func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	view, err := views.GetTemplate(tmpl)
	if err != nil {
		log.Println(err)
	}

	err = view.Execute(w, data, nil)
	if err != nil {
		log.Println(err)
	}
	return nil
}
