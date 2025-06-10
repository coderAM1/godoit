package godoit

import "context"

// interface

type Chronicler interface {
	// SetUpChronicle sets up the database being used for godoit to properly work.
	// I.E. creating tables with the needed columns.
	SetUpChronicle(ctx context.Context) error

	// UpsertOverseerInfo adds or updates the overseer info to the db, if needed.
	// TODO figure out if needed
	// UpsertOverseerInfo(ctx context.Context, info OverseerInfo)

	// RecordTask adds a specified Task to the database to be run at a later time.
	RecordTask(ctx context.Context, task Task) error

	// QueryTasks queries the database for Task(s) that need to be run based on a limit of
	// go routines to run.
	QueryTasks(ctx context.Context, limit int) ([]Task, error)

	// UpdateTask updates the db on whether the tasks
	// TODO decide if this should be a single task to update or if a bulk update function should be used
	UpdateTask(ctx context.Context, task Task) error
}
