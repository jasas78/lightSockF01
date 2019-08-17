//package lightsocks
package main

import (
	"io"
	"log"
	"net"
)

const (
	_VbufSize = 1024
)

// 加密传输的 TCP Socket
type _STsecureTCPConn struct {
	io.ReadWriteCloser
	Cipher *_STcipher
}

// 从输入流里读取加密过的数据，解密后把原数据放到bs里
func (secureSocket *_STsecureTCPConn) DecodeRead(bs []byte) (n int, err error) {
	n, err = secureSocket.Read(bs)
	if err != nil {
		return
	}
	secureSocket.Cipher._Fdecode(bs[:n])
	return
}

// 把放在bs里的数据加密后立即全部写入输出流
func (secureSocket *_STsecureTCPConn) EncodeWrite(bs []byte) (int, error) {
	secureSocket.Cipher._Fencode(bs)
	return secureSocket.Write(bs)
}

// 从src中源源不断的读取原数据加密后写入到dst，直到src中没有数据可以再读取
func (secureSocket *_STsecureTCPConn) EncodeCopy(dst io.ReadWriteCloser) error {
	buf := make([]byte, _VbufSize)
	for {
		readCount, errRead := secureSocket.Read(buf)
		if errRead != nil {
			if errRead != io.EOF {
				return errRead
			} else {
				return nil
			}
		}
		if readCount > 0 {
			writeCount, errWrite := (&_STsecureTCPConn{
				ReadWriteCloser: dst,
				Cipher:          secureSocket.Cipher,
			}).EncodeWrite(buf[0:readCount])
			if errWrite != nil {
				return errWrite
			}
			if readCount != writeCount {
				return io.ErrShortWrite
			}
		}
	}
}

// 从src中源源不断的读取加密后的数据解密后写入到dst，直到src中没有数据可以再读取
func (secureSocket *_STsecureTCPConn) DecodeCopy(dst io.Writer) error {
	buf := make([]byte, _VbufSize)
	for {
		readCount, errRead := secureSocket.DecodeRead(buf)
		if errRead != nil {
			if errRead != io.EOF {
				return errRead
			} else {
				return nil
			}
		}
		if readCount > 0 {
			writeCount, errWrite := dst.Write(buf[0:readCount])
			if errWrite != nil {
				return errWrite
			}
			if readCount != writeCount {
				return io.ErrShortWrite
			}
		}
	}
}

// see net.DialTCP
func DialTCPSecure(raddr *net.TCPAddr, ___Vcipher1 *_STcipher) (*_STsecureTCPConn, error) {
	remoteConn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		return nil, err
	}
	return &_STsecureTCPConn{
		ReadWriteCloser: remoteConn,
		Cipher:          ___Vcipher1,
	}, nil
}

// see net.ListenTCP
func ListenSecureTCP(laddr *net.TCPAddr, ___Vcipher2 *_STcipher, handleConn func(localConn *_STsecureTCPConn), didListen func(listenAddr net.Addr)) error {
	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		return err
	}

	defer listener.Close()

	if didListen != nil {
		didListen(listener.Addr())
	}

	for {
		localConn, err := listener.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}
		// localConn被关闭时直接清除所有数据 不管没有发送的数据
		localConn.SetLinger(0)
		go handleConn(&_STsecureTCPConn{
			ReadWriteCloser: localConn,
			Cipher:          ___Vcipher2,
		})
	}
}
