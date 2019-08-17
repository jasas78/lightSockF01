package main

import (
	"fmt"
	"github.com/phayes/freeport"
	"log"
	"net"
) // import

var version = "master"

func main() {
	log.SetFlags(log.Lshortfile)

	// 服务端监听端口随机生成
	port, err := freeport.GetFreePort()
	if err != nil {
		// 随机端口失败就采用 7448
		port = 7448
	}
	// 默认配置
	__Vsrv_config := &_ST_Config{
		ListenAddr: fmt.Sprintf(":%d", port),
		// 密码随机生成
		//Password: lightsocks.RandPassword(),
		Password: RandPassword(),
	}
	__Vsrv_config._Fcommon_ReadConfig()
	__Vsrv_config._Fcommon_SaveConfig()

	// 启动 server 端并监听
	lsServer, err := _Fserver_NewLsServer(__Vsrv_config.Password, __Vsrv_config.ListenAddr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(lsServer._Fserver_Listen(func(listenAddr net.Addr) {
		log.Println("使用配置：", fmt.Sprintf(`
本地监听地址 listen：
%s
密码 password：
%s
	`, listenAddr, __Vsrv_config.Password))
		log.Printf("lightsocks-server:%s 启动成功 监听在 %s\n", version, listenAddr.String())
	}))
} // main
