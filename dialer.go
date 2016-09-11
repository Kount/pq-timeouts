package pqtimeouts

import (
	"net"
	"time"
)

type timeoutDialer struct {
	netDial        func(string, string) (net.Conn, error)                // Allow this to be stubbed for testing
	netDialTimeout func(string, string, time.Duration) (net.Conn, error) // Allow this to be stubbed for testing
	readTimeout    time.Duration
	writeTimeout   time.Duration
}

func (t timeoutDialer) Dial(network string, address string) (net.Conn, error) {
	// If we don't have any timeouts set, just return a normal connection
	if t.readTimeout == 0 && t.writeTimeout == 0 {
		return t.netDial(network, address)
	}

	// Otherwise we want a timeoutConn to handle the read and write deadlines for us.
	c, err := t.netDial(network, address)
	if err != nil || c == nil {
		return c, err
	}

	return &timeoutConn{conn: c, readTimeout: t.readTimeout, writeTimeout: t.writeTimeout}, nil
}

func (t timeoutDialer) DialTimeout(network string, address string, timeout time.Duration) (net.Conn, error) {
	// If we don't have any timeouts set, just return a normal connection
	if t.readTimeout == 0 && t.writeTimeout == 0 {
		return t.netDialTimeout(network, address, timeout)
	}

	// Otherwise we want a timeoutConn to handle the read and write deadlines for us.
	c, err := t.netDialTimeout(network, address, timeout)
	if err != nil || c == nil {
		return c, err
	}

	return &timeoutConn{conn: c, readTimeout: t.readTimeout, writeTimeout: t.writeTimeout}, nil
}
