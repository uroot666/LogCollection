package TailFile

import (
	"github.com/hpcloud/tail"
)

func NewTailFile(path string) (TileObj *tail.Tail, err error) {
	config := tail.Config{
		ReOpen:    true,                                 // 重新打开
		Follow:    true,                                 // 是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 从文件的那个地方开始读
		MustExist: false,                                // 文件不存在不报错
		Poll:      true,
	}

	TileObj, err = tail.TailFile(path, config)
	return
}
