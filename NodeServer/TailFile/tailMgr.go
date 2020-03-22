package TailFile

import (
	"LogCollection/EtcdReadWatch"
	"fmt"
)

// 用于管理Tail对象
type TailMgr struct {
	ConfigList  []EtcdReadWatch.FilePath      // 从etcd拿到的需要监控的日志文件路径以及其对应的kafka的tpoic
	Instance    map[string]*TailFile          // 监控的日志对象的集合
	NewConfigch chan []EtcdReadWatch.FilePath // 用于接收变动信息的通道，给etcd使用
}

// 创建一个全局的管理者管理Tail对象
var Mgr TailMgr

// 初始化需要监控的日志文件，并全部添加到map中管理起来
func Init(fileConfg []EtcdReadWatch.FilePath) (err error) {
	Mgr = TailMgr{
		ConfigList:  fileConfg,
		Instance:    make(map[string]*TailFile, 20),
		NewConfigch: make(chan []EtcdReadWatch.FilePath),
	}
	for _, conf := range fileConfg {
		fm := fmt.Sprintf("%s_%s", conf.Topic, conf.Path)
		obj, _ := NewTailFile(conf.Topic, conf.Path)
		Mgr.Instance[fm] = obj
	}
	go Mgr.run()
	return
}

// 返回管理者用于接收etcd变动信息的通道
func GetNewConfCh() (NewConfCh chan []EtcdReadWatch.FilePath) {
	return Mgr.NewConfigch
}

// 监控etcd的配置信息变更
func (m TailMgr) run() {
	for newconf := range m.NewConfigch {
		// 将新配置和旧配置做对比，如果有新增就增加goroutine监控日志文件
		for _, c1 := range newconf {
			existence := false
			c1Str := fmt.Sprintf("%s_%s", c1.Topic, c1.Path)
			for c2Str, _ := range m.Instance {
				if c2Str == c1Str {
					existence = true
					break
				}
			}
			if !existence {
				obj, _ := NewTailFile(c1.Topic, c1.Path)
				m.Instance[c1Str] = obj
				m.ConfigList = append(m.ConfigList, c1)
			}
		}
		// 将新配置和旧配置做对比，如果有减少就停止对应的goroutine
		for c2Str, _ := range m.Instance {
			existence := false
			for _, c1 := range newconf {
				c1Str := fmt.Sprintf("%s_%s", c1.Topic, c1.Path)
				if c2Str == c1Str {
					existence = true
					break
				}
			}
			if !existence {
				obj := m.Instance[c2Str]
				obj.cancel()
				delete(m.Instance, c2Str)
			}
		}

	}
}
