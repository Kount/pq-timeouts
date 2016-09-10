package pqtimeouts

import (
	"fmt"
	"net"
	"testing"
	"time"
)

type testNetConn struct {
	readCalled             int
	readError              error
	writeCalled            int
	writeError             error
	closeCalled            int
	closeError             error
	setReadDeadlineCalled  int
	setReadDeadlineTime    time.Time
	setReadDeadlineError   error
	setWriteDeadlineCalled int
	setWriteDeadlineTime   time.Time
	setWriteDeadlineError  error
}

func (t *testNetConn) Read(b []byte) (int, error) {
	t.readCalled++

	return 0, t.readError
}

func (t *testNetConn) Write(b []byte) (int, error) {
	t.writeCalled++

	return 0, t.writeError
}

func (t *testNetConn) Close() error {
	t.closeCalled++

	return t.closeError
}

func (t *testNetConn) LocalAddr() net.Addr {
	return nil
}

func (t *testNetConn) RemoteAddr() net.Addr {
	return nil
}

func (t *testNetConn) SetDeadline(time time.Time) error {
	return nil
}

func (t *testNetConn) SetReadDeadline(time time.Time) error {
	t.setReadDeadlineCalled++
	t.setReadDeadlineTime = time

	return t.setReadDeadlineError
}

func (t *testNetConn) SetWriteDeadline(time time.Time) error {
	t.setWriteDeadlineCalled++
	t.setWriteDeadlineTime = time

	return t.setWriteDeadlineError
}

func TestReadConnNil(t *testing.T) {
	conn := &timeoutConn{}

	b := make([]byte, 5)
	_, err := conn.Read(b)

	if err == nil {
		t.Error("An error was expected")
	}

	if err.Error() != "Connection is nil" {
		t.Errorf("Error was not as expected: %q", err.Error())
	}
}

func TestReadTimeoutNotSet(t *testing.T) {
	testConn := &testNetConn{}

	conn := &timeoutConn{conn: testConn}

	b := make([]byte, 5)
	_, err := conn.Read(b)

	if err != nil {
		t.Error("Unexpected error")
	}

	if testConn.readCalled != 1 {
		t.Error("Read should have been called but was not")
	}

	if testConn.setReadDeadlineCalled != 0 {
		t.Error("SetReadDeadline should not have been called and was")
	}
}

func TestReadTimeoutSet(t *testing.T) {
	testConn := &testNetConn{}

	conn := &timeoutConn{conn: testConn, readTimeout: time.Duration(500) * time.Millisecond}

	b := make([]byte, 5)
	_, err := conn.Read(b)

	if err != nil {
		t.Error("Unexpected error")
	}

	if testConn.readCalled != 1 {
		t.Error("Read should have been called but was not")
	}

	if testConn.setReadDeadlineCalled != 2 {
		t.Error("SetReadDeadline should have been called twice and was not")
	}

	emptyTime := time.Time{}
	if testConn.setReadDeadlineTime != emptyTime {
		t.Error("Deadline time should have been cleared and was not")
	}
}

func TestReadError(t *testing.T) {
	testConn := &testNetConn{readError: fmt.Errorf("i/o timeout")}

	conn := &timeoutConn{conn: testConn, readTimeout: time.Duration(500) * time.Millisecond}

	b := make([]byte, 5)
	_, err := conn.Read(b)

	if err == nil {
		t.Error("An error was expected")
	}

	if err.Error() != "i/o timeout" {
		t.Errorf("The error was not as expected: %q", err.Error())
	}

	if testConn.setReadDeadlineCalled != 2 {
		t.Error("SetReadDeadline should have been called twice and was not")
	}

	emptyTime := time.Time{}
	if testConn.setReadDeadlineTime != emptyTime {
		t.Error("Deadline time should have been cleared and was not")
	}
}

func TestWriteConnNil(t *testing.T) {
	conn := &timeoutConn{}

	b := []byte{'t', 'e', 's', 't'}
	_, err := conn.Write(b)

	if err == nil {
		t.Error("An error was expected")
	}

	if err.Error() != "Connection is nil" {
		t.Errorf("Error was not as expected: %q", err.Error())
	}
}

func TestWriteTimeoutNotSet(t *testing.T) {
	testConn := &testNetConn{}

	conn := &timeoutConn{conn: testConn}

	b := []byte{'t', 'e', 's', 't'}
	_, err := conn.Write(b)

	if err != nil {
		t.Error("Unexpected error")
	}

	if testConn.writeCalled != 1 {
		t.Error("Write should have been called but was not")
	}

	if testConn.setWriteDeadlineCalled != 0 {
		t.Error("SetWriteDeadline should not have been called and was")
	}
}

func TestWriteTimeoutSet(t *testing.T) {
	testConn := &testNetConn{}

	conn := &timeoutConn{conn: testConn, writeTimeout: time.Duration(500) * time.Millisecond}

	b := []byte{'t', 'e', 's', 't'}
	_, err := conn.Write(b)

	if err != nil {
		t.Error("Unexpected error")
	}

	if testConn.writeCalled != 1 {
		t.Error("Write should have been called but was not")
	}

	if testConn.setWriteDeadlineCalled != 2 {
		t.Error("SetWriteDeadline should have been called twice and was not")
	}

	emptyTime := time.Time{}
	if testConn.setWriteDeadlineTime != emptyTime {
		t.Error("Deadline time should have been cleared and was not")
	}
}

func TestWriteError(t *testing.T) {
	testConn := &testNetConn{writeError: fmt.Errorf("i/o timeout")}

	conn := &timeoutConn{conn: testConn, writeTimeout: time.Duration(500) * time.Millisecond}

	b := []byte{'t', 'e', 's', 't'}
	_, err := conn.Write(b)

	if err == nil {
		t.Error("An error was expected")
	}

	if err.Error() != "i/o timeout" {
		t.Errorf("The error was not as expected: %q", err.Error())
	}

	if testConn.setWriteDeadlineCalled != 2 {
		t.Error("SetWriteDeadline should have been called twice and was not")
	}

	emptyTime := time.Time{}
	if testConn.setWriteDeadlineTime != emptyTime {
		t.Error("Deadline time should have been cleared and was not")
	}
}

func TestClose(t *testing.T) {
	testConn := &testNetConn{}

	conn := &timeoutConn{conn: testConn}

	err := conn.Close()

	if err != nil {
		t.Error("Unexpected error")
	}

	if testConn.closeCalled != 1 {
		t.Error("Close should have been called and was not")
	}

	if conn.conn != nil {
		t.Error("The connection should have been set to nil and was not")
	}
}

func TestCloseConnNil(t *testing.T) {
	conn := &timeoutConn{}

	err := conn.Close()

	if err == nil {
		t.Error("An error was expected")
	}

	if err.Error() != "Connection is nil" {
		t.Errorf("Error was not as expected: %q", err.Error())
	}
}

func TestCloseError(t *testing.T) {
	testConn := &testNetConn{closeError: fmt.Errorf("i/o timeout")}

	conn := &timeoutConn{conn: testConn}

	err := conn.Close()

	if err == nil {
		t.Error("An error was expected")
	}

	if err.Error() != "i/o timeout" {
		t.Errorf("Error was not as expected: %q", err.Error())
	}

	if conn.conn == nil {
		t.Error("The connection should not be nil since there was an error")
	}

}
