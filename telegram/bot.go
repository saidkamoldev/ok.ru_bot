package telegram

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"example.com/ok_bot/scraper"
	"gopkg.in/telebot.v3"
)

var (
	bot        *telebot.Bot
	lastPostID string
	mu         sync.Mutex
	postQueue  = make(map[string]scraper.Post) // Postlar navbati
)

const (
	adminID   = 1460974373 // Adminning Telegram ID'si
	channelID = -1001871864851 // Kanal ID'si
)

// 🔹 **Postni tekshirish va adminga yuborish**
func CheckAndSendPost(post scraper.Post) {
	mu.Lock()
	defer mu.Unlock()

	if post.ID == lastPostID {
		log.Println("⚠️ Post allaqachon yuborilgan:", post.ID)
		return
	}

	lastPostID = post.ID
	postQueue[post.ID] = post
	SendPostToAdmin(post)
}


// **📨 Adminga postni yuborish**
func SendPostToAdmin(post scraper.Post) {
	msg := fmt.Sprintf(
		"🆕 *Yangi post:* \n\n📜 *%s*\n🔗 [Postni ko‘rish](%s)",
		escapeMarkdown(post.Text), post.PostURL,
	)

	photo := &telebot.Photo{
		File:    telebot.FromURL(post.Image),
		Caption: msg,
	}

	inlineKeys := &telebot.ReplyMarkup{}
	btnApprove := inlineKeys.Data("✅ Tasdiqlash", "approve_"+post.ID)
	btnReject := inlineKeys.Data("❌ Bekor qilish", "reject_"+post.ID)
	inlineKeys.Inline(inlineKeys.Row(btnApprove, btnReject))

	_, err := bot.Send(&telebot.User{ID: adminID}, photo, inlineKeys)
	if err != nil {
		log.Println("❌ Adminga post yuborishda xatolik:", err)
	} else {
		log.Println("✅ Post adminga yuborildi:", post.ID)
	}
}

// **📢 Kanalga postni yuborish**
func ApprovePost(post scraper.Post) {
	msg := fmt.Sprintf(
		"📢 <b>Yangi post!</b>\n\n📜 %s\n🔗 <a href='%s'>Havola</a>",
		post.Text, post.PostURL,
	)

	photo := &telebot.Photo{
		File:    telebot.FromURL(post.Image),
		Caption: msg,
	}

	options := &telebot.SendOptions{ParseMode: telebot.ModeHTML}
	_, err := bot.Send(&telebot.Chat{ID: channelID}, photo, options)

	if err != nil {
		log.Println("❌ Kanalga post yuborishda xatolik:", err)
	} else {
		log.Println("✅ Post kanalga yuborildi:", post.ID)
	}
}

// **🛠 Botni sozlash**
func InitBot(token string) {
	pref := telebot.Settings{
		Token: token,
		Poller: &telebot.LongPoller{
			Timeout: 10,
		},
	}

	var err error
	bot, err = telebot.NewBot(pref)
	if err != nil {
		log.Fatal("❌ Telegram botni ishga tushirib bo‘lmadi:", err)
	} else {
		log.Println("✅ Telegram bot ishga tushdi!")
	}
}

// **🎛 Admin tugmalari bilan ishlash**
func SetupHandlers() {
	bot.Handle(telebot.OnCallback, func(c telebot.Context) error {
		data := c.Callback().Data
		postID := extractID(data)
		log.Println("postQueue tarkibi:", postQueue)

		post, exists := postQueue[postID]
		if !exists {
			return c.Send("❌ Xatolik: Post topilmadi!")
		}

		if strings.HasPrefix(data, "approve_") {
			ApprovePost(post)
			delete(postQueue, postID)
			return c.Edit("✅ Post tasdiqlandi va kanalga yuborildi!")
		} else if strings.HasPrefix(data, "reject_") {
			delete(postQueue, postID)
			return c.Edit("❌ Post rad etildi!")
		}

		return nil
	})
}

// **📌 Callback tugmachadan ID ajratib olish**
func extractID(data string) string {
    parts := strings.Split(data, "_")
    if len(parts) < 2 {
        return "" // noto‘g‘ri format bo‘lsa, bo‘sh string qaytariladi
    }
	fmt.Println(parts[1])
    return parts[1] // ID qaytariladi
}


// **🔠 Markdown matnni tozalash**
func escapeMarkdown(text string) string {
	symbols := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
	for _, sym := range symbols {
		text = strings.ReplaceAll(text, sym, "\\"+sym)
	}
	return text
}

// **🚀 Botni ishga tushirish**
func StartTelegramBot() {
	if bot == nil {
		log.Fatal("❌ Bot hali yaratilmagan! `InitBot(token)` funksiyasini chaqirishni unutmang!")
	}

	bot.Handle(telebot.OnCallback, func(c telebot.Context) error {
		data := c.Callback().Data
		postID := extractID(data)
	
		log.Println("Callback ma’lumotlari:", data)
		log.Println("Ajratilgan postID:", postID)
	
		post, exists := postQueue[postID]
		if !exists {
			log.Println("❌ Post topilmadi! postQueue tarkibi:", postQueue)
			return c.Send("❌ Xatolik: Post topilmadi!")
		}
	
		return nil
	})
	
	SetupHandlers()
	bot.Start()
}
