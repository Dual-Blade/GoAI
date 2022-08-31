package main

import (
	"log"
	"net/http"
	"net/url"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func proxyClient() *http.Client {
	// Set Proxy
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("http://127.0.0.1:7890")
	}
	transport := &http.Transport{Proxy: proxy}
	return &http.Client{Transport: transport}
}

func CopyPhotoSend(bot *tgbotapi.BotAPI, update tgbotapi.Update, from int64, to int64) {
	if update.Message != nil {
		if update.Message.Photo != nil && update.Message.Chat.ID == from {
			msg := tgbotapi.NewCopyMessage(to, update.Message.Chat.ID, update.Message.MessageID)
			//msg := tgbotapi.NewPhoto(-642879405, update.Message.Text)
			_, err := bot.Send(msg)
			if err != nil {
				log.Printf("Send copy photo fail [group name]:%s", update.Message.Chat.Title)
			}
		}
	}
}
func CopyUpiSend(bot *tgbotapi.BotAPI, update tgbotapi.Update, from int64, to int64) {
	if update.Message != nil {
		if update.Message.Chat.ID == from {
			text := update.Message.Text
			upiList := ExtractUPI(text, `[\w-]+@[\w]+`)
			if len(upiList) == 0 {
				return
			}
			for _, upi := range upiList {
				msg := tgbotapi.NewMessage(to, upi)
				_, err := bot.Send(msg)
				if err != nil {
					log.Printf("Send copy upi fail [group name]:%s", update.Message.Chat.Title)
				}
			}
			MergeSql(len(upiList))
			// @user to recv upi to test
			msg := tgbotapi.NewMessage(to, "@Sanju_Tiger")
			_, err := bot.Send(msg)
			if err != nil {
				log.Printf("Send @user fail [group name]:%s", update.Message.Chat.Title)
			}
		}
	}
}

func ExtractUPI(text string, p string) []string {
	r := regexp.MustCompile(p)
	return r.FindAllString(text, -1)
}

func main() {
	// Use proxy client to new a bot instance
	bot, err := tgbotapi.NewBotAPIWithClient("5391990902:AAFqZORAme3GWQrv0e4yYizJ_VIRHETLWas", tgbotapi.APIEndpoint, proxyClient())
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		// UPIGBP:-633478699  jiang:-665365433
		CopyUpiSend(bot, update, -642879405, -749897072)
		CopyPhotoSend(bot, update, -749897072, -642879405)
	}
}
