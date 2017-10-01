package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"../bk"

	"gopkg.in/telegram-bot-api.v4"
)

func getEnv(env, fallback string) string {
	if value, ok := os.LookupEnv(env); ok {
		return value
	}
	return fallback
}

func getEnvInt(env string, fallback int) int {
	if value, ok := os.LookupEnv(env); ok {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}

type sessions struct {
	cache         map[int64]*bk.Session
	createSession func() *bk.Session
}

func initSessions(createSession func() *bk.Session) sessions {
	return sessions{
		cache:         make(map[int64]*bk.Session),
		createSession: createSession,
	}
}

func (s *sessions) get(chatID int64) *bk.Session {
	if session, ok := s.cache[chatID]; ok {
		return session
	}

	log.Println("New session for ", chatID)

	session := s.createSession()
	s.cache[chatID] = session
	return session
}

func handleUpdate(sessions sessions, update tgbotapi.Update, out chan tgbotapi.Chattable) {
	// todo write only
	// updateId := update.UpdateID
	message := update.Message

	if message == nil {
		return
	}

	fromID := message.From.ID
	chatID := message.Chat.ID

	if int64(fromID) != chatID {
		return
	}

	session := sessions.get(chatID)

	command := []byte(strings.TrimSpace(message.Text))
	msg := bk.InMessage{Command: command}
	reply := session.HandleMessage(msg)

	out <- tgbotapi.NewMessage(chatID, reply.Text)
	// bot.Send(replyMsg)

	// s, _ := json.MarshalIndent(message, "", "  ")
	// fmt.Println("%s", string(s))
}

func replyReceiver(bot *tgbotapi.BotAPI, replies <-chan tgbotapi.Chattable) {
	for msg := range replies {
		bot.Send(msg)
	}
}

func updateReceiver(bot *tgbotapi.BotAPI, updates <-chan tgbotapi.Update, out chan tgbotapi.Chattable) {
	// todo leaveChat
	secretGenerator := bk.NewRandomSecretGenerator(4)
	sessions := initSessions(func() *bk.Session {
		session := bk.NewSession(secretGenerator)
		return &session
	})

	for update := range updates {
		go handleUpdate(sessions, update, out)
	}
}

func main() {
	var (
		domain  = getEnv("TG_DOMAINNAME", "example.com")
		token   = getEnv("TG_TOKEN", "")
		extPort = getEnvInt("TG_EXTPORT", 443)
		intPort = getEnvInt("TG_INTPORT", 9000)
	)

	if token == "" {
		log.Fatal("No token specified")
	}

	log.Println("Initializing bot")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal("Auth failed", err)
	}

	updateChannel := createWebhookUpdateChannel(bot, domain, extPort, intPort)
	updates, err := updateChannel.createChannel()

	if err != nil {
		log.Fatal("Cannot init webhook: ", err)
	}

	out := make(chan tgbotapi.Chattable)
	go replyReceiver(bot, out)
	go updateReceiver(bot, updates, out)

	log.Println(updateChannel)
	updateChannel.start()
}
