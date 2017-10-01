package main

import (
	"fmt"
	"log"
	"net/http"

	"gopkg.in/telegram-bot-api.v4"
)

type updateChannel interface {
	createChannel() (<-chan tgbotapi.Update, error)
	start()
}

type webhookUpdateChannel struct {
	bot     *tgbotapi.BotAPI
	extURL  string
	intAddr string
}

func createWebhookUpdateChannel(bot *tgbotapi.BotAPI, domain string, extPort, intPort int) updateChannel {
	url := fmt.Sprintf("https://%s:%d/%s", domain, extPort, bot.Token)
	addr := fmt.Sprintf("127.0.0.1:%d", intPort)

	return webhookUpdateChannel{bot, url, addr}
}

func (w webhookUpdateChannel) String() string {
	return fmt.Sprintf(
		"@%s listens for incoming updated on %s (%s)",
		w.bot.Self.UserName, w.extURL, w.intAddr,
	)
}

func (w webhookUpdateChannel) createChannel() (<-chan tgbotapi.Update, error) {
	webhook := tgbotapi.NewWebhook(w.extURL)
	log.Println("Webhook URL: ", w.extURL)
	_, err := w.bot.SetWebhook(webhook)
	if err != nil {
		return nil, err
		// log.Fatal("setWebhook failed", err)
	}

	updates := w.bot.ListenForWebhook("/" + w.bot.Token)
	return updates, nil
}

func (w webhookUpdateChannel) start() {
	http.ListenAndServe(w.intAddr, nil)
}
