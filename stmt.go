package sqlmetrics

import (
	"database/sql/driver"
	"time"
)

type metricsStmt struct {
	driver.Stmt
	query string
}

// Close close stmt
func (ms *metricsStmt) Close() error {
	startTime := time.Now()
	err := ms.Stmt.Close()
	reportMetrics("", "stmt_close", startTime, err)
	return err
}

// NumInput num input
func (ms *metricsStmt) NumInput() int {
	return ms.Stmt.NumInput()
}

// Exec exec
func (ms *metricsStmt) Exec(args []driver.Value) (driver.Result, error) {
	startTime := time.Now()
	res, err := ms.Stmt.Exec(args)
	reportMetrics(ms.query, "exec", startTime, err)
	return res, err
}

// Query query
func (ms *metricsStmt) Query(args []driver.Value) (driver.Rows, error) {
	startTime := time.Now()
	rows, err := ms.Stmt.Query(args)
	reportMetrics(ms.query, "query", startTime, err)
	return rows, err
}
