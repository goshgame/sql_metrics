package sqlmetrics

import (
	"database/sql/driver"
	"time"
)

type metricsConn struct {
	driver.Conn
}

func (mc *metricsConn) Prepare(query string) (driver.Stmt, error) {
	startTime := time.Now()
	stmt, err := mc.Conn.Prepare(query)
	reportMetrics(query, "prepare", startTime, err)
	if err != nil {
		return nil, err
	}
	return &metricsStmt{Stmt: stmt, query: query}, nil
}

func (mc *metricsConn) Close() error {
	startTime := time.Now()
	err := mc.Conn.Close()
	reportMetrics("", "close", startTime, err)
	return err
}

func (mc *metricsConn) Begin() (driver.Tx, error) {
	startTime := time.Now()
	tx, err := mc.Conn.Begin()
	reportMetrics("", "begin", startTime, err)
	return tx, err
}
