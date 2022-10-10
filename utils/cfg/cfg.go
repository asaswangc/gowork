package cfg

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/asaswangc/gowork/result"
	"github.com/asaswangc/gowork/variable"
	"log"
	"os"
)

var T Toml

func Init() {
	log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate)
	// 根据启动模式获取配置文件
	var path string
	if variable.Global.Get(variable.ConfPath) == "" {
		path = result.Result(os.Getwd()).Unwrap().(string)
	} else {
		path = variable.Global.Get(variable.ConfPath).(string)
	}
	var file = fmt.Sprintf("cfg%s.toml", variable.Global.Get(variable.RunMode))
	if _, err := toml.DecodeFile(fmt.Sprintf("%s/%s", path, file), &T); err != nil {
		log.Fatal("加载配置文件失败")
	}
}
