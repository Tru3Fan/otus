package bot

import (
	"otus/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api         *tgbotapi.BotAPI
	taskService service.TaskService
	userService service.UserService
}

func NewBot(token string, taskSvc service.TaskService, userSvc service.UserService) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &Bot{api: api, taskService: taskSvc, userService: userSvc}, nil
}

func (b *Bot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.api.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		b.handleMessage(update.Message)
	}
	return nil
}
