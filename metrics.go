package sqlmetrics

import (
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// SQLMetricsHandleTotal - SQL metrics handle total
	SQLMetricsHandleTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "sql_handle_total",
		Help: "Total number of sql handle make.",
	}, []string{"table", "method", "bcode"})
	// SQLMetricsHandleDuration - SQL metrics handle latencies in seconds
	SQLMetricsHandleDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "sql_handle_duration",
		Help: "The sql handle latencies in seconds",
	}, []string{"table", "method", "bcode"})
	// SQLDBMaxOpenConns - SQL metrics max open conns
	SQLDBMaxOpenConns = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "sql_db_max_open_conns",
		Help: "Maximum number of open connections to the database.",
	}, []string{"db", "master_slave"})
	// SQLDBOpenConns - SQL metrics open conns(use and idle)
	SQLDBOpenConns = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "sql_db_open_conns",
		Help: "The number of established connections both in use and idle.",
	}, []string{"db", "master_slave"})
	// SQLDBInUseConns - SQL metrics inuse conns
	SQLDBInUseConns = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "sql_db_in_use_conns",
		Help: "The number of connections currently in use.",
	}, []string{"db", "master_slave"})
	// SQLDBIdleConns - SQL metrics idel conns
	SQLDBIdleConns = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "sql_db_idle_conns",
		Help: "The number of idle connections.",
	}, []string{"db", "master_slave"})
)

func getMetricCode(err error) string {
	if err != nil {
		return "1"
	}
	return "0"
}

// 提取 SQL 动作 (select/insert/update/delete/other)
func extractOp(query string) string {
	if query == "" {
		return "other"
	}
	q := strings.TrimSpace(strings.ToLower(query))
	if q == "" {
		return "other"
	}
	// 拿第一个单词
	fields := strings.Fields(q)
	if len(fields) == 0 {
		return "other"
	}
	switch fields[0] {
	case "select":
		return "select"
	case "update":
		return "update"
	case "delete":
		return "delete"
	case "insert":
		return "insert"
	default:
		return "other"
	}
}

// 从 SQL 里提取 table 名（简单版：取 FROM/INTO/UPDATE 后第一个单词）
func extractTable(query string) string {
	q := strings.ToLower(query)
	if strings.Contains(q, "from ") {
		parts := strings.Split(q, "from ")
		if len(parts) > 1 {
			words := strings.Fields(parts[1])
			if len(words) > 0 {
				return words[0]
			}
		}
	}
	if strings.Contains(q, "into ") {
		parts := strings.Split(q, "into ")
		if len(parts) > 1 {
			words := strings.Fields(parts[1])
			if len(words) > 0 {
				return words[0]
			}
		}
	}
	if strings.HasPrefix(q, "update ") {
		words := strings.Fields(q)
		if len(words) > 1 {
			return words[1]
		}
	}
	return "unknown"
}

func reportMetrics(query string, fallbackMethod string, startTime time.Time, err error) {
	op := extractOp(query)
	method := op
	if method == "other" {
		method = fallbackMethod
	}
	bcode := getMetricCode(err)
	table := extractTable(query)
	duration := float64(time.Since(startTime).Milliseconds()) / 1000
	SQLMetricsHandleTotal.With(prometheus.Labels{"table": table, "method": method, "bcode": bcode}).Inc()
	SQLMetricsHandleDuration.With(prometheus.Labels{"table": table, "method": method, "bcode": bcode}).Observe(duration)
}
