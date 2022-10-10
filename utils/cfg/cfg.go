package cfg

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/asaswangc/gowork/result"
	"log"
	"os"
)

func Init(confPath string, runMode string, cfgStruct interface{}) {
	log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate)
	var path string
	if confPath == "" {
		path = result.Result(os.Getwd()).Unwrap().(string)
	} else {
		path = confPath
	}
	var file = fmt.Sprintf("cfg%s.toml", runMode)
	if _, err := toml.DecodeFile(fmt.Sprintf("%s/%s", path, file), cfgStruct); err != nil {
		log.Fatal("加载配置文件失败")
	}
}
