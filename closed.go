// Package neterrs exports functions and errors that not exported by official
// Go net package.
package neterrs

import "net"

// ErrClosed equals to net.errClosing.
var ErrClosed = makeErrClosed()

func acceptOne(l *net.TCPListener, done chan struct{}) {
	conn, err := l.Accept()
	if err != nil {
		panic(err)
	}
	l.Close()
	<-done
	conn.Close()
}

func listenRand() (*net.TCPAddr, chan struct{}) {
	var addr net.TCPAddr
	addr.IP = net.IPv4(127, 0, 0, 1)
	l, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		panic(err)
	}
	addr.Port = l.Addr().(*net.TCPAddr).Port
	done := make(chan struct{})
	go acceptOne(l, done)
	return &addr, done
}

func triggerErrClosed(conn *net.TCPConn) error {
	var buf [1]byte
	_, err := conn.Read(buf[:])
	if opErr, ok := err.(*net.OpError); ok {
		return opErr.Err
	}
	panic("neterrs: unexpected error for reading after closed")
}

func makeErrClosed() error {
	addr, done := listenRand()
	defer close(done)
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		panic(err)
	}
	conn.Close()
	return triggerErrClosed(conn)
}

// IsClosed returns a boolean indicating whether the error is caused by
// closed connection.
func IsClosed(err error) bool {
	switch err := err.(type) {
	case *net.OpError:
		return err.Err == ErrClosed
	default:
		return err == ErrClosed
	}
}
