package pqtimeouts

import (
	"fmt"
	"net"
	"testing"
	"time"
)

type testAddr struct {
	network string
	address string
}

func (t testAddr) Network() string {
	return t.network
}

func (t testAddr) String() string {
	return t.address
}

type testNetConn struct {
	readCalled               int
	readError                error
	writeCalled              int
	writeError               error
	closeCalled              int
	closeError               error
	localAddrCalled          int
	localAddr                net.Addr
	remoteAddrCalled         int
	remoteAddr               net.Addr
	setDeadlineCalled        int
	setDeadlineError         error
	setDeadlineTime          time.Time
	setReadDeadlineCalled    int
	setReadDeadlineTimePrev  time.Time
	setReadDeadlineTime      time.Time
	setReadDeadlineError     error
	setWriteDeadlineCalled   int
	setWriteDeadlineTimePrev time.Time
	setWriteDeadlineTime     time.Time
	setWriteDeadlineError    error
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
	t.localAddrCalled++

	return t.localAddr
}

func (t *testNetConn) RemoteAddr() net.Addr {
	t.remoteAddrCalled++

	return t.remoteAddr
}

func (t *testNetConn) SetDeadline(time time.Time) error {
	t.setDeadlineCalled++
	t.setDeadlineTime = time

	return t.setDeadlineError
}

func (t *testNetConn) SetReadDeadline(time time.Time) error {
	t.setReadDeadlineCalled++

	if t.setReadDeadlineTime != time {
		t.setReadDeadlineTimePrev = t.setReadDeadlineTime
	}
	t.setReadDeadlineTime = time

	return t.setReadDeadlineError
}

