package register

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type registerMgr struct {
	Path string     // 注册文件存储路径
	Reg  map[string]int64 // 单个文件对应的ID以及偏移量,ID为日志路径与topic的拼接
	Lock sync.RWMutex	// 用于并发安全
}

// 定义全局的注册表对象
var regMgr registerMgr

// 初始化注册表对象
func Init(path string) {
	// 判断文件是否存在，如不存在则创建一个空文件
	fmt.Println(path)

	_, err := os.Lstat(path)

	if os.IsNotExist(err) {
		fmt.Println("注册表文件不存在,创建...")
		err := RegWrite(path, make(map[string]int64))
		if err != nil {
			fmt.Println("创建注册文件失败")
		}
	}

	filePtr, err := os.Open(path)
	if err != nil {
		fmt.Println("Open file failed [Err:%s]", err.Error())
		return
	}
	defer filePtr.Close()

	reg := make(map[string]int64)

	// 创建json解码器
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&reg)
	if err != nil {
		fmt.Println("Decoder failed", err.Error())
	} else {
		fmt.Println("Decoder success")
		fmt.Println(reg)
	}

	regMgr = registerMgr{
		Path: path,
		Reg:  reg,
	}
}

// 写入文件函数
func RegWrite(path string, rl map[string]int64) (err error) {
	// 创建文件
	filePtr, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Create file failed", err.Error())
		return
	}
	defer filePtr.Close()

	// 创建Json编码器
	encoder := json.NewEncoder(filePtr)

	err = encoder.Encode(rl)
	if err != nil {
		fmt.Println("Encoder failed", err.Error())

	} else {
		fmt.Println("Encoder success")
	}
	return
}

// 将注册表对象持久化
func (r registerMgr) RegPersistence() (err error) {
	err = RegWrite(r.Path, r.Reg)
	return
}

// 根据ID确认是否存在，如果不存在则创建然后offset设置为0
func (r registerMgr) SelectOffset(id string)(offset int64){
	if _, ok := r.Reg[id]; !ok{
		offset = int64(0)
		r.Reg[id] = int64(0)
	} else{
		offset = r.Reg[id]
	}
	return
}

//// 对外开放一个获取注册表对象的函数
func GetRegMgr() *registerMgr {
	return &regMgr
}
