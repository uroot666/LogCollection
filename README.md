# LogCallection
用于收集日志发送到kafka.

启动流程

1. 载入配置文件
2. 初始化etcd连接，获取配置信息，生成一个包含所有需要监听文件和topic信息的对象的切片。
3. 初始化kafka连接，创建一个通道给tail模块发送日志。goroutine一个函数监听这个通道，将日志发送到kafka
4. 创建一个管理tail日志goroutine的管理者（tailMgr)，通过etcd获取的信息创建对应的对象并开始监听
5. tailMgr应该有一个通道用于接收etcd的变更信息，初始化这个通道，并对外提供这个通道。获取变更后做对应处理。
6. 创建启动一个etcd的watch，将变更信息发送到tailMgr变更信息的通道中



目录：

1. NodeServer是运行在节点监听日志的
2. WebServer用于提供web服务和接口，动态修改etcd的key