package main

import (
	"fmt"
	"io"
	"net"

	"github.com/gethinyan/go-proxy/pkg"
	"github.com/gethinyan/go-proxy/socks"
)

func main() {
	ln, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println("Server listen fail")
		return
	}
	for {
		sc, err := ln.Accept()
		if err != nil {
			fmt.Println("Server accept fail")
			continue
		}
		go handleServer(sc)
	}
}

func handleServer(sc net.Conn) {
	defer sc.Close()
	cc := pkg.CipherConn{Conn: sc}

	addr, err := socks.ReadAddr(sc)
	if err != nil {
		return
	}

	rc, err := net.Dial("tcp", addr.String())
	if err != nil {
		fmt.Println("server dial fail")
		return
	}
	defer rc.Close()

	go func() {
		io.Copy(rc, cc)
	}()
	io.Copy(cc, rc)
}
