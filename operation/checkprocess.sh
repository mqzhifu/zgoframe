#! /bin/bash

#通过shell检测一个进程是否已启动，如果未启动(挂了)，自动重启，跟superVisor略像

script_name="checkprocess"
process_name=$1

if  [ ! -n "$process_name" ] ;then
    echo "process_name is empty string."
    exit
fi

pidFile=$2

if  [ ! -n "$pidFile" ] ;then
    echo "pidFile is empty string."
    exit
fi

echo "process_name:$process_name"
echo "pidFile : $pidFile"

pid=0
#根据 进程名 关键字 获取(模糊搜索)进程ID
#缺点：因为是模糊搜索，如果其它进程名称相似，会识破认为也是对的
#精确搜索：grep - w , 如果一个脚本启动后，再开个窗口启动一样一样的，依然有漏洞
function getProcessPidByPS(){
  find_process_info=`ps -ef |grep $process_name|grep -v grep|grep -v $script_name`
  echo "find_process_info:$find_process_info"

  if [ "$find_process_info" !=  "" ];then
    echo processexist-need kill ;

    find_process_num=`ps -ef |grep $process_name|grep -v grep|grep -v $script_name|wc -l`
    echo "find_process_num:$find_process_num"

    if [ $find_process_num != 1 ];then
      echo "err:find_process_num != 1"
    else
      pid=`ps -ef |grep $process_name|grep -v grep|grep -v $script_name|awk '{print $2}'`
    fi
  else
    echo process-not-exist;
  fi

}
#获取一个已经启动进程的ID

#根据 进程名 关键字 获取(模糊搜索)进程ID ，精准查找
function getProcessPidByPidof(){
  pidofIds=`pidof $process_name`
  echo "pidof rs:$pidofId"
  if [ $pidofIds !=  "" ];then
    find_space=`echo $pidofIds |grep " "`
    if [ $find_space != ""];then
      echo "waning:have multi pids"
    fi
    pid=$pidofId
  fi
}

function getProcessPidByFile(){
  pidFile=$1
  if [ -f $pidFile ];then
    fileContentPid=`cat $pidFile`
    echo "fileContentPid:$fileContentPid"
    if [ "$fileContentPid" != "" ];then
      fileContentPidWc=`ps -p $fileContentPid|wc -l`
      echo "fileContentPidWc:$fileContentPidWc"
      if [ $fileContentPidWc -le 1 ];then
        echo "pid not exec."
      else
        pid=$fileContentPid
      fi
    else
      echo "file content is empty"
    fi
  else
    echo "file not exist"
  fi
}
getProcessPidByFile $pidFile
#getProcessPidByPS
#getProcessPidByPidof
echo $pid
if [ $pid != 0 ];then
  echo "kill $pid"
  `kill $pid`
  sleep 2
else
  echo "pid = 0 "
fi