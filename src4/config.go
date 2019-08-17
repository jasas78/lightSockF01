//package cmd
package main

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"log"
	"os"
	"path"
)

var (
	// 配置文件路径
	_VconfigPath string
)

type _ST_Config struct {
	ListenAddr string `json:"listen"`
	RemoteAddr string `json:"remote"`
	Password   string `json:"password"`
}

func _Finit2(___VfileNameExt string) {
	__Vhome1, _ := homedir.Dir()
	// 默认的配置文件名称
	__VconfigFilename1 := ".lightsocks.json"
	// 如果用户有传配置文件，就使用用户传入的配置文件
	if len(os.Args) == 2 {
		__VconfigFilename1 = os.Args[1]
	}
	_VconfigPath = path.Join(__Vhome1, __VconfigFilename1+___VfileNameExt)
}

// 保存配置到配置文件
func (___Vconfig1 *_ST_Config) _Fcommon_SaveConfig() {
	__VconfigJson, _ := json.MarshalIndent(___Vconfig1, "", "	")
	__Verr1 := ioutil.WriteFile(_VconfigPath, __VconfigJson, 0644)
	if __Verr1 != nil {
		fmt.Errorf("保存配置到文件 %s 出错: %s", _VconfigPath, __Verr1)
	}
	log.Printf("保存配置到文件 %s 成功\n", _VconfigPath)
}

func (___Vconfig2 *_ST_Config) _Fcommon_ReadConfig() {
	// 如果配置文件存在，就读取配置文件中的配置 assign 到 ___Vconfig2
	if _, __Verr2 := os.Stat(_VconfigPath); !os.IsNotExist(__Verr2) {
		log.Printf("从文件 %s 中读取配置\n", _VconfigPath)
		__Vfile1, __Verr3 := os.Open(_VconfigPath)
		if __Verr3 != nil {
			log.Fatalf("打开配置文件 %s 出错:%s", _VconfigPath, __Verr3)
		}
		defer __Vfile1.Close()

		__Verr3 = json.NewDecoder(__Vfile1).Decode(___Vconfig2)
		if __Verr3 != nil {
			log.Fatalf("格式不合法的 JSON 配置文件:\n%s", __Vfile1.Name())
		}
	}
}
