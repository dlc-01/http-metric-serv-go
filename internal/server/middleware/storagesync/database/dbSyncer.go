package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dlc-01/http-metric-serv-go/internal/general/config"
	"github.com/dlc-01/http-metric-serv-go/internal/general/logging"
	"github.com/dlc-01/http-metric-serv-go/internal/server/storage"
	"github.com/jackc/pgx/v5"
	"time"
)

var conn *pgx.Conn
var conf *config.ServerConfig
var ctx context.Context
var initialDB bool

func RunDumperDB() {
	dumpTicker := time.NewTicker(time.Duration(conf.StoreInterval) * time.Second)
	defer dumpTicker.Stop()
	for range dumpTicker.C {
		if err := DumpDB(); err != nil {
			logging.Fatalf("cannot dump db metrics to file: %s", err)
		}
	}
}

func RestoreDB() error {
	rows, err := conn.Query(ctx, `SELECT * FROM metrics`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id      string
			mtype   string
			deltaDB sql.NullInt64
			valueDB sql.NullFloat64
		)
		if err := rows.Scan(&id, &mtype, &deltaDB, &valueDB); err != nil {
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
		storage.SetMetric(id, mtype, value, delta)
	}
	return nil
}

func DumpDB() error {

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	for _, metric := range storage.GetMetrics() {
		query := `INSERT INTO metrics (id, mtype, mdelta, mvalue) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO UPDATE SET mvalue = EXCLUDED.mvalue, mdelta = EXCLUDED.mdelta;`
		if _, err := tx.Exec(ctx, query, metric.ID, metric.MType, metric.Delta, metric.Value); err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("cann't dump metric %q: %w", metric.ID, err)
		}
	}

	return tx.Commit(ctx)
}

func InitDB(cont context.Context, cnf *config.ServerConfig) error {
	ctx = cont
	conf = cnf
	var err error
	conn, err = pgx.Connect(ctx, conf.DatabaseAddress)
	if err != nil {
		return fmt.Errorf("cann't connect to db: %w", err)
	}
	logging.Info("connected to db")
	_, err = conn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS metrics (id TEXT PRIMARY KEY, mtype TEXT, mdelta BIGINT, mvalue DOUBLE PRECISION)")
	if err != nil {
		return fmt.Errorf("cann't create table: %w", err)
	}
	return nil
}
func ConnectDB() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
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
