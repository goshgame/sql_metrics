package sqlmetrics

import (
	"database/sql/driver"

	"github.com/go-sql-driver/mysql"
)

type metricsDriver struct {
	mysql.MySQLDriver
}

// Open driver open
func (md *metricsDriver) Open(name string) (driver.Conn, error) {
	conn, err := md.MySQLDriver.Open(name)
	if err != nil {
		return nil, err
	}
	return &metricsConn{Conn: conn}, nil
}
