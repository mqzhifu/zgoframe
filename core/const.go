package core

const (
	//全局(公共) 配置文件的类型
	DEFAULT_GLOBAL_CONFIG_FILE_TYPE = "toml"   //toml
	DEFAULT_GLOBAL_CONFIG_FILE_NAME = "config" //name

)

//@parse 全局(公共) 配置类型
const (
	DEFAULT_GLOBAL_CONFIG_TYPE_FILE   = "file"   //从文件中读取配置信息
	DEFAULT_GLOBAL_CONFIG_TYPE_ETCD   = "etcd"   //从ETCD中读取配置信息(暂未使用)
	DEFAULT_GLOBAL_CONFIG_TYPE_CENTER = "center" //从配置吣读取配置信息(暂未使用)
)

//@parse 全局配置文件中，每个模块的：开关选项
const (
	GLOBAL_CONFIG_MODEL_STATUS_OPEN = "open" //打开
	GLOBAL_CONFIG_MODEL_STATUS_OFF  = "off"  //关闭
)

//@parse HTTP 公共响应：自定义 状态码
const (
	HTTP_RES_COMM_ERROR   = 4   //失败
	HTTP_RES_COMM_SUCCESS = 200 //成功
)
