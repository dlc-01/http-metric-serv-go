package storagesync

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"github.com/jackc/pgx/v5"
	"time"
)

type dbConfig struct {
	ctx  context.Context
	conn *pgx.Conn
}

var db dbConfig

func ConnectDB() bool {
	ctx, cancel := context.WithTimeout(db.ctx, 1*time.Second)
	defer cancel()
	conn, err := pgx.Connect(ctx, conf.DatabaseAddress)
	if err != nil {
		logging.Errorf("can't connect to db: %s", err)
		return false
	}
	defer conn.Close(context.Background())
	logging.Info("connected to db")
	return true
}

func initDB() error {
	var err error
	db.conn, err = pgx.Connect(db.ctx, conf.DatabaseAddress)
	if err != nil {
		return fmt.Errorf("cann't connect to db: %w", err)
	}
	logging.Info("connected to db")
	_, err = db.conn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS metrics (id TEXT PRIMARY KEY, mtype TEXT, mdelta BIGINT, mvalue DOUBLE PRECISION)")
	if err != nil {
		return fmt.Errorf("cann't create table: %w", err)
	}
	return nil
}

func runDumperDB() {
	dumpTicker := time.NewTicker(time.Duration(conf.StoreInterval) * time.Second)
	defer dumpTicker.Stop()
	for range dumpTicker.C {
		if err := dumpDB(); err != nil {
			logging.Fatalf("cannot dump db metrics to file: %s", err)
		}
	}
}

func restoreDB() error {
	rows, err := db.conn.Query(db.ctx, `SELECT * FROM metrics`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			idDB    string
			mtypeDB string
			deltaDB sql.NullInt64
			valueDB sql.NullFloat64
		)
		if err := rows.Scan(&idDB, &mtypeDB, &deltaDB, &valueDB); err != nil {
			return err
		}
		var delta *int64
		if deltaDB.Valid {
			delta = &deltaDB.Int64
		}
		var value *float64
		if valueDB.Valid {
			value = &valueDB.Float64
		}
		storage.SetMetric(idDB, mtypeDB, value, delta)
	}
	return nil
}

func dumpDB() error {
	//TODO What kind insertion better with transaction or bulk insert
	tx, err := db.conn.Begin(db.ctx)
	if err != nil {
		return err
	}
	for _, metric := range storage.GetMetrics() {
		query := `INSERT INTO metrics (id, mtype, mdelta, mvalue) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO UPDATE SET mvalue = EXCLUDED.mvalue, mdelta =  coalesce(metrics.mdelta, 0) + EXCLUDED.mdelta;`
		if _, err := tx.Exec(db.ctx, query, metric.ID, metric.MType, metric.Delta, metric.Value); err != nil {
			tx.Rollback(db.ctx)
			return fmt.Errorf("cann't dump metric %q: %w", metric.ID, err)
		}
	}

	return tx.Commit(db.ctx)
}
