package sqlmetrics

import (
	"database/sql"
)

func init() {
	sql.Register("sqlmetrics", &metricsDriver{})
}
