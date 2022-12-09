<?php
/*
脚本执行：
php makepbservice.php pb ~/data/www/golang/zgoframe/protobuf/proto ~/data/www/golang/zgoframe/protobuf/pbservice
php makepbservice.php pb /data/www/golang/zgoframe/protobuf/proto /data/www/golang/zgoframe/protobuf/pbservice

功能描述：快速生成protobuf中间文件等工具集(比较懒，就不用GOLANG写了)
大体上一共是2类：
1、用 protoc 工具，编译生成 GO 和 JS 的所有 可执行文件
2、工具附加代码
    1. 编译 目录下，所有的.proto ，生成  pb.go (暂不在这里执行)
    2. 对每个服务的每个函数生成对应ID，供长连接使用(生成一个txt文件)
    3. 生成文件: (1) 所有服务的 grpc client 创建(负载的寻找 IP+PORT)  , (2) 根据服务名函数名-动态调用服务的方法
    4. 生成: 每个服务一个文件，是对该服务的具体实现(类似抽象类，无实现)

    CallGrpc -> CallServiceFuncFrameSync-> GetFrameSyncClient

>ps:以上所有功能，均依赖：.proto 描述文件

注：正则匹配一个 service 块时，结尾必须是：} ，上一行必须是\n结束
*/

define("DEBUG",1);

//引入 静态模块类(字符串)，用于动态生成 GO 文件、连接DB等
include 'makepbservice_template.php';
$template = new Template();
$db = new DB();

$argcNoticeMsg = "packageName=pb protoFilePath=/data/www/golang/zgoframe/protobuf/proto  outPath=/data/www/golang/zgoframe/protobuf/pbservice  ";
if (count($argv) < 4){
    exit($argcNoticeMsg);
}

$packageName = $argv[1];//包名,目前就一个：pb
$protoFilePath = $argv[2];//.proto文件路径
$outPath = $argv[3];//生成输出文件的路径

//生成函数映射ID时使用
//$GLOBALS["mapFuncIdNo"] = "1";
//$GLOBALS["mapServiceIdNo"] = "1";
$GLOBALS["map"] = "";
$mapIdSeparate = "|";
//快速生成一个服务的具体实现类（包名）
$serviceImplementPackage = "pbservice";
$fast_call_file_name = "ast_call.go.tmp";//最终 的 一个文件；网关动态调用 服务的函数
$callServiceFuncTotalStr = "";


//编译 proto 生成 PB 文件的SHELL脚本 ，这个分开执行，先不在这里执行了
// $compileCommand = "export PATH=\$PATH:/Users/mayanyan/go/bin; cd /Users/mayanyan/data/www/golang/src/zgoframe/protobuf; protoc --go_out=plugins=grpc:./pb ./proto/#proto_file_name#";
//pp("compileCommand:$compileCommand");

pp("packageName:$packageName , protoFilePath:$protoFilePath , outPath:$outPath serviceImplementPackage:$serviceImplementPackage");


//读取 proto 目录下的所有文件 .proto 文件
//$match = null;
$protoPathFileList = getDirFiles($protoFilePath);
if (count($protoPathFileList) <=0 ){
    exit("count(protoPathFileList) <=0 , 目录下的：.proto 一个没有 ");
}

$callGrpcServiceCase = "";

