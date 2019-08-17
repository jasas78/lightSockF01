//package lightsocks
package main

import (
	"log"
	"net"
)

type _STlsLocal struct {
	Cipher     *_STcipher
	ListenAddr *net.TCPAddr
	RemoteAddr *net.TCPAddr
}

// 新建一个本地端
// 本地端的职责是:
// 1. 监听来自本机浏览器的代理请求
// 2. 转发前加密数据
// 3. 转发socket数据到墙外代理服务端
// 4. 把服务端返回的数据转发给用户的浏览器
func NewLsLocal(___Vpassword string, ___VlistenAddr, ___VremoteAddr string) (*_STlsLocal, error) {
	__VbsPassword1, __Verr3 := _FparsePassword(___Vpassword)
	if __Verr3 != nil {
		return nil, __Verr3
	}
	__VstructListenAddr1, __Verr4 := net.ResolveTCPAddr("tcp", ___VlistenAddr)
	if __Verr4 != nil {
		return nil, __Verr4
	}
	__VstructRemoteAddr1, __Verr5 := net.ResolveTCPAddr("tcp", ___VremoteAddr)
	if __Verr5 != nil {
		return nil, __Verr5
	}
	return &_STlsLocal{
		Cipher:     _FnewCipher(__VbsPassword1),
		ListenAddr: __VstructListenAddr1,
		RemoteAddr: __VstructRemoteAddr1,
	}, nil
}

// 本地端启动监听，接收来自本机浏览器的连接
func (__Vlocal1 *_STlsLocal) _Fclient_Listen(___VdidListen func(___VlistenAddr net.Addr)) error {
	return _FlistenSecureTCP(__Vlocal1.ListenAddr, __Vlocal1.Cipher, __Vlocal1._FhandleConn, ___VdidListen)
}

func (__Vlocal2 *_STlsLocal) _FhandleConn(___VuserConn *_STsecureTCPConn) {
	defer ___VuserConn.Close()

	__VproxyServer, __Verr1 := DialTCPSecure(__Vlocal2.RemoteAddr, __Vlocal2.Cipher)
	if __Verr1 != nil {
		log.Println(__Verr1)
		return
	}
	defer __VproxyServer.Close()

	// Conn被关闭时直接清除所有数据 不管没有发送的数据
	//__VproxyServer.SetLinger(0)

	// 进行转发
	// 从 __VproxyServer 读取数据发送到 localUser
	go func() {
		__Verr2 := __VproxyServer.DecodeCopy(___VuserConn)
		if __Verr2 != nil {
			// 在 copy 的过程中可能会存在网络超时等 error 被 return，只要有一个发生了错误就退出本次工作
			___VuserConn.Close()
			__VproxyServer.Close()
		}
	}()
	// 从 localUser 发送数据发送到 __VproxyServer，这里因为处在翻墙阶段出现网络错误的概率更大
	___VuserConn.EncodeCopy(__VproxyServer)
}
