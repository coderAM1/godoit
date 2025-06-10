package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coderAM1/godoit/godoit"
	"github.com/coderAM1/godoit/pgchronicler"
	"github.com/jackc/pgx/v5"
	"strings"
	"time"
)

func main() {
	pgUrl := "postgres://myuser:mypassword@localhost:5432/mydb"
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, pgUrl)
	if err != nil {
		panic(err)
	}
	pgChronicler, err2 := pgchronicler.NewChronicler(ctx, conn, nil)
	if err2 != nil {
		panic(err2)
	}
	err = pgChronicler.SetUpChronicle(ctx)
	if err != nil {
		panic(err)
	}
	//task := createTestTask()
	//err = pgChronicler.RecordTask(ctx, task)
	//if err != nil {
	//	panic(err)
	//}
	selectString := fmt.Sprintf(pgchronicler.SELECT_TASKS_TO_RUN, pgchronicler.DEFAULT_TASK_TABLE_NAME, strings.Split(time.Now().UTC().Add(10*time.Minute).String(), ".")[0])
	fmt.Println(selectString)
	ovsr, ovsrErr := godoit.CreateOverseer(ctx, pgChronicler, nil, godoit.DefaultIdMaker, 5)
	if ovsrErr != nil {
		panic(ovsrErr)
	}
	ovsr.PutTaskInfo("test", func(ctx context.Context, args json.RawMessage) error {
		fmt.Println(string(args))
		return nil
	})
	// ovsr.Setup(ctx)
	err = ovsr.Start(ctx, 2*time.Second, 30*time.Second)
	panic(err)
}

func createTestTask() godoit.Task {
	args := TestArgs{
		TestName:   "test",
		NumberTest: 123,
	}
	jsonBytes, _ := json.Marshal(&args)
	now := time.Now().UTC()
	scheduled := now.Add(10 * time.Minute)
	return godoit.Task{
		Id:        "test-task",
		Name:      "test",
		Created:   now,
		Scheduled: scheduled,
		Updated:   now,
		Status:    godoit.PENDING,
		Args:      jsonBytes,
	}
}

type TestArgs struct {
	TestName   string `json:"testName"`
	NumberTest int    `json:"numberTest"`
}
