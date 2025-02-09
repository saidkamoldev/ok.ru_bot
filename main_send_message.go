// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"strings"
// 	"time"
// )

// // Telegram bot token va kanal ID
// const botToken = "6739271319:AAFZ2WzM1wI5CPb7qL7XqU3lKgnLHhQoalg"
// const chatID = "-1001871864851"

// // OK.ru guruh sahifasi va access token
// const okGroupURL = "https://api.ok.ru/group/70000034205196/messages"
// const accessToken = "-n-2ymYnH6EEo3Xs9hIgvk8Ku7udTeT5Z2M8E0Q0BgTsRI03wBey8RLPwc2iu07WsBMG4dUHnT48ptX4iA90:CMBDJQLGDIHBABABA"

// // OK.ru guruhidan yangi postlarni olish
// type OKPost struct {
// 	ID   string `json:"id"`
// 	Text string `json:"text"`
// }

// // OK.ru guruhidan postlarni olish
// func getOKPosts() ([]OKPost, error) {
// 	url := fmt.Sprintf("%s?access_token=%s", okGroupURL, accessToken)

// 	resp, err := http.Get(url)
// 	if err != nil {
// 		log.Println("OK.ru sahifasini ochib bo‘lmadi:", err)
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		log.Printf("API so‘rovi xato berdi: %d %s", resp.StatusCode, resp.Status)
// 		return nil, fmt.Errorf("API xatosi: %s", resp.Status)
// 	}

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Println("Javobni o‘qishda xato:", err)
// 		return nil, err
// 	}

// 	// API javobini konsolga chiqarish
// 	fmt.Println("API javobi:", string(body))

// 	var posts []OKPost
// 	err = json.Unmarshal(body, &posts)
// 	if err != nil {
// 		log.Println("JSONni parse qilishda xato:", err)
// 		return nil, err
// 	}

// 	return posts, nil
// }

// // Telegramga xabar yuborish
// func sendToTelegram(message string) {
// 	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
// 	data := fmt.Sprintf(`{"chat_id":"%s","text":"%s"}`, chatID, message)

// 	req, err := http.NewRequest("POST", url, strings.NewReader(data))
// 	if err != nil {
// 		log.Println("So‘rovni yaratishda xato:", err)
// 		return
// 	}
// 	req.Header.Set("Content-Type", "application/json")

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		log.Println("Telegramga yuborib bo‘lmadi:", err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		log.Printf("Telegramga yuborilgan xato: %s", resp.Status)
// 	} else {
// 		log.Println("Xabar Telegramga muvaffaqiyatli yuborildi!")
// 	}
// }

// func main() {
// 	for {
// 		// OK.ru guruhidan postlarni olish
// 		posts, err := getOKPosts()
// 		if err != nil {
// 			log.Println("Postlarni olishda xato:", err)
// 			time.Sleep(5 * time.Minute) // Xatolik bo‘lsa, 5 daqiqa kutish
// 			continue
// 		}

// 		// Yangi postlarni tekshirish
// 		for _, post := range posts {
// 			if strings.Contains(post.Text, "yangi post") { // Bu yerda o‘zingiz xohlagan filtrni qo‘shing
// 				log.Println("Yangi post topildi:", post.Text)
// 				sendToTelegram(fmt.Sprintf("OK.ru guruhida yangi post bor! 🔥\n%s", post.Text))
// 			} else {
// 				log.Println("Yangi post topilmadi:", post.Text)
// 			}
// 		}

//			// Har 5 daqiqada yana tekshirib ko‘rish
//			time.Sleep(5 * time.Minute)
//		}
//	}
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	// "os"
	"strings"
	"time"
)

// Telegram bot token va kanal ID
const botToken = "6739271319:AAFZ2WzM1wI5CPb7qL7XqU3lKgnLHhQoalg"
const chatID = "-1001871864851"
const accessToken = "-n-2ymYnH6EEo3Xs9hIgvk8Ku7udTeT5Z2M8E0Q0BgTsRI03wBey8RLPwc2iu07WsBMG4dUHnT48ptX4iA90:CMBDJQLGDIHBABABA"

// OK.ru guruh sahifasi va access token
const okGroupURL = "https://api.ok.ru/group/70000034205196/messages"

// OK.ru guruhidan yangi postlarni olish
type OKPost struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

// OK.ru guruhidan postlarni olish
func getOKPosts() ([]OKPost, error) {
	url := fmt.Sprintf("%s?access_token=%s&count=1", okGroupURL, accessToken)

	resp, err := http.Get(url)
	if err != nil {
		log.Println("OK.ru sahifasini ochib bo‘lmadi:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("API so‘rovi xato berdi: %d %s", resp.StatusCode, resp.Status)
		return nil, fmt.Errorf("API xatosi: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Javobni o‘qishda xato:", err)
		return nil, err
	}

	// API javobini konsolga chiqarish
	fmt.Println("API javobi:", string(body))

	var posts []OKPost
	err = json.Unmarshal(body, &posts)
	if err != nil {
		log.Println("JSONni parse qilishda xato:", err)
		return nil, err
	}

	return posts, nil
}

// Telegramga xabar yuborish
func sendToTelegram(message string) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	data := fmt.Sprintf(`{"chat_id":"%s","text":"%s"}`, chatID, message)

	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		log.Println("So‘rovni yaratishda xato:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Telegramga yuborib bo‘lmadi:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Telegramga yuborilgan xato: %s", resp.Status)
	} else {
		log.Println("Xabar Telegramga muvaffaqiyatli yuborildi!")
	}
}

// Postlar orasida o'zgarishni aniqlash uchun ID saqlash
var lastPostID string

// Yangi postlarni tekshirish va yuborish
func checkAndSendPosts(posts []OKPost) {
	for _, post := range posts {
		if post.ID != lastPostID {
			log.Println("Yangi post topildi:", post.Text)
			sendToTelegram(fmt.Sprintf("OK.ru guruhida yangi post bor! 🔥\n%s", post.Text))
			lastPostID = post.ID // So'nggi post ID'sini saqlash
		} else {
			log.Println("Yangi post topilmadi:", post.Text)
		}
	}
}

func main() {
	for {
		// OK.ru guruhidan postlarni olish
		posts, err := getOKPosts()
		if err != nil {
			log.Println("Postlarni olishda xato:", err)
			time.Sleep(5 * time.Minute) // Xatolik bo‘lsa, 5 daqiqa kutish
			continue
		}

		// Yangi postlarni tekshirish
		checkAndSendPosts(posts)

		// Har 5 daqiqada yana tekshirib ko‘rish
		time.Sleep(5 * time.Minute)
	}
}
