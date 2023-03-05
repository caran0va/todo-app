package main

import (
	"fmt"
	"time"
)

func (t Task) format() string {
	return fmt.Sprintf("%-20s %v\n%-20s %v\n%-20s %v\n", "Task: ", t.title, "Description:", t.descripion, "Creation Time:", t.creationTime)
}

// Create new taskList
func newTaskList(owner string) *TaskList {
	var tl TaskList
	tl.owner = owner
	return &tl
}

// create new task
func (tl *TaskList) newTask(title, desc string) {

	var t = Task{
		title:        title,
		descripion:   desc,
		status:       PENDING,
		creationTime: time.Now(),
	}

	tl.tasks = append(tl.tasks, t)
}
