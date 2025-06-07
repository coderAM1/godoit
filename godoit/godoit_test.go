package godoit_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/coderAM1/godoit/godoit"
	"github.com/stretchr/testify/assert"
	"sync/atomic"
	"testing"
	"time"
)

func TestManager_Setup(t *testing.T) {
	ctx := t.Context()
	functionCalled := atomic.Bool{}
	dbTest := &DbTaskerTester{
		SetUpDbFunc: func(ctx context.Context) error {
			functionCalled.Store(true)
			return nil
		},
	}

	man, err := godoit.CreateManager(ctx, dbTest, nil, nil, 2)
	assert.Nil(t, err)
	assert.NotNil(t, man)
	assert.False(t, functionCalled.Load())
	err = man.Setup(ctx)
	assert.Nil(t, err)
	assert.True(t, functionCalled.Load())
}

func TestManager_SetupError(t *testing.T) {
	ctx := t.Context()
	functionCalled := atomic.Bool{}
	expErr := errors.New("expected error")
	dbTest := &DbTaskerTester{
		SetUpDbFunc: func(ctx context.Context) error {
			functionCalled.Store(true)
			return expErr
		},
	}

	man, err := godoit.CreateManager(ctx, dbTest, nil, nil, 2)
	assert.Nil(t, err)
	assert.NotNil(t, man)
	assert.False(t, functionCalled.Load())
	err = man.Setup(ctx)
	assert.NotNil(t, err)
	assert.True(t, functionCalled.Load())
	assert.Equal(t, expErr, err)
}

func TestManager_PutTaskInfo(t *testing.T) {
	ctx := t.Context()
	dbTest := &DbTaskerTester{}
	man, err := godoit.CreateManager(ctx, dbTest, nil, nil, 2)
	assert.Nil(t, err)
	assert.NotNil(t, man)
	taskName := "test"
	err = man.PutTaskInfo(taskName, func(ctx context.Context, args json.RawMessage) error {
		return nil
	})
	assert.Nil(t, err)
	testFunc, ok := man.GetTask(taskName)
	assert.True(t, ok)
	assert.NotNil(t, testFunc)
	jsonBytes, jsonErr := json.Marshal(struct{}{})
	assert.Nil(t, jsonErr)
	err = testFunc(ctx, jsonBytes)
	assert.Nil(t, jsonErr)
}

func TestManager_PutTaskInfoError(t *testing.T) {
	ctx := t.Context()
	dbTest := &DbTaskerTester{}
	man, err := godoit.CreateManager(ctx, dbTest, nil, nil, 2)
	assert.Nil(t, err)
	assert.NotNil(t, man)
	taskName := "test"
	taskFunc := func(ctx context.Context, args json.RawMessage) error {
		return nil
	}
	err = man.PutTaskInfo(taskName, taskFunc)
	assert.Nil(t, err)
	err = man.PutTaskInfo(taskName, taskFunc)
	assert.NotNil(t, err)
}

func TestManager_GetTask(t *testing.T) {
	ctx := t.Context()
	dbTest := &DbTaskerTester{}
	man, err := godoit.CreateManager(ctx, dbTest, nil, nil, 2)
	assert.Nil(t, err)
	assert.NotNil(t, man)
	taskNameNoErr := "noErr"
	taskNameErr := "err"
	err = man.PutTaskInfo(taskNameNoErr, NoErrorTaskFunc())
	assert.Nil(t, err)
	errToGet := errors.New("test error")
	err = man.PutTaskInfo(taskNameErr, ErrorTaskFunc(errToGet))
	assert.Nil(t, err)

	noErrFunc, ok := man.GetTask(taskNameNoErr)
	assert.True(t, ok)
	assert.NotNil(t, noErrFunc)
	assert.Nil(t, noErrFunc(ctx, []byte{}))

	errFunc, ok := man.GetTask(taskNameErr)
	assert.True(t, ok)
	assert.NotNil(t, errFunc)
	returnErr := errFunc(ctx, []byte{})
	assert.NotNil(t, returnErr)
	assert.Equal(t, errToGet, returnErr)
}

func TestManager_BookTask(t *testing.T) {
	ctx := t.Context()
	functionCalled := atomic.Bool{}
	assert.False(t, functionCalled.Load())
	dbTest := &DbTaskerTester{
		BookTaskFunc: func(ctx context.Context, task godoit.Task) error {
			functionCalled.Store(true)
			return nil
		},
	}
	man, err := godoit.CreateManager(ctx, dbTest, nil, nil, 2)
	assert.Nil(t, err)
	assert.NotNil(t, man)

	taskName := "testTask"

	err = man.PutTaskInfo(taskName, NoErrorTaskFunc())
	assert.Nil(t, err)
	err = man.BookTask(ctx, taskName, time.Now(), []byte{})
	assert.Nil(t, err)
	assert.True(t, functionCalled.Load())
}

func TestManager_BookTaskNoTask(t *testing.T) {
	ctx := t.Context()
	dbTest := &DbTaskerTester{}
	man, err := godoit.CreateManager(ctx, dbTest, nil, nil, 2)
	assert.Nil(t, err)
	assert.NotNil(t, man)

	err = man.BookTask(ctx, "doesNotExist", time.Now(), []byte{})
	assert.NotNil(t, err)
	msg := err.Error()
	assert.Equal(t, godoit.TaskDoesNotExist, msg)
}

func TestManager_Start(t *testing.T) {
	ctx, cnc := context.WithTimeout(t.Context(), 5*time.Second)

	taskName := "testFunc"

	queuedIdOne := "taskOne"
	taskOneDone := atomic.Bool{}
	taskOne := CreateTask(queuedIdOne, taskName, []byte{})

	queuedIdTwo := "taskTwo"
	taskTwoDone := atomic.Bool{}
	taskTwo := CreateTask(queuedIdTwo, taskName, []byte{})

	testFunc := func(ctx context.Context, args json.RawMessage) error {
		if taskOneDone.Load() {
			taskTwoDone.Store(true)
			defer cnc()
			return errors.New("error happened")
		}
		taskOneDone.Store(true)
		return nil
	}

	num := atomic.Int32{}

	dbTest := &DbTaskerTester{
		QueryTasksFunc: func(ctx context.Context, limit int) ([]godoit.Task, error) {
			if !taskOneDone.Load() {
				return []godoit.Task{taskOne}, nil
			}
			return []godoit.Task{taskTwo}, nil
		},
		UpdateTaskFunc: func(ctx context.Context, task godoit.Task) error {
			num.Add(1)
			return nil
		},
	}

	man, err := godoit.CreateManager(ctx, dbTest, nil, nil, 2)
	assert.Nil(t, err)
	assert.NotNil(t, man)

	err = man.PutTaskInfo(taskName, testFunc)
	assert.Nil(t, err)

	err = man.Start(ctx, time.Millisecond*100, time.Second)
	assert.NotNil(t, err)
	assert.Equal(t, int32(2), num.Load())
	assert.True(t, taskOneDone.Load())
	assert.True(t, taskTwoDone.Load())
}
