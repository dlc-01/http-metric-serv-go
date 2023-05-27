package storage

import (
	"context"
	"fmt"
	"github.com/dlc-01/http-metric-serv-go/internal/general/config"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/general/metrics"
	"github.com/jackc/pgx/v5"
)

type dbStorage struct {
	conn *pgx.Conn
}

var db storage = &dbStorage{}

func (db *dbStorage) Сreate(ctx context.Context, cfg *config.ServerConfig) storage {
	var err error
	db.conn, err = pgx.Connect(ctx, cfg.DatabaseAddress)
	if err != nil {
		logging.Fatalf("cannot connected to db: %s", err)
	}
	logging.Info("connected to db")
	if err := db.PingStorage(ctx); err != nil {
		logging.Fatalf("error while try to ping db: %s", err)
	}
	if err = db.migrationDB(ctx); err != nil {
		logging.Fatalf("cannot create migration: %w", err)
	}
	return db
}

func (db *dbStorage) SetMetric(ctx context.Context, metric metrics.Metric) error {

	if err := db.PingStorage(ctx); err != nil {
		logging.Fatalf("error while try to ping db: %s", err)
	}

	tx, err := db.conn.Begin(ctx)
	defer tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error connected to db before dump: %w", err)
	}
	query := `INSERT INTO metrics (id, mtype, mdelta, mvalue) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO UPDATE SET mvalue = EXCLUDED.mvalue, mdelta =  coalesce(metrics.mdelta, 0) + EXCLUDED.mdelta;`
	if _, err := tx.Exec(ctx, query, metric.ID, metric.MType, metric.Delta, metric.Value); err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("cann't dump metric %q: %w", metric.ID, err)
	}
	return nil
}

func (db *dbStorage) SetMetricsBatch(ctx context.Context, metricsAll []metrics.Metric) error {

	if err := db.PingStorage(ctx); err != nil {
		return fmt.Errorf("error while try to ping db: %w", err)
	}

	batch := &pgx.Batch{}
	query := `INSERT INTO metrics (id, mtype, mdelta, mvalue) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO UPDATE SET mvalue = EXCLUDED.mvalue, mdelta =  coalesce(metrics.mdelta, 0) + EXCLUDED.mdelta;`
	for _, m := range metricsAll {
		batch.Queue(query, m.ID, m.MType, m.Delta, m.Value)
	}

	br := db.conn.SendBatch(ctx, batch)
	defer br.Close()

	if _, err := br.Exec(); err != nil {
		return fmt.Errorf("cannot send butch metric %w", err)
	}

	return nil
}

func (db *dbStorage) GetMetric(ctx context.Context, metric metrics.Metric) (metrics.Metric, error) {
	if err := db.PingStorage(ctx); err != nil {
		return metric, fmt.Errorf("error while try to ping db: %w", err)
	}

	row := db.conn.QueryRow(ctx, "SELECT * FROM metrics WHERE id = $1", metric.ID)
	err := row.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value)
	if err != nil {
		return metric, fmt.Errorf("cannot scan row: %w", err)
	}
	return metric, nil
}

func (db *dbStorage) GetAllMetrics(ctx context.Context) ([]metrics.Metric, error) {

	if err := db.PingStorage(ctx); err != nil {
		return []metrics.Metric{}, fmt.Errorf("error while try to ping db: %w", err)
	}

	rows, err := db.conn.Query(ctx, `SELECT * FROM metrics`)
	if err != nil {
		return nil, fmt.Errorf("error while sending query to db: %w", err)
	}
	defer rows.Close()
	metricsAll := []metrics.Metric{}
	for rows.Next() {
		metric := metrics.Metric{}
		if err := rows.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value); err != nil {
			return nil, err
		}

		metricsAll = append(metricsAll, metric)
	}
	return metricsAll, nil
}

func (db *dbStorage) GetAll(ctx context.Context) ([]string, error) {

	if err := db.PingStorage(ctx); err != nil {
		return []string{}, fmt.Errorf("error while try to ping db: %w", err)
	}

	str := []string{}
	rows, err := db.conn.Query(ctx, "SELECT id FROM metrics")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		str = append(str, id)
	}
	return str, nil

}

func (db *dbStorage) PingStorage(ctx context.Context) error {
	if db.conn != nil {
		err := db.conn.Ping(ctx)
		if err != nil {
			return fmt.Errorf("can't ping to db: %w", err)
		}
	}
	return nil
}

func (db *dbStorage) Close(ctx context.Context) {
	db.conn.Close(ctx)
}

func (db *dbStorage) migrationDB(ctx context.Context) error {
	_, err := db.conn.Exec(ctx, "CREATE TABLE IF NOT EXISTS metrics (id TEXT PRIMARY KEY, mtype TEXT, mdelta BIGINT, mvalue DOUBLE PRECISION)")
	if err != nil {
		return fmt.Errorf("cann't create table: %w", err)
	}
	return nil
}