func (t *testNetConn) SetWriteDeadline(time time.Time) error {
	t.setWriteDeadlineCalled++

	if t.setWriteDeadlineTime != time {
		t.setWriteDeadlineTimePrev = t.setWriteDeadlineTime
	}
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
		t.Errorf("Deadline time was not as expected: %+v", testConn.setReadDeadlineTime)
	}

	if testConn.setReadDeadlineTimePrev == emptyTime {
		t.Errorf("Previous deadline time was unexpected: %+v", testConn.setReadDeadlineTimePrev)
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
		t.Errorf("Deadline time was not as expected: %+v", testConn.setReadDeadlineTime)
	}

	if testConn.setReadDeadlineTimePrev == emptyTime {
		t.Errorf("Previous deadline time was unexpected: %+v", testConn.setReadDeadlineTimePrev)
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
		t.Errorf("Deadline time was not as expected: %+v", testConn.setWriteDeadlineTime)
	}

	if testConn.setWriteDeadlineTimePrev == emptyTime {
		t.Errorf("Previous deadline time was unexpected: %+v", testConn.setWriteDeadlineTimePrev)
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
		t.Errorf("Deadline time was not as expected: %+v", testConn.setWriteDeadlineTime)
	}

	if testConn.setWriteDeadlineTimePrev == emptyTime {
		t.Errorf("Previous deadline time was unexpected: %+v", testConn.setWriteDeadlineTimePrev)
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

	if testConn.closeCalled != 1 {
		t.Error("Close should have been called and was not")
	}

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

func TestLocalAddr(t *testing.T) {
	testConn := &testNetConn{localAddr: testAddr{network: "test", address: "testtest"}}

	conn := &timeoutConn{conn: testConn}

	addr := conn.LocalAddr()

	if testConn.localAddrCalled != 1 {
		t.Error("LocalAddr should have been called and was not")
	}

	if addr.Network() != "test" && addr.String() != "testtest" {
		t.Errorf("Local Addr was not as expected: %+v", addr)
	}
}

func TestLocalAddrConnNil(t *testing.T) {
	conn := &timeoutConn{}

	addr := conn.LocalAddr()

	if addr != nil {
		t.Error("Local Addr should have been nil")
	}
}

func TestRemoteAddr(t *testing.T) {
	testConn := &testNetConn{remoteAddr: testAddr{network: "test", address: "testtest"}}

	conn := &timeoutConn{conn: testConn}

	addr := conn.RemoteAddr()

	if testConn.remoteAddrCalled != 1 {
		t.Error("RemoteAddr should have been called and was not")
	}

	if addr.Network() != "test" && addr.String() != "testtest" {
		t.Errorf("Remote Addr was not as expected: %+v", addr)
	}
}

func TestLocalRemoteConnNil(t *testing.T) {
	conn := &timeoutConn{}

	addr := conn.RemoteAddr()

	if addr != nil {
		t.Error("Remote Addr should have been nil")
	}
}

func TestSetDeadline(t *testing.T) {
	testConn := &testNetConn{}

	conn := &timeoutConn{conn: testConn}

	err := conn.SetDeadline(time.Time{})

	if err != nil {
		t.Error("Unexpected error")
	}

	if testConn.setDeadlineCalled != 1 {
		t.Error("SetDeadline should have been called and was not")
	}

	emptyTime := time.Time{}
	if testConn.setDeadlineTime != emptyTime {
		t.Errorf("Deadline time was not as expected: %+v", testConn.setDeadlineTime)
	}
}

func TestSetDeadlineError(t *testing.T) {
	testConn := &testNetConn{setDeadlineError: fmt.Errorf("Can't set deadline")}

	conn := &timeoutConn{conn: testConn}

	err := conn.SetDeadline(time.Time{})

	if testConn.setDeadlineCalled != 1 {
		t.Error("SetDeadline should have been called and was not")
	}

	if err == nil {
		t.Error("An error was expected")
	}

	if err.Error() != "Can't set deadline" {
		t.Errorf("Error was not as expected: %q", err.Error())
	}
}

func TestSetDeadlineConnNil(t *testing.T) {
	conn := &timeoutConn{}

	err := conn.SetDeadline(time.Time{})

	if err == nil {
		t.Error("An error was expected.")
	}

	if err.Error() != "Connection is nil" {
		t.Errorf("Error was not as expected: %q", err.Error())
	}
}

func TestSetReadDeadline(t *testing.T) {
	testConn := &testNetConn{}

	conn := &timeoutConn{conn: testConn}

	err := conn.SetReadDeadline(time.Time{})

	if err != nil {
		t.Error("Unexpected error")
	}

	if testConn.setReadDeadlineCalled != 1 {
		t.Error("SetReadDeadline should have been called and was not")
	}

	emptyTime := time.Time{}
	if testConn.setReadDeadlineTime != emptyTime {
		t.Errorf("Deadline time was not as expected: %+v", testConn.setReadDeadlineTime)
	}
}

func TestSetReadDeadlineError(t *testing.T) {
	testConn := &testNetConn{setReadDeadlineError: fmt.Errorf("Can't set deadline")}

	conn := &timeoutConn{conn: testConn}

	err := conn.SetReadDeadline(time.Time{})

	if testConn.setReadDeadlineCalled != 1 {
		t.Error("SetReadDeadline should have been called and was not")
	}

	if err == nil {
		t.Error("An error was expected")
	}

	if err.Error() != "Can't set deadline" {
		t.Errorf("Error was not as expected: %q", err.Error())
	}
}

func TestSetReadDeadlineConnNil(t *testing.T) {
	conn := &timeoutConn{}

	err := conn.SetReadDeadline(time.Time{})

	if err == nil {
		t.Error("An error was expected.")
	}

	if err.Error() != "Connection is nil" {
		t.Errorf("Error was not as expected: %q", err.Error())
	}
}

func TestSetWriteDeadline(t *testing.T) {
	testConn := &testNetConn{}

	conn := &timeoutConn{conn: testConn}

	err := conn.SetWriteDeadline(time.Time{})

	if err != nil {
		t.Error("Unexpected error")
	}

	if testConn.setWriteDeadlineCalled != 1 {
		t.Error("SetWriteDeadline should have been called and was not")
	}

	emptyTime := time.Time{}
	if testConn.setWriteDeadlineTime != emptyTime {
		t.Errorf("Deadline time was not as expected: %+v", testConn.setWriteDeadlineTime)
	}
}

func TestSetWriteDeadlineError(t *testing.T) {
	testConn := &testNetConn{setWriteDeadlineError: fmt.Errorf("Can't set deadline")}

	conn := &timeoutConn{conn: testConn}

	err := conn.SetWriteDeadline(time.Time{})

	if testConn.setWriteDeadlineCalled != 1 {
		t.Error("SetWriteDeadline should have been called and was not")
	}

	if err == nil {
		t.Error("An error was expected")
	}

	if err.Error() != "Can't set deadline" {
		t.Errorf("Error was not as expected: %q", err.Error())
	}
}

func TestSetWriteDeadlineConnNil(t *testing.T) {
	conn := &timeoutConn{}

	err := conn.SetWriteDeadline(time.Time{})

	if err == nil {
		t.Error("An error was expected.")
	}

	if err.Error() != "Connection is nil" {
		t.Errorf("Error was not as expected: %q", err.Error())
	}
}
