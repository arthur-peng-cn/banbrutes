//MIT License
//
//Copyright (c) 2023 arthur
//
//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//copies of the Software, and to permit persons to whom the Software is
//furnished to do so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in all
//copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//SOFTWARE.

package main

import (
	"encoding/json"
	"flag"

	// "fmt"
	// "os"
	"path/filepath"

	"github.com/arthur/banbrutes/config"
	"github.com/arthur/banbrutes/rules"
	"github.com/arthur/banbrutes/server"
	"github.com/zngw/log"
	"github.com/zngw/zipinfo/ipinfo"
)

func main() {
	// 读取命令行配置文件参数
	c := flag.String("c", "./config.yml", "默认配置为 config.yml")
	s := flag.String("s", "", "默认配置为空")
	flag.Parse()

	if *s == "reload" {
		// 如果是reload，发送reload指令后退出
		config.SendReload()
		return
	}

	// 初始化配置
	err := config.Init(*c)
	if err != nil {
		panic(err)
	}

	// 初始化日志
	// _, file := filepath.Split(os.Args[0])
	// logFile, _ := filepath.Abs(config.Cfg.Logs)
	logFile := filepath.Join(config.Cfg.Logs, "log.txt")

	log.InitLog("all", logFile, "info", 7, true, []string{"add", "link", "net", "sys", "cmd", "init"})

	var ipCfg []interface{}
	err = json.Unmarshal([]byte(config.Cfg.IpInfo), &ipCfg)
	if err == nil {
		ipinfo.Init(ipCfg)
	} else {
		log.Error("init", "ipinfo init failed %s", ipCfg)
	}

	// 初始化规则
	rules.Init()

	// 启动用tail监听
	frpLog, _ := filepath.Abs(config.Cfg.FrpsLog)

	log.Debug("init", "启动日志监听-> %s", frpLog)
	srv, err := server.NewService(config.Cfg)
	if err != nil {
		log.Error("init", "create service failed \n")

	}
	srv.Run()
}
