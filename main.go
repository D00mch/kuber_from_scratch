package main

import (
	"dumch/cube/manager"
	"dumch/cube/task"
	"dumch/cube/worker"
	"fmt"
	"os"
	"strconv"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

func main() {
	whost := os.Getenv("CUBE_WORKER_HOST")
	wport, _ := strconv.Atoi(os.Getenv("CUBE_WORKER_PORT"))
	fmt.Printf("Worker host:port -> %s:%d", whost, wport)

	mhost := os.Getenv("CUBE_MANAGER_HOST")
	mport, _ := strconv.Atoi(os.Getenv("CUBE_MANAGER_PORT"))

	fmt.Println("Starting Cube worker")

	w := worker.Worker{
		Queue: *queue.New(),
		Db:    make(map[uuid.UUID]*task.Task),
	}
	wapi := worker.Api{Address: whost, Port: wport, Worker: &w}

	go w.RunTasks()
	go w.CollectStats()
	go wapi.Start()

	fmt.Println("Starting Cube manager")

	workers := []string{fmt.Sprintf("%s:%d", whost, wport)}
	fmt.Printf("Workers: %v", workers)
	m := manager.New(workers)
	mapi := manager.Api{Address: mhost, Port: mport, Manager: m}

	go m.ProcessTasks()
	go m.UpdateTasks()
	mapi.Start()
}

/*
export CUBE_WORKER_HOST=localhost
export CUBE_WORKER_PORT=5555 
export CUBE_MANAGER_HOST=localhost 
export CUBE_MANAGER_PORT=5556 
*/
