package main

import (
	"fmt"
	"strconv"
	"time"
)

func (t Task) format() string {
	return fmt.Sprintf("%-20s %v\n%-20s %v\n%-20s %v\n", "Task: ", t.Title, "Description:", t.Descripion, "Creation Time:", t.CreationTime)
}

// Create new taskList
func newTaskList(title, owner string) *TaskList {
	var tl = TaskList{
		ID:    strconv.Itoa(len(taskBook) + 1),
		Title: title,
		Owner: owner,
	}
	return &tl
}

// create new task
func (tl *TaskList) newTask(title, desc string) {

	var t = Task{
		ID:           strconv.Itoa(len(tl.Tasks) + 1),
		Title:        title,
		Descripion:   desc,
		Status:       PENDING,
		CreationTime: time.Now(),
	}

	tl.Tasks = append(tl.Tasks, t)
}
