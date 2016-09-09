package pqtimeouts

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
)

func init() {
	sql.Register("pqTimeouts", timeoutDriver{})
}

type timeoutDriver struct {
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func (t timeoutDriver) Open(connection string) (driver.Conn, error) {
	// Look for read_timeout and write_timeout in the connection string and extract the values.
	// read_timeout and write_timeout need to be removed from the connection string before calling pq as well.
	var newConnectionSettings []string
	var readTimeout time.Duration
	var writeTimeout time.Duration

	for _, setting := range strings.Fields(connection) {
		s := strings.Split(setting, "=")
		if s[0] == "read_timeout" {
			val, err := strconv.Atoi(s[1])
			if err != nil {
				return nil, fmt.Errorf("Error interpreting value for read_timeout")
			}
			readTimeout = time.Duration(val) * time.Millisecond // timeout is in milliseconds
		} else if s[0] == "write_timeout" {
			val, err := strconv.Atoi(s[1])
			if err != nil {
				return nil, fmt.Errorf("Error interpreting value for write_timeout")
			}
			writeTimeout = time.Duration(val) * time.Millisecond // timeout is in milliseconds
		} else {
			newConnectionSettings = append(newConnectionSettings, setting)
		}
	}

	newConnectionStr := strings.Join(newConnectionSettings, " ")

	return pq.DialOpen(timeoutDriver{readTimeout: readTimeout, writeTimeout: writeTimeout}, newConnectionStr)
}

func (t timeoutDriver) Dial(network string, address string) (net.Conn, error) {
	// If we don't have any timeouts set, just return a normal connection
	if t.readTimeout == 0 && t.writeTimeout == 0 {
		return net.Dial(network, address)
	}

	// Otherwise we want a timeoutConn to handle the read and write deadlines for us.
	c, err := net.Dial(network, address)
	if err != nil || c == nil {
		return c, err
	}

	return &timeoutConn{conn: c, readTimeout: t.readTimeout, writeTimeout: t.writeTimeout}, nil
}

func (t timeoutDriver) DialTimeout(network string, address string, timeout time.Duration) (net.Conn, error) {
	// If we don't have any timeouts set, just return a normal connection
	if t.readTimeout == 0 && t.writeTimeout == 0 {
		return net.DialTimeout(network, address, timeout)
	}

	// Otherwise we want a timeoutConn to handle the read and write deadlines for us.
	c, err := net.DialTimeout(network, address, timeout)
	if err != nil || c == nil {
		return c, err
	}

	return &timeoutConn{conn: c, readTimeout: t.readTimeout, writeTimeout: t.writeTimeout}, nil
}
