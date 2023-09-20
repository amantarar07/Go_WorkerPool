// server.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// var upgrader = websocket.Upgrader{} // use default options

type RequestData struct {
	JobCount    int `json:"jobCount"`
	WorkerCount int `json:"workerCount"`
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade our raw HTTP connection to a websocket based one
	fmt.Println("Inside Socket Handler")

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow connections from any origin
		},
	}
	
	// Upgrade the HTTP connection to a WebSocket connection

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error during connection upgradation:", err)
		return
	}
	JobCountSTRING := r.URL.Query().Get("JobCount")
	jobsCount, _ := strconv.Atoi(JobCountSTRING)

	workerCountSTRING := r.URL.Query().Get("WorkerCount")
	workerCount, _ := strconv.Atoi(workerCountSTRING)
	var jobs = make(chan int)

	go func() {

		for i := 1; i <= jobsCount; i++ {

			jobs <- i
		}

	}()

	GiveToWorkers(conn, jobs, workerCount, jobsCount)

}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Index Page")
}

func GiveToWorkers(conn *websocket.Conn, Jobs chan int, workerCount int, JobCount int) {

	var ProgressMutex = &sync.Mutex{}

	for i := 1; i <= workerCount; i++ {

		go Work(conn, Jobs, JobCount,ProgressMutex)

	}
}

func Work(conn *websocket.Conn, Jobs chan int, JobCount int, ProgressMutex *sync.Mutex) {

	fmt.Println("Now Working .....")
	for i := 1; i <= JobCount; i++ {

		time.Sleep(time.Second)
		ProgressMutex.Lock()
		var message = []byte(fmt.Sprintf("Job %d is Done", <-Jobs))
		
		conn.WriteMessage(websocket.TextMessage, message)
		ProgressMutex.Unlock()
	}

	defer conn.Close()
}

func main() {

	http.HandleFunc("/socket", socketHandler)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
