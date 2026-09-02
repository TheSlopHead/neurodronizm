package collector

import (
	"context"
	"fmt"
	"log"
	"neurodronizm/internal/generator"
	"neurodronizm/internal/store"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Collector(ctx context.Context, s *store.Store, gen *generator.Generator) {

	bot_token := os.Getenv("COLLECTOR_BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(bot_token)

	if err != nil {
		log.Printf("Cannot validate telegram bot: %v", err)
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
			} else if post.ForwardFromChat != nil || post.ForwardFrom != nil {
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

		}
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			if update.Message.IsCommand() && update.Message.Command() == "generate" {

				cmdCtx, cancel := context.WithTimeout(ctx, 50*time.Second)
				examples, err := s.GetLastPost(cmdCtx, 3)
				if err != nil {
					log.Printf("Cannot get posts from database: %v", err)
				}
				topic := update.Message.CommandArguments()

				if topic == "" {
					topic = "Сгенерируй 3 разных варианта поста в моем стиле, иронично и со стебом. Разделяй варианты строкой [POST_SPLIT]. Внутри самих постов этот маркер не используй"
				} else {
					topic += "Разделяй варианты строкой [POST_SPLIT]. Внутри самих постов этот маркер не используй"
				}

				variants, err := gen.GeneratePost(cmdCtx, examples, topic)
				if err != nil {
					log.Printf("Cannot generate post: %v", err)
				}

				var responseText string
				var draftIDs []int
				for i, variant := range variants {
					id, err := s.SaveDraft(cmdCtx, variant)
					if err != nil {
						log.Printf("Cannot use savedraft: %v", err)
						break
					}
					draftIDs = append(draftIDs, id)
					responseText += fmt.Sprintf("<b>Variant %d: </b>\n%s\n\n", i+1, variant)
				}

				var row []tgbotapi.InlineKeyboardButton
				for index, id := range draftIDs {

					callbackData := fmt.Sprintf("publish: %d", id)
					buttonText := fmt.Sprintf("Variant: %d", index+1)

					btn := tgbotapi.NewInlineKeyboardButtonData(buttonText, callbackData)
					row = append(row, btn)
				}
				keyboard := tgbotapi.NewInlineKeyboardMarkup(row)

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, responseText)
				msg.ParseMode = "HTML"
				msg.ReplyMarkup = keyboard
				_, err = bot.Send(msg)
				if err != nil {
					log.Printf("Cannot send message: %v", err)
				}
				cancel()

			}
		}
	}
}
