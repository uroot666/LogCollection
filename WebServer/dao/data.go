package dao

import (
	"LogCollectionWeb/Log"
	"context"
	"encoding/json"
	"go.etcd.io/etcd/clientv3"
	"go.uber.org/zap"
	"time"
)

// 用于序列化从etcd取出的value值
type PathTopic struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

// 用于管理从etcd取出的值
type KeyValue struct {
	Key   string      `json:"key"`
	Value []PathTopic `json:"value"`
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
		LogObj.Panicf("连接ETCD失败: %s", err)
		return
	}
	return
}

// 获取etcd中key对应的值
func GetKey(key string) (value *KeyValue, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	resp, err := cli.Get(ctx, key)
	if err != nil {
		LogObj.Panicf("从etcd获取value失败...")
		// fmt.Println("从etcd获取value失败...")
		return
	}
	cancel()
	value = &KeyValue{
		Key:   key,
		Value: make([]PathTopic, 10),
	}
	for _, v := range resp.Kvs {
		// fmt.Println("etcd value", string(v.Value))
		LogObj.Debugf("etcd value: %s", string(v.Value))
		err = json.Unmarshal(v.Value, &value.Value)
		if err != nil {
			// fmt.Println("etcd取出的结果反系列化失败...")
			// fmt.Println(err)
			LogObj.Panicf("etcd取出的结果反系列化失败: %v", err)
			return
		}
	}
	return
}

// 删除etcd的key
func DeletKey(key string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	_, err = cli.Delete(ctx, key)
	if err != nil {
		// fmt.Println("从etcd删除key失败...")
		LogObj.Errorf("从etcd删除key失败: %s", err)
		return
	}
	cancel()
	return
}

// 更新etcd对应key的值
func PutKey(kv *KeyValue) (err error) {
	v, err := json.Marshal(kv.Value)
	vStr := string(v)
	if err != nil {
		// fmt.Println("更新步骤序列化失败...", err)
		LogObj.Errorf("更新步骤序列化失败...", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	_, err = cli.Put(ctx, kv.Key, vStr)
	if err != nil {
		// fmt.Println("从etcd更新key失败...", err)
		LogObj.Errorf("从etcd更新key失败...", err)
		return
	}
	cancel()
	return
}

// 获取所有以某个字符串开头的key的值
func GetAllKey(p string) (allKv map[string][]PathTopic, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	getResp, err := cli.Get(ctx, p, clientv3.WithPrefix())
	if err != nil {
		LogObj.Errorf("获取所有key失败...", err)
		// fmt.Println("获取所有key失败...", err)
		return
	}
	cancel()
	allKv = make(map[string][]PathTopic, 10)
	for _, v := range getResp.Kvs {
		kv := make([]PathTopic, 10)
		err = json.Unmarshal(v.Value, &kv)
		if err != nil {
			LogObj.Errorf("获取所有key序列化失败...", err)
			// fmt.Println("获取所有key序列化失败...", err)
			return
		}
		allKv[string(v.Key)] = kv
	}
	return
}
