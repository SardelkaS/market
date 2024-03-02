package tg_bot_usecase

import (
	"core/config"
	"core/internal/tg_bot"
	tg_bot_model "core/internal/tg_bot/model"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"strings"
	"time"
)

type uc struct {
	bot *tele.Bot
	cfg *config.Config
}

func New(cfg *config.Config) (tg_bot.UC, error) {
	settings := tele.Settings{
		Token:  cfg.TgBot.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	core, err := tele.NewBot(settings)
	if err != nil {
		return nil, err
	}

	return &uc{
		bot: core,
		cfg: cfg,
	}, nil
}

func (u *uc) NotifyNewOrder(input tg_bot_model.NotifyNewOrderLogicInput) {
	internalId := input.InternalId
	if len(internalId) > 10 {
		internalId = internalId[:10]
	}
	message := fmt.Sprintf(`Новый заказ!
ID: %s
Контактные данные: %s
Товары: `,
		internalId, input.ContactData)

	for i, elem := range input.Products {
		if i == 0 {
			message += elem
		} else {
			message += ", " + elem
		}
	}

	for i := 0; i < 10; i++ {
		_, err := u.bot.Send(&tele.Chat{ID: u.cfg.TgBot.ChatId}, markdownMessage(message), tele.ModeMarkdownV2)
		if err == nil {
			break
		} else {
			fmt.Printf("Error to send message %v: %s", input, err.Error())
		}

		time.Sleep(time.Second * 2)
	}
}

func markdownMessage(message string) string {
	chars := []string{"_", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
	for _, char := range chars {
		message = strings.ReplaceAll(message, char, "\\"+char)
	}
	return message
}
