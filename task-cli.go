package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

type Task struct {
	TaskID      int    `json:"task_id"`
	Description string `json:"description"`
	CreateTime  string `json:"create_time"`
	UpdateTime  string `json:"update_time"`
	Status      string `json:"status"`
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
		if updateCmd.NArg() > 0 {
			update_state := updateCmd.Arg(0)
			add(update_state)
		} else {
			fmt.Println("Usage: task-cli update [task description]")
		}

	default:
		fmt.Println("Unknown command. Expected 'add', 'list', 'update', 'delete', 'mark-in-progress', 'mark-done' or 'help' subcommands")
	}
}

func add(description string) {
	// Define the file to store tasks
	fileName := "tasks.json"

	// Read existing tasks from file
	var tasks []Task
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&tasks); err != nil && err.Error() != "EOF" {
		fmt.Println("Error reading tasks:", err)
		return
	}

	// Create new task
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

	// Append new task
	tasks = append(tasks, newTask)

	// Write updated tasks back to file
	file.Truncate(0)
	file.Seek(0, 0)
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(tasks); err != nil {
		fmt.Println("Error writing tasks:", err)
		return
	}

	fmt.Println("Task added successfully:", newTask.TaskID)
}

func list(status string) {
	// List tasks based on their status
	fmt.Println("Listing Tasks")
	fmt.Println("-------------")
	fileName := "tasks.json"

	// Read existing tasks from the file
	var tasks []Task
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&tasks); err != nil && err.Error() != "EOF" {
		fmt.Println("Error reading tasks:", err)
		return
	}

	// Check if there are any tasks
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}

	// Print tasks
	filterStatus := status // Replace with the status you want to match

	if filterStatus == "all" {
		for _, task := range tasks {
			fmt.Printf(
				"Task ID: %d\nDescription: %s\nCreated: %s\nUpdated: %s\nStatus: %s\n\n",
				task.TaskID, task.Description, task.CreateTime, task.UpdateTime, task.Status,
			)
		}
	} else {
		for _, task := range tasks {
			if task.Status == filterStatus {
				fmt.Printf(
					"Task ID: %d\nDescription: %s\nCreated: %s\nUpdated: %s\nStatus: %s\n\n",
					task.TaskID, task.Description, task.CreateTime, task.UpdateTime, task.Status,
				)
			}
		}
	}
}