//处理每一个.proto 文件,目前是：一个 .proto 代表是一个 微服务
$ServiceFastCallSwitchCase = array();
foreach ($protoPathFileList as $k=>$fileName){
    //验证proto文件名
    $fileArr = explode(".",$fileName);
    if (count($fileArr) != 2){
        exit("处理 proto 目录文件时出错：file name err:$fileName");
    }
    if ( $fileArr[1] != "proto" ) {
        exit("处理 proto 目录文件时出错：file exit name must = .proto");
    }
    $serviceName = $fileArr[0];//文件名，即服务名
    if ( $serviceName == "common" ) {//忽略 common 服务，它没有任何函数，只是单纯的结构体定义
        echo "ignore common service\n";
        continue;
    }

    //编译proto生成pb 文件
    # compileProtoFile($compileCommand,$fileName);
    //开始具体处理一个文件里的内容，做:编译、分析等处理，注：一个文件里可能包括多个服务
    $serviceList = parserOneServiceFileContent($outPath,$protoFilePath,$serviceName,$packageName,$mapIdSeparate);
//    if (!$serviceList || count($serviceList) <= 0){
//        continue;
//    }
//    //生成快捷调用代码 + 函数映射ID
//    foreach ($serviceList as $serviceName=>$info){
//        dynamicCallGrpcService($serviceName,$info,$packageName);
//
//        $MountClientSwitchCaseStr = $template->MountClientSwitchCase();
//        $MountClientSwitchCaseStr = str_replace("#service_name#",$serviceName,$MountClientSwitchCaseStr);
//        $MountClientSwitchCaseStr = str_replace("#package_name#",$packageName,$MountClientSwitchCaseStr);
//
//        $GetClientStr = $template->GetClient();
//        $GetClientStr = str_replace("#service_name#",$serviceName,$GetClientStr);
//        $GetClientStr = str_replace("#package_name#",$packageName,$GetClientStr);
//
//        $ServiceFastCallSwitchCase[$serviceName] = array('MountClientSwitchCase'=>$MountClientSwitchCaseStr,"GetClient"=>$GetClientStr);
//
//        mapFunctionId($serviceName,$info,$mapIdSeparate);
//    }
}

createServiceFastCallSwitch($ServiceFastCallSwitchCase,$outPath);
createMapFile($outPath);

$outFile = $outPath . "/" . $fast_call_file_name;
$fd = fopen($outFile,"a+");
fwrite($fd,$callServiceFuncTotalStr);
//file_put_contents($outFile,$callServiceFuncTotalStr);


$callGrpcService = $template->CallGrpc();
$callGrpcService = str_replace("#case#",$callGrpcServiceCase,$callGrpcService);
fwrite($fd,$callGrpcService);

