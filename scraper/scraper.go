package scraper

import (
	"log"
	// "strings"
	"time"

	"github.com/gocolly/colly"
)

// Post strukturasi
type Post struct {
	ID      string
	Text    string
	PostURL string
	Image   string
}

// Global o‘zgaruvchi: so‘nggi post
var lastPost Post

// Sahifaning eng yuqori postini olish
func GetTopPost() (Post, error) {
	c := colly.NewCollector()
	var topPost Post

	c.OnHTML(".feed-w:first-child", func(e *colly.HTMLElement) {
		topPost = Post{
			Text:    e.ChildText(".media-text"),
			Image:   e.ChildAttr("img", "src"),
			PostURL: e.Request.URL.String(),
		}
	})

	err := c.Visit("https://ok.ru/group/70000034205196")
	if err != nil {
		log.Println("Scrapingda xatolik:", err)
		return topPost, err
	}

	if topPost.Text == "" {
		return topPost, logError("❌ Hech qanday yangi post topilmadi!", "")
	}

	return topPost, nil
}

// Har 10 sekundda postni tekshirish va yangisini yuborish
func StartScraper(newPostChan chan<- Post) {
	for {
		topPost, err := GetTopPost()
		if err == nil && topPost.Text != "" {
			// Agar post yangi bo‘lsa, uni saqlaymiz va botga yuboramiz
			if lastPost.Text != topPost.Text {
				lastPost = topPost
				newPostChan <- topPost
			}
		}
		time.Sleep(10 * time.Second)
	}
}

// Xatolarni log qilish
func logError(msg, data string) error {
	log.Printf("%s %s\n", msg, data)
	return nil
}
