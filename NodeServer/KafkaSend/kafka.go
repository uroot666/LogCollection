package KafkaSend

import (
	"LogCollection/Log"
	"github.com/Shopify/sarama"
)

// 定义一个全局的kafka连接
var kclient sarama.SyncProducer

var LogObj,_ = Log.GetLogObj()

// 初始化kafka连接
func Init(addr []string) (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true                   // 成功交付的信息将在 success channel 返回

	// 连接kafka
	kclient, err = sarama.NewSyncProducer(addr, config)
	if err != nil {
		// fmt.Println("producer closed, err:", err)
		LogObj.Errorf("producer closed, err: %s", err)
		return
	}
	return
}

// 对外开放一个往kafka写数据的函数
func SendToKafka(log, topic string) (err error) {
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(log)
	// 发送消息
	// pid, offset, err := kclient.SendMessage(msg)
	_, _, err = kclient.SendMessage(msg)
	if err != nil {
		// fmt.Println("发送消息失败:", err)
		LogObj.Errorf("发送消息失败: %s", err)
		return
	}
	return
	//fmt.Printf("发送消息结果: pid: %v offset: %v\n", pid, offset)
}
