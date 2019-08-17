package main

import (
	"fmt"
	//"github.com/gwuhaolin/lightsocks"
	//"github.com/gwuhaolin/lightsocks/cmd"
	"log"
	"net"
)

const (
	_Vserver_DefaultListenAddr = ":7448"
)

var _Vserver_version = "master"

func main() {
	log.SetFlags(log.Lshortfile)

	// 默认配置
	__Vclient_Config := &_ST_Config{
		ListenAddr: _Vserver_DefaultListenAddr,
	}
	__Vclient_Config._Fcommon_ReadConfig()
	__Vclient_Config._Fcommon_SaveConfig()

	// 启动 local 端并监听
	lsLocal, err := NewLsLocal(__Vclient_Config.Password, __Vclient_Config.ListenAddr, __Vclient_Config.RemoteAddr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(lsLocal._Fclient_Listen(func(___VclientListenAddr net.Addr) {
		log.Println("使用配置：", fmt.Sprintf(`
本地监听地址 listen：
%s
远程服务地址 remote：
%s
密码 password：
%s
	`, ___VclientListenAddr, __Vclient_Config.RemoteAddr, __Vclient_Config.Password))
		log.Printf("lightsocks-local:%s 启动成功 监听在 %s\n", _Vserver_version, ___VclientListenAddr.String())
	}))
}
