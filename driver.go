package sqlmetrics

import (
	"database/sql/driver"
)

type metricsDriver struct {
	parent driver.Driver
}

// Open driver open
func (md *metricsDriver) Open(name string) (driver.Conn, error) {
	conn, err := md.parent.Open(name)
	if err != nil {
		return nil, err
	}
	return &metricsConn{Conn: conn}, nil
}
