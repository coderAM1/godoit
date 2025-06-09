package godoit_test

import (
	"context"
	"encoding/json"
	"github.com/coderAM1/godoit/godoit"
	"time"
)

// ChroniclerMock is used as a way to help test the overseer based off of how the functions are configured
type ChroniclerMock struct {
	SetUpChronicleFunc     func(ctx context.Context) error
	UpsertOverseerInfoFunc func(ctx context.Context, info godoit.OverseerInfo)
	RecordTaskFunc         func(ctx context.Context, task godoit.Task) error
	QueryTasksFunc         func(ctx context.Context, limit int) ([]godoit.Task, error)
	UpdateTaskFunc         func(ctx context.Context, task godoit.Task) error
}

func (chron *ChroniclerMock) SetUpChronicle(ctx context.Context) error {
	return chron.SetUpChronicleFunc(ctx)
}

func (chron *ChroniclerMock) UpsertOverseerInfo(ctx context.Context, info godoit.OverseerInfo) {
	chron.UpsertOverseerInfoFunc(ctx, info)
}

func (chron *ChroniclerMock) RecordTask(ctx context.Context, task godoit.Task) error {
	return chron.RecordTaskFunc(ctx, task)
}

func (chron *ChroniclerMock) QueryTasks(ctx context.Context, limit int) ([]godoit.Task, error) {
	return chron.QueryTasksFunc(ctx, limit)
}

func (chron *ChroniclerMock) UpdateTask(ctx context.Context, task godoit.Task) error {
	return chron.UpdateTaskFunc(ctx, task)
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
		Id:        taskId,
		Name:      taskName,
		Created:   time.Now().UTC(),
		Scheduled: time.Now().UTC(),
		Updated:   time.Now().UTC(),
		Status:    godoit.PENDING,
		Args:      args,
	}
}
