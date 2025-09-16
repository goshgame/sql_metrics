package sqlmetrics

import (
	"database/sql/driver"

	"github.com/go-sql-driver/mysql"
)

type metricsDriver struct {
	parent *mysql.MySQLDriver
}

// Open driver open
func (md *metricsDriver) Open(dsn string) (driver.Conn, error) {
	conn, err := md.parent.Open(dsn)
	if err != nil {
		return nil, err
	}
	return &metricsConn{Conn: conn}, nil
}
