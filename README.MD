拼装一个GO的基础框架方便日常使用，主要是想把经常用的类统一化，像：log 链路追踪 etcd等，保证项目高可用

主要包说明 ：

| 类库                       | 说明         | 版本  |
|--------------------------|------------|-----|
| viper                    | 解析配置文件     ||     |
| zap                      | 日志         |     |
| casbin                   | 权限控制       ||     |
| jwt                      | 登陆验证       ||     |
| gin                      | http 容器    ||     |
| swaager                  | API管理及注解   ||     |
| etcd                     | 分布式存储      ||     |
| gorm                     | DB-MYSQL   ||     |
| grpc                     | 异步调用       ||     |
| protobuf                 | 传输内容压缩     ||     |
| uuid                     | uniq-uid   ||     |
| go-redis                 | redis      ||     |
| gorilla/websocket        | websocket  ||     |
| base64Captcha            | 图片二进制转码    ||     |
| sigs.k8s.io/yaml         | 解析yaml配置文件 ||     |
| prometheus/client_golang | metrics    ||     |

>详细的package，请去godoc里查看


辅助软件：

| key           | desc           |
|---------------|----------------|
| filebeat      | 收集日志并推送        |
| prometheus    | TSDB 收集metrics |
| grafana       | UI展示           |
| superVisor    | 控制进程           |
| git           | 版本控制           |
| rabbitmq      | 队列消息           |
| redis         | 缓存/容器          |
| mysql         | RDBS           |
| zipkin        | 链路追踪           |
| push_gateway  | 收集metrics      |
| node_exporter | 服务器metrics     |
| etcd          | 分布式存储          |

http request header:

| key           | desc           | require |
|---------------|----------------|---------|
| X-Source-Type | 来源载体           | 是       |
| X-Project-Id  | 项目ID           | 是       |
| X-Access      | 项目访问KEY        | 是       |
| X-Request-Id  | 请求ID           | 否       |
| X-Trace-Id    | 追踪 ID          | 否       |
| Client-Info   | 客户端信息,详细去godoc | 否       |

http response body:

```azure
{
    "code": 200,
    "data": {}，
    "msg" : "",
}
```

启动doc
```
go get golang.org/x/tools/cmd/godoc
godoc -http=:6060
```

生成swagger
```
go get -u github.com/swaggo/swag/cmd/swag@v1.7.9
~/go/bin/swag init --parseDependency --parseInternal --parseDepth 1
```

http://127.0.0.1:1111/swagger/index.html


注意：
这里最恶心的是ETCD，不兼容MOD模式。  
pkg.go.com 上的包是:v3.3.25+incompatible  
而官方是：3.4.7  
而不使用mod ,直接 go get go.etcd.io/etc/client 是可以的.  

现在都是MOD模式，只能用3.3.25，而下载过程中就是一堆的BUG...  

1. bbolt包已经换了位置 ，代码中还是使用旧的地址
>replace github.com/coreos/bbolt v1.3.4 => go.etcd.io/bbolt v1.3.4

2. grpc包得降级
>replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

这个是最恶心的，导致如果你项目中使用grpc 也得跟着降级protoc 的版本

go get -u github.com/golang/protobuf/protoc-gen-go@v1.2.0
go get github.com/golang/protobuf/protoc-gen-go@v1.3.2


必须开启组件/类库：
配置类->viper
APP->所有项目必须得有APPID ，统计管理
日志->zap
