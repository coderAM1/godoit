package main

import (
	"encoding/json"
	"fmt"
	"github.com/coderAM1/godoit/godoit"
	"strings"
	"time"
)

func main() {
	fmt.Println(time.Now().UTC().String())
	array := strings.Split(time.Now().UTC().Add(10*time.Minute).String(), " +")
	fmt.Println(array[0])
	fmt.Println(array[1])
	//pgUrl := "postgres://myuser:mypassword@localhost:5432/mydb"
	//ctx := context.Background()
	//conn, err := pgx.Connect(ctx, pgUrl)
	//if err != nil {
	//	panic(err)
	//}
	//pgChronicler, err2 := pgchronicler.NewChronicler(ctx, conn, nil)
	//if err2 != nil {
	//	panic(err2)
	//}
	//err = pgChronicler.SetUpChronicle(ctx)
	//if err != nil {
	//	panic(err)
	//}
	////task := createTestTask()
	////err = pgChronicler.RecordTask(ctx, task)
	////if err != nil {
	////	panic(err)
	////}
	//selectString := fmt.Sprintf(pgchronicler.SELECT_TASKS_TO_RUN, pgchronicler.DEFAULT_TASK_TABLE_NAME, strings.Split(time.Now().UTC().Add(10*time.Minute).String(), ".")[0])
	//fmt.Println(selectString)
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
