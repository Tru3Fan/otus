package bot

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.api.GetUpdatesChan(u)
	for update := range updates {
		if update.CallbackQuery != nil {
			b.handleCallback(update.CallbackQuery)
			continue
		}
		if update.Message != nil {
			b.handleMessage(update.Message)
		}
	}
	return nil
}

func (b *Bot) handleMessage(msg *tgbotapi.Message) {
	if msg.IsCommand() && msg.Command() == "start" {
		b.handleStart(msg)
		return
	}
	if msg.IsCommand() && msg.Command() == "adduser" && msg.From.ID == b.adminID {
		b.handleAddUser(msg)
		return
	}
	if !b.isRegistered(msg.From.ID) {
		b.send(msg.Chat.ID, "Нет доступа. Обратитесь к администратору.")
		return
	}

	state := b.getState(msg.From.ID)
	switch state.Step {
	case StepWaitingTaskTitle:
		state.Draft.Title = strings.TrimSpace(msg.Text)
		state.Step = StepWaitingTaskBody
		b.sendWithKeyboard(msg.Chat.ID, "Описание задачи (или пропустить):", b.skipKeyboard())
	case StepWaitingTaskBody:
		state.Draft.Body = strings.TrimSpace(msg.Text)
		state.Step = StepWaitingAssignee
		b.sendAssigneeChoice(msg.Chat.ID, msg.From.ID)
	default:
		b.sendMainMenu(msg.Chat.ID)
	}
}

func (b *Bot) handleStart(msg *tgbotapi.Message) {
	_, err := b.userService.GetUserByTelegramID(msg.Chat.ID)
	if err == nil {
		b.sendMainMenu(msg.Chat.ID)
		return
	}

	ok, err := b.userService.IsPendingUser(msg.From.UserName)
	if err != nil || !ok {
		b.send(msg.Chat.ID, "Нет доступа. Обратитесь к администратору")
		return
	}
	_, err = b.userService.ConfirmPendingUser(msg.From.ID, msg.From.UserName)
	if err != nil {
		b.send(msg.Chat.ID, "Ошибка регистрации: "+err.Error())
		return
	}
	b.send(msg.Chat.ID, "Добро пожаловать, @"+msg.From.FirstName+"!")
	b.sendMainMenu(msg.Chat.ID)
}
func (b *Bot) handleAddUser(msg *tgbotapi.Message) {
	username := strings.TrimPrefix(msg.CommandArguments(), "@")
	if username == "" {
		b.send(msg.Chat.ID, "Формат: /adduser @username")
		return
	}
	err := b.userService.AddPendingUser(username)
	if err != nil {
		b.send(msg.Chat.ID, "Ошибка: "+err.Error())
		return
	}
	b.send(msg.Chat.ID, "Пользователь @"+username+" добавлен в список ожидания.")
}

func (b *Bot) handleCallback(cb *tgbotapi.CallbackQuery) {
	b.api.Request(tgbotapi.NewCallback(cb.ID, ""))
	chatID := cb.Message.Chat.ID
	userID := cb.From.ID

	if !b.isRegistered(userID) {
		b.send(chatID, "Нет доступа.")
		return
	}
	switch cb.Data {
	case "menu":
		b.sendMainMenu(chatID)
	case "new_task":
		b.getState(userID).Step = StepWaitingTaskTitle
		b.send(chatID, "Введите название задачи")
	case "my_tasks":
		b.showMyTasks(chatID, userID)
	case "assigned_by_me":
		b.showAssignedByMe(chatID, userID)
	case "assigned_to_me":
		b.showAssignedToMe(chatID, userID)
	case "archive":
		b.showArchive(chatID, userID)
	case "skip_body":
		state := b.getState(userID)
		state.Draft.Body = ""
		state.Step = StepWaitingAssignee
		b.sendAssigneeChoice(chatID, userID)
	default:
		switch {
		case strings.HasPrefix(cb.Data, "assignee_"):
			b.handleAssigneeChoice(chatID, userID, cb.Data)
		case strings.HasPrefix(cb.Data, "deadline_"):
			b.handleDeadlineChoice(chatID, userID, cb.Data)
		case strings.HasPrefix(cb.Data, "task_"):
			id, _ := strconv.Atoi(strings.TrimPrefix(cb.Data, "task_"))
			b.showTaskDetail(chatID, id)
		case strings.HasPrefix(cb.Data, "close_task_"):
			id, _ := strconv.Atoi(strings.TrimPrefix(cb.Data, "close_task_"))
			b.handleCloseTask(chatID, userID, id)
		case strings.HasPrefix(cb.Data, "inprogress_task_"):
			id, _ := strconv.Atoi(strings.TrimPrefix(cb.Data, "inprogress_task_"))
			b.handleInProgressTask(chatID, userID, id)
		case strings.HasPrefix(cb.Data, "cancel_task_"):
			id, _ := strconv.Atoi(strings.TrimPrefix(cb.Data, "cancel_task_"))
			b.handleCancelTask(chatID, userID, id)
		}
	}
}

