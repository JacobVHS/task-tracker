package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
			update(taskID, newTitle) // Replace `add` with `update` or similar logic
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

func update(expectedTaskID string, newTitle string) {
	// Open the JSON file
	file, err := os.Open("tasks.json")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// Read the file contents
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Parse the JSON data
	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	// Convert expectedTaskID to integer
	searchTaskID, err := strconv.Atoi(expectedTaskID)
	if err != nil {
		log.Fatalf("Invalid task ID: %v", err)
	}

	// Find and update the task
	found := false
	for i := range tasks {
		if tasks[i].TaskID == searchTaskID {
			tasks[i].Description = newTitle // Update description
			found = true
			break
		}
	}

	// Check if a match was found
	if found {
		// Serialize updated tasks back to JSON
		updatedData, err := json.MarshalIndent(tasks, "", "  ")
		if err != nil {
			log.Fatalf("Failed to serialize updated tasks: %v", err)
		}

		// Write updated tasks back to the file
		if err := ioutil.WriteFile("tasks.json", updatedData, 0644); err != nil {
			log.Fatalf("Failed to write updated tasks to file: %v", err)
		}

		fmt.Println("Updated tasks saved to file.")
	} else {
		fmt.Printf("Task with TaskID %d not found\n", searchTaskID)
	}
}

func delete(expectedTaskID string) {
	// open the JSON file
	file, err := os.Open("tasks.json")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// read the contents
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// parse the json data
	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	// Convert expectedTaskID to integer
	searchTaskID, err := strconv.Atoi(expectedTaskID)
	if err != nil {
		log.Fatalf("Invalid task ID: %v", err)
	}

	// Create a new slice excluding the task with matching TaskID
	var updatedTasks []Task
	found := false
	for _, task := range tasks {
		if task.TaskID == searchTaskID {
			found = true // Mark as found
			continue     // Skip this task
		}
		updatedTasks = append(updatedTasks, task)
	}

	// Check if a match was found
	if found {
		// Serialize updated tasks back to JSON
		updatedData, err := json.MarshalIndent(updatedTasks, "", "  ")
		if err != nil {
			log.Fatalf("Failed to serialize updated tasks: %v", err)
		}

		// Write updated tasks back to the file
		if err := ioutil.WriteFile("tasks.json", updatedData, 0644); err != nil {
			log.Fatalf("Failed to write updated tasks to file: %v", err)
		}

		fmt.Println("Task deleted successfully and updated tasks saved to file.")
	} else {
		fmt.Printf("Task with TaskID %d not found\n", searchTaskID)
	}

}

func status(expectedTaskID string, newStatus string) {
	// Open the JSON file
	file, err := os.Open("tasks.json")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// Read the file contents
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Parse the JSON data
	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	// Convert expectedTaskID to integer
	searchTaskID, err := strconv.Atoi(expectedTaskID)
	if err != nil {
		log.Fatalf("Invalid task ID: %v", err)
	}

	// Find and update the task
	found := false
	for i := range tasks {
		if tasks[i].TaskID == searchTaskID {
			tasks[i].Status = newStatus // Update Status
			found = true
			break
		}
	}
	// Check if a match was found
	if found {
		// Serialize updated tasks back to JSON
		updatedData, err := json.MarshalIndent(tasks, "", "  ")
		if err != nil {
			log.Fatalf("Failed to serialize updated tasks: %v", err)
		}

		// Write updated tasks back to the file
		if err := ioutil.WriteFile("tasks.json", updatedData, 0644); err != nil {
			log.Fatalf("Failed to write updated tasks to file: %v", err)
		}

		fmt.Println("Updated tasks saved to file.")
	} else {
		fmt.Printf("Task with TaskID %d not found\n", searchTaskID)
	}
}
