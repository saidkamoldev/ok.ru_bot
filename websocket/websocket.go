package websocket

import (
	"encoding/json"
	"log"
	"net/http"

	"example.com/ok_bot/scraper"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // Ulangan WebSocket mijozlari
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// WebSocket ulanishlarini boshqarish
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket ulanish xatosi:", err)
		return
	}
	defer conn.Close()
	clients[conn] = true

	// WebSocket mijozlari faol bo‘lsa, ularga xabar yuborish
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			delete(clients, conn)
			break
		}
	}
}

// Yangi postni barcha mijozlarga yuborish
func BroadcastPost(post scraper.Post) {
	postJSON, _ := json.Marshal(post)
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, postJSON)
		if err != nil {
			log.Println("Mijozga yozishda xatolik:", err)
			client.Close()
			delete(clients, client)
		}
	}
}

// WebSocket serverni ishga tushirish
func StartWebSocketServer() {
	http.HandleFunc("/ws", HandleConnections)
	log.Println("WebSocket server 8080-portda ishlamoqda...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
