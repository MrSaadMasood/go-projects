// Package manager manages tasks
package manager

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Logger interface {
	Log(message string)
}

type status string

var taskStatusEnum = struct {
	queued    status
	running   status
	cancelled status
	completed status
	failed    status
}{
	queued:    "QUEUED",
	running:   "RUNNING",
	cancelled: "CANCELLED",
	completed: "COMPLETED",
	failed:    "FAILED",
}

type Task interface {
	Execute() error
	GetTimeout() *time.Duration
}

type EnqueueTask struct {
	Task  Task
	Delay time.Duration
}

type ExtendedTask struct {
	task   Task
	status status
	timer  *time.Timer
	delay  time.Duration
}
type ExtendedTaskMap map[string]ExtendedTask

type taskManager struct {
	maxConcurrency int
	activeTasks    ExtendedTaskMap
	waitingTasks   ExtendedTaskMap
	logger         Logger
	running        bool
	shutdown       bool
	done           chan bool
	rwMutex        sync.RWMutex
}

func (tm *taskManager) Enqueue(et ...EnqueueTask) {
	if tm.shutdown {
		tm.logger.Log("Task Manager Shutting Down. Cannot Accept More Tasks")
		return
	}

	for _, v := range et {
		tm.safeWrite(func() {
			id := uuid.NewString()
			tm.waitingTasks[id] = ExtendedTask{task: v.Task, status: taskStatusEnum.queued, delay: v.Delay}
		})
	}

	if !tm.running {
		tm.running = true
		tm.Run()
	}
}

func (tm *taskManager) ListTasks() {
	tm.rwMutex.RLock()
	defer tm.rwMutex.RUnlock()

	tm.logger.Log("Listing Active Tasks: " + strconv.Itoa(len(tm.activeTasks)))
	for k, v := range tm.activeTasks {
		fmt.Println("Active task ID: "+k, " VALUE: ", v.task, " Status: ", v.status, " Delay", v.delay)
	}

	tm.logger.Log("Listing Waiting Tasks: " + strconv.Itoa(len(tm.waitingTasks)))
	for k, v := range tm.waitingTasks {
		if v.status != taskStatusEnum.queued {
			continue
		}
		fmt.Println("Waiting task ID: "+k, " VALUE: ", v.task, " Status: ", v.status, " Delay", v.delay)
	}
}

func (tm *taskManager) CancelTask(id string) {
	extendedTask, ok := tm.waitingTasks[id]
	if ok {
		tm.RemoveTask(id)
		tm.logger.Log("Waiting Or Cancelled Task Deleted" + id)
		return
	}
	extendedTask, ok = tm.activeTasks[id]
	if !ok {
		tm.logger.Log("Incorrect Task Id Provided To Cancel Task")
		return
	}

	extendedTask.timer.Stop()

	delete(tm.activeTasks, id)
	extendedTask.status = taskStatusEnum.cancelled
	tm.waitingTasks[id] = extendedTask

	tm.logger.Log("Active Task Stopped and Cancelled: " + id)
}

func (tm *taskManager) RemoveTask(id string) {
	_, wok := tm.waitingTasks[id]
	activeTask, aok := tm.activeTasks[id]
	if wok {
		delete(tm.waitingTasks, id)
		return
	}

	if aok {
		activeTask.timer.Stop()
		delete(tm.activeTasks, id)
		return
	}

	tm.logger.Log("Incorect task id provided to remove task")
}

func (tm *taskManager) GracefulShutdown() {
	tm.shutdown = true
	for _, v := range tm.activeTasks {
		v.timer.Stop()
	}

	defer func() {
		tm.activeTasks = make(ExtendedTaskMap)
		tm.waitingTasks = make(ExtendedTaskMap)
	}()

	tm.logger.Log("Active and waiting task cancelled and removed for graceful shutdown")
}

func (tm *taskManager) Run() {
	if tm.shutdown {
		return
	}

	var wg sync.WaitGroup
	for len(tm.waitingTasks) != 0 {
		safeConcurrency := int(math.Min(float64(tm.maxConcurrency), float64(len(tm.waitingTasks))))
		count := 0
		for k, v := range tm.waitingTasks {
			if count > safeConcurrency {
				break
			}
			count++

			wg.Add(1)

			tm.safeDelete(tm.waitingTasks, k)
			v.status = taskStatusEnum.running
			timer := time.AfterFunc(v.delay, func() {
				tm.execute(k, &v, &wg)
			})
			v.timer = timer
			tm.safeWrite(func() {
				tm.activeTasks[k] = v
			})

		}
		wg.Wait()
		tm.running = false

	}
}

func (tm *taskManager) safeDelete(etm ExtendedTaskMap, id string) {
	tm.rwMutex.Lock()
	defer tm.rwMutex.Unlock()
	delete(etm, id)
}

func (tm *taskManager) safeWrite(f func()) {
	tm.rwMutex.Lock()
	defer tm.rwMutex.Unlock()
	f()
}

func (tm *taskManager) execute(id string, et *ExtendedTask, wg *sync.WaitGroup) {
	done := make(chan bool)

	timeout := et.task.GetTimeout()
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)

	go func() {
		defer func() {
			cancel()
			wg.Done()
		}()

		select {
		case <-ctx.Done():
			tm.safeWrite(func() {
				tm.CancelTask(id)
			})
		case success := <-done:
			if !success {
				et.status = taskStatusEnum.failed
				tm.logger.Log("FAILED: Active Task with id: " + id + " " + string(et.status))
			} else {
				et.status = taskStatusEnum.completed
				tm.logger.Log("Active Task with id: " + id + " " + string(et.status))
			}
		}
	}()

	err := et.task.Execute()
	if err != nil {
		done <- false
		return
	}
	done <- true
}

func CreateNewTaskManager(logger Logger, maxConcurrency int) taskManager {
	return taskManager{activeTasks: make(ExtendedTaskMap), waitingTasks: make(ExtendedTaskMap), logger: logger, maxConcurrency: maxConcurrency, done: make(chan bool)}
}
