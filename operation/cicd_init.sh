#! /bin/bash

#远程服务器IP
SERVER_HOST=$1
#所有项目的基础目录
PROJECT_ORI_DIR=$2
#git clone下来的所有项目的基目录
TARGET_GIT_PROJECT_DIR=$3
#项目：配置，目录
CONFIG_PROJTECT_NAME=$4
#当前项目名称
PROJECT_NAME=$5
#计算出:最终环境变量，输出到文件中，供调用者引入
ENV_FILE=$6

#git clone下来的单个项目的目录
TARGET_GIT_PROJECT_BUILD_ORI_DIR=$TARGET_GIT_PROJECT_DIR/$PROJECT_NAME
#项目最终的总目录
PROJECT_DIR=$PROJECT_ORI_DIR/$PROJECT_NAME
#远程服务器IP，生成新变量
SERVER_HOST=$TEST_SERVER_HOST

echo "debug SERVER_HOST:$SERVER_HOST PROJECT_NAME:$PROJECT_NAME CONFIG_PROJTECT_NAME:$CONFIG_PROJTECT_NAME "
echo "debug TARGET_GIT_PROJECT_DIR:$TARGET_GIT_PROJECT_DIR TARGET_GIT_PROJECT_BUILD_ORI_DIR:$TARGET_GIT_PROJECT_BUILD_ORI_DIR PROJECT_DIR:$PROJECT_DIR"

#根据git生成一个项目的目录名称
CI_COMMIT_TIME=$(git show -s --format=%ct)
CI_COMMIT_TIME_FORMATTED=`TZ='Asia/Shanghai' date -d @"$CI_COMMIT_TIME" "+%Y%m%d_%H%M%S"`
CI_COMMIT_ID=$(git rev-parse --short HEAD)
APP_NAME_FULL="$CI_COMMIT_TIME_FORMATTED-$CI_COMMIT_ID"

#处理SSH连接相关
echo "debug process ssh:"
mkdir -p /root/.ssh
echo "$KNOWN_HOST_202" > /root/.ssh/known_hosts
echo "$ID_RSA_202" > /root/.ssh/id_rsa
chmod 600 /root/.ssh/id_rsa && chmod 700 /root/.ssh
ls -l /root/.ssh

#项目最终的目录
PROJECT_DIR_FINAL=$PROJECT_DIR/$APP_NAME_FULL
#项目最终的编译目录
BUILD_DIR="$TARGET_GIT_PROJECT_BUILD_ORI_DIR/$APP_NAME_FULL"

echo "debug APP_NAME_FULL:$APP_NAME_FULL PROJECT_DIR_FINAL:PROJECT_DIR_FINAL BUILD_DIR:$BUILD_DIR"

echo "export TARGET_GIT_PROJECT_BUILD_ORI_DIR=$TARGET_GIT_PROJECT_BUILD_ORI_DIR" >> $ENV_FILE
echo "export PROJECT_DIR=$PROJECT_DIR" >> $ENV_FILE
echo "export APP_NAME_FULL=$APP_NAME_FULL" >> $ENV_FILE
echo "export PROJECT_DIR_FINAL=$PROJECT_DIR_FINAL" >> $ENV_FILE
echo "export BUILD_DIR=$BUILD_DIR" >> $ENV_FILE
echo "export SERVER_HOST=$SERVER_HOST" >> $ENV_FILE

chmod 777 $ENV_FILE
