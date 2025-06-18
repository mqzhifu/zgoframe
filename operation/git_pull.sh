#!/bin/bash

BASE_DIR=/Users/clarissechamley/data/gambl/
PROJECT_CNT=0


do_git_func(){
  cd $1
  echo -n  $2" "$3
  echo ""
#  git checkout prod
#  git checkout dev
  git describe --abbrev=0 --tags
#  git pull
#  rm -rf go.sum;go mod tidy
  cd ..
}

for FIRST_LEVERL_DIR in *
do
  dir=$BASE_DIR$FIRST_LEVERL_DIR
  for project_dir in $dir/*
  do
    if [ ! -d $project_dir ];then
        continue
    fi

    arr=(${project_dir//\// })
    module=${arr[4]}
    if [ "$module" = "test" ];then
      continue
    elif [ "$module" = "model" ];then
      continue
    else
      if [ "$module" = "thirdgame" ];then
        do_git_func $project_dir $module ${arr[5]}
        PROJECT_CNT=$[1+$PROJECT_CNT]
      fi
    fi

  done
done


echo "PROJECT count :$PROJECT_CNT"
