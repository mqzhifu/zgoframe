source /etc/profile
cd /data/www/golang/src/zgoframe/protobuf
protoc --go_out=plugins=grpc:./pb ./proto/zgoframe.proto

php makepbservice.php pbservice proto zgoframe.proto pbservice pb.