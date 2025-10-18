// Package task defines task and provides contructors to initialize tasks
package task

import (
	"context"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Logger interface {
	Log(message string)
}
type task struct {
	logger     Logger
	id         string
	message    string
	timeout    *time.Duration
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func (t task) execute() {
	t.logger.Log(t.message)
}

func (t task) GetTimeout() *time.Duration {
	timeout := t.timeout
	if timeout == nil {
		fallback := time.Second * 300
		return &fallback
	} else {
		value := time.Second * *timeout
		return &value
	}
}

func (t task) Cancel() {
	t.cancelFunc()
}

type PrintTask struct {
	task
	repeat int
}

func (pt PrintTask) Execute() error {
	for i := range pt.repeat {
		select {
		case <-pt.ctx.Done():
			pt.logger.Log("Cancelled In Progress Print Task: ----->" + pt.id)
			return nil
		default:
			pt.logger.Log("Task with id: " + pt.id + "repeated " + strconv.Itoa(i) + "times")
			pt.execute()
		}
	}
	return nil
}

type SleepTask struct {
	task
	delay time.Duration
}

func (st SleepTask) Execute() error {
	timer := time.NewTimer(st.delay)
	select {
	case <-timer.C:
		st.execute()
		return nil
	case <-st.ctx.Done():
		st.logger.Log("Cancelled In Progress Sleep Task: ----->" + st.id)
		return nil
	}
}

func NewPrintTask(message string, logger Logger, repeat int, timeout *time.Duration) PrintTask {
	ctx, cancelFunc := context.WithCancel(context.Background())
	return PrintTask{task: task{message: message, logger: logger, id: uuid.NewString(), timeout: timeout, ctx: ctx, cancelFunc: cancelFunc}, repeat: repeat}
}

func NewSleepTask(message string, logger Logger, delay time.Duration, timeout *time.Duration) SleepTask {
	ctx, cancelFunc := context.WithCancel(context.Background())
	return SleepTask{task: task{message: message, logger: logger, id: uuid.NewString(), timeout: timeout, ctx: ctx, cancelFunc: cancelFunc}, delay: delay}
}
