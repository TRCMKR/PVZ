package tx_manager

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// Database ...
type database interface {
	Select(context.Context, interface{}, string, ...interface{}) error
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	ExecQueryRow(context.Context, string, ...interface{}) pgx.Row
	SendBatch(context.Context, *pgx.Batch) pgx.BatchResults
	BeginTx(context.Context, pgx.TxOptions) (pgx.Tx, error)
}

// TxManager ...
type TxManager struct {
	db database
}

// NewTxManager ...
func NewTxManager(db database) *TxManager {
	return &TxManager{
		db: db,
	}
}

// RunSerializable ...
func (m *TxManager) RunSerializable(ctx context.Context, fn func(context.Context, pgx.Tx) error) error {
	opts := pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadWrite,
	}

	return m.beginFunc(ctx, opts, fn)
}

// RunRepeatableRead ...
func (m *TxManager) RunRepeatableRead(ctx context.Context, fn func(context.Context, pgx.Tx) error) error {
	opts := pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadWrite,
	}

	return m.beginFunc(ctx, opts, fn)
}

// RunReadCommitted ...
func (m *TxManager) RunReadCommitted(ctx context.Context, fn func(context.Context, pgx.Tx) error) error {
	opts := pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	}

	return m.beginFunc(ctx, opts, fn)
}

func (m *TxManager) beginFunc(ctx context.Context, opts pgx.TxOptions, fn func(context.Context, pgx.Tx) error) error {
	tx, err := m.db.BeginTx(ctx, opts)

	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if err = fn(ctx, tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
