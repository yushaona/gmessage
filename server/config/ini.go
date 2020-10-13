/*
	读取配置文件
*/

package config

import (
	"fmt"
	"os"

	"gopkg.in/ini.v1"
)

var (
	TCP_HOST string
)

func init() {
	cfg, err := ini.Load("./conf/server.ini")
	if err != nil {
		fmt.Printf("Fail to read server.ini: %v", err)
		os.Exit(1)
	}
	section := cfg.Section("server")
	TCP_HOST = section.Key("TCP_HOST").String()
}
