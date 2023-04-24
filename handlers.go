package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Index
func getTaskLists(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(taskBook)
}

// Show
func getTaskList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, taskList := range taskBook {
		if taskList.ID == params["listID"] {
			json.NewEncoder(w).Encode(taskList)
			return
		}
	}
}

// Create
func createTaskList(w http.ResponseWriter, r *http.Request) {
	if !readOnlyMode {
		w.Header().Set("Content-Type", "application/json")
		var newTaskList TaskList
		json.NewDecoder(r.Body).Decode(&newTaskList) //create new tasklist based on request parameters
		newTaskList.ID = strconv.Itoa(len(taskBook) + 1)
		taskBook = append(taskBook, newTaskList)
		json.NewEncoder(w).Encode(newTaskList)
	}
}

// Update
func updateTaskList(w http.ResponseWriter, r *http.Request) {
	if !readOnlyMode {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		for i, tasklist := range taskBook {
			if tasklist.ID == params["listID"] {
				var editTaskList TaskList
				json.NewDecoder(r.Body).Decode(&editTaskList)
				editTaskList.ID = tasklist.ID
				editTaskList.Tasks = tasklist.Tasks // Tasks are null after edit if i don't have this. not sure why its needed
				taskBook = append(taskBook[:i], taskBook[i+1:]...)
				taskBook = append(taskBook, editTaskList)
				return
			}
		}

	}
}

// Delete
func deleteTaskList(w http.ResponseWriter, r *http.Request) {
	if !readOnlyMode {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		for i, taskList := range taskBook {
			if taskList.ID == params["listID"] {
				taskBook = append(taskBook[:i], taskBook[i+1:]...)
				break
			}
		}
		json.NewEncoder(w).Encode(taskBook)
	}
}

// Index
func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, taskList := range taskBook {
		if taskList.ID == params["listID"] {
			json.NewEncoder(w).Encode(taskList.Tasks)
		}
	}

}

// Show
func getTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, tasklist := range taskBook {
		if tasklist.ID == params["listID"] {
			for _, task := range tasklist.Tasks {
				if task.ID == params["taskID"] {
					json.NewEncoder(w).Encode(task)
				}

			}
		}
	}
}

// Create
func createTask(w http.ResponseWriter, r *http.Request) {
	if !readOnlyMode {
		w.Header().Set("Content-Type", "application/json")

		params := mux.Vars(r)
		for i, _ := range taskBook {
			if taskBook[i].ID == params["listID"] {
				var newTask Task
				json.NewDecoder(r.Body).Decode(&newTask) //create new task based on request body parameters
				newTask.CreationTime = time.Now()
				newTask.ID = strconv.Itoa(len(taskBook[i].Tasks) + 1)
				taskBook[i].Tasks = append(taskBook[i].Tasks, newTask)
				fmt.Printf("Tasks:%v", taskBook)
				json.NewEncoder(w).Encode(newTask)
			}

		}
	}
}

// Update
func updateTask(w http.ResponseWriter, r *http.Request) {
	if !readOnlyMode {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)

		for i, _ := range taskBook {
			if taskBook[i].ID == params["listID"] {
				for j, _ := range taskBook[i].Tasks {
					if taskBook[i].Tasks[j].ID == params["taskID"] {
						var editTask Task
						editTask.CreationTime = taskBook[i].Tasks[j].CreationTime
						json.NewDecoder(r.Body).Decode(&editTask)                                     //edit existing task based on request body parameters
						taskBook[i].Tasks = append(taskBook[i].Tasks[:j], taskBook[i].Tasks[j+1:]...) //remove the task to be editted
						editTask.ID = params["taskID"]                                                //copy over details from OG task
						taskBook[i].Tasks = append(taskBook[i].Tasks, editTask)                       //add edited task into the tasklist
						return
					}
				}
			}
		}
	}
}

// Delete
func deleteTask(w http.ResponseWriter, r *http.Request) {
	if !readOnlyMode {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		for i, taskList := range taskBook {
			if taskList.ID == params["listID"] {
				for j, _ := range taskBook[i].Tasks {
					if taskBook[i].Tasks[j].ID == params["taskID"] {
						taskBook[i].Tasks = append(taskBook[i].Tasks[:j], taskBook[i].Tasks[j+1:]...)
						break
					}
				}
			}
		}
		json.NewEncoder(w).Encode(taskBook)
	}
}
