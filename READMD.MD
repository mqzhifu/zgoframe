
| 类库 | 说明 |版本  |
| --- | --- | --- |
| viper | 解析配置文件 ||
| zap | 日志
| casbin | 权限控制 ||
| jwt | 登陆验证 ||
| gin | http 容器  ||
| swaager | API管理及注解   ||
| etcd | 分布式存储 ||
| gorm | DB-MYSQL ||
| grpc | 异步调用  ||
| protobuf | 传输内容压缩||
| uuid | uniq-uid ||
| go-redis | redis ||
| gorilla/websocket| websocket ||
| base64Captcha| 图片二进制转码 ||
| sigs.k8s.io/yaml| 解析yaml配置文件 ||
| prometheus/client_golang| metrics ||


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


必须有的组件：
配置类->viper
APP->所有项目必须得有APPID ，统计管理
日志->zap