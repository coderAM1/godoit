package godoit

import "context"

// interface

type DbTasker interface {
	// SetUpDb sets up the database being used for godoit to properly work.
	// I.E. creating tables with the needed columns.
	SetUpDb(ctx context.Context) error

	// UpsertManagerInfo adds or updates the manager info to the db, if needed.
	UpsertManagerInfo(ctx context.Context, info ManagerInfo)

	// BookTask adds a specified Task to the database to be run at a later time.
	BookTask(ctx context.Context, task Task) error

	// QueryTasks queries the database for Task(s) that need to be run based on a limit of
	// go routines to run.
	QueryTasks(ctx context.Context, limit int) ([]Task, error)

	// UpdateTask updates the db on whether the tasks
	UpdateTask(ctx context.Context, task Task) error
}
