package EtcdReadWatch

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

type FilePath struct {
	Path string `json:"path"`
}

func Init(addr string) (cli *clientv3.Client, err error) {
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{addr},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("连接ETCD失败...")
		return
	}
	return
}

func GetFilePath(key string, cli *clientv3.Client) (filepath []*FilePath, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	resp, err := cli.Get(ctx, key)
	if err != nil {
		fmt.Println("从etcd获取value失败...")
		return
	}
	cancel()
	for _, v := range resp.Kvs {
		fmt.Println("etcd value", string(v.Value))
		err = json.Unmarshal(v.Value, &filepath)
		if err != nil {
			fmt.Println("etcd取出的结果反系列化失败...")
			fmt.Println(err)
			return
		}
	}
	return
}
