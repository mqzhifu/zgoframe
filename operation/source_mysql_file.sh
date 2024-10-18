#!/bin/bash

if [ ! -n "$1" ]; then
    echo "DB_FILE_DIR empty!"
    exit
fi

if [ ! -n "$2" ]; then
    echo "DB_NAME empty!"
    exit
fi


echo "DB_FILE_DIR:$1 , DB_NAME:$2"
DB_FILE_DIR=$1
DB_NAME=$2

dir=`ls $DB_FILE_DIR`
echo "" > all.sql

for i in $dir
do
    cat "$DB_FILE_DIR/$i;" >> all.sql
done

mysql -uroot -p123456 -h127.0.0.1 -e "drop database $DB_NAME"
mysql -uroot -p123456 -h127.0.0.1 -e "create database $DB_NAME charset=utf8;show databases;use $DB_NAME;all.sql;"