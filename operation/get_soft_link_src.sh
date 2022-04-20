#/bin/bash

set -e

pathFile=$1

rs=`ls -l $pathFile|awk -F '->' '{print $2}'`
echo $rs
