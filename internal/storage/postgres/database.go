package postgres

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Database struct {
	cluster *pgxpool.Pool
}

func newDatabase(cluster *pgxpool.Pool) *Database {
	return &Database{cluster: cluster}
}

func (db Database) GetPool() *pgxpool.Pool {
	return db.cluster
}

func (db Database) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Get(ctx, db.cluster, dest, query, args...)
}

func (db Database) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Select(ctx, db.cluster, dest, query, args...)
}

func (db Database) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	return db.cluster.Exec(ctx, query, args...)
}

func (db Database) ExecQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return db.cluster.QueryRow(ctx, query, args...)
}

func (db Database) SendBatch(ctx context.Context, batch *pgx.Batch) pgx.BatchResults {
	return db.cluster.SendBatch(ctx, batch)
}

func (db Database) BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error) {
	return db.cluster.BeginTx(ctx, opts)
}

func (db Database) Close() {
	db.cluster.Close()
}
