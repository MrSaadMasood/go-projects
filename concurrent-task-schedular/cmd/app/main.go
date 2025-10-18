package main

import (
	"fmt"
	"log/slog"
	"main/internal/env"
	"main/internal/logger"
	"main/internal/manager"
	"main/internal/task"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func randomTimeoutGenerator() *time.Duration {
	randomValue := rand.Intn(10)
	if randomValue > 5 {
		return nil
	}
	duration := time.Second * time.Duration(randomValue)
	return &duration
}

func main() {
	printTaskLogger := logger.NewLogHandler(logger.NewJSONLogger())
	sleepTaskLogger := logger.NewLogHandler(logger.NewJSONLogger())
	taskManagerLogger := logger.NewLogHandler(logger.NewJSONLogger())

	printTaskLogger.AddAttribute([]slog.Attr{slog.String("logger_type", "PRINT_LOGGER")})
	sleepTaskLogger.AddAttribute([]slog.Attr{slog.String("logger_type", "SLEEP_LOGGER")})
	taskManagerLogger.AddAttribute([]slog.Attr{slog.String("logger_type", "TASK_MANAGER_LOGGER")})

	taskManager := manager.CreateNewTaskManager(&taskManagerLogger, env.MaxConcurrency)

	TaskCount := 20

	tasks := make([]manager.EnqueueTask, 0)
	for i := range TaskCount {
		var t manager.Task
		random := rand.Intn(10)

		if random > 5 {
			t = task.NewPrintTask(
				"Print Task:"+strconv.Itoa(i),
				&printTaskLogger,
				rand.Intn(5),
				randomTimeoutGenerator(),
			)
		} else {
			t = task.NewSleepTask(
				"Sleep Task"+strconv.Itoa(i),
				&sleepTaskLogger,
				time.Second*time.Duration(rand.Intn(20)),
				randomTimeoutGenerator(),
			)
		}
		tasks = append(tasks, manager.EnqueueTask{
			Task: t, Delay: time.Second * time.Duration(rand.Intn(10)),
		},
		)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool)
	go func() {
		closeSignal := <-sig
		fmt.Println("the close signal is", closeSignal)
		taskManager.GracefulShutdown()
		done <- true
	}()

	go func() {
		taskManager.Enqueue(tasks...)
		done <- true
	}()

	<-done
}
