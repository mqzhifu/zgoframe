#! /bin/bash
#执行 compile.sh ~/data/www/golang/zgoframe 项目目录(绝对路径)
echo $PATH
#引入环境变量，确保执行指令 protoc 及 protoc-gen-go 没问题
source /etc/profile
base_dir=$1
protobuf_base_dir="$base_dir/protobuf"
package_name="pb"

function CheckDirExist(){
  dir=$1
  if [ ! -d "$dir" ]; then
    echo "dir( $dir )  is not exist."
    exit
  fi
}


CheckDirExist $base_dir

echo "base_dir:$base_dir"

protobuf_pb_dir="$protobuf_base_dir/$package_name"
protobuf_proto_dir="$protobuf_base_dir/proto"
static_dir="$base_dir/static"
static_proto_dir="$static_dir/proto"
pbservice_dir="$protobuf_base_dir/pbservice"

CheckDirExist $protobuf_base_dir
CheckDirExist $protobuf_pb_dir
CheckDirExist $protobuf_proto_dir
CheckDirExist $static_dir
CheckDirExist $static_proto_dir
CheckDirExist $pbservice_dir


#进入到.proto目录下
cd $protobuf_proto_dir
#编译GO，生成 GRPC protobuf的： PB 文件(项目根目录/protobuf/pb/)
echo "protoc  --go_out=plugins=grpc:$protobuf_pb_dir ./*.proto"
protoc  --go_out=plugins=grpc:$protobuf_pb_dir ./*.proto
#protoc --go_out=.=../pb  --go-grpc_out=../pb ./*.proto  ，这个好像是新版本可以使用
#protoc --go_out=../pb/ ./*.proto //仅生成 pb 文件



#生成 js 的pb文件
echo "protoc --js_out=import_style=commonjs,binary:$protobuf_pb_dir/js *.proto"
protoc --js_out=import_style=commonjs,binary:$protobuf_pb_dir/js *.proto

#下面是让该 js-pb 文件支持浏览器的方式，需要安装 browserify
cd $protobuf_pb_dir/js
echo "browserify all exports.js to exports.pb.js"



browserify exports_common.js > exports_common_pb.js
browserify exports_frame_sync.js > exports_frame_sync_pb.js
browserify exports_gateway.js > exports_gateway_pb.js
browserify exports_game_match.js > exports_game_match_pb.js
browserify exports_twin_agora.js > exports_twin_agora_pb.js

echo "mv exports.pb.js to pash static"

mv exports_common_pb.js     $static_dir/js/pb/common_pb.js
mv exports_frame_sync_pb.js $static_dir/js/pb/frame_sync_pb.js
mv exports_gateway_pb.js    $static_dir/js/pb/gateway_pb.js
mv exports_game_match_pb.js $static_dir/js/pb/game_match_pb.js
mv exports_twin_agora_pb.js $static_dir/js/pb/twin_agora_pb.js

echo "finish clear *_pb.js"
rm *_pb.js

cd $protobuf_base_dir
#生成 id 映射 proto 的文件.txt ，以及动态网关需要的类文件
/usr/local/Cellar/php@7.4/7.4.33_1/bin/php  makepbservice.php pb $protobuf_proto_dir $pbservice_dir


cp -r $protobuf_proto_dir/* $static_proto_dir
cp -r $pbservice_dir/map.txt $static_proto_dir

