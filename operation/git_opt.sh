#!/bin/bash

BASE_DIR=/Users/clarissechamley/data/gambl/

PROJECT_CNT=0
for FIRST_DIR in *
do
  dir=$BASE_DIR$FIRST_DIR
  for project_dir in $dir/*
  do
    if test -d $project_dir
    then
      arr=(${project_dir//\// })
      module=${arr[4]}
      if [ "$module" = "test" ];then
	      continue
      elif [ "$module" = "model" ];then
	      continue
      else
        if [ "$module" = "pay" ];then
          PROJECT_CNT=$[1+$PROJECT_CNT]
          cd $project_dir
          echo -n  ${arr[5]}" "
          #git checkout prod
          git describe --abbrev=0 --tags
          # rm -rf go.sum;go mod tidy
          #git pull
        fi
        cd ..
       fi
    fi
  done
done
echo "PROJECT count :$PROJECT_CNT"
