#golang 环境
FROM golang:1.18-alpine AS builder
#FROM golang:1.18 AS build

#使用 alpine ，可以减少镜像大小。但是 alpine 默认的源在国内访问不了，需要修改为国内的源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

#安装编译需要的环境gcc等
RUN apk add build-base

#设置代码的工作目录，容器启动直接进入此目录
WORKDIR /app

#设置GOLANG的环境变量
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

#复制项目代码
COPY . .

#下载 goland 依赖包
#RUN go version;go env -w GO111MODULE=on;go env -w GOPROXY=https://goproxy.cn,direct;
RUN go mod tidy;

#编译项目代码
RUN go build -o ar120

#帮助文档
#RUN go install github.com/swaggo/swag/cmd/swag@v1.7.9;
#RUN $HOME/go/bin/swag init --parseDependency --parseInternal --parseDepth 3;
#RUN swag -v 

#开放的端口号，注：这里需要看一下项目中的配置文件，要保持一致
EXPOSE 3333 5555

CMD [ "./ar120","-e","5"]
