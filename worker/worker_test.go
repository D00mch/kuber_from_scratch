package worker

import (
	"dumch/cube/task"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

func TestWoker(test *testing.T) {

	w1, w2 := newWorker(), newWorker()
	t1, t2 := newTask(1), newTask(2)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go startTaskOnWorker(w1, t1, &wg)
	go startTaskOnWorker(w2, t2, &wg)

	wg.Wait()
	fmt.Println("Finished")
}

func TestApi(test *testing.T) {
	host := "localhost" // os.Getenv("CUBE_HOST")
	port := 5555

	fmt.Println("Starting Cube worker")
	w := newWorker()

	api := Api{Address: host, Port: port, Worker: &w}

	go runTasks(&w)
	api.Start()
}

func runTasks(w *Worker) {
	for {
		if w.Queue.Len() != 0 {
			result := w.RunTask()
			if result.Error != nil {
				log.Printf("Error running task: %v\n", result.Error)
			}
		} else {
			log.Println("No tasks to process currently.")
		}
		log.Println("Sleeping for 10 seconds.")
		time.Sleep(10 * time.Second)
	}
}

func startTaskOnWorker(w Worker, t task.Task, wg *sync.WaitGroup) {
	fmt.Println("starting task")
	w.AddTask(t)
	result := w.RunTask()
	if result.Error != nil {
		panic(result.Error)
	}

	t.ContainerID = result.ContainerId
	fmt.Printf("task %s is running in container %s\n", t.ID, t.ContainerID)
	fmt.Println("Sleepy time")
	time.Sleep(time.Second * 10)

	fmt.Printf("stopping task %s\n", t.ID)
	t.State = task.Completed
	w.AddTask(t)
	result2 := w.RunTask()
	if result2.Error != nil {
		panic(result2.Error)
	}
	wg.Done()
}

func newWorker() Worker {
	db := make(map[uuid.UUID]*task.Task)
	return Worker{
		Queue: *queue.New(),
		Db:    db,
	}
}

func newTask(number int) task.Task {
	return task.Task{
		ID:    uuid.New(),
		Name:  fmt.Sprintf("test-container-%d", number),
		State: task.Scheduled,
		Image: "strm/helloworld-http",
	}
}
