<?php
//正则匹配一个service 块时，结尾必须是：} ，上一行必须是\n结束
define("DEBUG",0);


$argcNoticeMsg = "servicePackage=xxxx protoFilePath=xxx protoFileName=xxxx outPath=xxx pbClassName=xxxx";
if (count($argv) < 6){
    exit($argcNoticeMsg);
}

// $servicePackage = "pbservice";
// $protoFilePath = "proto";
// $protoFileName = "demo.proto";
// $outPath = "$servicePackage";
// $pbClassName = "pb.";
$servicePackage = $argv[1];
$protoFilePath = $argv[2];
$protoFileName = $argv[3];
$outPath = "$servicePackage";
$pbClassName = $argv[5];



$outFile = $outPath . "/" . $protoFileName. ".go";
$protoFileContent = file_get_contents($protoFilePath . "/".$protoFileName);
pp($protoFileContent);

$match = null;


preg_match_all('/package(.*);/isU',$protoFileContent,$match);
$package = trim($match[1][0]);

pp("\n\n\n\n\n");

preg_match_all('/service(.*){(.*)\n}/isU',$protoFileContent,$match);
pp($match);

if (count($match[0]) == 0){
    exit("no match any service block~");
}

$service = null;
foreach ($match[1] as $k =>$v){
    $serviceName = trim($v);
    $rpcFuncMatch = null;
    preg_match_all("/rpc(.*)\((.*)\)(.*)returns(.*)\((.*)\)/isU",$match[0][$k],$rpcFuncMatch);
    foreach ($rpcFuncMatch[0] as $k2=>$v2){
        $arr = array(
            'name'=>trim($rpcFuncMatch[1][$k2]),
             'in'=>trim($rpcFuncMatch[2][$k2]),
            'out'=>trim($rpcFuncMatch[5][$k2]),
        );
        $service[$serviceName][] = $arr;
    }
//    var_dump($service);exit;
//    exit(1111);

}
$s = "#";


$go_file_content = get_new_go_file();

$go_file_content = str_replace($s."package".$s,$servicePackage,$go_file_content);



foreach ($service as $serviceName=>$funs){
    $go_class_str = go_class_str();
    $go_class_str = str_replace($s."service_class_name".$s,$serviceName,$go_class_str);
    $func_total_str = "\n";
    foreach ($funs as $k=>$v){
        $func_str = get_class_func_str();
        $func_str = str_replace($s."class_name".$s,lcfirst($serviceName),$func_str);
        $func_str = str_replace($s."class_type".$s,$serviceName,$func_str);
        $func_str = str_replace($s."func_name".$s,$v['name'],$func_str);
        $func_str = str_replace($s."para_in_type".$s,$pbClassName.$v['in'],$func_str);
        $func_str = str_replace($s."para_in_name".$s,lcfirst($v['in']),$func_str);
        $func_str = str_replace($s."return_type".$s,$pbClassName.$v['out'],$func_str);
        $func_str = str_replace($s."return_name".$s,lcfirst($v['out']),$func_str);
        $func_total_str .= $func_str . "\n";
    }
    $go_class_str = str_replace($s."funcs".$s,$func_total_str,$go_class_str);
    $go_file_content .=  $go_class_str . "\n" ;
}



file_put_contents($outFile,$go_file_content);

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

function get_class_func_str(){
$go_func =<<<EOF
func (#class_name# *#class_type#)#func_name#(ctx context.Context,#para_in_name# *#para_in_type#) (*#return_type#,error){
    #return_name# := &#return_type#{}
    return #return_name#,nil
}
EOF;
    return $go_func;
}




function pp($info){
    if(DEBUG == 1){
        var_dump($info);
    }
}
