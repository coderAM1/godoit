package pgchronicler

import (
	"fmt"
	"strings"
	"time"
)

const DEFAULT_TASK_TABLE_NAME = "tasks"

// CREATE_TASK_TABLE TODO: add recurring/retry logic
const CREATE_TASK_TABLE = "CREATE TABLE IF NOT EXISTS %s (" +
	"id TEXT PRIMARY KEY," +
	"taskname TEXT NOT NULL," +
	"created TIMESTAMP WITHOUT TIME ZONE NOT NULL," +
	"scheduled TIMESTAMP WITHOUT TIME ZONE NOT NULL," +
	"updated TIMESTAMP WITHOUT TIME ZONE NOT NULL," +
	"status TEXT NOT NULL," +
	"args JSONB NOT NULL" +
	")"

const SELECT_TASKS_TO_RUN = "SELECT * FROM %s where (scheduled <= '%s' AND status = 'PENDING')"

const INSERT_TASK = "INSERT INTO %s(id, taskname, created, scheduled, updated, status, args) VALUES($1, $2, $3, $4, $5, $6, $7)"

type PgNamingOverrides struct {
	tableName string
}

func createTaskTableCommand(name string) string {
	nameToUse := name
	if nameToUse == "" {
		nameToUse = DEFAULT_TASK_TABLE_NAME
	}
	return fmt.Sprintf(CREATE_TASK_TABLE, DEFAULT_TASK_TABLE_NAME)
}

func createInsertTaskCommand(name string) string {
	nameToUse := name
	if nameToUse == "" {
		nameToUse = DEFAULT_TASK_TABLE_NAME
	}
	return fmt.Sprintf(INSERT_TASK, DEFAULT_TASK_TABLE_NAME)
}

func getProperTimeCheckString() string {
	return strings.Split(time.Now().UTC().String(), ".")[0]
}
