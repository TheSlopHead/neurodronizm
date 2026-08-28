package collector

import (
	"context"
	"log"
	"neurodronizm/cmd/store"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Collector(ctx context.Context, s *store.Store) {

	bot_token := os.Getenv("COLLECTOR_BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(bot_token)

	if err != nil {
		panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)

	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.ChannelPost != nil {
			post := update.ChannelPost
			if post.Text == "" && post.Caption == "" {
				continue
			}
			if post.Text == "" && post.Caption != "" {
				s.SavePost(ctx, int64(post.MessageID), post.Caption, time.Unix(int64(post.Date), 0))
				//TODO - сделать нормально по условиям, создать переменную  для текста и обработать все возможные ошибки
			} else {
				err := s.SavePost(ctx, int64(post.MessageID), post.Text, time.Unix(int64(post.Date), 0))
				if err != nil {
					log.Printf("failed to save post: %v", err)
				} else {
					log.Printf("Channel post %d saved succesfully", post.MessageID)
				}
			}

			if update.Message != nil {
				log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
				msg.ReplyToMessageID = update.Message.MessageID

				bot.Send(msg)
			}

		}
	}
}
