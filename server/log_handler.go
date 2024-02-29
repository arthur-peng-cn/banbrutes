package server

import (
	"path/filepath"
	"time"

	"github.com/arthur/banbrutes/config"
	"github.com/arthur/banbrutes/rules"
	"github.com/hpcloud/tail"
	"github.com/zngw/log"
)

func logMonitorServer(filename string, rule interface{}) {
	frpLog, _ := filepath.Abs(filename)
	tails, err := tail.TailFile(frpLog, tail.Config{
		ReOpen:    true,                                 // 重新打开
		Follow:    true,                                 // 是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 从文件的哪个地方开始读
		MustExist: false,                                // 文件不存在不报错
		Poll:      true,
	})

	if err != nil {
		log.Error("sys", "tail file failed, err:%v", err)
		return
	}

	log.Trace("sys", "banbrutes 已启动，正在监听日志文件：%s", frpLog)
	var line *tail.Line
	var ok bool
	go func() {
		for {
			line, ok = <-tails.Lines
			if !ok {
				log.Error("sys", "tail file close reopen, filename:%s\n", tails.Filename)
				time.Sleep(time.Second)
				continue
			}
			// rule := "^* \\[I] \\[.*] \\[.*] \\[(.*?)] get a user connection \\[(.*?)]"
			switch rule := rule.(type) {
			case string:
				rules.Check(line.Text, rule)
			case []config.RegFilter:
				rules.CheckFilter(line.Text, rule)
			}
		}
	}()
}
