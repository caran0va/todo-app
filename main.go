package main

import (
	"bufio"
	"fmt"
	"os"
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
	taskBook = append(taskBook, *newTaskList("Kaycees Hella Todo List", "Kaycee"))
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
	version := "Version 0.2.0"
	fmt.Print(centerPaddedString(fmt.Sprintf("Welcome to Cara's Todo list application !!\n\n%v", version), '#', 100))
}

func mainMenu() {
	var command []string
	//var command[1] string
	//var command[1]2 int
	scanner := bufio.NewScanner(os.Stdin)
	var selectedItem string = ""
	var selectedIdx int = -1
	var selector string = "  "
menuloop:
	for {
		// reset input fields
		for i, _ := range command {
			command[i] = ""
		}

		fmt.Printf("\nPlease type a command> ")
		scanner.Scan()
		if scanner.Text() != "" {
			command = strings.Fields(scanner.Text())

			switch {
			default:
				fmt.Println()
				greeting()
				fmt.Println()
				fmt.Println("HELP - Display list of available commands.")
				fmt.Println("CREATE -  Create a tasklist or task.")
				fmt.Println("LIST -  Display a list of objects.")
				fmt.Println("DELETE - Unimplemented at this time.")
				fmt.Println("SELECT - Shift focus to an object.")
				fmt.Println("SELECTED -  Display the selected object.")
			case strings.EqualFold(command[0], "create") && len(command) == 2:
				switch {
				case strings.EqualFold(command[1], "task") && selectedItem == "Task List" && selectedIdx != -1:
					fmt.Printf("Task Name: ")
					scanner.Scan()
					taskName := scanner.Text()
					fmt.Printf("Task Description: ")
					scanner.Scan()
					taskDescription := scanner.Text()
					taskBook[selectedIdx].newTask(taskName, taskDescription)
					fmt.Printf("Task %s Created", taskName)
				case strings.EqualFold(command[1], "tasklist"):
					fmt.Printf("Task  List Owner: ")
					scanner.Scan()
					taskListOwner := scanner.Text()
					fmt.Printf("Task List Title: ")
					scanner.Scan()
					taskListName := scanner.Text()
					taskBook = append(taskBook, *newTaskList(taskListName, taskListOwner))
					fmt.Printf("Task List %s Created for %s", taskListName, taskListOwner)
				}

			case strings.EqualFold(command[0], "list") && len(command) == 2:
				switch {
				case strings.EqualFold(command[1], "tasklist"):
					fmt.Printf("%-20v %-40v %-20v %v", "\n  Task List ###:", "Title:", "Owner:", "# of Tasks:\n\n")
					for i, tl := range taskBook {
						if selectedItem == "Task List" && selectedIdx == i {
							selector = "* "
						} else {
							selector = "  "
						}
						fmt.Printf("%-20v%-40v %-20v %v\n", selector+"Tasklist "+strconv.Itoa(i), tl.title, tl.owner, len(tl.tasks))
						selector = "  "
					}
				case strings.EqualFold(command[1], "task"):
					fmt.Printf("%-20v %-40v %v", "\n  Task ###:", "Title:", "Description:\n\n\n")
					for i, task := range taskBook[selectedIdx].tasks {
						fmt.Printf("%-20v%-40v %-20v\n\n", selector+"Task "+strconv.Itoa(i), task.title, task.descripion)
					}

				}
			case strings.EqualFold(command[0], "select") && len(command) == 3:
				switch {
				case strings.EqualFold(command[1], "tasklist"):
					selectedItem = "Task List"
					if arg2, err := strconv.Atoi(command[2]); err == nil && arg2 >= 0 && arg2 <= len(taskBook)-1 {
						selectedIdx = arg2
						fmt.Printf("Task List %d is the selected %s\n", selectedIdx, selectedItem)
					} else if arg2 > len(taskBook)-1 {
						fmt.Printf("Task List %d does not exist, please select a valid index", arg2)
					} else {
						fmt.Printf("Error occured: %v", err)
					}

				}
			case strings.EqualFold(command[0], "delete"):
				// TODO implement deleting thingies
			case strings.EqualFold(command[0], "exit") && len(command) == 1:
				break menuloop
			case strings.EqualFold(command[0], "selected") && len(command) == 1:
				if selectedIdx >= 0 && selectedItem != "" {
					fmt.Printf("%s %d is currently selected\n", selectedItem, selectedIdx)
				} else {
					fmt.Printf("Nothing is selected.")
				}
			}
		}
	}
}
