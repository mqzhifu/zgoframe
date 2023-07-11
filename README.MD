拼装一个GO的基础框架，主要是想把经常用的类/包 统一化/SDK化，统一管理，节省重复造轮子。  
像：log 链路追踪 etcd等，保证项目高可用/高扩展

主要 3方 package 说明 ：

| 类库                       | 说明            | 版本  |
|--------------------------|---------------|-----|
| viper                    | 解析配置文件        ||     |
| zap                      | 日志            |     |
| casbin                   | 权限控制          ||     |
| jwt                      | 登陆验证          ||     |
| gin                      | http 容器       ||     |
| swaager                  | API管理及注解      ||     |
| etcd                     | 分布式存储         ||     |
| gorm                     | DB-MYSQL      ||     |
| grpc                     | 远程 调用         ||     |
| protobuf                 | 传输内容压缩        ||     |
| uuid                     | uniq-uid      ||     |
| go-redis                 | redis         ||     |
| redigo                   | redis         ||     |
| gorilla/websocket        | websocket     ||     |
| base64Captcha            | 图片二进制转码       ||     |
| sigs.k8s.io/yaml         | 解析yaml配置文件    ||     |
| prometheus/client_golang | metrics       ||     |
| go-supervisord           | 管理 supervisor ||     |
| alibabacloud-go          | 阿里云SDK 主要是OSS ||     |

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


目录说明：
![目录说明](https://github.com/mqzhifu/zgoframe/blob/master/dir_desc.png)

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
~/go/bin/godoc -http=:6060
```

安装swag指令行工具
```
go get -u github.com/swaggo/swag/cmd/swag@v1.7.9
```
生成swagger 文档，这里得参数，因为有些结构体用的3方库
``
~/go/bin/swag init --parseDependency --parseInternal --parseDepth 3
``
访问地址：
>http://127.0.0.1:1111/swagger/index.html

启动:
>go run main.go -e 1

查看指令行参数：
>go run main.go -h
```
  -bs string
    	BuildStatic (default "off")
  -cfn string
    	configFileName (default "config")
  -cs string
    	configSource:file or etcd (default "file")
  -ct string
    	configFileType (default "toml")
  -debug int
    	startup debug mode level
  -e int
    	must require , 1本地2开发3测试4预发布5线上
  -etl string
    	get etcd config url (default "http://127.0.0.1/getEtcdCluster/Ip/Port")
  -t string
    	testFlag:empty or 1
```
>-e 是使用最多的，且是必须的。主要是区分环境，如：本地、开发、测试、预发布、线上 

>-t 是做测试使用


必须开启组件/类库：

| key   | desc                   |
|-------|------------------------|
| viper | 读取配置文件                 |
| mysql | projectInfo目前是存在MYSQL中 |
| zap   | 日志                     |
| redis | http限流使用               |



配置文件说明

目前仅使用了toml格式，也可以兼容yaml 等，自行研究    
每个模块里都有status: open  | close ，自行设置  
基本上大部分模块，从配置中分析一下，也能知道个大概  

CICD/容器部署

自行参考 Dockfile 文件

