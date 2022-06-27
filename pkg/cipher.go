package pkg

import (
	"fmt"
	"net"

	"golang.org/x/crypto/chacha20"
)

// EncryptKey 加密 key
const EncryptKey = "12345678123456781234567812345678"

// CipherConn 流数据加解密连接
type CipherConn struct {
	Conn net.Conn
}

var nonce []byte

func init() {
	nonce = []byte{199, 85, 121, 195, 59, 196, 122, 51, 254, 131, 195, 121, 19, 94, 150, 197, 246, 80, 119, 133, 127, 211, 235, 90}
}

// Encode 加密
func (cc CipherConn) Encode(src []byte, dst []byte) {
	encoder, err := chacha20.NewUnauthenticatedCipher([]byte(EncryptKey), nonce)
	if err != nil {
		fmt.Println("chacha20 NewUnauthenticatedCipher fail")
		fmt.Println(err)
		return
	}

	// fmt.Println("EncodeBefore：", src)
	// 一个字节一个字节加密
	for i := 0; i < len(src); i++ {
		encoder.XORKeyStream(dst[i:i+1], src[i:i+1])
	}
	// fmt.Println("EncodeAfter：", dst)

	return
}

// Decode 解密
func (cc CipherConn) Decode(dst []byte, src []byte, n int) {
	decoder, err := chacha20.NewUnauthenticatedCipher([]byte(EncryptKey), nonce)
	if err != nil {
		fmt.Println("chacha20 NewUnauthenticatedCipher fail")
		fmt.Println(err)
		return
	}

	// fmt.Println("DecodeBefore：", dst)
	// 一个字节一个字节解密
	for i := 0; i < n; i++ {
		decoder.XORKeyStream(src[i:i+1], dst[i:i+1])
	}
	// fmt.Println("DecodeAfter：", src)

	return
}

// Read 读字节流
func (cc CipherConn) Read(b []byte) (n int, err error) {
	buf := make([]byte, len(b))
	n, err = cc.Conn.Read(buf)
	if err != nil {
		return 0, err
	}
	// 加密
	cc.Decode(buf, b, n)

	return
}

// Write 写字节流
func (cc CipherConn) Write(b []byte) (n int, err error) {
	buf := make([]byte, len(b))
	// 加密
	cc.Encode(b, buf)
	n, err = cc.Conn.Write(buf)
	if err != nil {
		return 0, err
	}

	return
}
