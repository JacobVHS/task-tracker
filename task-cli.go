package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Task struct {
	TaskID      int    `json:"task_id"`
	Description string `json:"description"`
	CreateTime  string `json:"create_time"`
	UpdateTime  string `json:"update_time"`
	Status      string `json:"status"`
}

func getTaskFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get user home directory: %v", err)
	}
	return filepath.Join(homeDir, ".tasks.json")
}

func readTasks() ([]Task, error) {
	filePath := getTaskFilePath()
	var tasks []Task

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return tasks, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&tasks); err != nil && err.Error() != "EOF" {
		return tasks, err
	}

	return tasks, nil
}

func writeTasks(tasks []Task) error {
	filePath := getTaskFilePath()

	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	helpCmd := flag.NewFlagSet("help", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	mipCmd := flag.NewFlagSet("mark-in-progress", flag.ExitOnError)
	mdCmd := flag.NewFlagSet("mark-done", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("Expected 'add', 'list', 'update', 'delete', 'mark-in-progress', 'mark-done' or 'help' subcommands")
		return
	}

	switch os.Args[1] {
	case "help":
		helpCmd.Parse(os.Args[2:])
		if helpCmd.NArg() > 0 {
			fmt.Printf("Help topic: %s\n", helpCmd.Arg(0))
		} else {
			fmt.Println("Usage: task-cli help [topic]")
		}

	case "add":
		addCmd.Parse(os.Args[2:])
		if addCmd.NArg() > 0 {
			taskDescription := addCmd.Arg(0)
			add(taskDescription)
		} else {
			fmt.Println("Usage: task-cli add [task description]")
		}

	case "list":
		listCmd.Parse(os.Args[2:])
		var status string
		if listCmd.NArg() > 0 {
			status = listCmd.Arg(0)
		} else {
			status = "all"
		}
		list(status)

	case "update":
		updateCmd.Parse(os.Args[2:])
		if updateCmd.NArg() == 2 { // Expecting two arguments: ID and description
			taskID := updateCmd.Arg(0)
			newTitle := updateCmd.Arg(1)
			update(taskID, newTitle)
		} else {
			fmt.Println("Usage: task-cli update [task description]")
		}

	case "delete":
		deleteCmd.Parse(os.Args[2:])
		if deleteCmd.NArg() > 0 {
			taskID := deleteCmd.Arg(0)
			delete(taskID)
		} else {
			fmt.Println("Usage: task-cli delete [taskID]")
		}

	case "mark-in-progress":
		mipCmd.Parse(os.Args[2:])
		if mipCmd.NArg() > 0 {
			taskID := mipCmd.Arg(0)
			status(taskID, "in-progress")
		} else {
			fmt.Println("Usage: task-cli mark-in-progress [taskID]")
		}

	case "mark-done":
		mdCmd.Parse(os.Args[2:])
		if mdCmd.NArg() > 0 {
			taskID := mdCmd.Arg(0)
			status(taskID, "done")
		} else {
			fmt.Println("Usage: task-cli mark-done [taskID]")
		}

	default:
		fmt.Println("Unknown command. Expected 'add', 'list', 'update', 'delete', 'mark-in-progress', 'mark-done' or 'help' subcommands")
	}
}

func add(description string) {
	tasks, err := readTasks()
	if err != nil {
		fmt.Println("Error reading tasks:", err)
		return
	}

	newTaskID := 1
	if len(tasks) > 0 {
		newTaskID = tasks[len(tasks)-1].TaskID + 1
	}

	currentTime := time.Now().Format("2006-01-02 15:04:05")
	newTask := Task{
		TaskID:      newTaskID,
		Description: description,
		CreateTime:  currentTime,
		UpdateTime:  currentTime,
		Status:      "new",
	}

	tasks = append(tasks, newTask)

	if err := writeTasks(tasks); err != nil {
		fmt.Println("Error writing tasks:", err)
		return
	}

	fmt.Println("Task added successfully:", newTask.TaskID)
}

func list(status string) {
	tasks, err := readTasks()
	if err != nil {
		fmt.Println("Error reading tasks:", err)
		return
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}

	for _, task := range tasks {
		if status == "all" || task.Status == status {
			fmt.Printf(
				"Task ID: %d\nDescription: %s\nCreated: %s\nUpdated: %s\nStatus: %s\n\n",
				task.TaskID, task.Description, task.CreateTime, task.UpdateTime, task.Status,
			)
		}
	}
}

func update(expectedTaskID string, newTitle string) {
	tasks, err := readTasks()
	if err != nil {
		fmt.Println("Error reading tasks:", err)
		return
	}

	taskID, err := strconv.Atoi(expectedTaskID)
	if err != nil {
		fmt.Println("Invalid Task ID:", err)
		return
	}

	updated := false
	for i := range tasks {
		if tasks[i].TaskID == taskID {
			tasks[i].Description = newTitle
			tasks[i].UpdateTime = time.Now().Format("2006-01-02 15:04:05")
			updated = true
			break
		}
	}

	if !updated {
		fmt.Printf("Task with ID %d not found.\n", taskID)
		return
	}

	if err := writeTasks(tasks); err != nil {
		fmt.Println("Error writing tasks:", err)
		return
	}

	fmt.Println("Task updated successfully.")
}

func delete(expectedTaskID string) {
	tasks, err := readTasks()
	if err != nil {
		fmt.Println("Error reading tasks:", err)
		return
	}

	taskID, err := strconv.Atoi(expectedTaskID)
	if err != nil {
		fmt.Println("Invalid Task ID:", err)
		return
	}

	updatedTasks := []Task{}
	deleted := false
	for _, task := range tasks {
		if task.TaskID == taskID {
			deleted = true
			continue
		}
		updatedTasks = append(updatedTasks, task)
	}

	if !deleted {
		fmt.Printf("Task with ID %d not found.\n", taskID)
		return
	}

	if err := writeTasks(updatedTasks); err != nil {
		fmt.Println("Error writing tasks:", err)
		return
	}

	fmt.Println("Task deleted successfully.")
}

func status(expectedTaskID string, newStatus string) {
	tasks, err := readTasks()
	if err != nil {
		fmt.Println("Error reading tasks:", err)
		return
	}

	taskID, err := strconv.Atoi(expectedTaskID)
	if err != nil {
		fmt.Println("Invalid Task ID:", err)
		return
	}

	updated := false
	for i := range tasks {
		if tasks[i].TaskID == taskID {
			tasks[i].Status = newStatus
			tasks[i].UpdateTime = time.Now().Format("2006-01-02 15:04:05")
			updated = true
			break
		}
	}

	if !updated {
		fmt.Printf("Task with ID %d not found.\n", taskID)
		return
	}

	if err := writeTasks(tasks); err != nil {
		fmt.Println("Error writing tasks:", err)
		return
	}

	fmt.Println("Task status updated successfully.")
}
