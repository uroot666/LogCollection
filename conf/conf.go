package conf

type KafkaConf struct {
	AddrPort string
}

type EtcdConf struct {
	AddrPort string
	Key      string
}
