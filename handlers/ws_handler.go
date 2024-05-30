package handlers

import (
	"context"
	"log"
	"net/http"
	"sync"
	"website/database"
	"website/models"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	clients   = make(map[*websocket.Conn]bool)
	mutex     sync.Mutex
	broadcast = make(chan models.Message)
)

func HandleConnection(conn *websocket.Conn, username string, chatID string) {
	defer conn.Close()

	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()

	defer func() {
		mutex.Lock()
		delete(clients, conn)
		mutex.Unlock()
	}()

	loadPreviousMessages(conn, chatID)

	for {
		var msg models.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("Error reading message: %v", err)
			}
			break
		}
		msg.Username = username
		msg.ChatID = chatID
		saveMessage(msg, chatID)
		broadcast <- msg
	}
}

func HandleMessages() {
	for msg := range broadcast {
		mutex.Lock()
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Error writing message: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}

func saveMessage(msg models.Message, chatID string) {
	collection := database.Client.Database("project").Collection("messages_" + chatID)
	_, err := collection.InsertOne(context.Background(), msg)
	if err != nil {
		log.Printf("Error saving message to MongoDB: %v", err)
	}
}

func loadPreviousMessages(conn *websocket.Conn, chatID string) {
	collection := database.Client.Database("project").Collection("messages_" + chatID)
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Printf("Error loading previous messages: %v", err)
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var msg models.Message
		if err := cursor.Decode(&msg); err != nil {
			log.Printf("Error decoding message: %v", err)
			continue
		}
		if err := conn.WriteJSON(msg); err != nil {
			log.Printf("Error sending previous message: %v", err)
		} else {
			log.Printf("Sent previous message: %+v", msg) // Add this line for logging
		}
	}
	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
	}
}
=======
	"log"
	"net/http"
	"sync"
	"website/models"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

var (
    clients = make(map[*websocket.Conn]bool)
    chats   = make(map[string]*models.Chat)
    mutex   sync.Mutex
    broadcast = make(chan models.Message)
)

func HandleConnection(conn *websocket.Conn) {
	var newChat models.Chat
	// Generate a unique chat ID and create a new chat session
	newChat.ID = generateUniqueID()
	newChat.ClientID = "client1" // Replace with actual client ID
	newChat.Active = true
	chats[newChat.ID] = &newChat

	chatNotification := struct {
			ClientID string
			ChatID   string
	}{
			ClientID: newChat.ClientID,
			ChatID:   newChat.ID,
	}

	for client := range clients {
			err := client.WriteJSON(chatNotification)
			if err != nil {
					log.Println("Error sending chat notification:", err)
			}
	}

	// Handle messages for the new chat
	clients[conn] = true
	defer func() {
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
			conn.Close()
	}()

	for {
			var msg models.Message
			err := conn.ReadJSON(&msg)
			if err != nil {
					log.Println("Error reading message:", err)
					break
			}
			broadcast <- msg
	}
}


func HandleMessages() {
    for {
        msg := <-broadcast
        mutex.Lock()
        for client := range clients {
            err := client.WriteJSON(msg)
            if err != nil {
                log.Println("Error writing message:", err)
                client.Close()
                delete(clients, client)
            }
        }
        mutex.Unlock()
    }
}

func generateUniqueID() string {
    return uuid.New().String()
}

