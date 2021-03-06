package utils

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
)

var (
	G_server_name  string // 项目名称
	G_server_addr  string // 服务器IP地址
	G_server_port  string // 服务器端口
	G_redis_addr   string // redis ip地址
	G_redis_port   string // redis port
	G_redis_dbnum  string // redis db 编号
	G_mysql_addr   string // mysql ip地址
	G_mysql_port   string // mysql 端口
	G_mysql_dbname string // mysql db name
	G_fastdfs_port string // fastdfs 端口
	G_fastdfs_addr string // fastdfs ip
)

func InitConfig() {
	appconf, err := config.NewConfig("ini", "./conf/app.conf")
	if err != nil {
		beego.Debug(err)
		return
	}

	G_server_name = appconf.String("appname")
	G_server_addr = appconf.String("httpaddr")
	G_server_port = appconf.String("httpport")
	G_redis_addr = appconf.String("redisaddr")
	G_redis_port = appconf.String("redisport")
	G_redis_dbnum = appconf.String("redisdbnum")
	G_mysql_addr = appconf.String("mysqladdr")
	G_mysql_port = appconf.String("mysqlport")
	G_mysql_dbname = appconf.String("mysqldbbane")
	G_fastdfs_addr = appconf.String("fastdfsaddr")
	G_fastdfs_port = appconf.String("fastdfsport")
	return
}

func init() {
	InitConfig()
}
