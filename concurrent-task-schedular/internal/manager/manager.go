// Package manager manages tasks
package manager

import (
	"context"
	"math"
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
	Cancel()
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
	shutdownCtx    context.Context
	shutDown       context.CancelFunc
	done           chan bool
	rwMutex        sync.RWMutex
	wg             sync.WaitGroup
}

func (tm *taskManager) Enqueue(et ...EnqueueTask) {
	select {
	case <-tm.shutdownCtx.Done():
		tm.logger.Log("Task Manager Shutting Down. Cannot Accept More Tasks")
		return
	default:
		for _, v := range et {
			id := uuid.NewString()
			tm.waitingTasks[id] = ExtendedTask{task: v.Task, status: taskStatusEnum.queued, delay: v.Delay}
		}

		tm.Run()
		tm.wg.Wait()
	}

}

func (tm *taskManager) cancelTask(id string) {
	extendedTask, ok := tm.waitingTasks[id]
	if ok {
		tm.removeTask(id)
		tm.logger.Log("Waiting Or Cancelled Task Deleted" + id)
		return
	}
	extendedTask, ok = tm.activeTasks[id]
	if !ok {
		tm.logger.Log("Incorrect Task Id Provided To Cancel Task")
		return
	}

	tm.removeTask(id)
	extendedTask.status = taskStatusEnum.cancelled
	tm.waitingTasks[id] = extendedTask

	tm.logger.Log("Active Task Stopped and Cancelled: " + id)
}

func (tm *taskManager) removeTask(id string) {
	_, wok := tm.waitingTasks[id]
	activeTask, aok := tm.activeTasks[id]
	if wok {
		delete(tm.waitingTasks, id)
		return
	}

	if aok {
		activeTask.timer.Stop()
		activeTask.task.Cancel()
		delete(tm.activeTasks, id)
		return
	}

	tm.logger.Log("Incorect task id provided to remove task")
}

func (tm *taskManager) GracefulShutdown() {
	tm.shutDown()

	tm.wg.Wait()

	tm.rwMutex.Lock()
	defer tm.rwMutex.Unlock()

	for id := range tm.activeTasks {
		tm.cancelTask(id)
	}

	defer func() {
		tm.activeTasks = make(ExtendedTaskMap)
		tm.waitingTasks = make(ExtendedTaskMap)
	}()

	tm.logger.Log("Active and waiting task cancelled and removed for graceful shutdown")
}

func (tm *taskManager) Run() {

	tm.rwMutex.Lock()
	defer tm.rwMutex.Unlock()

	safeConcurrency := int(math.Min(float64(tm.maxConcurrency), float64(len(tm.waitingTasks))))
	count := 0

	for k, v := range tm.waitingTasks {

		select {
		case <-tm.shutdownCtx.Done():
			return
		default:
			if count > safeConcurrency {
				break
			}
			count++

			tm.wg.Add(1)
			delete(tm.waitingTasks, k)
			v.status = taskStatusEnum.running
			timer := time.AfterFunc(v.delay, func() {
				tm.execute(k, &v)
			})
			v.timer = timer
			tm.activeTasks[k] = v

		}
	}

}

func (tm *taskManager) execute(id string, et *ExtendedTask) {
	done := make(chan bool)

	timeout := et.task.GetTimeout()
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)

	go func() {

		tm.rwMutex.Lock()
		defer tm.rwMutex.Unlock()

		defer func() {
			cancel()
			tm.wg.Done()
		}()

		select {
		case <-tm.shutdownCtx.Done():
			et.task.Cancel()
			return
		case <-ctx.Done():
			tm.cancelTask(id)
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
	ctx, cancelFunc := context.WithCancel(context.Background())
	return taskManager{activeTasks: make(ExtendedTaskMap), waitingTasks: make(ExtendedTaskMap), logger: logger, maxConcurrency: maxConcurrency, done: make(chan bool), shutdownCtx: ctx, shutDown: cancelFunc}
}
