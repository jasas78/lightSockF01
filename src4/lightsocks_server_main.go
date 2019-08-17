package main

import (
	"fmt"
	//"github.com/phayes/freeport"
	"log"
	"net"
) // import

var _Vsrv_Version = "master"

func main() {
	_Finit1()
	_Finit2(".server")
	log.SetFlags(log.Lshortfile)

	// 服务端监听端口: 采用 7448
	__Vsrv_port := 7448

	// 默认配置
	__Vsrv_config := &_ST_Config{
		ListenAddr: fmt.Sprintf(":%d", __Vsrv_port),
		// 密码随机生成
		Password: _FrandPassword(),
	}
	__Vsrv_config._Fcommon_ReadConfig()
	__Vsrv_config._Fcommon_SaveConfig()

	// 启动 server 端并监听
	__Vsrv_lsServer, __VsrvErr2 :=
		_Fserver_NewLsServer(__Vsrv_config.Password, __Vsrv_config.ListenAddr)
	if __VsrvErr2 != nil {
		log.Fatalln(__VsrvErr2)
	}
	log.Fatalln(__Vsrv_lsServer._Fserver_Listen(func(___Vsrv_listenAddr net.Addr) {
		log.Println("使用配置：", fmt.Sprintf(`
本地监听地址 listen：
%s
密码 password：
%s
	`, ___Vsrv_listenAddr, __Vsrv_config.Password))
		log.Printf("lightsocks-server:%s 启动成功 监听在 %s\n", _Vsrv_Version, ___Vsrv_listenAddr.String())
	}))
} // main
