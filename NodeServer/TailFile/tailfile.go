package TailFile

import (
	"LogCollection/KafkaSend"
	"context"
	"fmt"
	"github.com/hpcloud/tail"
	"time"
)

// 存储tail对象的结构体
type TailFile struct {
	Topic    string
	Path     string
	Instance *tail.Tail
	ctx      context.Context // 用于后面结束对应的goroutine
	cancel   context.CancelFunc
}

// 监听日志日志文件，返回文件对象的指针
func NewTailFile(topic, path string) (TailObj *TailFile, err error) {
	config := tail.Config{
		ReOpen:    true,                                 // 重新打开
		Follow:    true,                                 // 是否跟随
		Location:  &tail.SeekInfo{Offset: 294, Whence: 0}, // 从文件的哪个地方开始读
		MustExist: false,                                // 文件不存在不报错
		Poll:      true,
		//Location:  &tail.SeekInfo{Offset: 0, Whence: io.SeekStart}, // 从文件最开始读取
	}
	fObj, err := tail.TailFile(path, config)
	ctx, cancel := context.WithCancel(context.Background())
	TailObj = &TailFile{
		Topic:    topic,
		Path:     path,
		Instance: fObj,
		ctx:      ctx,
		cancel:   cancel,
	}
	offset, _:= TailObj.Instance.Tell()
	fmt.Println("开始的offset", offset)
	go TailObj.run()
	return
}

// 将新增日志发送到kafka
func (f TailFile) run() {
	for {
		select {
		case line := <-f.Instance.Lines:
			KafkaSend.SendToKafka(line.Text, f.Topic)
			offset, _:= f.Instance.Tell()
			fmt.Println("运行中", offset)
		case <-f.ctx.Done():
			fmt.Println("退出了...", f.Topic)
			offset, _:= f.Instance.Tell()
			fmt.Println("结束的", offset)
			return
		default:
			time.Sleep(500 * time.Millisecond)
		}
	}

}
