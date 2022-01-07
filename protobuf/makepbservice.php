<?php
/*
脚本执行：
php makepbservice.php pb /data/www/golang/src/zgoframe/protobuf/proto /data/www/golang/src/zgoframe/protobuf/pbservice

功能描述：快速生成protobuf中间文件等工具集

1. 编译 目录下，所有的.proto ，生成  pb.go
2. 生成grpc 服务的 快捷实现go文件
3. 生成一个服务的快捷调用方法go文件
4. 对每个服务的每个函数生成对应ID，供长连接使用(生成一个txt文件)

以上所有功能，均依赖：.proto 描述文件

注：正则匹配一个service 块时，结尾必须是：} ，上一行必须是\n结束
*/
define("DEBUG",1);


$argcNoticeMsg = "packageName=xxxx protoFilePath=xxx  outPath=xxx  ";
if (count($argv) < 4){
    exit($argcNoticeMsg);
}

$packageName = $argv[1];//包名
$protoFilePath = $argv[2];//.proto文件路径
$outPath = $argv[3];//生成输出文件的路径

//生成函数映射ID时使用
$GLOBALS["mapFuncIdNo"] = "0";
$GLOBALS["mapServiceIdNo"] = "1";
$GLOBALS["map"] = "";
$mapIdSeparate = "|";

$compileCommand = "export PATH=\$PATH:/var/root/go/bin; cd /data/www/golang/src/zgoframe/protobuf; protoc --go_out=plugins=grpc:./pb ./proto/#proto_file_name#";

pp("packageName:$packageName , protoFilePath:$protoFilePath , outPath:$outPath");
pp("compileCommand:$compileCommand");


