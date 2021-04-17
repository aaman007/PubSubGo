package pubsub

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	"log"
)

const (
	PUBLISH     = "publish"
	SUBSCRIBE   = "subscribe"
	UNSUBSCRIBE = "unsubscribe"
)

func randomID() string {
	return uuid.NewV4().String()
}

type PubSub struct {
	Clients       []Client
	Subscriptions []Subscription
}

func (ps *PubSub) NewClient(conn *websocket.Conn) *Client {
	client := Client{ID: randomID(), Conn: conn}
	ps.Clients = append(ps.Clients, client)
	
	log.Println("Added new client with ID:", client.ID)
	log.Println("Total Client Now:", len(ps.Clients))
	
	payload := []byte("Hello Client. Your ID is:" + client.ID)
	_ = client.send(payload)
	
	return &client
}

func (ps *PubSub) RemoveClient(client *Client) {
	// Remove client's subscriptions
	for index, subscription := range ps.Subscriptions {
		if subscription.Client.ID == client.ID {
			ps.Subscriptions = append(ps.Subscriptions[:index], ps.Subscriptions[index+1:]...)
		}
	}
	log.Println("Removed Client's Subscriptions. Total Subscriptions Now:", len(ps.Subscriptions))
	
	// Remove Client
	for index, _client := range ps.Clients {
		if _client.ID == client.ID {
			ps.Clients = append(ps.Clients[:index], ps.Clients[index+1:]...)
		}
	}
	
	log.Println("Removed Client. Total Clients Now:", len(ps.Clients))
}

func (ps *PubSub) GetSubscriptions(client *Client, topic string) []Subscription {
	var subscriptions []Subscription
	for _, subscription := range ps.Subscriptions {
		if client != nil {
			if client.ID == subscription.Client.ID && topic == subscription.Topic {
				subscriptions = append(subscriptions, subscription)
			}
		} else if subscription.Topic == topic {
			subscriptions = append(subscriptions, subscription)
		}
	}
	return subscriptions
}

func (ps *PubSub) Subscribe(client *Client, topic string) {
	if len(ps.GetSubscriptions(client, topic)) > 0 {
		log.Println("Already Subscribed")
		return
	}
	subscription := Subscription{Client: client, Topic: topic}
	ps.Subscriptions = append(ps.Subscriptions, subscription)
	log.Printf("Client %v subscribed to %v\n", client.ID, topic)
	log.Println("Total Subscriptions Now:", len(ps.Subscriptions))
}

func (ps *PubSub) Unsubscribe(client *Client, topic string) {
	unsubscribed := false
	for index, subscription := range ps.Subscriptions {
		if subscription.Client.ID == client.ID && subscription.Topic == topic {
			unsubscribed = true
			ps.Subscriptions = append(ps.Subscriptions[:index], ps.Subscriptions[index+1:]...)
			log.Println("Unsubscribed from topic")
			log.Println("Total Subscriptions Now:", len(ps.Subscriptions))
		}
	}
	
	if !unsubscribed {
		log.Println("You did not subscribed to this topic")
	}
}

func (ps *PubSub) Publish(topic string, message []byte, exclude *Client) {
	subscriptions := ps.GetSubscriptions(nil, topic)
	for _, subscription := range subscriptions {
		if exclude != nil && subscription.Client.ID == exclude.ID {
			continue
		}
		client := subscription.Client
		_ = client.send(message)
		
		log.Println("Sent message to client", client.ID)
	}
}

func (ps *PubSub) HandleReceivedMessage(client *Client, messageType int, p []byte) {
	message := Message{}
	err := json.Unmarshal(p, &message)
	if err != nil {
		log.Println("Message not in correct format", err)
		return
	}
	
	switch message.Action {
	case PUBLISH:
		ps.Publish(message.Topic, message.Message, client)
	case SUBSCRIBE:
		ps.Subscribe(client, message.Topic)
	case UNSUBSCRIBE:
		ps.Unsubscribe(client, message.Topic)
	default:
		log.Println("Unknown Action Type")
	}
}

type Subscription struct {
	Topic  string
	Client *Client
}

type Message struct {
	Action  string          `json:"action"`
	Topic   string          `json:"topic"`
	Message json.RawMessage `json:"message"`
}

type Client struct {
	ID   string
	Conn *websocket.Conn
}

func (c *Client) send(message []byte) error {
	return c.Conn.WriteMessage(websocket.TextMessage, message)
}
