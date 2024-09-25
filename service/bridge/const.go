package bridge

// @parse 请求3方服务的协议方法
const (
	REQ_SERVICE_METHOD_HTTP   = 1 //http
	REQ_SERVICE_METHOD_GRPC   = 2 //grpc
	REQ_SERVICE_METHOD_NATIVE = 3 //本地
)

const (
	BRIDGE_SLEEP_TIME = 5     //内部调用时，每个服务要监听 管道里的消息，睡眠时间
	GATEWAY_ADMIN_UID = 99999 //后端反向给前端推送消息时，最好加上一个来源UID
)
