#! /bin/bash
#执行 compile.sh ~/data/www/golang/zgoframe 项目目录

#引入环境变量，确保执行指令 protoc 及 protoc-gen-go 没问题
source /etc/profile
base_dir=$1

if  [ ! -n "$base_dir" ] ;then
    echo "project base_dir is empty string."
    exit
fi

echo "base_dir:$base_dir"

#进入到.proto目录下
cd $base_dir/protobuf/proto
#编译GO，生成 GRPC protobuf的： PB 文件
echo "protoc  --go_out=plugins=grpc:../pb ./*.proto"
protoc  --go_out=plugins=grpc:../pb ./*.proto
#protoc --go_out=.=../pb  --go-grpc_out=../pb ./*.proto  ，这个好像是新版本可以使用
#protoc --go_out=../pb/ ./*.proto //仅生成 pb 文件
#生成 id 映射 proto 的文件.txt ，以及动态网关需要的类文件
#php makepbservice.php pbservice proto *.proto pbservice pb.

#生成js pb文件，这里注意下目录，JS支持浏览器的方式依赖node几个类包
echo "protoc --js_out=import_style=commonjs,binary:./../pb/js *.proto"
protoc --js_out=import_style=commonjs,binary:./../pb/js *.proto
cd $base_dir/protobuf/pb/js


echo "browserify all exports.js to exports.pb.js"

browserify exports_common.js > exports_common_pb.js
browserify exports_frame_sync.js > exports_frame_sync_pb.js
browserify exports_gateway.js > exports_gateway_pb.js
browserify exports_game_match.js > exports_game_match_pb.js
browserify exports_twin_agora.js > exports_twin_agora_pb.js

echo "mv  exports.pb.js to pash static"

mv exports_common_pb.js $base_dir/static/js/pb/common_pb.js
mv exports_frame_sync_pb.js $base_dir/static/js/pb/frame_sync_pb.js
mv exports_gateway_pb.js $base_dir/static/js/pb/gateway_pb.js
mv exports_game_match_pb.js $base_dir/static/js/pb/game_match_pb.js
mv exports_twin_agora_pb.js $base_dir/static/js/pb/twin_agora_pb.js

echo "finish clear *_pb.js"
rm *_pb.js
