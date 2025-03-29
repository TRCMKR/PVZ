package tx_manager

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type Database interface {
	Get(context.Context, interface{}, string, ...interface{}) error
	Select(context.Context, interface{}, string, ...interface{}) error
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	ExecQueryRow(context.Context, string, ...interface{}) pgx.Row
	SendBatch(context.Context, *pgx.Batch) pgx.BatchResults
	BeginTx(context.Context, pgx.TxOptions) (pgx.Tx, error)
}

type txManagerKey struct{}

type TxManager struct {
	db Database
}

func NewTxManager(db Database) *TxManager {
	return &TxManager{
		db: db,
	}
}

func (m *TxManager) RunSerializable(ctx context.Context, fn func(ctxTx context.Context) error) error {
	opts := pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadOnly,
	}

	return m.beginFunc(ctx, opts, fn)
}

func (m *TxManager) RunRepeatableRead(ctx context.Context, fn func(ctxTx context.Context) error) error {
	opts := pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadOnly,
	}

	return m.beginFunc(ctx, opts, fn)
}

func (m *TxManager) RunReadCommitted(ctx context.Context, fn func(ctxTx context.Context) error) error {
	opts := pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadOnly,
	}

	return m.beginFunc(ctx, opts, fn)
}

func (m *TxManager) beginFunc(ctx context.Context, opts pgx.TxOptions, fn func(ctxTx context.Context) error) error {
	tx, err := m.db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	ctx = context.WithValue(ctx, txManagerKey{}, tx)
	if err = fn(ctx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (m *TxManager) GetQueryEngine(ctx context.Context) Database {
	v, ok := ctx.Value(txManagerKey{}).(Database)
	if ok && v != nil {
		return v
	}

	return m.db
}
