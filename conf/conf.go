package conf

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

var conf_txt = `host : 0.0.0.0
port : 8088
user : admin
password : admin
debug : False
rootPath : .
auth : False
randomRoute : True`

type Conf struct {
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	Debug       bool   `yaml:"debug"`
	RootPath    string `yaml:"rootPath"`
	Auth        bool   `yaml:"auth"`
	RandomRoute bool   `yaml:"randomRoute"`
}

func Init_config() Conf {
	cfgFile := "config.yml"
	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		fmt.Printf("文件%s不存在\n", cfgFile)
		fmt.Println("正在初始化配置文件")
		fmt.Println("请修改配置文件后重新启动")
		err := os.WriteFile(cfgFile, []byte(conf_txt), 0644)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Scanln()
		os.Exit(0)
	}
	return load(cfgFile)

}

func load(cfgFile string) Conf {
	fmt.Printf("正在读取%s文件\n", cfgFile)
	dataBytes, err := os.ReadFile(cfgFile)
	if err != nil {
		fmt.Println("读取文件失败：", err)
		os.Exit(1)
	}
	config := Conf{}
	err = yaml.Unmarshal(dataBytes, &config)
	if !strings.HasSuffix(config.RootPath, "/") {
		config.RootPath += "/"
	}
	if err != nil {
		fmt.Println("解析 yaml 文件失败：", err)
		os.Exit(1)
	}
	return config
}
