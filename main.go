package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
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
	ID           string    `json:"taskID"`
	Title        string    `json:"title"`
	Descripion   string    `json:"description"`
	Status       Status    `json:"status"`
	CreationTime time.Time `json:"creationTime"`
	//progressTime time.Duration `json: "progress`
}
type TaskList struct {
	ID    string `json:"listID"`
	Title string `json:"title"`
	Owner string `json:"owner"`
	Tasks []Task `json:"tasks"`
}

func main() {

	//greeting()
	taskBook = append(taskBook, *newTaskList("Caras First Todo List", "Cara"))
	taskBook = append(taskBook, *newTaskList("Kaycees Hella Todo List", "Kaycee"))
	//time.Sleep(2 * time.Second)
	taskBook[0].newTask("New Task #2", "This is something I should also do when I have time. Or rather, I should make time to do this lol")

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	} else if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
	//temporary loop to print tasks from the TaskList carasTasks
	//TODO: implement list function
	//for i, _ := range carasTasks.tasks {
	//	fmt.Println(carasTasks.tasks[i].format())

}

func run(ctx context.Context) error {
	tun, err := ngrok.Listen(ctx,
		config.HTTPEndpoint(),
		ngrok.WithAuthtokenFromEnv(),
	)
	if err != nil {
		return err
	}

	log.Println("tunnel created:", tun.URL())

	router := mux.NewRouter()
	//
	router.HandleFunc("/", func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(200)
		greeting(writer)

	})
	//tasklist CRUD
	router.HandleFunc("/tasklist", getTaskLists).Methods("GET")
	router.HandleFunc("/tasklist/{listID}", getTaskList).Methods("GET")
	router.HandleFunc("/tasklist", createTaskList).Methods("POST")
	router.HandleFunc("/tasklist/{listID}", updateTaskList).Methods("POST")
	router.HandleFunc("/tasklist/{listID}", deleteTaskList).Methods("DELETE")
	//task CRUD
	router.HandleFunc("/tasklist/{listID}/task", getTasks).Methods("GET")
	router.HandleFunc("/tasklist/{listID}/task/{taskID}", getTask).Methods("GET")
	router.HandleFunc("/tasklist/{listID}/task", createTask).Methods("POST")
	router.HandleFunc("/tasklist/{listID}/task/{taskID}", updateTask).Methods("POST")
	router.HandleFunc("/tasklist/{listID}/task/{taskID}", deleteTask).Methods("DELETE")

	return http.Serve(tun, http.HandlerFunc(router.ServeHTTP))
}

func greeting(writer http.ResponseWriter) {
	//screen.Clear()
	version := "Version 0.2.0-ngrok"
	fmt.Fprintf(writer, centerPaddedString(fmt.Sprintf("Welcome to Cara's Todo list application !!\n\n%v", version), '#', 100))
}
