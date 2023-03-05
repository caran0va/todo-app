package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Status = int

var taskBook = []TaskList{}

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
	title string
	owner string
	tasks []Task
}

func main() {
	greeting()
	taskBook = append(taskBook, *newTaskList("Caras First Todo List", "Cara"))
	//time.Sleep(2 * time.Second)
	taskBook[0].newTask("New Task #2", "This is something I should also do when I have time. Or rather, I should make time to do this lol")

	//temporary loop to print tasks from the TaskList carasTasks
	//TODO: implement list function
	//for i, _ := range carasTasks.tasks {
	//	fmt.Println(carasTasks.tasks[i].format())
	//}
	mainMenu()
}

func greeting() {
	//screen.Clear()
	version := "Version 0.1.0"
	fmt.Print(centerPaddedString(fmt.Sprintf("Welcome to Cara's Todo list application !!\n\n%v", version), '#', 100))
}

func mainMenu() {
	var command string
	var arg string
	var arg2 int

	var selectedItem string = ""
	var selectedIdx int = -1
menuloop:
	for {
		fmt.Printf("\nPlease type a command> ")
		fmt.Scanf("%s %s %d", &command, &arg, &arg2)
		//fmt.Printf("%v %v\n", command, arg)

		//fmt.Println(arg)
		switch {
		case strings.EqualFold(command, "help"):

		case strings.EqualFold(command, "list"):
			switch {
			case strings.EqualFold(arg, "tasklist"):
				fmt.Printf("%-20v %-40v %-20v %v", "\nTask List ###:", "Title:", "Owner:", "# of Tasks:\n")
				for i, tl := range taskBook {
					fmt.Printf("%-20s %-40v %-20v %v\n\n", "Tasklist "+strconv.Itoa(i), tl.title, tl.owner, len(tl.tasks))
				}
			case strings.EqualFold(arg, "task"):
			}
		case strings.EqualFold(command, "select"):
			switch {
			case strings.EqualFold(arg, "tasklist"):
				selectedItem = "tasklist"
				if arg2 >= 0 {
					selectedIdx = arg2
					fmt.Printf("Task List %d is the selected %s\n", selectedIdx, selectedItem)
				}
			}
		case strings.EqualFold(command, "delete"):

		case strings.EqualFold(command, "exit"):
			break menuloop
		}
		command = ""
		arg = ""
		arg2 = -1
	}

}
