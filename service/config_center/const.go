package config_center

// @parse 配置中心-数据持久化类型
const (
	PERSISTENCE_TYPE_OFF     = 0 //关闭
	PERSISTENCE_TYPE_MYSQL   = 1 //mysql数据库
	PERSISTENCE_TYPE_REDIS   = 2 //redis缓存
	PERSISTENCE_TYPE_FILE    = 3 //文件
	PERSISTENCE_TYPE_ETCD    = 4 //etcd
	PERSISTENCE_TYPE_CONSULE = 5 //consul
)
