package pgtasker

import (
	"context"
	"errors"
	"github.com/coderAM1/godoit/godoit"
	"github.com/jackc/pgx/v5"
)

type PgIt struct {
	conn *pgx.Conn
	ctx  context.Context
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

func (pg *PgIt) SetUpDb(ctx context.Context) error {
	return nil
}

func (pg *PgIt) UpsertManagerInfo(ctx context.Context, info godoit.ManagerInfo) {

}

func (pg *PgIt) BookTask(ctx context.Context, task godoit.Task) error {
	return nil
}

func (pg *PgIt) QueryTasks(ctx context.Context, limit int) ([]godoit.Task, error) {
	return []godoit.Task{}, nil
}

func (pg *PgIt) UpdateTask(ctx context.Context, task godoit.Task) error {
	return nil
}
