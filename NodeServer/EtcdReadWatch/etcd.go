package EtcdReadWatch

import (
	"LogCollection/Log"
	"context"
	"encoding/json"
	"go.etcd.io/etcd/clientv3"
	"go.uber.org/zap"
	"time"
)

type FilePath struct {
	Topic string `json:topic`
	Path  string `json:"path"`
}

// 给定一个全局的etcd连接对象
var cli *clientv3.Client

var LogObj *zap.SugaredLogger

// 初始化etcd连接
func Init(addr string) (err error) {
	LogObj, _ = Log.GetLogObj()
	if LogObj == nil {
		return
	}

	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{addr},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// fmt.Println("连接ETCD失败...")
		LogObj.Errorf("连接ETCD失败...")
		return
	}
	return
}

// 通过给定主机名称到etcd中获取对应的值，值表示有哪些日志需要收集，以及在kafka中对应的topic
func GetFilePath(key string) (filepath []FilePath, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	resp, err := cli.Get(ctx, key)
	if err != nil {
		// fmt.Println("从etcd获取value失败...")
		LogObj.Errorf("从etcd获取value失败...")
		return
	}
	cancel()
	for _, v := range resp.Kvs {
		// fmt.Println("etcd value", string(v.Value))
		LogObj.Debugf("etcd value: %v", string(v.Value))
		err = json.Unmarshal(v.Value, &filepath)
		if err != nil {
			// fmt.Println("etcd取出的结果反系列化失败...")
			// fmt.Println(err)
			LogObj.Errorf("etcd取出的结果反系列化失败, %s", err)
			return
		}
	}
	return
}

// watch etcd的变更信息，有变动时发送到对应通道
func WatchChang(key string, NewConfch chan<- []FilePath) (err error) {
	ch := cli.Watch(context.Background(), key)
	for change := range ch {
		for _, evt := range change.Events {
			var NewConf []FilePath
			//  fmt.Println(string(evt.Kv.Value))
			LogObj.Debugf("etcd get value: %s", evt.Kv.Value)
			if evt.Type != clientv3.EventTypeDelete {
				err = json.Unmarshal(evt.Kv.Value, &NewConf)
				if err != nil {
					// fmt.Println("etcd变更信息反序列化失败...", err)
					LogObj.Errorf("etcd变更信息反序列化失败: %s", err)
					return
				}
				NewConfch <- NewConf
			}
		}
	}
	return
}
