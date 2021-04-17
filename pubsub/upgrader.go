package pubsub

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var pubSubManager = PubSub{}

func readPump(client *Client) {
	for {
		messageType, p, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println("Something went wrong")
			pubSubManager.RemoveClient(client)
			break
		}
		
		pubSubManager.HandleReceivedMessage(client, messageType, p)
	}
}

func ServeWS(w http.ResponseWriter, req *http.Request) {
	upgrader.CheckOrigin = func(req *http.Request) bool {
		return true
	}
	
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := pubSubManager.NewClient(conn)
	readPump(client)
}