//var_dump($callServiceFuncTotalStr);exit;
exit(111);
//具体分析/处理一个文件中的内容
function parserOneServiceFileContent($outPath,$protoFilePath,$protoFileNamePrefix,$packName,$mapIdSeparate){
    //公共变量，最后外层函数可能得用
    global $template;
    global $serviceImplementPackage ;
    global $ServiceFastCallSwitchCase;

    //.proto文件名
    $protoFileName = $protoFilePath . "/". $protoFileNamePrefix.".proto";
    //输出路径
    $outFile = $outPath . "/" . $protoFileNamePrefix. ".go";
    //打开该文件，获取文件内容
    $protoFileContent = file_get_contents($protoFileName);
    //正则：获取包名
//    preg_match_all('/package(.*);/isU',$protoFileContent,$match);
//    $package = trim($match[1][0]);


    pp("parserOneServiceFileContent , protoFileName: $protoFileName , outFile:$outFile " );

    pp("\n\n\n");
    //读取一个文件中的，若干个service 块
    preg_match_all('/service(.*){(.*)\n}/isU',$protoFileContent,$match);

    if (count($match[0]) == 0){
        pp("parserOneServiceFileContent err : 正则未匹配到 service 块，即：该 proto 文件内没有任何函数.");
        return array();
    }
    //读取一个service中的所有rpc 函数名
    $service = null;
    foreach ($match[1] as $k =>$v){
        $serviceName = trim($v);
        $rpcFuncMatch = null;
        //正则：解析一个rpc 函数 的详细信息，名称 输入 输出 描述信息
        preg_match_all("/rpc(.*)\((.*)\)(.*)returns(.*)\((.*)\)(.*)\/\/(.*)\n/isU",$match[0][$k],$rpcFuncMatch);
        foreach ($rpcFuncMatch[0] as $k2=>$v2){
            $arr = array(
                'name'=>trim($rpcFuncMatch[1][$k2]),
                'in'=>trim($rpcFuncMatch[2][$k2]),
                'out'=>trim($rpcFuncMatch[5][$k2]),
                'desc'=>trim($rpcFuncMatch[7][$k2]),
            );
            pp("name:".$arr["name"] . " in:".$arr["in"]. " out:".$arr["out"]. " desc:".$arr["desc"]);

            $service[$serviceName][] = $arr;
        }
    }
    pp("一个 proto 文件，服务下的所有 函数-参数/结构体 分析出来了");
    pp("给每个函数生成 go 调用 代码，网关在反向调用时，使用");
    $s = $template->separator;
//    $serviceImplementPackage = "pbservice";
    //读取静态模板:一个新GO文件的，头部信息
    $go_file_content = $template->serviceImplementHeader();
    //开始替换模板中的动态变量
    $go_file_content = str_replace($s."serviceImplementPackage".$s,$serviceImplementPackage,$go_file_content);
    $go_file_content = str_replace($s."package".$s,$packName,$go_file_content);

    //生成： 每个服务的 grpc client，每个服务中包含的函数的调用代码
    foreach ($service as $serviceName=>$funs){
        //每个服务的具体实现，得有一个结构体，再将所有函数挂在该结构体
        $go_class_str = $template->serviceImplementStruct();
        $go_class_str = str_replace($s."service_class_name".$s,$serviceName,$go_class_str);
        $func_total_str = "\n";
        foreach ($funs as $k=>$v){
            $func_str = $template->serviceImplementFunc();
            $func_str = str_replace($s."class_name".$s,lcfirst($serviceName),$func_str);
            $func_str = str_replace($s."class_type".$s,$serviceName,$func_str);
            $func_str = str_replace($s."func_name".$s,$v['name'],$func_str);
            $func_str = str_replace($s."para_in_type".$s,$packName.".".$v['in'] ,$func_str);
            $func_str = str_replace($s."para_in_name".$s,lcfirst($v['in']),$func_str);
            $func_str = str_replace($s."return_type".$s,$packName.".".$v['out'],$func_str);
            $func_str = str_replace($s."return_name".$s,lcfirst($v['out']),$func_str);
            $func_total_str .= $func_str . "\n";
        }
        //最后将一个服务下的所有函数，统一替换进一个文件中
        $go_class_str = str_replace($s."funcs".$s,$func_total_str,$go_class_str);
        $go_file_content .=  $go_class_str . "\n" ;
    }
    //每个服务，生成一个文件
    file_put_contents($outFile,$go_file_content);



    //生成： 网关 反射 动态调用 服务 函数的快捷调用代码 + 函数映射ID表。 这里最终是统一生成 2个文件，外部处理。
    foreach ($service as $serviceName=>$info){
        dynamicCallGrpcService($serviceName,$info,$packName);

        $MountClientSwitchCaseStr = $template->MountClientSwitchCase();
        $MountClientSwitchCaseStr = str_replace("#service_name#",$serviceName,$MountClientSwitchCaseStr);
        $MountClientSwitchCaseStr = str_replace("#package_name#",$packName,$MountClientSwitchCaseStr);

        $GetClientStr = $template->GetClient();
        $GetClientStr = str_replace("#service_name#",$serviceName,$GetClientStr);
        $GetClientStr = str_replace("#package_name#",$packName,$GetClientStr);

        $ServiceFastCallSwitchCase[$serviceName] = array('MountClientSwitchCase'=>$MountClientSwitchCaseStr,"GetClient"=>$GetClientStr);

        mapFunctionId($serviceName,$info,$mapIdSeparate);
    }


    pp("parserOneServiceFileContent finish.");
    return $service;
}

//动态调用grpc 服务 方法
function dynamicCallGrpcService($serviceName, $serviceInfo,$packageName){
    global $template,$callServiceFuncTotalStr,$callGrpcServiceCase;

    $callServiceFuncStr = $template->CallServiceFunc();
    $callServiceFuncStr = str_replace("#service_name#",$serviceName,$callServiceFuncStr);

    $serviceFuncListStr = "";
    foreach ($serviceInfo as $k=>$service){
        $callServiceFuncCaseStr = $template->CallServiceFuncCase();
        $callServiceFuncCaseStr = str_replace("#service_name#",$serviceName,$callServiceFuncCaseStr);
        $callServiceFuncCaseStr = str_replace("#package_name#",$packageName,$callServiceFuncCaseStr);
        $callServiceFuncCaseStr = str_replace("#request#",$service["in"],$callServiceFuncCaseStr);
        $callServiceFuncCaseStr = str_replace("#func_name#",$service["name"],$callServiceFuncCaseStr);

        $serviceFuncListStr .= $callServiceFuncCaseStr . "\n";

    }
    $callServiceFuncStr = str_replace("#case#",$serviceFuncListStr,$callServiceFuncStr);
    $callServiceFuncTotalStr .= $callServiceFuncStr . "\n";

    $CallGrpcCase = $template->CallGrpcCase();
    $CallGrpcCase = str_replace("#service_name#",$serviceName,$CallGrpcCase);
    $callGrpcServiceCase .= $CallGrpcCase . "\n";
}

