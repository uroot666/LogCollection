package main

import (
	"LogCollectionWeb/dao"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 定义全局的web对象
var r *gin.Engine

func Init() {
	r = gin.Default()
	r.Static("/static/", "./static")
	r.LoadHTMLGlob("templates/*")
	// 用于检查服务是否可访问
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// 首页
	r.GET("/", func(c *gin.Context) {
		all, _ := dao.GetAllKey("/logagent/")
		var keymap = make(map[string]int)
		for k, v := range all{
			keymap[k] = len(v)
		}
		c.HTML(http.StatusOK, "index.html", gin.H{"title": "LogCollectionWeb管理", "keymap": keymap})
	})

	// 获取所有key的信息
	r.GET("/value/allKey", GetAllKey)
	// 操作单个key(get put del)
	r.GET("/value/key", GetKey)
	r.POST("/value/key", PutKey)
	r.PUT("/value/key", PutKey)
	r.DELETE("/value/key", DeleteKey)
}

func main() {
	// 初始化web连接
	Init()
	// 初始化etcd连接
	dao.Init("192.168.1.111:2380")
	// 启动web对象
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}

func GetAllKey(c *gin.Context) {
	key := c.DefaultQuery("key", "/err")
	all, _ := dao.GetAllKey(key)
	c.JSON(http.StatusOK, all)
}

func GetKey(c *gin.Context) {
	key := c.DefaultQuery("key", "/err")
	t, _ := dao.GetKey(key)
	c.JSON(http.StatusOK, t)
}

func PutKey(c *gin.Context) {
	//kv := &dao.KeyValue{
	//	Key: "/logagent/192.168.1.105/collect_config",
	//	Value: []dao.PathTopic{
	//		{Path: "c:/tmp/test1/nginx.log", Topic: "web_log"},
	//		{Path: "c:/tmp/test2/redis.log", Topic: "redis_log"},
	//		{Path: "c:/tmp/test3/mysql.log", Topic: "mysql_log"},
	//	},
	//}
	var kv dao.KeyValue
	if err := c.ShouldBindJSON(&kv); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_ = dao.PutKey(&kv)
	fmt.Println("123", kv)
	c.JSON(http.StatusOK, "put ok")
}

func DeleteKey(c *gin.Context) {
	key := c.DefaultQuery("key", "/err")
	dao.DeletKey(key)
	c.JSON(http.StatusOK, "get ok")
}
