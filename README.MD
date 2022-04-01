拼装一个GO的基础框架，主要是想把经常用的类/包 统一化/SDK化，统一管理，节省重复造轮子。  
像：log 链路追踪 etcd等，保证项目高可用

主要package说明 ：

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
| protobuf                 | 传输内容压缩     ||     |g
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

安装swag指令行工具
```
go get -u github.com/swaggo/swag/cmd/swag@v1.7.9
```
生成swagger 文档，这里得参数，因为有些结构体用的3方库
``
~/go/bin/swag init --parseDependency --parseInternal --parseDepth 1
``
访问地址：
>http://127.0.0.1:1111/swagger/index.html


包管理 注意：  
这里最恶心的是ETCD，它不是完全兼容MOD模式。官方的版本比较高，给MOD的是略低些的版本，有些新的功能不可以用，同时：连带着GRPC也不能用高版本  
pkg.go.com 上的包是:v3.3.25+incompatible  

官方：3.4.7  
MOD：3.3.25

1. bbolt包已经换了位置 ，代码中还是使用旧的地址
>replace github.com/coreos/bbolt v1.3.4 => go.etcd.io/bbolt v1.3.4

2. grpc包得降级
>replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

这个是最恶心的，导致如果你项目中使用grpc 也得跟着降级protoc 的版本

如果使用protobuf，需要先下载工具包
```azure
go get -u github.com/golang/protobuf/protoc-gen-go@v1.2.0
go get github.com/golang/protobuf/protoc-gen-go@v1.3.2
```


启动:
>go run main.go -e local

指令行参数：
>go run main.go -h
```
  -cfn string
    	configFileName (default "config")
  -cs string
    	configSource:file or etcd (default "file")
  -ct string
    	configFileType (default "toml")
  -debug int
    	startup debug mode level
  -e string
    	must require , env:local test pre dev online
  -etl string
    	get etcd config url (default "http://127.0.0.1/getEtcdCluster/Ip/Port")
```

必须开启组件/类库：

| key   | desc                   |
|-------|------------------------|
| viper | 读取配置文件                 |
| mysql | projectInfo目前是存在MYSQL中 |
| zap   | 日志                     |
| redis | http限流使用               |

配置文件说明：  
目前仅使用了toml格式，也可以兼容yaml 等，自行研究    
每个模块里都有status: open  | close ，自行设置  
基本上大部分模块，从配置中分析一下，也能知道个大概  
