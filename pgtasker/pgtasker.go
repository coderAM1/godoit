package pgtasker

import (
	"context"
	"errors"
	"github.com/coderAM1/godoit/godoit"
	"github.com/jackc/pgx/v5"
)

type PgIt struct {
	conn   *pgx.Conn
	logger godoit.LogIt
	ctx    context.Context
}

func NewPgIt(ctx context.Context, conn *pgx.Conn) (*PgIt, error) {
	if conn == nil {
		return nil, errors.New("conn cannot be nil")
	}
	return &PgIt{
		conn: conn,
		ctx:  ctx,
	}, nil
}

func (pg *PgIt) SetUpChronicle(ctx context.Context) error {
	pg.logger.InfoLog(ctx, "starting setting up postgres db")
	return nil
}

func (pg *PgIt) UpsertOverseerInfo(ctx context.Context, info godoit.OverseerInfo) {

}

func (pg *PgIt) RecordTask(ctx context.Context, task godoit.Task) error {
	pg.logger.InfoLog(ctx, "starting setting up postgres db")
	return nil
}

func (pg *PgIt) QueryTasks(ctx context.Context, limit int) ([]godoit.Task, error) {
	if limit <= 0 {
		return []godoit.Task{}, nil
	}
	return []godoit.Task{}, nil
}

func (pg *PgIt) UpdateTask(ctx context.Context, task godoit.Task) error {
	return nil
}
