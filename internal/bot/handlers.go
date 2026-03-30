package bot

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleMessage(msg *tgbotapi.Message) {
	if !b.isAllowed(msg.From.ID) {
		b.send(msg.Chat.ID, "Доступ запрещён")
		return
	}

	switch msg.Command() {
	case "start":
		b.send(msg.Chat.ID, "Привет! /help - список команд")
	case "help":
		b.send(msg.Chat.ID, "/tasks — мои задачи\n/newtask <название> — создать\n/status <id> <статус> — обновить статус\n/done <id> — выполнена")
	case "tasks":
		b.handleTasks(msg)
	case "newtask":
		b.handleNewTask(msg)
	case "status":
		b.handleStatus(msg)
	case "done":
		b.handleDone(msg)
	default:
		b.send(msg.Chat.ID, "Неизвестная команда. /help - список команд")
	}
}

func (b *Bot) handleTasks(msg *tgbotapi.Message) {
	user, err := b.userService.GetUserByTelegramID(msg.From.ID)
	if err != nil {
		b.send(msg.Chat.ID, "Пользователь не найден.")
		return
	}
	tasks, err := b.taskService.GetTasksByUser(user.UserID)
	if err != nil {
		b.send(msg.Chat.ID, "Ошибка получения задач.")
		return
	}
	if len(tasks) == 0 {
		b.send(msg.Chat.ID, "У вас нет задач.")
		return
	}
	var sb strings.Builder
	for _, t := range tasks {
		sb.WriteString(fmt.Sprintf("#%d %s [%s]\n", t.TaskID, t.Title, t.Status))
	}
	b.send(msg.Chat.ID, sb.String())
}

func (b *Bot) handleNewTask(msg *tgbotapi.Message) {
	title := msg.CommandArguments()
	if title == "" {
		b.send(msg.Chat.ID, "Укажи название: /newtask <название>")
		return
	}
	user, err := b.userService.GetUserByTelegramID(msg.From.ID)
	if err != nil {
		b.send(msg.Chat.ID, "Пользователь не найден.")
		return
	}
	t, err := b.taskService.CreateTask(title, user.UserID)
	if err != nil {
		b.send(msg.Chat.ID, "Ошибка создания задачи.")
		return
	}
	b.send(msg.Chat.ID, fmt.Sprintf("Задача #%d создана: %s", t.TaskID, t.Title))
}

func (b *Bot) handleStatus(msg *tgbotapi.Message) {
	args := strings.Fields(msg.CommandArguments())
	if len(args) != 2 {
		b.send(msg.Chat.ID, "Формат: /status <id> <pending|in_progress|done>")
		return
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		b.send(msg.Chat.ID, "Неверный id задачи.")
		return
	}
	t, err := b.taskService.UpdateTaskStatus(id, args[1])
	if err != nil {
		b.send(msg.Chat.ID, "Ошибка: "+err.Error())
		return
	}
	b.send(msg.Chat.ID, fmt.Sprintf("Задача #%d: статус обновлен на %s", t.TaskID, t.Status))
}

func (b *Bot) handleDone(msg *tgbotapi.Message) {
	id, err := strconv.Atoi(msg.CommandArguments())
	if err != nil {
		b.send(msg.Chat.ID, "Формат: /done <id>")
		return
	}
	t, err := b.taskService.UpdateTaskStatus(id, "done")
	if err != nil {
		b.send(msg.Chat.ID, "Ошибка: "+err.Error())
		return
	}
	b.send(msg.Chat.ID, fmt.Sprintf("Задача #%d выполнена", t.TaskID))
}

func (b *Bot) isAllowed(telegramID int64) bool {
	_, err := b.userService.GetUserByTelegramID(telegramID)
	return err == nil
}
func (b *Bot) send(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	b.api.Send(msg)
}
