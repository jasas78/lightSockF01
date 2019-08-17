package main

import (
	"fmt"
	//"github.com/gwuhaolin/lightsocks"
	//"github.com/gwuhaolin/lightsocks/cmd"
	"log"
	"net"
)

const (
	_Vclient_DefaultListenAddr = ":7448"
)

var _Vclient_version = "master"

func main() { // main clinet
    _Finit()
	log.SetFlags(log.Lshortfile)

	// 默认配置
	__Vclient_Config := &_ST_Config{
		ListenAddr: _Vclient_DefaultListenAddr,
	}
	__Vclient_Config._Fcommon_ReadConfig()
	__Vclient_Config._Fcommon_SaveConfig()

	// 启动 local 端并监听
	__Vclient_lsLocal, __VclientErr := NewLsLocal(__Vclient_Config.Password,
		__Vclient_Config.ListenAddr, __Vclient_Config.RemoteAddr)
	if __VclientErr != nil {
		log.Fatalln(__VclientErr)
	}
	log.Fatalln(__Vclient_lsLocal._Fclient_Listen(func(___VclientListenAddr net.Addr) {
		log.Println("使用配置：", fmt.Sprintf(`
本地监听地址 listen：
%s
远程服务地址 remote：
%s
密码 password：
%s
	`, ___VclientListenAddr, __Vclient_Config.RemoteAddr, __Vclient_Config.Password))
		log.Printf("lightsocks-local:%s 启动成功 监听在 %s\n", _Vclient_version, ___VclientListenAddr.String())
	}))
} // main clinet
