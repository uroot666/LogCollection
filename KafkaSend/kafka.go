package KafkaSend

import (
	"fmt"
	"github.com/Shopify/sarama"
)

var kclient sarama.SyncProducer

func Init(addr []string) (err error) {
	config := sarama.NewConfig()
	// tailf 包使用
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true                   // 成功交付的信息将在 success channel 返回

	// 连接kafka
	kclient, err = sarama.NewSyncProducer(addr, config)
	if err != nil {
		fmt.Println("producer closed, err:", err)
		return
	}
	return
}

func SendToKafka(log, topic string) {
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(log)

	// 发送消息
	pid, offset, err := kclient.SendMessage(msg)
	if err != nil {
		fmt.Println("发送消息失败:", err)
	}
	fmt.Printf("发送消息结果: pid: %v offset: %v\n", pid, offset)
}
