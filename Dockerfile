#golang 编译环境
FROM golang:1.18-alpine AS builder

#使用 alpine-OS ，可以减少镜像大小。但是 alpine 默认的源在国内访问不了，需要修改为国内的源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

#给 OS 安装 GOLANG 编译时，需要的类库: gcc 等
RUN apk add build-base

#设置代码的工作目录，容器启动直接进入此目录
WORKDIR /app

#设置GOLANG的环境变量
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

#将代码统一复制项目代码中
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


#二阶段部署，上面阶段如果直接运行，镜像大概是：1.5GB，使用 alpine-runner 更小
FROM alpine AS runner
WORKDIR /app


#RUN mkdir -p  /app/static

#COPY . .
COPY static ./static
#COPY protobuf ./protobuf
#COPY config.toml .
COPY --from=builder /app/ar120 .


#开放的端口号，注：这里需要看一下项目中的配置文件，要保持一致
EXPOSE 3333 5555

CMD [ "./ar120","-e","5"]
#CMD [ "./ar120","-e","5","-bs","on"]
