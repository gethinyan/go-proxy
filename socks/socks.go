package socks

import (
	"errors"
	"io"
	"net"
	"strconv"
)

// MaxAddrLen 最大的地址长度
const MaxAddrLen = 1 + 1 + 255 + 2

// MaxReqLen 最大的请求长度
const MaxReqLen = 1 + 1 + 1 + MaxAddrLen

// 客户端请求的类型
const (
	CmdConnect = 1
)

// 地址类型
const (
	ATypeDomain = 3
	ATypeIPV4   = 1
	ATypeIPV6   = 4
)

// 错误信息
var (
	ErrInvalidCmd   = errors.New("invalid cmd")
	ErrInvalidAType = errors.New("invalid aType")
)

// Addr 请求地址字节流
type Addr []byte

func (a Addr) String() string {
	host, port := "", ""

	switch a[0] {
	case ATypeDomain:
		host = string(a[2 : a[1]+2])
		port = strconv.Itoa((int(a[a[1]+2]) << 8) | int(a[a[1]+2+1]))
	case ATypeIPV4:
	case ATypeIPV6:
	}

	return net.JoinHostPort(host, port)
}

// ReadAddr 读取请求地址
func ReadAddr(r io.Reader) (Addr, error) {
	buf := make([]byte, MaxAddrLen)
	_, err := io.ReadFull(r, buf[:1])
	if err != nil {
		return nil, err
	}

	switch buf[0] {
	case ATypeDomain:
		_, err = io.ReadFull(r, buf[1:2])
		if err != nil {
			return nil, err
		}
		_, err = io.ReadFull(r, buf[2:2+buf[1]+2])
		return buf[:1+1+buf[1]+2], err
	case ATypeIPV4:
		_, err = io.ReadFull(r, buf[1:1+net.IPv4len+2])
		return buf[:1+net.IPv4len+2], err
	case ATypeIPV6:
		_, err = io.ReadFull(r, buf[1:1+net.IPv6len+2])
		return buf[:1+net.IPv6len+2], err
	}

	return nil, ErrInvalidAType
}

// HandShake 建立 socks5 连接
func HandShake(sc net.Conn) (Addr, error) {
	buf := make([]byte, MaxReqLen)
	// 读取 VER NMETHODS METHODS
	if _, err := io.ReadFull(sc, buf[:2]); err != nil {
		return nil, err
	}
	nMethods := buf[1]
	if _, err := io.ReadFull(sc, buf[:nMethods]); err != nil {
		return nil, err
	}
	// 通知客户端使用的 VER METHOD
	if _, err := sc.Write([]byte{5, 0}); err != nil {
		return nil, err
	}
	// 读取 VER CMD RSV ATYP DST.ADDR DST.PORT
	if _, err := io.ReadFull(sc, buf[:3]); err != nil {
		return nil, err
	}

	addr, err := ReadAddr(sc)
	if err != nil {
		return nil, err
	}

	if buf[1] != CmdConnect {
		return nil, ErrInvalidCmd
	}
	if _, err := sc.Write([]byte{5, 0, 0, 1, 0, 0, 0, 0, 0, 0}); err != nil {
		return nil, err
	}

	return addr, nil
}