func (b *Bot) showMyTasks(chatID int64, telegramID int64) {
	user, err := b.userService.GetUserByTelegramID(telegramID)
	if err != nil {
		b.send(chatID, "Пользователь не найден.")
		return
	}
	task, err := b.taskService.GetTasksByUser(user.UserID)
	if err != nil || len(task) == 0 {
		b.send(chatID, "Задач нет.")
		return
	}
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, t := range task {

		if t.Status == "done" || t.Status == "cancelled" {
			continue
		}

		label := fmt.Sprintf("#%d %s [%s]", t.TaskID, t.Title, statusText(t.Status))
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(label, fmt.Sprintf("task_%d", t.TaskID)),
		))
	}
	if len(rows) == 0 {
		b.send(chatID, "Активных задач нет.")
		return
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("« Меню", "menu"),
	))
	b.sendWithKeyboard(chatID, "Ваши задачи:", tgbotapi.NewInlineKeyboardMarkup(rows...))
}

func (b *Bot) showAssignedByMe(chatID int64, telegramID int64) {
	user, err := b.userService.GetUserByTelegramID(telegramID)
	if err != nil {
		b.send(chatID, "Пользователь не найден.")
		return
	}
	task, err := b.taskService.GetTasksByAuthor(user.UserID)
	if err != nil || len(task) == 0 {
		b.send(chatID, "Задач нет.")
		return
	}
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, t := range task {
		label := fmt.Sprintf("#%d %s [%s]", t.TaskID, t.Title, statusText(t.Status))
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(label, fmt.Sprintf("task_%d", t.TaskID)),
		))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("« Меню", "menu"),
	))
	b.sendWithKeyboard(chatID, "Задачи, которые вы назначили:", tgbotapi.NewInlineKeyboardMarkup(rows...))
}

func (b *Bot) showAssignedToMe(chatID int64, telegramID int64) {
	user, err := b.userService.GetUserByTelegramID(telegramID)
	if err != nil {
		b.send(chatID, "Пользователь не найден.")
		return
	}
	all, err := b.taskService.GetTasksByUser(user.UserID)
	if err != nil {
		b.send(chatID, "Ошибка получения задач.")
		return
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, t := range all {
		if t.AssignedBy != 0 && t.AssignedBy != user.UserID {
			label := fmt.Sprintf("#%d %s [%s]", t.TaskID, t.Title, statusText(t.Status))
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(label, fmt.Sprintf("task_%d", t.TaskID)),
			))
		}
	}
	if len(rows) == 0 {
		b.send(chatID, "Задач нет.")
		return
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("« Меню", "menu"),
	))
	b.sendWithKeyboard(chatID, "Задачи, назначенные вам:", tgbotapi.NewInlineKeyboardMarkup(rows...))
}

func (b *Bot) showArchive(chatID int64, telegramID int64) {
	user, err := b.userService.GetUserByTelegramID(telegramID)
	if err != nil {
		b.send(chatID, "Пользователь не найден.")
		return
	}
	all, err := b.taskService.GetTasksByUser(user.UserID)
	if err != nil {
		b.send(chatID, "Ошибка получения задач.")
		return
	}
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, t := range all {
		if t.Status == "done" || t.Status == "cancelled" {
			label := fmt.Sprintf("#%d %s [%s]", t.TaskID, t.Title, statusText(t.Status))
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(label, fmt.Sprintf("task_%d", t.TaskID)),
			))
		}
	}
	if len(rows) == 0 {
		b.send(chatID, "Архив пуст.")
		return
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("« Меню", "menu"),
	))
	b.sendWithKeyboard(chatID, "Архив задач:", tgbotapi.NewInlineKeyboardMarkup(rows...))
}

func (b *Bot) showTaskDetail(chatID int64, taskID int) {
	t, err := b.taskService.GetTask(taskID)
	if err != nil {
		b.send(chatID, "Задача не найдена.")
		return
	}
	deadline := "не указан"
	if t.Deadline != nil {
		deadline = t.Deadline.Format("02.01.2006")
	}
	text := fmt.Sprintf("📌 #%d %s\nСтатус: %s\nДедлайн: %s\n%s", t.TaskID, t.Title, statusText(t.Status), deadline, t.Body)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("▶️ В работу", fmt.Sprintf("inprogress_task_%d", t.TaskID)),
			tgbotapi.NewInlineKeyboardButtonData("✅ Закрыть", fmt.Sprintf("close_task_%d", t.TaskID)),
			tgbotapi.NewInlineKeyboardButtonData("❌ Отменить", fmt.Sprintf("cancel_task_%d", t.TaskID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("« Назад", "menu"),
		),
	)
	b.sendWithKeyboard(chatID, text, keyboard)
}

func (b *Bot) sendAssigneeChoice(chatID int64, telegramID int64) {
	users, err := b.userService.GetUsers()
	if err != nil {
		b.send(chatID, "Ошибка получения пользователей.")
		return
	}
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, u := range users {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("@"+u.Username, fmt.Sprintf("assignee_%d", u.UserID)),
		))
	}
	b.sendWithKeyboard(chatID, "Выберете исполнителя:", tgbotapi.NewInlineKeyboardMarkup(rows...))
}

