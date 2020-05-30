package TailFile

import (
	"LogCollection/KafkaSend"
	"LogCollection/Log"
	"LogCollection/register"
	"context"
	"fmt"
	"github.com/hpcloud/tail"
	"go.uber.org/zap"
	"time"
)

// 存储tail对象的结构体
type TailFile struct {
	Topic    string             // kafka topic
	Path     string             // 被监听日志路径
	Instance *tail.Tail         // 被监听日志文件 Tail对象
	ctx      context.Context    // 用于后面结束对应的goroutine
	cancel   context.CancelFunc // 用于后面结束对应的goroutine
}

var LogObj *zap.SugaredLogger

// 监听日志文件，返回文件对象的指针
func NewTailFile(topic, path string) (TailObj *TailFile, err error) {
	LogObj, _ = Log.GetLogObj()
	if LogObj == nil {
		return
	}

	regMgr := register.GetRegMgr()
	regMgr.Lock.RLock()
	fm := fmt.Sprintf("%s_%s", topic, path)
	offset := regMgr.SelectOffset(fm)
	regMgr.Lock.RUnlock()

	config := tail.Config{
		ReOpen:    true,                                      // 重新打开
		Follow:    true,                                      // 是否跟随
		Location:  &tail.SeekInfo{Offset: offset, Whence: 0}, // 从文件的哪个地方开始读
		MustExist: false,                                     // 文件不存在不报错
		Poll:      true,
		Logger:    tail.DiscardingLogger, // 禁用日志输出
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
	// fmt.Println("开始的offset", offset)
	LogObj.Debugf("开始的offset %v", offset)
	go TailObj.run()
	return
}

// 将新增日志发送到kafka
func (f TailFile) run() {
	regMgr := register.GetRegMgr()
	fm := fmt.Sprintf("%s_%s", f.Topic, f.Path)
	for {
		select {
		case line := <-f.Instance.Lines:
			err := KafkaSend.SendToKafka(line.Text, f.Topic)
			// 如果消息发送失败则不记录文件偏移量
			if err != nil {
				return
			}
			regMgr.Lock.Lock()
			regMgr.Reg[fm], _ = f.Instance.Tell()
			_ = register.RegWrite(regMgr.Path, regMgr.Reg)
			regMgr.Lock.Unlock()

		case <-f.ctx.Done():
			// 日志文件读取配置删除之后，记录此次文件的偏移量，以便下次可以直接使用
			regMgr.Lock.Lock()
			regMgr.Reg[fm], _ = f.Instance.Tell()
			_ = register.RegWrite(regMgr.Path, regMgr.Reg)
			regMgr.Lock.Unlock()
			return
		default:
			time.Sleep(500 * time.Millisecond)
		}
	}
}
