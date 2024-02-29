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

package rules

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/arthur/banbrutes/config"
)

// 解析日志
func parse(text string, rule string) (err error, ip, name string, port int) {
	// 从frp日志中获取tcp连接信息
	// 2024/01/24 20:37:51 [I] [proxy.go:204] [de369b802e44e3f9] [S0-SSH] get a user connection [185.226.106.34:40432]
	if !strings.Contains(text, "get a user connection") {
		err = fmt.Errorf("not tcp link")
		return
	}

	// 正则表达式获取转发名和请求ID
	compileRegex := regexp.MustCompile(rule)
	matchArr := compileRegex.FindStringSubmatch(text)

	if len(matchArr) <= 2 {
		err = fmt.Errorf("not tcp link")
		return
	}

	// 转发名
	name = matchArr[1]
	addr := matchArr[2]
	addrArray := strings.Split(addr, ":")
	if len(addrArray) != 2 {
		err = fmt.Errorf(addr + " addr error")
		return
	}

	// 请求IP
	ip = addrArray[0]

	if v, ok := config.Cfg.NamePort[name]; ok {
		port = v
	} else {
		port = -1
	}

	return
}

// 解析日志
func parseFilter(text string, rules []config.RegFilter) (err error, ip, name string, port int) {

	name = ""
	for _, rule := range rules {
		// fmt.Println(rule)
		re := regexp.MustCompile(rule.Expression)
		result := re.FindStringSubmatch(text)
		offset := []int{1, 2}
		port = -1
		if rule.Offset != nil {
			offset = rule.Offset
		}
		if len(result) >= len(offset) {
			ip = result[offset[0]]

			if len(offset) > 1 {
				port, err = strconv.Atoi(result[offset[1]])
			}
			// fmt.Printf("IP地址: %s, 端口号: %d\n", ip, port)
			return
		}
	}
	err = fmt.Errorf("not valid line")
	return
}
