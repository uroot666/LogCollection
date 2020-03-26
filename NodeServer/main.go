package main

import (
	"LogCollection/EtcdReadWatch"
	"LogCollection/KafkaSend"
	"LogCollection/TailFile"
	"LogCollection/conf"
	"LogCollection/register"
	"fmt"
	"gopkg.in/ini.v1"
	"sync"
)

var (
	kafkaConf conf.KafkaConf
	etcdConf  conf.EtcdConf
	wg        sync.WaitGroup
)

func main() {
	// 加载配置文件
	cfg, err := ini.Load("./conf/config.ini")
	if err != nil {
		fmt.Printf("加载配置文件失败, err:%#v \n", err)
		return
	}
	err = cfg.Section("kafka").MapTo(&kafkaConf)
	err = cfg.Section("etcd").MapTo(&etcdConf)

	// 初始化注册表文件，用于记录文件的偏移量
	register.Init("./register.json")

	// 初始化kafka连接
	KafkaSend.Init([]string{kafkaConf.AddrPort})

	// 初始化etcd连接
	err = EtcdReadWatch.Init(etcdConf.AddrPort)
	if err != nil {
		fmt.Println("初始化etcd失败")
	}

	filePathConf, _ := EtcdReadWatch.GetFilePath(etcdConf.Key)
	TailFile.Init(filePathConf)
	NewConfCh := TailFile.GetNewConfCh()
	wg.Add(1)
	EtcdReadWatch.WatchChang(etcdConf.Key, NewConfCh)
	wg.Wait()

}
