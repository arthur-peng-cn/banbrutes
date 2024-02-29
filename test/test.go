package main

import (
	"fmt"
	"regexp"

	"github.com/arthur/banbrutes/config"
)

func main1() {
	text := "[I] [proxy.go:204] [65e59a1bb269f263] [家里电脑] get a user connection [113.215.189.96:6866]"

	// 匹配ip的正则表达式
	re := regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`)

	// 查找ip地址
	ip := re.FindString(text)

	fmt.Println(ip)
}

func main() {
	err := config.Init("./config.yml")
	if err != nil {
		fmt.Println("load config failed %v", err)
	}

	fmt.Printf("(*%s*)\n", config.Cfg.ListenAddr)
	str := "[I] [proxy.go:204] [65e59a1bb269f263] [家里电脑] get a user connection [113.215.189.96:6866]"
	for _, filter := range config.Cfg.Filters {
		for _, reg := range filter.RegFilters {
			fmt.Println(reg)
			re := regexp.MustCompile(reg.Expression)
			result := re.FindStringSubmatch(str)
			offset := []int{1, 2}
			port := ""
			if reg.Offset != nil {
				offset = reg.Offset
			}
			if len(result) >= len(offset) {
				ip := result[offset[0]]

				if len(offset) > 1 {
					port = result[offset[1]]
				}
				fmt.Printf("IP地址: %s, 端口号: %s\n", ip, port)
			}
		}
	}

}

func main2() {
	text := "Failed password for root from 218.92.0.55 port 55778 ssh2"

	re := regexp.MustCompile(`Failed password.* (\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})\s+port\s+(\d+)`)
	match := re.FindStringSubmatch(text)

	if len(match) > 0 {
		ipAddress := match[1]
		port := match[2]

		fmt.Println("IP地址:", ipAddress)
		fmt.Println("端口号:", port)
	}
}