func (b *Bot) handleAssigneeChoice(chatID int64, telegramID int64, data string) {
	assigneeID, err := strconv.Atoi(strings.TrimPrefix(data, "assignee_"))
	if err != nil {
		b.send(chatID, "Ошибка.")
		return
	}
	state := b.getState(telegramID)
	state.Draft.AssigneeID = assigneeID
	state.Step = StepWaitingDeadline
	b.sendDeadlineChoice(chatID)
}

func (b *Bot) sendDeadlineChoice(chatID int64) {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Сегодня", "deadline_today"),
			tgbotapi.NewInlineKeyboardButtonData("Завтра", "deadline_tomorrow"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("+3дня", "deadline_3days"),
			tgbotapi.NewInlineKeyboardButtonData("+неделя", "deadline_week"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Без срока", "deadline_none"),
		),
	)
	b.sendWithKeyboard(chatID, "Выберите дедлайн:", keyboard)
}

func (b *Bot) handleDeadlineChoice(chatID int64, telegramID int64, data string) {
	now := time.Now()
	var deadline *time.Time
	switch data {
	case "deadline_today":
		d := now
		deadline = &d
	case "deadline_tomorrow":
		d := now.AddDate(0, 0, 1)
		deadline = &d
	case "deadline_3days":
		d := now.AddDate(0, 0, 3)
		deadline = &d
	case "deadline_week":
		d := now.AddDate(0, 0, 7)
		deadline = &d
	case "deadline_none":
		deadline = nil
	}

	state := b.getState(telegramID)
	author, err := b.userService.GetUserByTelegramID(telegramID)
	if err != nil {
		b.send(chatID, "Ошибка: пользователь не найден")
		return
	}
	t, err := b.taskService.CreateTaskFull(state.Draft.Title, state.Draft.Body, state.Draft.AssigneeID, author.UserID, deadline)
	if err != nil {
		b.send(chatID, "Ошибка создания задачи: "+err.Error())
		return
	}
	b.send(chatID, fmt.Sprintf("Задача #%d создана!", t.TaskID))
	b.notifyAssignee(state.Draft.AssigneeID, t.TaskID, t.Title)
	b.sendMainMenu(chatID)
	b.resetState(telegramID)
}

func (b *Bot) handleCloseTask(chatID int64, telegramID int64, taskID int) {
	t, err := b.taskService.CloseTask(taskID)
	if err != nil {
		b.send(chatID, "Ошибка: "+err.Error())
		return
	}
	b.send(chatID, fmt.Sprintf("Задача #%d закрыта.", t.TaskID))
	b.notifyAuthor(t.AssignedBy, t.TaskID, t.Title)
	b.sendMainMenu(chatID)
}

func (b *Bot) handleCancelTask(chatID int64, telegramID int64, taskID int) {
	t, err := b.taskService.UpdateTaskStatus(taskID, "cancelled")
	if err != nil {
		b.send(chatID, "Ошибка: "+err.Error())
		return
	}
	b.send(chatID, fmt.Sprintf("Задача #%d отменена.", t.TaskID))
	b.sendMainMenu(chatID)
}

func (b *Bot) handleInProgressTask(chatID int64, telegramID int64, taskID int) {
	t, err := b.taskService.UpdateTaskStatus(taskID, "in_progress")
	if err != nil {
		b.send(chatID, "Ошибка: "+err.Error())
		return
	}
	b.send(chatID, fmt.Sprintf("Задача #%d взята в работу.", t.TaskID))
	b.sendMainMenu(chatID)
}

func (b *Bot) notifyAssignee(assigneeID int, taskID int, title string) {
}
func (b *Bot) notifyAuthor(authorID int, taskID int, title string) {

}

func (b *Bot) sendMainMenu(chatID int64) {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📋 Мои задачи", "my_tasks"),
			tgbotapi.NewInlineKeyboardButtonData("✏️ Новая задача", "new_task"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📤 Я назначил", "assigned_by_me"),
			tgbotapi.NewInlineKeyboardButtonData("📥 Мне назначили", "assigned_to_me"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📦 Архив", "archive"),
		),
	)
	b.sendWithKeyboard(chatID, "Главное меню:", keyboard)
}

func (b *Bot) sendWithKeyboard(chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	b.api.Send(msg)
}
func (b *Bot) skipKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Пропустить", "skip_body"),
		),
	)
}

func (b *Bot) isRegistered(telegramID int64) bool {
	_, err := b.userService.GetUserByTelegramID(telegramID)
	return err == nil
}
func (b *Bot) send(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	b.api.Send(msg)
}

func statusText(s string) string {
	switch s {
	case "pending":
		return "ожидает"
	case "in_progress":
		return "в работе"
	case "done":
		return "выполнено"
	case "cancelled":
		return "отменена"
	default:
		return s
	}
}
