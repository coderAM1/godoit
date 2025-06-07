package godoit_test

import (
	"context"
	"encoding/json"
	"github.com/coderAM1/godoit/godoit"
	"time"
)

// DbTaskerTester is used as a way to help test the manager based off of how the functions are configured
type DbTaskerTester struct {
	SetUpDbFunc           func(ctx context.Context) error
	UpsertManagerInfoFunc func(ctx context.Context, info godoit.ManagerInfo)
	BookTaskFunc          func(ctx context.Context, task godoit.Task) error
	QueryTasksFunc        func(ctx context.Context, limit int) ([]godoit.Task, error)
	UpdateTaskFunc        func(ctx context.Context, task godoit.Task) error
}

func (tester *DbTaskerTester) SetUpDb(ctx context.Context) error {
	return tester.SetUpDbFunc(ctx)
}

func (tester *DbTaskerTester) UpsertManagerInfo(ctx context.Context, info godoit.ManagerInfo) {
	tester.UpsertManagerInfoFunc(ctx, info)
}

func (tester *DbTaskerTester) BookTask(ctx context.Context, task godoit.Task) error {
	return tester.BookTaskFunc(ctx, task)
}

func (tester *DbTaskerTester) QueryTasks(ctx context.Context, limit int) ([]godoit.Task, error) {
	return tester.QueryTasksFunc(ctx, limit)
}

func (tester *DbTaskerTester) UpdateTask(ctx context.Context, task godoit.Task) error {
	return tester.UpdateTaskFunc(ctx, task)
}

func NoErrorTaskFunc() func(ctx context.Context, args json.RawMessage) error {
	return func(ctx context.Context, args json.RawMessage) error {
		return nil
	}
}

func ErrorTaskFunc(err error) func(ctx context.Context, args json.RawMessage) error {
	return func(ctx context.Context, args json.RawMessage) error {
		return err
	}
}

func CreateTask(taskId string, taskName string, args []byte) godoit.Task {
	return godoit.Task{
		Id:      taskId,
		Name:    taskName,
		Created: time.Now(),
		When:    time.Now(),
		Updated: time.Now(),
		Status:  godoit.PENDING,
		Args:    args,
	}
}
