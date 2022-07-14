package main

import (
	"fmt"
	"io"
	"net"

	"github.com/gethinyan/go-proxy/pkg"
	"github.com/gethinyan/go-proxy/socks"
)

func main() {
	l, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println("Client listen fail")
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Client accept fail")
		}

		go handleClient(conn)
	}
}

func handleClient(sc net.Conn) {
	defer sc.Close()

	addr, err := socks.HandShake(sc)
	if err != nil {
		return
	}

	// rc 传输需要加解密
	rc, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("Dial server fail")
	}
	defer rc.Close()
	cc := pkg.CipherConn{Conn: rc}

	if _, err = cc.Write(addr); err != nil {
		fmt.Printf("failed to send address: %v", err)
		return
	}

	go io.Copy(cc, sc)
	io.Copy(sc, cc)
}
