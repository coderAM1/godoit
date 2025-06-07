package godoit

import (
	"context"
	"encoding/json"
	"errors"
	"sync/atomic"
	"time"
)

// structs

type TaskFunc func(ctx context.Context, args json.RawMessage) error

type Manager struct {
	dbTasker             DbTasker
	idMaker              IdMaker
	retryConfig          RetryConfig
	taskMap              map[string]TaskFunc
	ctx                  context.Context
	threadLimit          int
	durationBetweenQuery time.Duration
	taskDuration         time.Duration
}

type ManagerInfo struct {
	name    string `json:"name"`
	podName string `json:"podName"`
}

// functions

func CreateManager(ctx context.Context, tasker DbTasker, config RetryConfig, idMaker IdMaker, threadLimit int) (*Manager, error) {
	if tasker == nil {
		return nil, errors.New("cannot utilize nil tasker")
	}
	if config == nil {
		return nil, errors.New("cannot utilize nil config")
	}
	if threadLimit <= 0 {
		return nil, errors.New("cannot less than or equal to zero threads")
	}
	idGen := idMaker
	if idGen == nil {
		idGen = DefaultIdMaker
	}
	return &Manager{
		dbTasker:             tasker,
		retryConfig:          config,
		taskMap:              make(map[string]TaskFunc),
		idMaker:              idGen,
		ctx:                  ctx,
		threadLimit:          threadLimit,
		durationBetweenQuery: time.Second * 15,
		taskDuration:         time.Second * 30,
	}, nil
}

func (man *Manager) PutTaskInfo(taskName string, taskFunc TaskFunc) error {
	_, ok := man.taskMap[taskName]
	if ok {
		return errors.New("task name already utilized")
	}
	man.taskMap[taskName] = taskFunc
	return nil
}

func (man *Manager) BookTask(ctx context.Context, taskName string, when time.Time, args json.RawMessage) error {
	_, ok := man.taskMap[taskName]
	if !ok {
		return errors.New("task does not exist")
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
	err = man.dbTasker.BookTask(ctx, task)
	return err
}

func (man *Manager) Setup(ctx context.Context) error {
	return man.dbTasker.SetUpDb(ctx)
}

// Start starts infinite for loop to query and run tasks unless error occur or ctx gets cancelled.
// Recommended to use go routine when calling,
func (man *Manager) Start(ctx context.Context) error {
	ctx, cnc := context.WithCancel(ctx)
	defer cnc()
	queryTimer := time.NewTicker(man.durationBetweenQuery)
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
		tasks, err := man.dbTasker.QueryTasks(ctx, man.threadLimit-currentTasksRunning)
		if err != nil {
			return err
		}

		for _, task := range tasks {
			go func(task Task) {
				taskFunc, ok := man.taskMap[task.Name]
				if !ok {
					// TODO decide how to handle manager and db wise
					return
				}
				taskCtx, taskCnc := context.WithTimeout(ctx, man.taskDuration)
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
				man.dbTasker.UpdateTask(ctx, updatedTask)
			}(task)
		}
		tickCount++
	}
}
