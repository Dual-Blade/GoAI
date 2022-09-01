package main

import (
	"fmt"
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

func CopyPhotoSend(bot *tgbotapi.BotAPI, update tgbotapi.Update, from int64, to int64) error {
	if update.Message != nil {
		if update.Message.Photo != nil && update.Message.Chat.ID == from {
			msg := tgbotapi.NewCopyMessage(to, update.Message.Chat.ID, update.Message.MessageID)
			//msg := tgbotapi.NewPhoto(-642879405, update.Message.Text)
			_, err := bot.Send(msg)
			if err != nil {
				return fmt.Errorf("CopyPhotoSend fail [group name-%s]: %s", update.Message.Chat.Title, err)
			}
		}
	}
	return nil
}

func CopyUpiSend(bot *tgbotapi.BotAPI, update tgbotapi.Update, from int64, to int64) (err error) {
	if update.Message != nil {
		if update.Message.Chat.ID == from {
			text := update.Message.Text
			upiList, err := ExtractUPI(text, `[\w-]+@[\w]+`)
			if err != nil {
				return fmt.Errorf("CopyUpiSend fail [group name-%s]: %s", update.Message.Chat.Title, err)
			}
			if len(upiList) == 0 {
				return fmt.Errorf("未解析到UPI")
			}
			//发送upi列表到目标群，并统计成功数量
			n := 0
			for _, upi := range upiList {
				msg := tgbotapi.NewMessage(to, upi)
				if _, err = bot.Send(msg); err != nil {
					log.Printf("发送消息失败-upi：%s", err)
				} else {
					n++
				}
			}
			//成功数量0
			if n == 0 {
				return fmt.Errorf("CopyUpiSend fail: 0")
			}
			// @user to recv upi to test
			msg := tgbotapi.NewMessage(to, "@Sanju_Tiger")
			if _, err = bot.Send(msg); err != nil {
				log.Printf("发送消息失败-@user：%s", err)
			}
			//发送成功数量到目标群
			msg = tgbotapi.NewMessage(from, fmt.Sprintf("upi 发送成功，数量：%d", n))
			if _, err = bot.Send(msg); err != nil {
				log.Printf("发送消息失败-upi数量：%s", err)
			}
			//储存统计数量到mysql
			if err = MergeSql(n); err != nil {
				return fmt.Errorf("CopyUpiSend fail: %s", err.Error())
			}
		}
	}
	return nil
}

func ExtractUPI(text string, p string) ([]string, error) {
	r, err := regexp.Compile(p)
	if err != nil {
		return []string{}, fmt.Errorf("ExtractUPI fail: %s", err)
	}
	return r.FindAllString(text, -1), nil
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
		// UPIGBP:-633478699 -749897072  jiang:-665365433 -642879405
		err := CopyUpiSend(bot, update, -642879405, -749897072)
		if err != nil {
			log.Println(err)
		}
		err = CopyPhotoSend(bot, update, -749897072, -642879405)
		if err != nil {
			log.Println(err)
		}
	}
}
