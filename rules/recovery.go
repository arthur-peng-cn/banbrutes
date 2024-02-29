package rules

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/arthur/banbrutes/util"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zngw/log"
)

// 用线程安装的map保存ip记录
// 从sqlite数据库中读取没有recovered 的IP记录，启动1小时的清理，任务recovery，调用Release 函数释放recoverTime < time.Now()的记录
// release调用成功后需要更新数据库中的recovered状态
// 实现add(ip, port, banType)函数，添加一个新的ban item 记录到banlist，并插入数据库

var ban_list = sync.Map{}

type banItem struct {
	id          int
	ip          string
	port        int
	banType     string
	desc        string
	createTime  string
	recoverTime int64
	recovered   bool
}

type Recovery struct {
	db *sql.DB
}

func NewRecoverySrv() (srv *Recovery) {

	_db, _ := createDatabase()
	srv = &Recovery{
		db: _db,
	}
	return
}

func createDatabase() (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", "ban_list.db")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS ban_items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ip TEXT,
		port INTEGER,
		banType TEXT,
		desc TEXT,
		createTime TEXT,
		recoverTime INTEGER,
		recovered BOOLEAN
	)`)

	return
}

func (srv *Recovery) init() {

	if srv.db == nil {
		srv.db, _ = createDatabase()
	}

	rows, err := srv.db.Query("SELECT id, ip, port, banType, createTime, recoverTime, desc FROM ban_items WHERE recovered = ?", false)
	if err != nil {
		// 处理错误
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id, port int
		var recoverTime int64
		var ip, banType, desc, createTime string
		err = rows.Scan(&id, &ip, &port, &banType, &createTime, &recoverTime, &desc)
		if err != nil {
			// 处理错误
			continue
		}

		banItem := &banItem{
			ip:          ip,
			port:        port,
			banType:     banType,
			desc:        desc,
			createTime:  createTime,
			recoverTime: recoverTime,
			recovered:   false,
		}
		ban_list.Store(ip, banItem)
	}
}

func (srv *Recovery) Add(ip string, port int, banType string, desc string) {
	// 在ban_list中添加新的ban item记录
	banItem := &banItem{
		ip:          ip,
		port:        port,
		banType:     banType,
		desc:        desc,
		createTime:  time.Now().Format("2006-01-02 15:04:05"),
		recoverTime: time.Now().Unix() + 3600*24*10, // 10后尝试恢复
		recovered:   false,
	}
	ban_list.Store(ip, banItem)

	// 插入数据库
	srv.insert(banItem)
}

// InsertBanItemToDB 向SQLite数据库中插入一条记录
func (srv *Recovery) insert(it *banItem) {

	// 准备插入数据的SQL语句
	stmt, err := srv.db.Prepare("INSERT INTO ban_items(ip, port, banType, createTime, recoverTime, desc, recovered) VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Error("link", err.Error())
	}
	defer stmt.Close()

	// 执行插入操作
	_, err = stmt.Exec(it.ip, it.port, it.banType, it.createTime, it.recoverTime, it.desc, it.recovered)
	if err != nil {
		log.Error("link", err.Error())
	}

	log.Trace("link", "插入记录成功 %s:%d", it.ip, it.port)
}

func (srv *Recovery) UpdateRecoveredStatus(id int) {
	// 更新数据库中id对应的ban item的recovered状态
	stmt, err := srv.db.Prepare("UPDATE ban_items SET recovered = ? WHERE id = ?")
	if err != nil {
		// 处理错误
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(true, id)
	if err != nil {
		// 处理错误
		return
	}
}

func (srv *Recovery) Run() {
	// 异步清理
	go func() {
		for {
			delTime := time.Now().Unix() - 3600
			ban_list.Range(func(key, value interface{}) bool {
				// ip := key.(string)
				item := value.(*banItem)
				if !item.recovered && item.recoverTime < delTime {
					if Release(item.ip, item.port, item.banType) {
						item.recovered = true
						// 更新数据库中的recovered状态
						srv.UpdateRecoveredStatus(item.id)
					}
				}
				return true
			})

			// 休眠1小时
			time.Sleep(time.Hour)
		}
	}()
}

func Release(ip string, port int, banType string) (res bool) {
	cmd := ""
	switch banType {
	case "iptables":
		if port == -1 {
			cmd = fmt.Sprintf("iptables -D INPUT -s %s -j DROP", ip)
		} else {
			cmd = fmt.Sprintf("iptables -D INPUT -s %s -ptcp --dport %d -j DROP", ip, port)
		}
	case "firewall":
		if port == -1 {
			cmd = fmt.Sprintf("firewall-cmd --permanent --remove-rich-rule=\"rule family=\"ipv4\" source address=\"%s\" reject\"", ip)
		} else {
			cmd = fmt.Sprintf("firewall-cmd --permanent --remove-rich-rule=\"rule family=\"ipv4\" source address=\"%s\" port protocol=\"tcp\" port=\"%d\" reject\"", ip, port)
		}
		cmd += " && firewall-cmd --reload"
	case "md":
		if port == -1 {
			cmd = fmt.Sprintf("netsh advfirewall firewall delete rule name=brute-ban-%s", ip)
		} else {
			cmd = fmt.Sprintf("netsh advfirewall firewall delete rule name=brute-ban-%s-%d", ip, port)
		}
	}

	if cmd != "" {
		log.Info("cmd", cmd)
		result := util.Command(cmd)
		if result != "" {
			log.Trace("sys", result)
		}
	}
	res = true
	return
}
