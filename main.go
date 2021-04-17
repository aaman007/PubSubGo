package main

import (
	"github.com/aaman007/pubsubgo2/pubsub"
	"log"
	"net/http"
)

func websocketHandler(w http.ResponseWriter, req *http.Request) {
	pubsub.ServeWS(w, req)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, "static")
	})
	http.HandleFunc("/ws", websocketHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
