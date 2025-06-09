package pgchronicler

import (
	"context"
	"errors"
	"github.com/coderAM1/godoit/godoit"
	"github.com/jackc/pgx/v5"
)

type Chronicler struct {
	conn   *pgx.Conn
	naming PgNamingOverrides
	logger godoit.LogIt
	ctx    context.Context
}

func NewChronicler(ctx context.Context, conn *pgx.Conn, logger godoit.LogIt) (*Chronicler, error) {
	if conn == nil {
		return nil, errors.New("conn cannot be nil")
	}
	return &Chronicler{
		conn:   conn,
		logger: logger,
		// TODO: add ability to pass in
		naming: PgNamingOverrides{},
		ctx:    ctx,
	}, nil
}

func (pg *Chronicler) SetUpChronicle(ctx context.Context) error {
	// pg.logger.InfoLog(ctx, "starting set up for godoit for postgres")
	createTable := createTaskTableCommand(pg.naming.tableName)
	_, err := pg.conn.Exec(ctx, createTable)
	// TODO figure out if indexing is needed
	return err
}

func (pg *Chronicler) UpsertOverseerInfo(ctx context.Context, info godoit.OverseerInfo) {

}

func (pg *Chronicler) RecordTask(ctx context.Context, task godoit.Task) error {
	// pg.logger.InfoLog(ctx, "starting setting up postgres db")
	recordTask := createInsertTaskCommand(pg.naming.tableName)
	_, err := pg.conn.Exec(ctx, recordTask, task.Id, task.Name, task.Created, task.Scheduled, task.Updated, task.Status, task.Args)
	return err
}

func (pg *Chronicler) QueryTasks(ctx context.Context, limit int) ([]godoit.Task, error) {
	if limit <= 0 {
		return []godoit.Task{}, nil
	}
	return []godoit.Task{}, nil
}

func (pg *Chronicler) UpdateTask(ctx context.Context, task godoit.Task) error {
	return nil
}
