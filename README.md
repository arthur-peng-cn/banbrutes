[comment]: <> (dtapps)
[![GitHub Org's stars](https://img.shields.io/github/stars/zngw)](https://github.com/zngw)

[comment]: <> (go)
[![godoc](https://pkg.go.dev/badge/github.com/arthur/banbrutes?status.svg)](https://pkg.go.dev/github.com/arthur/banbrutes)
[![oproxy.cn](https://goproxy.cn/stats/github.com/arthur/banbrutes/badges/download-count.svg)](https://goproxy.cn/stats/github.com/arthur/banbrutes)
[![goreportcard.com](https://goreportcard.com/badge/github.com/arthur/banbrutes)](https://goreportcard.com/report/github.com/arthur/banbrutes)
[![deps.dev](https://img.shields.io/badge/deps-go-red.svg)](https://deps.dev/go/github.com%2Fdtapps%2Fgo-ssh-tunnel)

[comment]: <> (github.com)
[![watchers](https://badgen.net/github/watchers/zngw/banbrutes)](https://github.com/arthur/banbrutes/watchers)
[![stars](https://badgen.net/github/stars/zngw/banbrutes)](https://github.com/arthur/banbrutes/stargazers)
[![forks](https://badgen.net/github/forks/zngw/banbrutes)](https://github.com/arthur/banbrutes/network/members)
[![issues](https://badgen.net/github/issues/zngw/banbrutes)](https://github.com/arthur/banbrutes/issues)
[![branches](https://badgen.net/github/branches/zngw/banbrutes)](https://github.com/arthur/banbrutes/branches)
[![releases](https://badgen.net/github/releases/zngw/banbrutes)](https://github.com/arthur/banbrutes/releases)
[![tags](https://badgen.net/github/tags/zngw/banbrutes)](https://github.com/arthur/banbrutes/tags)
[![license](https://badgen.net/github/license/zngw/banbrutes)](https://github.com/arthur/banbrutes/blob/master/LICENSE)
[![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/zngw/banbrutes)](https://github.com/arthur/banbrutes)
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/zngw/banbrutes)](https://github.com/arthur/banbrutes/releases)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/zngw/banbrutes)](https://github.com/arthur/banbrutes/tags)
[![GitHub pull requests](https://img.shields.io/github/issues-pr/zngw/banbrutes)](https://github.com/arthur/banbrutes/pulls)
[![GitHub issues](https://img.shields.io/github/issues/zngw/banbrutes)](https://github.com/arthur/banbrutes/issues)
[![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/zngw/banbrutes)](https://github.com/arthur/banbrutes)
[![GitHub language count](https://img.shields.io/github/languages/count/zngw/banbrutes)](https://github.com/arthur/banbrutes)
[![GitHub search hit counter](https://img.shields.io/github/search/zngw/banbrutes/go)](https://github.com/arthur/banbrutes)
[![GitHub top language](https://img.shields.io/github/languages/top/zngw/banbrutes)](https://github.com/arthur/banbrutes)

# banbrutes
* 监控日志文件，自定义规则，使用系统自带的防火墙(iptables、firewall、Microsoft Defender)拦截tcp连接的ip，防止暴力破解
* 支持同时监控多个日志文件，支持同一日志文件用多种过滤规则来过滤，支持IP，PORT随机出现的规则
* 支持各类通过IP区域以及访问频次进行规则过滤的场景, 支持通过正则表达式方法设置过滤规则
* 支持作为Frps的访问控制plugin的方式工作
* 通过数据库记录所有被拦截的ip，在设定的天数之后自动解除拦截，可以通过设置recovery规则设置，默认是10天后解除


## 前提
1、因为需要用到命令行修改系统防火墙，所以运行程序需要root或管理员权限


## Frps plugin工作方式
```
[plugin.frp-info]
addr = 127.0.0.1:8888
path = /handler
ops = NewUserConn
```


## 配置

```yaml
# 防止黑客暴力破解
# 输出日志目录
logs: ./log/

filters:
  - # log
    log_file: /var/log/auth.log
    reg_filters: 
      - # ssh 过滤规则
        expression: 'Failed password.* (\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})\s+port\s+(\d+)'
      - # log

# 启用防火墙类型 iptables / firewall / md (Microsoft Defender)
tables_type: iptables

# ip白名单:
allow_ip:
  - 127.0.0.1

# 端口白名单
allow_port:
  - 80
  - 443

# 规则访问
rules:

  # 按数组顺序来，匹配到了就按匹配的规则执行，跳过此规则。
  # 地区 country-国家， regionName-省名，名字中不带省字， city-市名，名字中也不带市字
  # 端口: -1 所有端口
  # time: 时间区间
  # count: 访问次数，-1不限，0限制。其他为 time时间内访问count次，超出频率就限制

  - # 中国上海IP允许
    port: -1
    country: 中国
    regionName: 上海
    city: 上海
    time: 1
    count: -1

  - # 中国地区IP 10分钟3次，超出这频率添加防火墙
    port: -1
    country: 中国
    regionName: 浙江
    city:
    time: 600
    count: 3

  - # 其他地区IP 直接加入防火墙
    port: -1
    country:
    regionName:
    city:
    time: 1
    count: 0
```

## 启动
直接使用`nohup ./banbrutes -c config.yml &`启动

也可以新建`/etc/systemd/system/banbrutes.service`文件加入系统,以服务方式启动
```ini
[Unit]
Description=frps daemon
After=syslog.target  network.target
Wants=network.target

[Service]
Type=simple
ExecStart=/usr/local/frp/banbrutes -c /usr/local/frp/config.yml
Restart= always
RestartSec=1min

[Install]
WantedBy=multi-user.target

```

