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

// structs

type Task struct {
	Id             string          `json:"id"`
	Name           string          `json:"taskName"`
	Created        time.Time       `json:"created"`
	When           time.Time       `json:"when"`
	Updated        time.Time       `json:"updated"`
	Status         status          `json:"status"`
	Args           json.RawMessage `json:"args"`
	Retry          bool            `json:"retry,omitempty"`
	RetryAmounts   int             `json:"retryAmounts,omitempty"`
	AttemptedTimes []time.Time     `json:"attemptedTimes,omitempty"`
	Recurring      bool            `json:"recurringTask,omitempty"`
}

type RecurringTaskInfo struct {
	Name string `json:"name"`
	Cron string `json:"cron"`
}

func (t Task) CreateUpdatedTask(stat status, updated time.Time) Task {
	return Task{
		Id:      t.Id,
		Name:    t.Name,
		Created: t.Created,
		When:    t.When,
		Updated: updated,
		Status:  stat,
		Args:    t.Args,
	}
}
