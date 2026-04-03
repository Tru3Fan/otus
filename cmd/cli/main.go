package main

import (
	"bufio"
	"fmt"
	"os"
	"otus/internal/db"
	postgresRepo "otus/internal/repository/postgres"
	"otus/internal/service"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	if err := db.Connect(); err != nil {
		fmt.Println("Error connecting to database", err)
		os.Exit(1)
	}
	defer db.Disconnect()

	taskRepo := postgresRepo.NewTaskRepo()
	userRepo := postgresRepo.NewUserRepo()
	taskSvc := service.NewTaskService(taskRepo)
	userSvc := service.NewUserService(userRepo)

	runCLI(taskSvc, userSvc)
}

func runCLI(taskSvc service.TaskService, userSvc service.UserService) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		printMenu()
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())
		switch input {
		case "1":
			listTasks(taskSvc)
		case "2":
			createTask(scanner, taskSvc)
		case "3":
			updateTask(scanner, taskSvc)
		case "4":
			updateStatus(scanner, taskSvc)
		case "5":
			filterByStatus(scanner, taskSvc)
		case "6":
			deleteTask(scanner, taskSvc)
		case "7":
			listUsers(userSvc)
		case "0":
			fmt.Println("Выход.")
			return
		default:
			fmt.Println("Неверный ввод")

		}
	}
}

func printMenu() {
	fmt.Println("\n=== Task Manager ===")
	fmt.Println("1. Список задач")
	fmt.Println("2. Создать задачу")
	fmt.Println("3. Обновить задачу")
	fmt.Println("4. Изменить статус задачи")
	fmt.Println("5. Фильтр по статусу")
	fmt.Println("6. Удалить задачу")
	fmt.Println("7. Список пользователей")
	fmt.Println("0. Выход")
	fmt.Println("Выберите: ")
}

func listTasks(taskSvc service.TaskService) {
	tasks, err := taskSvc.GetTasks()
	if err != nil {
		fmt.Println("Ошибка", err)
		return
	}
	if len(tasks) == 0 {
		fmt.Println("Задач нет")
		return
	}
	for _, t := range tasks {
		fmt.Printf("[%d] %s | статус: %s\n", t.TaskID, t.Title, t.Status)
	}
}

func listUsers(userSvc service.UserService) {
	users, err := userSvc.GetUsers()
	if err != nil {
		fmt.Println("Ошибка", err)
		return
	}
	if len(users) == 0 {
		fmt.Println("Пользователей нет")
		return
	}
	for _, u := range users {
		fmt.Printf("[%d] %s\n", u.UserID, u.Username)
	}
}

func createTask(scanner *bufio.Scanner, taskSvc service.TaskService) {
	fmt.Print("Название задачи: ")
	scanner.Scan()
	title := strings.TrimSpace(scanner.Text())
	if title == "" {
		fmt.Println("Название не может быть пустым")
		return
	}
	task, err := taskSvc.CreateTask(title, 0)
	if err != nil {
		fmt.Println("Ошибка: ", err)
		return
	}
	fmt.Printf("Создана задача [%d] %s\n", task.TaskID, task.Title)
}

func updateTask(scanner *bufio.Scanner, taskSvc service.TaskService) {
	fmt.Print("ID задачи: ")
	scanner.Scan()
	id, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil {
		fmt.Println("Неверный ID")
		return
	}
	fmt.Print("Новое название: ")
	scanner.Scan()
	title := strings.TrimSpace(scanner.Text())
	if title == "" {
		fmt.Println("Название не может быть пустым")
		return
	}
	task, err := taskSvc.UpdateTask(id, title, 0)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	fmt.Printf("Обновлена задача [%d] %s\n", task.TaskID, task.Title)
}

func deleteTask(scanner *bufio.Scanner, taskSvc service.TaskService) {
	fmt.Print("ID задачи: ")
	scanner.Scan()
	id, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil {
		fmt.Println("Неверный ID")
		return
	}
	if err := taskSvc.DeleteTask(id); err != nil {
		fmt.Println("Ошибка: ", err)
		return
	}
	fmt.Println("Задача удалена")
}

func filterByStatus(scanner *bufio.Scanner, taskSvc service.TaskService) {
	fmt.Print("Статус (pending/in_progress/done): ")
	scanner.Scan()
	status := strings.TrimSpace(scanner.Text())
	tasks, err := taskSvc.GetTasksByStatus(status)
	if err != nil {
		fmt.Println("Ошибка: ", err)
		return
	}
	if len(tasks) == 0 {
		fmt.Println("Задач с таким статусом нет")
		return
	}
	for _, t := range tasks {
		fmt.Printf("[%d] %s | статус: %s\n", t.TaskID, t.Title, t.Status)
	}
}

func updateStatus(scanner *bufio.Scanner, taskSvc service.TaskService) {
	fmt.Print("ID задачи: ")
	scanner.Scan()
	id, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil {
		fmt.Println("Неверный ID")
		return
	}
	fmt.Print("Статус (pending/in_progress/done/cancelled): ")
	scanner.Scan()
	status := strings.TrimSpace(scanner.Text())
	task, err := taskSvc.UpdateTaskStatus(id, status)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	fmt.Printf("Статус задачи [%d] обновлён: %s\n", task.TaskID, task.Status)
}
