#/bin/bash

set -e

#./cp_project.sh /data/www/golang/src/zgoframe /data/www/golang/src/log_slave

ORI_PROJECT_DIR=$1
TAR_GET_DIR=$2

echo "ORI_PROJECT_DIR:$ORI_PROJECT_DIR TAR_GET_DIR:$TAR_GET_DIR"

#先把目标目录下的目录清空下
cd $TAR_GET_DIR

ls -l

rm -rf *
rm -rf .drone.yml .gitignore

cd $ORI_PROJECT_DIR

targetZipFileFullName="${TAR_GET_DIR}/git_project.zip"

echo "targetZipFileFullName:$targetZipFileFullName"

`git archive --format zip  --output "$targetZipFileFullName" master`

cd $TAR_GET_DIR
ls -l

unzip $targetZipFileFullName
rm -rf $targetZipFileFullName

cat ${ORI_PROJECT_DIR}/config.toml |sed 's/projectId = 6/projectId = 3/' > $TAR_GET_DIR/config.toml
cat $TAR_GET_DIR/config.toml |sed 's/projectId = 6/projectId = 3/' > $TAR_GET_DIR/config.toml