$match = null;
//读取一个目录下的所有文件
$protoPathFileList = getDirFiles($protoFilePath);
if (count($protoPathFileList) <=0 ){
    exit("count(protoPathFileList) <=0");
}
//处理每一个.proto 文件
$ServiceFastCallSwitchCase = array();
foreach ($protoPathFileList as $k=>$fileName){
    $fileArr = explode(".",$fileName);
    if (count($fileArr) != 2){
        exit("file name err:$fileName");
    }
    if ( $fileArr[1] != "proto" ) {
        exit("file exit name must = .proto");
    }

    compileProtoFile($compileCommand,$fileName);
    //开始具体处理一个文件里的内容，做:编译、分析等处理，注：一个文件里可能包括多个服务
    $serviceList = oneService($outPath,$protoFilePath,$fileArr[0],$packageName);
    if (!$serviceList || count($serviceList) <= 0){
        continue;
    }
    //生成快捷调用代码 + 函数映射ID
    foreach ($serviceList as $serviceName=>$info){
        $MountClientSwitchCaseStr = MountClientSwitchCase();
        $MountClientSwitchCaseStr = str_replace("#service_name#",$serviceName,$MountClientSwitchCaseStr);
        $MountClientSwitchCaseStr = str_replace("#package_name#",$packageName,$MountClientSwitchCaseStr);

        $GetClientStr = GetClient();
        $GetClientStr = str_replace("#service_name#",$serviceName,$GetClientStr);
        $GetClientStr = str_replace("#package_name#",$packageName,$GetClientStr);

        $ServiceFastCallSwitchCase[$serviceName] = array('MountClientSwitchCase'=>$MountClientSwitchCaseStr,"GetClient"=>$GetClientStr);

        mapFunctionId($serviceName,$info,$mapIdSeparate);
    }
}
createServiceFastCallSwitch($ServiceFastCallSwitchCase,$outPath);
createMapFile($outPath);
exit(111);
//编译proto文件，生成pb.go 文件
function compileProtoFile($compileCommand,$fileName){
    $compileCommandFile = str_replace("#proto_file_name#",$fileName,$compileCommand);
    pp("compileProtoFile: $compileCommandFile");
    $output = shell_exec($compileCommandFile);
    pp("output:$output");
}
//函数名映射ID
function mapFunctionId($serviceName,$serviceFuncListInfo,$mapIdSeparate){
    $serviceId = $GLOBALS["mapServiceIdNo"];
    if(strlen($serviceId) < 2 ){
        $serviceId = $serviceId . "0";
    }
    foreach ($serviceFuncListInfo as $k=>$funcInfo){
        $in = "empty";
        $out = "empty";
        if ($funcInfo["in"]){
            $in =  $funcInfo["in"];
        }
        if ($funcInfo["out"]){
            $out =  $funcInfo["out"];
        }

        $mapFuncIdNo = $GLOBALS["mapFuncIdNo"];
        if(strlen($mapFuncIdNo) == 1 ){
            $mapFuncIdNo = "00".$mapFuncIdNo;
        }else if(strlen($mapFuncIdNo) == 2 ){
            $mapFuncIdNo = "0".$mapFuncIdNo;
        }



//    $outContent .= $GLOBALS["mapIdNo"]  . $mapIdSeparate .$serviceInfo["name"] . $mapIdSeparate .  $in . $mapIdSeparate .  $out . $mapIdSeparate. $arr['desc'] . "\n";
        $outContent = $serviceId . $mapFuncIdNo . $mapIdSeparate .$serviceName . $mapIdSeparate. $funcInfo["name"] . $mapIdSeparate .  $in . $mapIdSeparate .  $out . $mapIdSeparate. $funcInfo['desc'] . "\n";
        $GLOBALS["mapFuncIdNo"] ++;
        $GLOBALS["map"] .= $outContent;
    }

    $GLOBALS["mapServiceIdNo"]++;

}
//具体分析/处理一个文件中的内容
function oneService($outPath,$protoFilePath,$protoFileNamePrefix,$packName){
    pp("oneService:".$protoFileNamePrefix);
    //.proto文件名
    $protoFileName = $protoFilePath . "/". $protoFileNamePrefix.".proto";
    //输出路径
    $outFile = $outPath . "/" . $protoFileNamePrefix. ".go";
    //打开该文件，交获取文件内容
    $protoFileContent = file_get_contents($protoFileName);
    //获取包名
    preg_match_all('/package(.*);/isU',$protoFileContent,$match);
    $package = trim($match[1][0]);

    pp("\n\n\n");
    //读取一个文件中的，若干个service 块
    preg_match_all('/service(.*){(.*)\n}/isU',$protoFileContent,$match);
//    pp($match);

    if (count($match[0]) == 0){
        pp("no match any service block~");
        return array();
    }
    //读取一个service中的所有rpc 函数名
    $service = null;
    foreach ($match[1] as $k =>$v){
        $serviceName = trim($v);
        $rpcFuncMatch = null;
        preg_match_all("/rpc(.*)\((.*)\)(.*)returns(.*)\((.*)\)(.*)\/\/(.*)\n/isU",$match[0][$k],$rpcFuncMatch);
//        var_dump($rpcFuncMatch[0]);
        foreach ($rpcFuncMatch[0] as $k2=>$v2){
            $arr = array(
                'name'=>trim($rpcFuncMatch[1][$k2]),
                 'in'=>trim($rpcFuncMatch[2][$k2]),
                'out'=>trim($rpcFuncMatch[5][$k2]),
                'desc'=>trim($rpcFuncMatch[7][$k2]),
            );
            $service[$serviceName][] = $arr;
        }
    //    var_dump($service);exit;
    //    exit(1111);

    }
    $s = "#";
    $serviceImplementPackage = "pbservice";

    $go_file_content = get_new_go_file();

    $go_file_content = str_replace($s."package".$s,$serviceImplementPackage,$go_file_content);




    foreach ($service as $serviceName=>$funs){
        $go_class_str = go_class_str();
        $go_class_str = str_replace($s."service_class_name".$s,$serviceName,$go_class_str);
        $func_total_str = "\n";
        foreach ($funs as $k=>$v){
            $func_str = get_class_func_str();
            $func_str = str_replace($s."class_name".$s,lcfirst($serviceName),$func_str);
            $func_str = str_replace($s."class_type".$s,$serviceName,$func_str);
            $func_str = str_replace($s."func_name".$s,$v['name'],$func_str);
            $func_str = str_replace($s."para_in_type".$s,$packName.".".$v['in'] ,$func_str);
            $func_str = str_replace($s."para_in_name".$s,lcfirst($v['in']),$func_str);
            $func_str = str_replace($s."return_type".$s,$packName.".".$v['out'],$func_str);
            $func_str = str_replace($s."return_name".$s,lcfirst($v['out']),$func_str);
            $func_total_str .= $func_str . "\n";
        }
        $go_class_str = str_replace($s."funcs".$s,$func_total_str,$go_class_str);
        $go_file_content .=  $go_class_str . "\n" ;
    }

    file_put_contents($outFile,$go_file_content);

    return $service;

}
function createMapFile($outPath){
    $outFile = $outPath . "/" . "map.txt";
    file_put_contents($outFile,$GLOBALS["map"]);
}
function createServiceFastCallSwitch($ServiceFastCallSwitchCase,$outPath){
    $cases = "";
    $getClient = "";
    foreach ($ServiceFastCallSwitchCase as $k=>$v){
        $cases .= $v['MountClientSwitchCase'] . "\n";
        $getClient .= $v['GetClient'] . "\n";
    }
    $MountClientSwitchStr = MountClientSwitch();
    $MountClientSwitchStr = str_replace("#switch_case#",$cases,$MountClientSwitchStr);

    $content = "\n\n" . $getClient . "\n\n" . $MountClientSwitchStr;

    $outFile = $outPath . "/" . "fast_call.go.tmp";
    file_put_contents($outFile,$content);
}

