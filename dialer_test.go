package pqtimeouts

import (
	"fmt"
	"net"
	"reflect"
	"testing"
	"time"
)

func TestDialNoTimeouts(t *testing.T) {
	testConn := &testNetConn{}

	testDial := func(network string, address string) (net.Conn, error) {
		return testConn, nil
	}

	testDialTimeout := func(network string, address string, timeout time.Duration) (net.Conn, error) {
		return testConn, nil
	}

	dialer := timeoutDialer{netDial: testDial, netDialTimeout: testDialTimeout}

	conn, err := dialer.Dial("testNetwork", "testAddress")

	if err != nil {
		t.Error("Unexpected error")
	}

	if conn == nil {
		t.Error("Connection should not be nil")
	}

	if reflect.TypeOf(conn).String() != "*pqtimeouts.testNetConn" {
		t.Errorf("Connection type was not as expected: %q", reflect.TypeOf(conn).String())
	}
}

func TestDialWithReadTimeout(t *testing.T) {
	testConn := &testNetConn{}

	testDial := func(network string, address string) (net.Conn, error) {
		return testConn, nil
	}

	testDialTimeout := func(network string, address string, timeout time.Duration) (net.Conn, error) {
		return testConn, nil
	}

	dialer := timeoutDialer{
		netDial:        testDial,
		netDialTimeout: testDialTimeout,
		readTimeout:    time.Duration(1000)}

	conn, err := dialer.Dial("testNetwork", "testAddress")

	if err != nil {
		t.Error("Unexpected error")
	}

	if conn == nil {
		t.Error("Connection should not be nil")
	}

	if reflect.TypeOf(conn).String() != "*pqtimeouts.timeoutConn" {
		t.Errorf("Connection type was not as expected: %q", reflect.TypeOf(conn).String())
	}
}

func TestDialWithWriteTimeout(t *testing.T) {
	testConn := &testNetConn{}

	testDial := func(network string, address string) (net.Conn, error) {
		return testConn, nil
	}

	testDialTimeout := func(network string, address string, timeout time.Duration) (net.Conn, error) {
		return testConn, nil
	}

	dialer := timeoutDialer{
		netDial:        testDial,
		netDialTimeout: testDialTimeout,
		writeTimeout:   time.Duration(1000)}

	conn, err := dialer.Dial("testNetwork", "testAddress")

	if err != nil {
		t.Error("Unexpected error")
	}

	if conn == nil {
		t.Error("Connection should not be nil")
	}

	if reflect.TypeOf(conn).String() != "*pqtimeouts.timeoutConn" {
		t.Errorf("Connection type was not as expected: %q", reflect.TypeOf(conn).String())
	}
}

func TestDialError(t *testing.T) {
	testConn := &testNetConn{}

	testDial := func(network string, address string) (net.Conn, error) {
		return nil, fmt.Errorf("Could not connect")
	}

	testDialTimeout := func(network string, address string, timeout time.Duration) (net.Conn, error) {
		return testConn, nil
	}

	dialer := timeoutDialer{
		netDial:        testDial,
		netDialTimeout: testDialTimeout,
		readTimeout:    time.Duration(1000),
		writeTimeout:   time.Duration(1000)}

	conn, err := dialer.Dial("testNetwork", "testAddress")

	if err == nil {
		t.Error("An error was expected")
	}

	if err.Error() != "Could not connect" {
		t.Errorf("Error was not as expected: %q", err.Error())
	}

	if conn != nil {
		t.Error("Connection should be nil")
	}
}

func TestDialTimeoutNoTimeouts(t *testing.T) {
	testConn := &testNetConn{}

	testDial := func(network string, address string) (net.Conn, error) {
		return testConn, nil
	}

	testDialTimeout := func(network string, address string, timeout time.Duration) (net.Conn, error) {
		return testConn, nil
	}

	dialer := timeoutDialer{netDial: testDial, netDialTimeout: testDialTimeout}

	conn, err := dialer.DialTimeout("testNetwork", "testAddress", time.Duration(500))

	if err != nil {
		t.Error("Unexpected error")
	}

	if conn == nil {
		t.Error("Connection should not be nil")
	}

	if reflect.TypeOf(conn).String() != "*pqtimeouts.testNetConn" {
		t.Errorf("Connection type was not as expected: %q", reflect.TypeOf(conn).String())
	}
}

func TestDialTimeoutWithReadTimeout(t *testing.T) {
	testConn := &testNetConn{}

	testDial := func(network string, address string) (net.Conn, error) {
		return testConn, nil
	}

	testDialTimeout := func(network string, address string, timeout time.Duration) (net.Conn, error) {
		return testConn, nil
	}

	dialer := timeoutDialer{
		netDial:        testDial,
		netDialTimeout: testDialTimeout,
		readTimeout:    time.Duration(1000)}

	conn, err := dialer.DialTimeout("testNetwork", "testAddress", time.Duration(500))

	if err != nil {
		t.Error("Unexpected error")
	}

	if conn == nil {
		t.Error("Connection should not be nil")
	}

	if reflect.TypeOf(conn).String() != "*pqtimeouts.timeoutConn" {
		t.Errorf("Connection type was not as expected: %q", reflect.TypeOf(conn).String())
	}
}

func TestDialTimeoutWithWriteTimeout(t *testing.T) {
	testConn := &testNetConn{}

	testDial := func(network string, address string) (net.Conn, error) {
		return testConn, nil
	}

	testDialTimeout := func(network string, address string, timeout time.Duration) (net.Conn, error) {
		return testConn, nil
	}

	dialer := timeoutDialer{
		netDial:        testDial,
		netDialTimeout: testDialTimeout,
		writeTimeout:   time.Duration(1000)}

	conn, err := dialer.DialTimeout("testNetwork", "testAddress", time.Duration(500))

	if err != nil {
		t.Error("Unexpected error")
	}

	if conn == nil {
		t.Error("Connection should not be nil")
	}

	if reflect.TypeOf(conn).String() != "*pqtimeouts.timeoutConn" {
		t.Errorf("Connection type was not as expected: %q", reflect.TypeOf(conn).String())
	}
}

func TestDialTimeoutError(t *testing.T) {
	testConn := &testNetConn{}

	testDial := func(network string, address string) (net.Conn, error) {
		return testConn, nil
	}

	testDialTimeout := func(network string, address string, timeout time.Duration) (net.Conn, error) {
		return nil, fmt.Errorf("Could not connect")
	}

	dialer := timeoutDialer{
		netDial:        testDial,
		netDialTimeout: testDialTimeout,
		readTimeout:    time.Duration(1000),
		writeTimeout:   time.Duration(1000)}

	conn, err := dialer.DialTimeout("testNetwork", "testAddress", time.Duration(500))

	if err == nil {
		t.Error("An error was expected")
	}

	if err.Error() != "Could not connect" {
		t.Errorf("Error was not as expected: %q", err.Error())
	}

	if conn != nil {
		t.Error("Connection should be nil")
	}
}
