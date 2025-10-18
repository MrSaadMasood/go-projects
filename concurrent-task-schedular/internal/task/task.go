// Package task defines task and provides contructors to initialize tasks
package task

import (
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Logger interface {
	Log(message string)
}
type task struct {
	logger  Logger
	id      string
	message string
	timeout *time.Duration
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

type PrintTask struct {
	task
	repeat int
}

func (pt PrintTask) Execute() error {
	for i := range pt.repeat {
		pt.logger.Log("Task with id: " + pt.id + "repeated " + strconv.Itoa(i) + "times")
		pt.execute()
	}
	return nil
}

type SleepTask struct {
	task
	delay time.Duration
}

func (st SleepTask) Execute() error {
	timer := time.NewTimer(st.delay)
	<-timer.C
	st.execute()
	return nil
}

func NewPrintTask(message string, logger Logger, repeat int, timeout *time.Duration) PrintTask {
	return PrintTask{task: task{message: message, logger: logger, id: uuid.NewString(), timeout: timeout}, repeat: repeat}
}

func NewSleepTask(message string, logger Logger, delay time.Duration, timeout *time.Duration) SleepTask {
	return SleepTask{task: task{message: message, logger: logger, id: uuid.NewString(), timeout: timeout}, delay: delay}
}
