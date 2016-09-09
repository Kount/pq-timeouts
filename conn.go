package pqtimeouts

import (
	"fmt"
	"net"
	"time"
)

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
		return
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
