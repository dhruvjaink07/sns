package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type Notification struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

var (
	notifications []Notification
	nextID        int = 1
	clients           = make(map[*websocket.Conn]bool)
	broadcast         = make(chan Notification)
	upgrader          = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func main() {
	http.HandleFunc("/ping", handlePing)
	http.HandleFunc("/notify", handleNotify)
	http.HandleFunc("/notifications", handleGetNotifications)
	http.HandleFunc("/notifications/", handleDeleteNotification)
	http.HandleFunc("/ws", handleWebSocketConnections)

	go handleBroadcast()

	go func() {
		for {
			time.Sleep(5 * time.Second)
			n := Notification{
				ID:      nextID,
				Title:   "Periodic Notification",
				Message: "You'll receive this every 5 seconds",
			}
			nextID++
			notifications = append(notifications, n)
			broadcast <- n
		}
	}()

	fmt.Println("Notification Service is running on port 8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Notification Service is running ✅")
}

func handleNotify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	clientIP := r.RemoteAddr
	log.Printf("Client IP: %s", clientIP)

	var n Notification
	err := json.NewDecoder(r.Body).Decode(&n)
	if err != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
		return
	}

	n.ID = nextID
	nextID++
	notifications = append(notifications, n)

	broadcast <- n

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Notification received: %s", n.Title)
}

func handleGetNotifications(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	clientIP := r.RemoteAddr
	log.Printf("Client IP: %s", clientIP)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

func handleDeleteNotification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	clientIP := r.RemoteAddr
	log.Printf("Client IP: %s", clientIP)

	id := r.URL.Path[len("/notifications/"):]
	notificationID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for i, n := range notifications {
		if n.ID == notificationID {
			notifications = append(notifications[:i], notifications[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Notification not found", http.StatusNotFound)
}

func handleWebSocketConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()
	
	clientIP := r.RemoteAddr
	log.Printf("WebSocket client connected: %s", clientIP)

	clients[conn] = true

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket client disconnected: %s", clientIP)
			delete(clients, conn)
			break
		}
	}
}

func handleBroadcast() {
	for {
		notification := <-broadcast
		for client := range clients {
			err := client.WriteJSON(notification)
			if err != nil {
				log.Println("WebSocket write error:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}