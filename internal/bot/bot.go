package bot

import (
	"otus/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api         *tgbotapi.BotAPI
	taskService service.TaskService
	userService service.UserService
	adminID     int64
	state       map[int64]*UserState
}

func NewBot(token string, taskSvc service.TaskService, userSvc service.UserService, adminID int64) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &Bot{api: api, taskService: taskSvc, userService: userSvc, adminID: adminID, state: make(map[int64]*UserState)}, nil
}
