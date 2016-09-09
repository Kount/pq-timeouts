package pqtimeouts

import (
	"net"
	"time"
)

type testConn struct {
}

func (t *testConn) Read(b []byte) (int, error) {
	return 0, nil
}

func (t *testConn) Write(b []byte) (int, error) {
	return 0, nil
}

func (t *testConn) Close() error {
	return nil
}

func (t *testConn) LocalAddr() net.Addr {
	return nil
}

func (t *testConn) RemoteAddr() net.Addr {
	return nil
}

func (t *testConn) SetDeadline(time time.Time) error {
	return nil
}

func (t *testConn) SetReadDeadline(time time.Time) error {
	return nil
}

func (t *testConn) SetWriteDeadline(time time.Time) error {
	return nil
}