//编译proto文件，生成pb.go 文件，这个单独用shell执行再好一些
function compileProtoFile($compileCommand,$fileName){
    $compileCommandFile = str_replace("#proto_file_name#",$fileName,$compileCommand);
    pp("compileProtoFile: $compileCommandFile");
    $output = shell_exec($compileCommandFile);
    pp("output:$output");
}
//函数名映射ID
function mapFunctionId($serviceName,$serviceFuncListInfo,$mapIdSeparate){
    $mapFuncIdNo = "1";
    pp("mapFunctionId:$serviceName");
//    $serviceId = $GLOBALS["mapServiceIdNo"];
    global $db;
    $projectList = $db->GetProjectList();
    if (!$projectList || count($projectList) <= 0 ){
        exit("db GetProjectList empty~");
    }

//    var_dump($projectList);exit;
    $searchProject = null;
    foreach ($projectList as $k=>$v){
        if ($v['name'] == $serviceName){
            $searchProject = $v;
        }
    }

    if (!$searchProject){
        exit("searchProject empty~");
    }
    $serviceId = $searchProject['id'];

    if(strlen($serviceId) < 2 ){
        $serviceId = $serviceId . "0";
    }
    $mapFuncIdNo = 100;
    foreach ($serviceFuncListInfo as $k=>$funcInfo){
        $in = "empty";
        $out = "empty";
        if ($funcInfo["in"]){
            $in =  $funcInfo["in"];
        }
        if ($funcInfo["out"]){
            $out =  $funcInfo["out"];
        }

//         if(strlen($mapFuncIdNo) == 1 ){
//             $mapFuncIdNo = "00".$mapFuncIdNo;
//         }else if(strlen($mapFuncIdNo) == 2 ){
//             $mapFuncIdNo = "0".$mapFuncIdNo;
//         }
//    $outContent .= $GLOBALS["mapIdNo"]  . $mapIdSeparate .$serviceInfo["name"] . $mapIdSeparate .  $in . $mapIdSeparate .  $out . $mapIdSeparate. $arr['desc'] . "\n";
        $outContent = $serviceId . $mapFuncIdNo . $mapIdSeparate .$serviceName . $mapIdSeparate. $funcInfo["name"] . $mapIdSeparate .  $in . $mapIdSeparate .  $out . $mapIdSeparate. $funcInfo['desc'] . "\n";
//        $GLOBALS["mapFuncIdNo"] ++;$GLOBALS["mapFuncIdNo"] ++;//PB生成的全是奇数，用于表示：client请求server ，长连接时会有server请求client，这些都是偶数
        $mapFuncIdNo++;$mapFuncIdNo++;
        $GLOBALS["map"] .= $outContent;
    }

//    $GLOBALS["mapServiceIdNo"]++;
}

//创建 函数 映射 txt 表
function createMapFile($outPath){
    $outFile = $outPath . "/" . "map.txt";
    file_put_contents($outFile,$GLOBALS["map"]);
}
//grpc switch service
function createServiceFastCallSwitch($ServiceFastCallSwitchCase,$outPath){
    global $template;
    global $fast_call_file_name;

    $cases = "";
    $getClient = "";
    foreach ($ServiceFastCallSwitchCase as $k=>$v){
        $cases .= $v['MountClientSwitchCase'] . "\n";
        $getClient .= $v['GetClient'] . "\n";
    }
    $MountClientSwitchStr = $template->MountClientSwitch();
    $MountClientSwitchStr = str_replace("#switch_case#",$cases,$MountClientSwitchStr);

    $content = $template->fastCallFileHeader() .   "\n\n" . $getClient . "\n\n" . $MountClientSwitchStr;

    $outFile = $outPath . "/" . $fast_call_file_name;
    file_put_contents($outFile,$content . "\n");
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

function pp($info,$is_br = 1 ){
    if(DEBUG == 1){
        $end = "";
        if ($is_br){
            $end = "\n";
        }
        echo $info.$end;
    }
}
