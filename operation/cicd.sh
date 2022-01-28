#/bin/bash

set -e

#这里主要是做CICD的SHELL脚本，但大部分功能转移到了go 里，这里仅剩下一些小部分功能，且不太好被GO替代的

PROJECT_GIT_URL=$1
SERVICE_GIT_CLONE_PATH=$2
SERVICE_NAME=$3

cd $SERVICE_GIT_CLONE_PATH

rm -rf ./*

#echo "clone $PROJECT_GIT_URL $SERVICE_NAME"
git clone $PROJECT_GIT_URL $SERVICE_NAME

#echo "cd $SERVICE_GIT_CLONE_PATH/$SERVICE_NAME"
cd $SERVICE_GIT_CLONE_PATH/$SERVICE_NAME

#根据git生成一个项目的目录名称
CI_COMMIT_TIME=$(git show -s --format=%ct)
CI_COMMIT_TIME_FORMATTED=`TZ='Asia/Shanghai' date -d "$CI_COMMIT_TIME" "+%Y%m%d_%H%M%S"`
CI_COMMIT_ID=$(git rev-parse --short HEAD)
#APP_NAME_FULL="$CI_COMMIT_TIME_FORMATTED-$CI_COMMIT_ID"


echo $CI_COMMIT_ID
