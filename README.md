# LogCallection
> 用于日志收集的工具。实现的功能是跟踪读取日志文件，将读取的日志条目发送到kafka中。依赖于etcd的watch功能，实现了在不重启客户端的情况下，动态的增加和删减需要监控的日志文件。

## 两个组件

### NodeServer
运行在需要日志收集的节点上，根据从etcd获取到的日志路径列表来收集日志。

#### 依赖环境
1. etcd: 存储需要监听的日志的列表
2. kafka: 用于接收日志

#### 大致流程
1. 初始化register对象，用于对读取日志的偏移量做记录与查找。初始化kafka连接，提供一个对外的函数用于发送日志。初始化etcd连接,用于根据配置文件中的key来查询对应的日志列表。同时维护一个for循环，将watch到的配置变动发送到一个channel中供tailmgr模块做相应的操作。
2. 根据从etcd中获取的日志列表的信息，初始化一个TailMgr的对象，用于管理实际监听文件变动的TailFile对象。TailMgr会将topic与文件路径拼接成一个字符串做唯一标识，来创建Tail对象。之后会维护一个循环来从etcd传输数据的channel中读取配置变动，根据日志条目的增删来调整TailFile对象的创建以及结束。
3. TailMgr管理者创建TailFile对象。创建之前会通过最开始register对象的一个方法，来获取该日志是否有被读取过。如果有则设置从记录的偏移位置开始读取，没有则从0开始读取。对象会有一个run方法维护一个循环来读取日志变动，有变动则通过kafka模块暴露的函数来发送到kafka中。发送成功之后会记录当前偏移量，避免后续客户端重启之后产生重复发送。

### WebServer
提供一个restful api接口以及一个web页面，用来修改和查看etcd中存储的NodeServer监听的所有日志的情况。

#### 依赖环境
1. etcd: 存储需要监听的日志的列表

#### api
```
GET("/value/allKey")

// 操作单个key
GET("/value/key")  
Query key=topic

POST("/value/key")
PUT("/value/key")
{
    Key: "/logagent/192.168.1.105/collect_config", 
    Value:[
        {Path: "c:/tmp/test1/nginx.log", Topic: "web_log"},
    ]
}

DELETE("/value/key")
Query key=topic
```
#### 页面
页面支持增删改查所有topic对应的日志条目。增加日志条目的格式为`topic=filepath`。
![image]https://raw.githubusercontent.com/uroot666/LogCollection/master/img/web.png

## 使用
创建好etcd以及kafka，然后修改NodeServer以及WebServer的conf文件夹中的conf.ini。
NodeServer运行到需要监听日志的节点，可以连接上etcd与kafka。
WebServer可以脸上etcd与kafka即可，不限制部署位置。默认是监听的8080端口。
