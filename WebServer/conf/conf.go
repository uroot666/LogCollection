package conf

type EtcdConf struct {
	AddrPort string // etcdip及端口
	Key      string // 用于过滤出相同前缀的key,同一个管理端只能管理一个key前缀
}
