package sqlmetrics

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

func init() {
	sql.Register("sqlmetrics", &metricsDriver{parent: &mysql.MySQLDriver{}})
}