function get_new_go_file(){
$str =<<<EOF
package #package#

import (
	"context"
	"zgoframe/protobuf/pb"
)

EOF;
    return $str;
}

function go_class_str(){
$go_class_str =<<<EOF
type #service_class_name# struct{}

#funcs#
EOF;
return $go_class_str;
}
//一个GRPC服务的，具体实现函数体
function get_class_func_str(){
$go_func =<<<EOF
func (#class_name# *#class_type#)#func_name#(ctx context.Context,#para_in_name# *#para_in_type#) (*#return_type#,error){
    #return_name# := &#return_type#{}
    return #return_name#,nil
}
EOF;
    return $go_func;
}

function getDirFiles($path){
    if(!is_dir($path)){
        return false;
    }
    //scandir方法
    $arr = array();
    $data = scandir($path);
    foreach ($data as $value){
        if($value != '.' && $value != '..'){
            $arr[] = $value;
        }
    }
    return $arr;
}




function pp($info){
    if(DEBUG == 1){
        var_dump($info);
    }
}

function MountClientSwitchCase(){
$str =<<<EOF
    case "#service_name#":
        incClient = #package_name#.New#service_name#Client(myGrpcClient.ClientConn)
EOF;
return $str;
}
function MountClientSwitch(){
    $str =<<<EOF
    func  (myGrpcClient *MyGrpcClient)MountClientToConnect(serviceName string){
        var incClient interface{}
        switch serviceName {
            #switch_case#
        }
    
        myGrpcClient.GrpcClientList[serviceName] = incClient
	}
EOF;
    return $str;
}

function GetClient(){
    $str =<<<EOF
    func (grpcManager *GrpcManager)Get#service_name#Client(name string)(#package_name#.#service_name#Client,error){
        client, err := grpcManager.GetClientByLoadBalance(name,0)
        if err != nil{
            return nil,err
        }
    
        return client.(pb.#service_name#Client),nil
    }
EOF;
    return $str;
}




