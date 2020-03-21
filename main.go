package main

import (
	"LogCollection/EtcdReadWatch"
	"LogCollection/KafkaSend"
	"LogCollection/TailFile"
	"fmt"
	"time"
)

func main() {
	KafkaSend.Init([]string{"192.168.198.163:9092"})

	addr := "192.168.198.163:2379"
	key := "filepath"
	cli, err := EtcdReadWatch.Init(addr)
	if err != nil {
		fmt.Println("初始化etcd失败")
	}
	filePathObj, _ := EtcdReadWatch.GetFilePath(key, cli)
	if filePathObj[0].Path == "" {
		return
	}
	fmt.Printf("开始监听: %v\n", filePathObj[0].Path)
	TailObj, _ := TailFile.NewTailFile(filePathObj[0].Path)
	for {
		select {
		case line := <-TailObj.Lines:
			fmt.Printf("日志: %v\n", line.Text)
			KafkaSend.SendToKafka(line.Text, "testTopic")
		default:
			time.Sleep(500 * time.Millisecond)
		}
	}

}
