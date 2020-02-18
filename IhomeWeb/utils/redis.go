package utils

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/garyburd/redigo/redis"
	_ "github.com/gomodule/redigo/redis"
)

func Redis(serverName, addr, port, dbNum string) (bm cache.Cache, err error) {

	redis_conf := map[string]string{
		"key":   serverName,
		"conn":  addr + ":" + port,
		"dbNum": dbNum,
	}
	beego.Info(redis_conf)
	//将map转换为json
	redis_conf_json, _ := json.Marshal(redis_conf)
	// 创建redis句柄
	bm, err = cache.NewCache("redis", string(redis_conf_json))
	return

}
