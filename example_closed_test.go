package neterrs_test

import (
	"fmt"
	"net"
	"reflect"

	"github.com/kezhuw/neterrs"
)

func getClosedConn() net.Conn {
	l, err := net.ListenTCP("tcp", &net.TCPAddr{})
	if err != nil {
		panic(err)
	}
	done := make(chan struct{})
	defer close(done)
	go func() {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		<-done
		conn.Close()
		l.Close()
	}()
	conn, err := net.DialTCP("tcp", nil, l.Addr().(*net.TCPAddr))
	if err != nil {
		panic(err)
	}
	conn.Close()
	return conn
}

func ExampleIsClosed_listenerAccept() {
	l, err := net.ListenTCP("tcp", &net.TCPAddr{})
	if err != nil {
		panic(err)
	}
	l.Close()
	_, err = l.Accept()
	fmt.Println(neterrs.IsClosed(err))
	if !neterrs.IsClosed(err) {
		fmt.Println(err)
		fmt.Println(reflect.TypeOf(err))
	}
	// Output: true
}

func ExampleIsClosed_connRead() {
	var buf [1]byte
	conn := getClosedConn()
	_, err := conn.Read(buf[:])
	fmt.Println(neterrs.IsClosed(err))
	// Output: true
}

func ExampleIsClosed_connWrite() {
	var buf [1]byte
	conn := getClosedConn()
	_, err := conn.Write(buf[:])
	fmt.Println(neterrs.IsClosed(err))
	// Output: true
}
