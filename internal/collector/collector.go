package collector

import (
	"context"
	"fmt"
	"log"
	"neurodronizm/internal/generator"
	"neurodronizm/internal/store"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Collector(ctx context.Context, s *store.Store, gen *generator.Generator) {

	bot_token := os.Getenv("COLLECTOR_BOT_TOKEN")
	channel_name := os.Getenv("CHANNEL_NAME")
	raw_my_id := os.Getenv("MY_ID")
	my_id, err := strconv.ParseInt(raw_my_id, 10, 64)
	if err != nil {
		log.Printf("Cannot parse to int: %v", err)
	}
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
			if update.Message.From.ID != my_id {
				continue
			}
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

					callbackData := fmt.Sprintf("publish:%d", id)
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
		if update.CallbackQuery != nil {
			if update.CallbackQuery.From.ID != my_id {
				continue
			}
			if strings.HasPrefix(update.CallbackQuery.Data, "publish") {
				callbackctx, cancel := context.WithTimeout(ctx, 30*time.Second)
				idS := strings.TrimPrefix(update.CallbackQuery.Data, "publish:")
				id, err := strconv.Atoi(idS)
				if err != nil {
					log.Printf("Cannot convert string to int: %v", err)
					cancel()
					continue
				}
				post, err := s.GetDraftByID(callbackctx, id)
				if err != nil {
					log.Printf("Cannot get draft by id: %v", err)
				}
				err = s.ChangeStatus(callbackctx, id)
				if err != nil {
					log.Printf("Cannot change draft status: %v", err)
				}
				channel_post := tgbotapi.NewMessageToChannel(channel_name, post)
				if _, err := bot.Send(channel_post); err != nil {
					log.Printf("Bot cannot publish post: %v", err)
				}
				callBackAnswer := tgbotapi.NewCallback(update.CallbackQuery.ID, "Пост успешно опубликован в канал!")
				if _, err := bot.Request(callBackAnswer); err != nil {
					log.Printf("Cannot answer to callback: %v", err)
				}
				cancel()
			}
		}
	}
}
