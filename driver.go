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
	sql.Register("pqTimeouts", &timeoutDriver{})
}

type timeoutDriver struct {
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func (t *timeoutDriver) Open(connection string) (driver.Conn, error) {
	// Look for read_timeout and write_timeout in the connection string and extract the values.
	var newConnectionSettings []string

	for _, setting := range strings.Fields(connection) {
		s := strings.Split(setting, "=")
		if s[0] == "read_timeout" {
			val, err := strconv.Atoi(s[1])
			if err != nil {
				return nil, fmt.Errorf("Error interpreting value for read_timeout")
			}
			t.readTimeout = time.Duration(val) * time.Millisecond // timeout is in milliseconds
		} else if s[0] == "write_timeout" {
			val, err := strconv.Atoi(s[1])
			if err != nil {
				return nil, fmt.Errorf("Error interpreting value for write_timeout")
			}
			t.writeTimeout = time.Duration(val) * time.Millisecond // timeout is in milliseconds
		} else {
			newConnectionSettings = append(newConnectionSettings, setting)
		}
	}

	newConnectionStr := strings.Join(newConnectionSettings, " ")

	return pq.DialOpen(&timeoutDriver{}, newConnectionStr)
}

func (t *timeoutDriver) Dial(network string, address string) (net.Conn, error) {
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

func (t *timeoutDriver) DialTimeout(network string, address string, timeout time.Duration) (net.Conn, error) {
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

type timeoutConn struct {
	conn         net.Conn
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func (t *timeoutConn) Read(b []byte) (n int, err error) {
	if t.conn != nil {
		if t.readTimeout != 0 {
			// Set a read deadline before we call read.
			t.conn.SetReadDeadline(time.Now().Add(t.readTimeout))
		}
		n, err = t.conn.Read(b)
		if t.readTimeout != 0 {
			// Clear the deadline if we have one set
			t.conn.SetReadDeadline(time.Time{})
		}
		return
	}
	return 0, fmt.Errorf("Connection is nil")
}

func (t *timeoutConn) Write(b []byte) (n int, err error) {
	if t.conn != nil {
		if t.writeTimeout != 0 {
			// Set a write deadline before we call write.
			t.conn.SetWriteDeadline(time.Now().Add(t.writeTimeout))
		}
		n, err = t.conn.Write(b)
		if t.writeTimeout != 0 {
			// Clear the deadline if we have one set
			t.conn.SetWriteDeadline(time.Time{})
		}
	}
	return 0, fmt.Errorf("Connection is nil")
}

func (t *timeoutConn) Close() (err error) {
	if t.conn != nil {
		err = t.conn.Close()
		if err != nil {
			// If the close looked successful, set the connection to nil
			t.conn = nil
		}
		return
	}
	return fmt.Errorf("Connection is nil")
}

func (t *timeoutConn) LocalAddr() net.Addr {
	if t.conn != nil {
		return t.conn.LocalAddr()
	}
	return nil
}

func (t *timeoutConn) RemoteAddr() net.Addr {
	if t.conn != nil {
		return t.conn.RemoteAddr()
	}
	return nil
}

func (t *timeoutConn) SetDeadline(time time.Time) error {
	if t.conn != nil {
		return t.conn.SetDeadline(time)
	}
	return fmt.Errorf("Connection is nil")
}

func (t *timeoutConn) SetReadDeadline(time time.Time) error {
	if t.conn != nil {
		return t.conn.SetReadDeadline(time)
	}
	return fmt.Errorf("Connection is nil")
}

func (t *timeoutConn) SetWriteDeadline(time time.Time) error {
	if t.conn != nil {
		return t.conn.SetWriteDeadline(time)
	}
	return fmt.Errorf("Connection is nil")
}
