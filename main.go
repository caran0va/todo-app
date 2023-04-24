package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/launchdarkly/go-sdk-common/v3/ldcontext"
	ld "github.com/launchdarkly/go-server-sdk/v6"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

type Status = int

var taskBook = []TaskList{}

var sdk_key = ""

const featureFlagKey = "sample-flag"

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
	taskBook[0].newTask("New Task", "This is something I should also do when I have time. Or rather, I should make time to do this lol")

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	} else if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}

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

	sdk_key = os.Getenv("LD_SDKKEY")
	if sdk_key == "" {
		log.Fatal("SDK Key is empty. please edit .env to have LD_SDKKEY with sdk key")
	}
	client, _ := ld.MakeClient(sdk_key, 5*time.Second)
	if client.Initialized() {
		log.Print("LaunchDarkly SDK successfully initialized!")
	} else {
		log.Fatal("LaunchDarkly SDK failed to initalize...")
	}

	context := ldcontext.NewBuilder("example-user-key").
		Name("Cara").
		Build()

	flagValue, err := client.BoolVariation(featureFlagKey, context, false)
	if err != nil {
		log.Printf("error: " + err.Error())
	}

	log.Printf("Feature flag '%s' is %t for this context", featureFlagKey, flagValue)

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

	if !flagValue { // if the flag is false, then allow editing and deleting entries
		router.HandleFunc("/tasklist/{listID}", updateTaskList).Methods("POST")
		router.HandleFunc("/tasklist/{listID}", deleteTaskList).Methods("DELETE")
		router.HandleFunc("/tasklist/{listID}/task/{taskID}", updateTask).Methods("POST")
		router.HandleFunc("/tasklist/{listID}/task/{taskID}", deleteTask).Methods("DELETE")
	}

	//task CRUD
	router.HandleFunc("/tasklist/{listID}/task", getTasks).Methods("GET")
	router.HandleFunc("/tasklist/{listID}/task/{taskID}", getTask).Methods("GET")
	router.HandleFunc("/tasklist/{listID}/task", createTask).Methods("POST")

	return http.Serve(tun, http.HandlerFunc(router.ServeHTTP))
}

func greeting(writer http.ResponseWriter) {
	//screen.Clear()
	version := "Version 0.2.0-ngrok"
	fmt.Fprintf(writer, centerPaddedString(fmt.Sprintf("Welcome to Cara's Todo list application !!\n\n%v", version), '#', 100))
}
