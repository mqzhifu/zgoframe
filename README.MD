
# 初衷

使用各种开源包，拼装一个 GO 的基础框架，把经常用的类/包，进行 SDK 化，统一管理。目的：
- SDK 可以用于公共使用，不用重复造轮子，且性能高，BUG 少
- CICD 部署

# 3方 package 说明 ：

>package 的详细信息，请去godoc里查看

| 类库                       | 说明                   | 版本  |
| ------------------------ | -------------------- | --- |
| viper                    | 解析配置文件               |     |
| fsnotify                 | 监听文件变化 viper 使用      |     |
| BurntSushi/toml          | 解析toml配置文件， viper 使用 |     |
| zap                      | 日志                   |     |
| casbin                   | 权限控制                 |     |
| jwt                      | 登陆验证                 |     |
| gin                      | http 容器              |     |
| swaager                  | API管理及注解             |     |
| etcd                     | 分布式存储                |     |
| gorm                     | DB-MYSQL             |     |
| grpc                     | 远程 调用                |     |
| protobuf                 | 传输内容压缩               |     |
| uuid                     | uniq-uid             |     |
| go-redis                 | redis                |     |
| redigo                   | redis                |     |
| gorilla/websocket        | websocket            |     |
| base64Captcha            | 图片二进制转码              |     |
| sigs.k8s.io/yaml         | 解析yaml配置文件           |     |
| prometheus/client_golang | metrics              |     |
| go-supervisord           | 管理 supervisor        |     |
| alibabacloud-go          | 阿里云SDK 主要是OSS        |     |
| go-gomail                | 发邮件                  |     |
| ......                   | ......               |     |


# 系统-辅助软件：

| key | desc |  
|---------------|----------------|  
| filebeat | 收集日志并推送 |  
| prometheus | TSDB 收集metrics |  
| grafana | UI展示 |  
| superVisor | 控制进程 |  
| git | 版本控制 |  
| rabbitmq | 队列消息 |  
| redis | 缓存/容器 |  
| mysql | RDBS |  
| zipkin | 链路追踪 |  
| push_gateway | 收集metrics |  
| node_exporter | 服务器metrics |  
| etcd | 分布式存储 |  


# 目录说明：
![目录说明](https://github.com/mqzhifu/zgoframe/blob/master/dir_desc.png)

# HTTP 请求规范


#### http request header:

| key | desc | require |  
|---------------|----------------|---------|  
| X-Source-Type | 来源载体 | 是 |  
| X-Project-Id | 项目ID | 是 |  
| X-Access | 项目访问KEY | 是 |  
| X-Request-Id | 请求ID | 否 |  
| X-Trace-Id | 追踪 ID | 否 |  
| Client-Info | 客户端信息,详细去godoc | 否 |  

#### http response body:

```javascript  
{  
"code": 200,  
"data": {}，  
"msg" : "",  
}  
```  


# 必须开启组件/类库：

| key   | desc                   |     |
| ----- | ---------------------- | --- |
| viper | 读取配置文件                 |     |
| mysql | projectInfo目前是存在MYSQL中 |     |
| zap   | 日志                     |     |
| redis | http限流使用               |     |

# 项目启动流程


#### 配置文件

cp config/config.toml.tmp config/config.toml

主要修改的：
- mysql
- redis
- 日志目录

#### go  doc

安装：
>go get golang.org/x/tools/cmd/godoc@v0.17.0
>go install golang.org/x/tools/cmd/godoc@v0.17.0

启动
>$HOME/go/bin/godoc -http=:6060

查看
>http://127.0.0.1:6060

#### swag

安装
>go get -u github.com/swaggo/swag/cmd/swag@v1.7.9

生成swagger 文档，这里得参数，因为有些结构体用的3方库
>$HOME/go/bin/swag init --parseDependency --parseInternal --parseDepth 3

访问地址：
>http://127.0.0.1:1111/swagger/index.html

#### 主程序

启动:
>go run main.go -e 1

查看指令行参数：
>go run main.go -h

| 参数     | 类型     | 描述                              | 默认      |     |
| :----- | :----- | :------------------------------ | :------ | --- |
| -bs    | string | BuildStatic                     | off     |     |
| -cfn   | string | configFileName                  | config  |     |
| -cs    | string | configSource:file or etcd       | file    |     |
| -ct    | string | configFileType                  | toml    |     |
| -debug | int    | startup debug mode level        | 0       |     |
| -e     | int    | must require , 1本地2开发3测试4预发布5线上 | 必填，没有默认 |     |
| etl    | string | get etcd config url             |         |     |
| -t     | string | 是否为测试，且测试的具体模块名                 | 空       |     |
|        |        |                                 |         |     |

>-e 是使用最多的，且是必选项    。 -t 是做测试使用

# 代码编写注意

1. 严禁使用 CGO(维护、编译都太复杂)，尤其引用 3方包做好测试
2. 严格使用 MSYQL 联表(尽量做汇总表，或者存 redis)
3. 每一个API接口，都要写上对应的 SWAG 描述信息
4. 每段代码，都要写注释
5. 函数名、变量名不要使用单字符、拼音之类的，最好使用有意义的单词
6. 代码逻辑最好要清晰，，不要随意添加/删除文件：
   - 接口就放在 API 目录下
   - DB定义就放到 model 下
   - 配置文件放到 config 目录下
   - ......
7. 有公共的SDK就必须使用，不要自己编写，重复造轮子
8. 开发项目，要有 CICD 的概念，加了哪些模块，如何部署，都要考虑进去，且有:
   - docker_makefile
   - docker_composer
9. 任何无用的文件不允许进GIT库
   - 加密相关的：SSH证书、3方密钥等
   - IDE 的配置文件
   - 日志文件
   - ......
10. godoc go-test pprof go-cover 都要用起来
11. 反射只用结构体的一些操作，节省代码，不要做复杂的反射，效率太低，阅读太难
12. 泛型只用函数参数，一些简单的操作，像|定义|之类的不要用泛型，阅读太难

# 配置文件说明

目前仅使用了toml格式，也可以兼容yaml 等，自行研究  
每个模块里都有 status:  open | close ，自行设置  
基本上大部分模块，从配置中分析一下，也能知道个大概

# CICD/容器部署

>详细的可参考 Dockfile 文件

3方依赖容器:
- mysql
- redis
- superVisor
- os (alpine-OS)

3方依赖软件：
- golang 1.18
- git


# 结语

项目的初衷：不是要很强大，而是最小化实现主要功能，保持 灵活 多变 少依赖 的特点
>从部署上看，也是非常的简单，能不依赖就尽量不依赖

本人很反感用一堆高大上的软件( spring | k8s )，或者高大上的技术(CGO)，看着挺牛B，一但：项目迁移/项目部署：
- 各种OS不兼容(主要是win兼容不友好)
- 各种相互依赖软件安装
- 依赖软件还有版本要求
- 代码均在容器里运行，但没有完整的日志、链路追踪，导致上线后，BUG短时间无法修改
- 程序员本地开发，环境无法完全安装，测试起来很麻烦
- ......

真的很让人很恼火。

注：如果真的项目很大、开发人员很多，使用的开源软件也很多，那，不推荐使用此框架，建议：使用 K8S/云原生 体系，做完整的容器化部署。