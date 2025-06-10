package godoit

import (
	"encoding/json"
	"time"
)

// status

type status string

const PENDING status = "PENDING"
const GOING status = "GOING"
const DONE status = "DONE"
const FAILED status = "FAILED"
const UNKNOWN status = "UNKNOWN"

// structs

type Task struct {
	Id        string          `json:"id" db:"id"`
	Name      string          `json:"taskName" db:"taskname"`
	Created   time.Time       `json:"created" db:"created"`
	Scheduled time.Time       `json:"scheduled" db:"scheduled"`
	Updated   time.Time       `json:"updated" db:"updated"`
	Status    status          `json:"status" db:"status"`
	Args      json.RawMessage `json:"args" db:"args"`
	// Retry          bool            `json:"retry,omitempty"`
	// RetryAmounts   int             `json:"retryAmounts,omitempty"`
	// AttemptedTimes []time.Time     `json:"attemptedTimes,omitempty"`
	// Recurring      bool            `json:"recurringTask,omitempty"`
}

type RecurringTaskInfo struct {
	Name string `json:"name"`
	Cron string `json:"cron"`
}

func (t Task) CreateUpdatedTask(stat status, updated time.Time) Task {
	return Task{
		Id:        t.Id,
		Name:      t.Name,
		Created:   t.Created,
		Scheduled: t.Scheduled,
		Updated:   updated,
		Status:    stat,
		Args:      t.Args,
	}
}
