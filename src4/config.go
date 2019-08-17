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

func init() {
	home, _ := homedir.Dir()
	// 默认的配置文件名称
	configFilename := ".lightsocks.json"
	// 如果用户有传配置文件，就使用用户传入的配置文件
	if len(os.Args) == 2 {
		configFilename = os.Args[1]
	}
	_VconfigPath = path.Join(home, configFilename)
}

// 保存配置到配置文件
func (config *_ST_Config) _Fcommon_SaveConfig() {
	__VconfigJson, _ := json.MarshalIndent(config, "", "	")
	__Verr1 := ioutil.WriteFile(_VconfigPath, __VconfigJson, 0644)
	if __Verr1 != nil {
		fmt.Errorf("保存配置到文件 %s 出错: %s", _VconfigPath, __Verr1)
	}
	log.Printf("保存配置到文件 %s 成功\n", _VconfigPath)
}

func (config *_ST_Config) _Fcommon_ReadConfig() {
	// 如果配置文件存在，就读取配置文件中的配置 assign 到 config
	if _, __Verr2 := os.Stat(_VconfigPath); !os.IsNotExist(__Verr2) {
		log.Printf("从文件 %s 中读取配置\n", _VconfigPath)
		file, __Verr3 := os.Open(_VconfigPath)
		if __Verr3 != nil {
			log.Fatalf("打开配置文件 %s 出错:%s", _VconfigPath, __Verr3)
		}
		defer file.Close()

		__Verr3 = json.NewDecoder(file).Decode(config)
		if __Verr3 != nil {
			log.Fatalf("格式不合法的 JSON 配置文件:\n%s", file.Name())
		}
	}
}
