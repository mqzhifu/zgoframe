set -e

#根据git生成一个项目的目录名称
CI_COMMIT_TIME=$(git show -s --format=%ct)
CI_COMMIT_TIME_FORMATTED=`TZ='Asia/Shanghai' date -d @"$CI_COMMIT_TIME" "+%Y%m%d_%H%M%S"`
CI_COMMIT_ID=$(git rev-parse --short HEAD)
APP_NAME_FULL="$CI_COMMIT_TIME_FORMATTED-$CI_COMMIT_ID"

echo $APP_NAME_FULL