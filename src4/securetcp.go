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
func (secureSocket *_STsecureTCPConn) DecodeRead(bs []byte) (___Vn int, ___Verr1 error) {
	___Vn, ___Verr1 = secureSocket.Read(bs)
	if ___Verr1 != nil {
		return
	}
	secureSocket.Cipher._Fdecode(bs[:___Vn])
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
func _FdialTCPSecure(raddr *net.TCPAddr, ___Vcipher1 *_STcipher) (*_STsecureTCPConn, error) {
	remoteConn, __Verr2 := net.DialTCP("tcp", nil, raddr)
	if __Verr2 != nil {
		return nil, __Verr2
	}
	return &_STsecureTCPConn{
		ReadWriteCloser: remoteConn,
		Cipher:          ___Vcipher1,
	}, nil
}

// see net.ListenTCP
func _FlistenSecureTCP(___VlistenAddr *net.TCPAddr, ___Vcipher2 *_STcipher,
	___FhandleConnCallback func(*_STsecureTCPConn),
	___VlistenInitShowConfig func(net.Addr)) error {

	___VacceptListenER2, __Verr3 := net.ListenTCP("tcp", ___VlistenAddr)
	if __Verr3 != nil {
		return __Verr3
	}

	defer ___VacceptListenER2.Close()

	if ___VlistenInitShowConfig != nil {
		___VlistenInitShowConfig(___VacceptListenER2.Addr())
	}

	for {
		___VlistenAcceptConnEach, __Verr4 := ___VacceptListenER2.AcceptTCP()
		if __Verr4 != nil {
			log.Println(__Verr4)
			continue
		}
		// ___VlistenAcceptConnEach被关闭时直接清除所有数据 不管没有发送的数据
		___VlistenAcceptConnEach.SetLinger(0)

		// _FsrvHandleConn (srv) , _FhandleConnClient
		go ___FhandleConnCallback(&_STsecureTCPConn{
			ReadWriteCloser: ___VlistenAcceptConnEach,
			Cipher:          ___Vcipher2,
		})
	}
}
