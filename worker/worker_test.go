package worker

import (
	"dumch/cube/task"
	"fmt"
	"testing"
	"time"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

func TestWoker(test *testing.T) {

	db := make(map[uuid.UUID]*task.Task)
	w := Worker{
		Queue: *queue.New(),
		Db:    db,
	}

	t := task.Task{
		ID:    uuid.New(),
		Name:  "test-container-1",
		State: task.Scheduled,
		Image: "strm/helloworld-http",
	}

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
}
