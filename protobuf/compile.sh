source /etc/profile
cd /data/www/golang/src/zgoframe/protobuf/proto

protoc  --go_out=plugins=grpc:../pb ./*.proto

php makepbservice.php pbservice proto *.proto pbservice pb.


#生成js pb文件，这里注意下目录，JS支持浏览器的方式依赖node几个类包
protoc --js_out=import_style=commonjs,binary:./../pb/js *.proto