package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

const featureFlagKey = "read-only-mode"

var readOnlyMode = false
var lastFlag = !readOnlyMode
var flagChan = make(chan bool)
var client = new(ld.LDClient)

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

	SetupCloseHandler()
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
	client, _ = ld.MakeClient(sdk_key, 5*time.Second)
	if client.Initialized() {
		log.Print("LaunchDarkly SDK successfully initialized!")
	} else {
		log.Fatal("LaunchDarkly SDK failed to initalize...")
	}

	context := ldcontext.NewBuilder("example-user-key").
		Name("Cara").
		Build()
	go checkFlag(context)

	go listenToFlag()

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
	router.HandleFunc("/tasklist/{listID}/task/{taskID}", updateTask).Methods("POST")
	router.HandleFunc("/tasklist/{listID}/task/{taskID}", deleteTask).Methods("DELETE")

	//task CRUD
	router.HandleFunc("/tasklist/{listID}/task", getTasks).Methods("GET")
	router.HandleFunc("/tasklist/{listID}/task/{taskID}", getTask).Methods("GET")
	router.HandleFunc("/tasklist/{listID}/task", createTask).Methods("POST")

	return http.Serve(tun, http.HandlerFunc(router.ServeHTTP))
}

func greeting(writer http.ResponseWriter) {
	//screen.Clear()
	version := "Version 0.2.0-ngrok"
	fmt.Fprint(writer, centerPaddedString(fmt.Sprintf("Welcome to Cara's Todo list application !!\n\n%v", version), '#', 100))
}

func listenToFlag() {
	for {
		readOnlyMode = <-flagChan
		if readOnlyMode != lastFlag {
			//when the flag channel changes, let us know
			log.Printf("Feature flag '%s' is %t for this context", featureFlagKey, readOnlyMode)
			lastFlag = readOnlyMode
		}

	}

}

func checkFlag(context ldcontext.Context) {
	for {
		time.Sleep(time.Second)
		flagValue, err := client.BoolVariation(featureFlagKey, context, false)
		if err != nil {
			log.Printf("error: " + err.Error())
		}

		//write new value to the channel
		flagChan <- flagValue
	}
}

func SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		onClose()
		os.Exit(0)
	}()
}

func onClose() {
	// Here we ensure that the SDK shuts down cleanly and has a chance to deliver analytics
	// events to LaunchDarkly before the program exits. If analytics events are not delivered,
	// the context attributes and flag usage statistics will not appear on your dashboard. In
	// a normal long-running application, the SDK would continue running and events would be
	// delivered automatically in the background.
	client.Close()
}
