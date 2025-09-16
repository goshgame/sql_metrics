package sqlmetrics

import (
	"context"
	"database/sql/driver"
	"strconv"
	"time"
)

type metricsConn struct {
	driver.Conn
}

// NamedValueChecker 实现参数转换
func (mc *metricsConn) CheckNamedValue(nv *driver.NamedValue) error {
	switch v := nv.Value.(type) {
	case uint64:
		if v >= 1<<63 {
			// 超过 int64 范围，转成字符串
			nv.Value = strconv.FormatUint(v, 10)
			return nil
		}
		// 转成 int64，避免触发 default 报错
		nv.Value = int64(v)
		return nil
	case int64:
		// 原样返回，允许负数
		return nil
	}
	return driver.ErrSkip // 交给默认逻辑处理
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

// ExecContext
func (mc *metricsConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	execer, ok := mc.Conn.(driver.ExecerContext)
	if !ok {
		// fallback 到老接口
		return nil, driver.ErrSkip
	}
	start := time.Now()
	res, err := execer.ExecContext(ctx, query, args)
	reportMetrics(query, "exec", start, err)
	return res, err
}

// QueryContext
func (mc *metricsConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	queryer, ok := mc.Conn.(driver.QueryerContext)
	if !ok {
		// fallback
		return nil, driver.ErrSkip
	}
	start := time.Now()
	rows, err := queryer.QueryContext(ctx, query, args)
	reportMetrics(query, "query", start, err)
	return rows, err
}

// PrepareContext
func (mc *metricsConn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	preparer, ok := mc.Conn.(driver.ConnPrepareContext)
	if !ok {
		// fallback
		return mc.Prepare(query)
	}
	start := time.Now()
	stmt, err := preparer.PrepareContext(ctx, query)
	reportMetrics(query, "prepare", start, err)
	if err != nil {
		return nil, err
	}
	return &metricsStmt{Stmt: stmt, query: query}, nil
}

// Ping
func (mc *metricsConn) Ping(ctx context.Context) error {
	pinger, ok := mc.Conn.(driver.Pinger)
	if !ok {
		return driver.ErrSkip
	}
	start := time.Now()
	err := pinger.Ping(ctx)
	reportMetrics("", "ping", start, err)
	return err
}
