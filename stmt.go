package sqlmetrics

import (
	"database/sql/driver"
	"fmt"
	"math"
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

// NamedValueChecker 实现参数转换
func (ms *metricsStmt) CheckNamedValue(nv *driver.NamedValue) error {
	switch v := nv.Value.(type) {
	case uint64:
		if v > math.MaxInt64 {
			// 超过 int64 范围，转成字符串
			nv.Value = fmt.Sprintf("%d", v)
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
