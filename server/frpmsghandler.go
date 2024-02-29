package server

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/arthur/banbrutes/rules"
	"github.com/zngw/log"
)

var executor = make(chan struct{}, 4)
var mu sync.Mutex
var SSH_IP_MODE = "ok"

// timestamp_to_str 格式化时间戳
func timestampToStr(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.Format("2006-01-02 15:04:05")
}

// loginOperation 处理Login操作：frpc登录frps
func loginOperation(data map[string]interface{}) string {
	strFmt := "frp-client登录\nfrp版本：%v\n系统类型：%v\n系统架构：%v\n登录时间：%v\n连接池大小：%v"
	txt := fmt.Sprintf(strFmt, data["version"], data["os"], data["arch"], timestampToStr(int64(data["timestamp"].(float64))), data["pool_count"])
	return txt
}

// newProxyOperation 处理NewProxy操作：frpc与frps之间建立通道用于内网穿透
func newProxyOperation(data map[string]interface{}) string {
	runID := data["user"].(map[string]interface{})["run_id"]
	proxyType := data["proxy_type"].(string)
	txt := fmt.Sprintf("frp-client建立穿透代理\n主机ID：%v\n代理名称：%v\n代理类型：%v\n", runID, data["proxy_name"], proxyType)
	if proxyType == "tcp" || proxyType == "udp" {
		txt += fmt.Sprintf("远程端口：%v\n", data["remote_port"])
	} else if proxyType == "http" || proxyType == "https" {
		txt += fmt.Sprintf("子域名：%v\n", data["subdomain"])
	}
	return txt
}

// newUserConnOperation 处理NewUserConn操作：用户连接内网机器
func newUserConnOperation(data map[string]interface{}) (string, string, bool) {
	runID := data["user"].(map[string]interface{})["run_id"]
	ip := data["remote_addr"].(string)
	parts := strings.Split(ip, ":")
	ip = parts[0]
	isAllow := true
	if SSH_IP_MODE != "no" {
		refuse, desc, _, _ := rules.CheckRules(ip, -1)
		if refuse {
			log.Info("link", "用户(%s)链接拒绝原因: %v", ip, desc)
			rules.Refuse(ip, -1)
			isAllow = false
		}
	}

	// strFmt := "用户连接内网机器\n内网主机ID：%v\n代理名称：%v\n代理类型：%v\n登录时间：%v\n用户IP和端口：%v"
	// txt := fmt.Sprintf(strFmt, runID, data["proxy_name"], data["proxy_type"], timestampToStr(int64(data["timestamp"].(float64))), data["remote_addr"])
	txt := " + " + runID.(string)
	return txt, ip, isAllow
}

// newWorkConnOperation 处理NewWorkConn操作
func newWorkConnOperation(data map[string]interface{}) {
	// Do something here for NewWorkConn operation
}

// asyncSend 异步发送消息
func asyncSend(txt string, isAllowSSH bool, ip string) {
	// mu.Lock()
	// defer mu.Unlock()

	// // 用户地理位置
	// position := ""
	// if ip != "" {
	// 	position = ip2geo(ip)
	// 	txt += fmt.Sprintf("\n用户地理位置：%v", position)
	// }

	// // 是否允许用户连接
	// txt += "\n允许连接：是"
	// if !isAllowSSH {
	// 	txt += "否"
	// }

	// // 发送消息给接收者
	// for _, receiver := range RECEIVERS {
	// 	switch receiver {
	// 	case "dingtalk":
	// 		sendTextDingtalk(txt)
	// 	case "feishu":
	// 		sendTextFeishu(txt)
	// 	default:
	// 		// Handle other receivers
	// 	}
	// }
}

// handleMsg 处理各种frps信息
func handleMsg(data map[string]interface{}) bool {
	// 当前建立frp的类型
	operation := data["op"].(string)
	// frp请求的具体信息
	content := data["content"].(map[string]interface{})

	// 发送给管理员用户的提示
	txt := ""
	ip := ""
	// 是否允许用户ssh连接
	isAllowSSH := true

	switch operation {
	case "Ping":
		return true
	case "Login":
		txt = loginOperation(content)
	case "NewProxy":
		txt = newProxyOperation(content)
	case "NewUserConn":
		content["timestamp"] = int(time.Now().Unix())
		txt, ip, isAllowSSH = newUserConnOperation(content)
	case "NewWorkConn":
		// newWorkConnOperation(content)
		return true
	default:
		// Handle other cases
		return true
	}

	// 通知用户
	go asyncSend(txt, isAllowSSH, ip)
	return isAllowSSH
}
