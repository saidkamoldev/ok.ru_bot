// package main

// import (
// 	"fmt"
// 	"log"
// 	"strings"

// 	"github.com/gocolly/colly"
// )

// // OK.ru sahifa URL'si
// const okGroupURL = "https://ok.ru/group/70000034205196" // Guruh ID'sini o‘zingizga moslang

// func main() {
// 	// Colly collector yaratish
// 	c := colly.NewCollector()

// 	// Sahifa ichidagi postlarni yig‘ish
// 	c.OnHTML("div.media-text", func(e *colly.HTMLElement) {
// 		postText := strings.TrimSpace(e.Text)
// 		fmt.Println("Post:", postText)
// 	})

// 	// Sahifani ochish
// 	err := c.Visit(okGroupURL)
// 	if err != nil {
// 		log.Fatal("Sahifani ochib bo‘lmadi:", err)
// 	}
// }
// package main

// import (
// 	"encoding/json"
// 	// "fmt"
// 	"log"
// 	"net/http"
// 	"strings"
// 	"time"

// 	"github.com/gocolly/colly"
// 	"github.com/gorilla/websocket"
// )

// // OK.ru sahifa URL'si
// const okGroupURL = "https://ok.ru/group/70000034205196"

// // Post tuzilmasi
// type Post struct {
// 	Text     string `json:"text"`
// 	ImageURL string `json:"image_url"`
// }

// // WebSocket uchun upgrader
// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool { return true },
// }

// var clients = make(map[*websocket.Conn]bool) // Barcha WebSocket ulanishlari

// // Oxirgi postni olish funksiyasi
// func getLatestPost() (Post, error) {
// 	c := colly.NewCollector()
// 	var latestPost Post

// 	// Post matnini olish
// 	c.OnHTML("div.media-text", func(e *colly.HTMLElement) {
// 		if latestPost.Text == "" {
// 			latestPost.Text = strings.TrimSpace(e.Text)
// 		}
// 	})

// 	// Postdagi rasmni olish
// 	c.OnHTML("img", func(e *colly.HTMLElement) {
// 		imgSrc := e.Attr("src")
// 		if latestPost.ImageURL == "" {
// 			latestPost.ImageURL = imgSrc
// 		}
// 	})

// 	// Sahifani ochish
// 	err := c.Visit(okGroupURL)
// 	if err != nil {
// 		return latestPost, err
// 	}

// 	return latestPost, nil
// }

// // WebSocket handler
// func handleConnections(w http.ResponseWriter, r *http.Request) {
// 	ws, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println("WebSocket ulanishda xato:", err)
// 		return
// 	}
// 	defer ws.Close()

// 	clients[ws] = true // Yangi mijozni qo'shish

// 	for {
// 		// Mijoz WebSocketni yopganda, uni o‘chiramiz
// 		_, _, err := ws.ReadMessage()
// 		if err != nil {
// 			delete(clients, ws)
// 			break
// 		}
// 	}
// }

// // WebSocket orqali yangi postlarni jo'natish
// func sendUpdates() {
// 	var lastPost Post

// 	for {
// 		post, err := getLatestPost()
// 		if err != nil {
// 			log.Println("Postlarni olishda xato:", err)
// 		} else if post.Text != lastPost.Text {
// 			log.Println("Yangi post aniqlandi:", post.Text)

// 			// WebSocket orqali barcha mijozlarga yangi postni jo‘natish
// 			postJSON, _ := json.Marshal(post)
// 			for client := range clients {
// 				err := client.WriteMessage(websocket.TextMessage, postJSON)
// 				if err != nil {
// 					log.Println("WebSocket jo‘natishda xato:", err)
// 					client.Close()
// 					delete(clients, client)
// 				}
// 			}
// 			lastPost = post
// 		}

// 		time.Sleep(10 * time.Second) // 10 soniyadan keyin yana tekshirish
// 	}
// }

// func main() {
// 	// WebSocket server
// 	http.HandleFunc("/ws", handleConnections)

// 	// WebSocket serverni ishga tushirish
// 	go sendUpdates()

// 	log.Println("WebSocket server 8080-portda ishlayapti...")
// 	err := http.ListenAndServe(":8080", nil)
// 	if err != nil {
// 		log.Fatal("Serverni ishga tushirib bo‘lmadi:", err)
// 	}
// }

package main

import (


	"example.com/ok_bot/scraper"
	"example.com/ok_bot/telegram"
	"example.com/ok_bot/websocket"
	// "OK_Bot_Project/scraper"
	// "OK_Bot_Project/telegram"
	// "OK_Bot_Project/websocket"
)

func main() {

	
	// Telegram botni ishga tushirish
	telegram.InitBot("6739271319:AAFZ2WzM1wI5CPb7qL7XqU3lKgnLHhQoalg")
	go telegram.StartTelegramBot()

	// WebSocket serverni ishga tushirish
	go websocket.StartWebSocketServer()

	// Scraperni ishga tushirish
	newPostChan := make(chan scraper.Post)
	go scraper.StartScraper(newPostChan)

	for post := range newPostChan {
		telegram.SendPostToAdmin(post)
	}
}
	