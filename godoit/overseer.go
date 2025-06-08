package godoit

import (
	"context"
	"encoding/json"
	"errors"
	"sync/atomic"
	"time"
)

const TaskDoesNotExist = "task does not exist"

// structs

type TaskFunc func(ctx context.Context, args json.RawMessage) error

type Overseer struct {
	chronicler  Chronicler
	idMaker     IdMaker
	retryConfig RetryConfig
	taskMap     map[string]TaskFunc
	ctx         context.Context
	threadLimit int
	started     atomic.Bool
}

type OverseerInfo struct {
	name    string `json:"name"`
	podName string `json:"podName"`
}

// functions

func CreateOverseer(ctx context.Context, tasker Chronicler, retry RetryConfig, idMaker IdMaker, threadLimit int) (*Overseer, error) {
	if tasker == nil {
		return nil, errors.New("cannot utilize nil tasker")
	}
	//if config == nil {
	//	return nil, errors.New("cannot utilize nil config")
	//}
	if threadLimit <= 0 {
		return nil, errors.New("cannot less than or equal to zero threads")
	}
	idGen := idMaker
	if idGen == nil {
		idGen = DefaultIdMaker
	}
	return &Overseer{
		chronicler:  tasker,
		retryConfig: retry,
		taskMap:     make(map[string]TaskFunc),
		idMaker:     idGen,
		ctx:         ctx,
		threadLimit: threadLimit,
		started:     atomic.Bool{},
	}, nil
}

func (man *Overseer) PutTaskInfo(taskName string, taskFunc TaskFunc) error {
	_, ok := man.taskMap[taskName]
	if ok {
		return errors.New("task name already utilized")
	}
	man.taskMap[taskName] = taskFunc
	return nil
}

func (man *Overseer) GetTask(taskName string) (TaskFunc, bool) {
	task, ok := man.taskMap[taskName]
	return task, ok
}

func (man *Overseer) BookTask(ctx context.Context, taskName string, when time.Time, args json.RawMessage) error {
	_, ok := man.taskMap[taskName]
	if !ok {
		return errors.New(TaskDoesNotExist)
	}
	id, err := man.idMaker(ctx, taskName, when)
	if err != nil {
		return err
	}
	now := time.Now()
	task := Task{
		Id:      id,
		Name:    taskName,
		Created: now,
		When:    when,
		Updated: now,
		Status:  PENDING,
		Args:    args,
		Retry:   false,
	}
	err = man.chronicler.RecordTask(ctx, task)
	return err
}

func (man *Overseer) Setup(ctx context.Context) error {
	return man.chronicler.SetUpChronicle(ctx)
}

// Start starts infinite for loop to query and run tasks unless error occur or ctx gets cancelled.
// Recommended to use go routine when calling,
func (man *Overseer) Start(ctx context.Context, durationBetweenQuery time.Duration, taskDuration time.Duration) error {
	if man.started.Swap(true) {
		return errors.New("already started")
	}
	ctx, cnc := context.WithCancel(ctx)
	defer cnc()
	queryTimer := time.NewTicker(durationBetweenQuery)
	defer queryTimer.Stop()
	tickCount := 0
	taskGoing := atomic.Int64{}
	// probably not needed but just being safe
	taskGoing.Store(0)
	for {
		if tickCount > 4 {
			tickCount = 0
			// TODO add logic for cleaning up bad state tasks and removing completed jobs after an amount of time
		}
		select {
		case <-queryTimer.C:
		case <-ctx.Done():
			return errors.New("context cancelled")
		case <-man.ctx.Done():
			return errors.New("context cancelled")
		}
		currentTasksRunning := int(taskGoing.Load())
		tasks, err := man.chronicler.QueryTasks(ctx, man.threadLimit-currentTasksRunning)
		if err != nil {
			return err
		}

		for _, task := range tasks {
			go func(task Task) {
				taskFunc, ok := man.taskMap[task.Name]
				if !ok {
					// TODO decide how to handle overseer and db wise
					return
				}
				taskCtx, taskCnc := context.WithTimeout(ctx, taskDuration)
				defer taskCnc()
				taskGoing.Add(1)
				// TODO figure out better way to handle contexts getting cancelled with task
				taskErr := taskFunc(taskCtx, task.Args)
				taskGoing.Add(-1)
				var updatedTask Task
				if taskErr != nil {
					updatedTask = task.CreateUpdatedTask(FAILED, time.Now())
					// TODO figure out what else to do with this error
				} else {
					updatedTask = task.CreateUpdatedTask(DONE, time.Now())
				}
				// TODO handle error also probably best to just batch these with a channel and update on ticks
				man.chronicler.UpdateTask(ctx, updatedTask)
			}(task)
		}
		tickCount++
	}
}
