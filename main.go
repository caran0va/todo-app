package main

import (
	"fmt"
	"time"
)

type Status = int

var taskBook = make(map[string]TaskList, 0)

// Status
const (
	PENDING Status = iota
	IN_PROGRESS
	DONE
	ARCHIVED
)

type Task struct {
	title        string
	descripion   string
	status       Status
	creationTime time.Time
	progressTime time.Duration
}
type TaskList struct {
	owner string
	tasks []Task
}

func main() {
	greeting()
	var carasTasks = *newTaskList("cara")

	carasTasks.newTask("New Task", "This is something I should do")
	time.Sleep(2 * time.Second)
	carasTasks.newTask("New Task #2", "This is something I should also do when I have time. Or rather, I should make time to do this lol")
	for i, _ := range carasTasks.tasks {
		fmt.Println(carasTasks.tasks[i].format())
	}

}

func greeting() {
	version := "0.0.0a"
	fmt.Print(centerPaddedString(fmt.Sprintf("Welcome to Cara's Todo list application !!\n\n%v", version), '#', 70))
	fmt.Println()
}
