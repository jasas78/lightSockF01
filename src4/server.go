//package lightsocks
package main

import (
	"encoding/binary"
	"net"
)

type _STlsServer struct {
	Cipher     *_STcipher
	ListenAddr *net.TCPAddr
}

// 新建一个服务端
// 服务端的职责是:
// 1. 监听来自本地代理客户端的请求
// 2. 解密本地代理客户端请求的数据，解析 SOCKS5 协议，连接用户浏览器真正想要连接的远程服务器
// 3. 转发用户浏览器真正想要连接的远程服务器返回的数据的加密后的内容到本地代理客户端
func _Fserver_NewLsServer(___Vpassword3 string, ___VlistenAddr1 string) (*_STlsServer, error) {
	__VbsPassword2, __Verr1 := _FparsePassword(___Vpassword3)
	if __Verr1 != nil {
		return nil, __Verr1
	}
	__VstructListenAddr2, __Verr2 := net.ResolveTCPAddr("tcp", ___VlistenAddr1)
	if __Verr2 != nil {
		return nil, __Verr2
	}
	return &_STlsServer{
		Cipher:     _FnewCipher(__VbsPassword2),
		ListenAddr: __VstructListenAddr2,
	}, nil

}

// 运行服务端并且监听来自本地代理客户端的请求
func (___VlsServer1 *_STlsServer) _Fserver_Listen(___VsrvDidListen func(__VlistenAddr2 net.Addr)) error {
	return _FlistenSecureTCP(___VlsServer1.ListenAddr, ___VlsServer1.Cipher,
		___VlsServer1._FsrvHandleConn, ___VsrvDidListen)
}

// 解 SOCKS5 协议
// https://www.ietf.org/rfc/rfc1928.txt
func (___VlsServer2 *_STlsServer) _FsrvHandleConn(___VlocalConn *_STsecureTCPConn) {
	defer ___VlocalConn.Close()
	__Vbuf7 := make([]byte, 256)

	/**
	   The ___VlocalConn connects to the __VdstServer, and sends a ver
	   identifier/method selection message:
		          +----+----------+----------+
		          |VER | NMETHODS | METHODS  |
		          +----+----------+----------+
		          | 1  |    1     | 1 to 255 |
		          +----+----------+----------+
	   The VER field is set to X'05' for this ver of the protocol.  The
	   NMETHODS field contains the number of method identifier octets that
	   appear in the METHODS field.
	*/
	// 第一个字段VER代表Socks的版本，Socks5默认为0x05，其固定长度为1个字节
	_, __Verr3 := ___VlocalConn.DecodeRead(__Vbuf7)
	// 只支持版本5
	if __Verr3 != nil || __Vbuf7[0] != 0x05 {
		return
	}

	/**
	   The __VdstServer selects from one of the methods given in METHODS, and
	   sends a METHOD selection message:

		          +----+--------+
		          |VER | METHOD |
		          +----+--------+
		          | 1  |   1    |
		          +----+--------+
	*/
	// 不需要验证，直接验证通过
	___VlocalConn.EncodeWrite([]byte{0x05, 0x00})

	/**
	  +----+-----+-------+------+----------+----------+
	  |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
	  +----+-----+-------+------+----------+----------+
	  | 1  |  1  | X'00' |  1   | Variable |    2     |
	  +----+-----+-------+------+----------+----------+
	*/

	// 获取真正的远程服务的地址
	__Vn4, __Verr4 := ___VlocalConn.DecodeRead(__Vbuf7)
	// __Vn4 最短的长度为7 情况为 ATYP=3 DST.ADDR占用1字节 值为0x0
	if __Verr4 != nil || __Vn4 < 7 {
		return
	}

	// CMD代表客户端请求的类型，值长度也是1个字节，有三种类型
	// CONNECT X'01'
	if __Vbuf7[1] != 0x01 {
		// 目前只支持 CONNECT
		return
	}

	var __VdIP []byte
	// aType 代表请求的远程服务器地址类型，值长度1个字节，有三种类型
	switch __Vbuf7[3] {
	case 0x01:
		//	IP V4 address: X'01'
		__VdIP = __Vbuf7[4 : 4+net.IPv4len]
	case 0x03:
		//	DOMAINNAME: X'03'
		ipAddr, __Verr5 := net.ResolveIPAddr("ip", string(__Vbuf7[5:__Vn4-2]))
		if __Verr5 != nil {
			return
		}
		__VdIP = ipAddr.IP
	case 0x04:
		//	IP V6 address: X'04'
		__VdIP = __Vbuf7[4 : 4+net.IPv6len]
	default:
		return
	}
	__VdPort4 := __Vbuf7[__Vn4-2:]
	__VdstAddr4 := &net.TCPAddr{
		IP:   __VdIP,
		Port: int(binary.BigEndian.Uint16(__VdPort4)),
	}

	// 连接真正的远程服务
	__VdstServer, __Verr6 := net.DialTCP("tcp", nil, __VdstAddr4)
	if __Verr6 != nil {
		return
	} else {
		defer __VdstServer.Close()
		// Conn被关闭时直接清除所有数据 不管没有发送的数据
		__VdstServer.SetLinger(0)

		// 响应客户端连接成功
		/**
		  +----+-----+-------+------+----------+----------+
		  |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
		  +----+-----+-------+------+----------+----------+
		  | 1  |  1  | X'00' |  1   | Variable |    2     |
		  +----+-----+-------+------+----------+----------+
		*/
		// 响应客户端连接成功
		___VlocalConn.EncodeWrite([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	}

	// 进行转发
	// 从 localUser 读取数据发送到 __VdstServer
	go func() {
		__Verr7 := ___VlocalConn.DecodeCopy(__VdstServer)
		if __Verr7 != nil {
			// 在 copy 的过程中可能会存在网络超时等 error 被 return，只要有一个发生了错误就退出本次工作
			___VlocalConn.Close()
			__VdstServer.Close()
		}
	}()
	// 从 __VdstServer 读取数据发送到 localUser，这里因为处在翻墙阶段出现网络错误的概率更大
	(&_STsecureTCPConn{
		Cipher:          ___VlocalConn.Cipher,
		ReadWriteCloser: __VdstServer,
	}).EncodeCopy(___VlocalConn)
} // _STlsServer . _FsrvHandleConn
