package pqtimeouts

import (
	"database/sql/driver"
	"reflect"
	"testing"
	"time"

	"github.com/lib/pq"
)

func TestOpenNoTimeoutsAdded(t *testing.T) {
	var connection string
	var dialer pq.Dialer

	testDialOpen := func(d pq.Dialer, name string) (_ driver.Conn, err error) {
		connection = name
		dialer = d
		return nil, nil
	}

	driver := TimeoutDriver{DialOpen: testDialOpen}

	testConnection := "user=pqtest dbname=pqtest sslmode=verify-full"
	_, err := driver.Open(testConnection)

	if err != nil {
		t.Error("Unexpected error")
	}

	if connection != testConnection {
		t.Errorf("The connection string was not as expected: %q", connection)
	}

	if reflect.TypeOf(dialer).String() != "pqtimeouts.timeoutDialer" {
		t.Errorf("The type of the dialer is not as expected: %q", reflect.TypeOf(dialer).String())
	}
}

func TestOpenReadTimeoutAdded(t *testing.T) {
	var connection string
	var dialer pq.Dialer

	testDialOpen := func(d pq.Dialer, name string) (_ driver.Conn, err error) {
		connection = name
		dialer = d
		return nil, nil
	}

	driver := TimeoutDriver{DialOpen: testDialOpen}

	_, err := driver.Open("user=pqtest read_timeout=700 dbname=pqtest sslmode=verify-full")

	if err != nil {
		t.Error("Unexpected error")
	}

	if connection != "user=pqtest dbname=pqtest sslmode=verify-full" {
		t.Errorf("The connection string was not as expected: %q", connection)
	}

	if reflect.TypeOf(dialer).String() != "pqtimeouts.timeoutDialer" {
		t.Errorf("The type of the dialer is not as expected: %q", reflect.TypeOf(dialer).String())
	}

	if toDialer, ok := dialer.(timeoutDialer); ok {
		if toDialer.readTimeout != time.Duration(700)*time.Millisecond {
			t.Error("Read timeout was not set to the correct duration")
		}
	} else {
		t.Error("The dialer is not a timeoutDialer")
	}
}

func TestOpenWriteTimeoutAdded(t *testing.T) {
	var connection string
	var dialer pq.Dialer

	testDialOpen := func(d pq.Dialer, name string) (_ driver.Conn, err error) {
		connection = name
		dialer = d
		return nil, nil
	}

	driver := TimeoutDriver{DialOpen: testDialOpen}

	_, err := driver.Open(" user=pqtest write_timeout=968      dbname=pqtest sslmode=verify-full		")

	if err != nil {
		t.Error("Unexpected error")
	}

	if connection != "user=pqtest dbname=pqtest sslmode=verify-full" {
		t.Errorf("The connection string was not as expected: %q", connection)
	}

	if reflect.TypeOf(dialer).String() != "pqtimeouts.timeoutDialer" {
		t.Errorf("The type of the dialer is not as expected: %q", reflect.TypeOf(dialer).String())
	}

	if toDialer, ok := dialer.(timeoutDialer); ok {
		if toDialer.writeTimeout != time.Duration(968)*time.Millisecond {
			t.Error("Read timeout was not set to the correct duration")
		}
	} else {
		t.Error("The dialer is not a timeoutDialer")
	}
}

func TestOpenTimeoutsAddedWriteError(t *testing.T) {
	dialOpenCalled := false

	testDialOpen := func(d pq.Dialer, name string) (_ driver.Conn, err error) {
		dialOpenCalled = true
		return nil, nil
	}

	driver := TimeoutDriver{DialOpen: testDialOpen}

	_, err := driver.Open(" user=pqtest write_timeout=seven read_timeout=7    dbname=pqtest sslmode=verify-full		")

	if err == nil {
		t.Error("An error was expected")
	}

	if err.Error() != "Error interpreting value for write_timeout" {
		t.Errorf("The error is unexpected: %q", err.Error())
	}

	if dialOpenCalled != false {
		t.Error("DialOpen should not have been called and was")
	}
}

func TestOpenTimeoutsAddedReadError(t *testing.T) {
	dialOpenCalled := false

	testDialOpen := func(d pq.Dialer, name string) (_ driver.Conn, err error) {
		dialOpenCalled = true
		return nil, nil
	}

	driver := TimeoutDriver{DialOpen: testDialOpen}

	_, err := driver.Open(" user=pqtest    write_timeout=680 read_timeout=    dbname=pqtest sslmode=verify-full		")

	if err == nil {
		t.Error("An error was expected")
	}

	if err.Error() != "Error interpreting value for read_timeout" {
		t.Errorf("The error is unexpected: %q", err.Error())
	}

	if dialOpenCalled != false {
		t.Error("DialOpen should not have been called and was")
	}
}

func TestPostgresURL(t *testing.T) {
	var connection string
	var dialer pq.Dialer

	testDialOpen := func(d pq.Dialer, name string) (_ driver.Conn, err error) {
		connection = name
		dialer = d
		return nil, nil
	}

	driver := TimeoutDriver{DialOpen: testDialOpen}

	_, err := driver.Open("postgres://pqtest:password@localhost/pqtest?read_timeout=500&sslmode=verify-full&write_timeout=100")

	if err != nil {
		t.Error("Unexpected error")
	}

	if connection != "dbname=pqtest host=localhost password=password sslmode=verify-full user=pqtest" {
		t.Errorf("The connection string was not as expected: %q", connection)
	}

	if toDialer, ok := dialer.(timeoutDialer); ok {
		if toDialer.readTimeout != time.Duration(500)*time.Millisecond {
			t.Error("Read timeout was not set to the correct duration")
		}

		if toDialer.writeTimeout != time.Duration(100)*time.Millisecond {
			t.Error("Write timeout was not set to the correct duration")
		}

	} else {
		t.Error("The dialer is not a timeoutDialer")
	}
}

func TestPostgresqlURLError(t *testing.T) {
	dialOpenCalled := false

	testDialOpen := func(d pq.Dialer, name string) (_ driver.Conn, err error) {
		dialOpenCalled = true
		return nil, nil
	}

	driver := TimeoutDriver{DialOpen: testDialOpen}

	_, err := driver.Open("postgresql://pqtest\\\\/:password@localhost/pqtest?read_timeout=500&sslmode=verify-full&write_timeout=100")

	if err == nil {
		t.Error("An error was expected")
	}

	if err.Error() != "parse postgresql://pqtest\\\\/:password@localhost/pqtest?read_timeout=500&sslmode=verify-full&write_timeout=100: invalid character \"\\\\\" in host name" {
		t.Errorf("The error was not as expected: %q", err.Error())
	}

	if dialOpenCalled {
		t.Error("DialOpen should not have been called")
	}
}
