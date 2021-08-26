<?php
//正则匹配一个service 块时，结尾必须是：} ，上一行必须是\n结束
define("DEBUG",1);

$protoFilePath = "proto";
$protoFileName = "api.proto";
$outPath = "./";
$outName = "map.txt";
$outContent = "";
$s = "|";

$outFile = $outPath . "/" . $outName;
$protoFileContent = file_get_contents($protoFilePath . "/".$protoFileName);
// pp($protoFileContent);

$match = null;


preg_match_all('/service(.*){(.*)\n}/isU',$protoFileContent,$match);
// pp($match);

if (count($match[0]) == 0){
    exit("no match any service block~");
}

$service = null;
$idNo = 1000;

foreach ($match[1] as $k =>$v){
    $serviceName = trim($v);
    $rpcFuncMatch = null;
//     var_dump($match[0][$k]);
    preg_match_all("/rpc(.*)\((.*)\)(.*)returns(.*)\((.*)\)(.*)\/\/(.*)\n/isU",$match[0][$k],$rpcFuncMatch);
//     var_dump($rpcFuncMatch);exit;
    foreach ($rpcFuncMatch[0] as $k2=>$v2){
//     var_dump($rpcFuncMatch[0]);
        $arr = array(
            'name'=>trim($rpcFuncMatch[1][$k2]),
            'in'=>trim($rpcFuncMatch[2][$k2]),
            'out'=>trim($rpcFuncMatch[5][$k2]),
            'desc'=>trim($rpcFuncMatch[7][$k2]),
        );

        $service[$serviceName][] = $arr;

        $in = "empty";
        $out = "empty";
        if ($arr["in"]){
            $in =  $arr["in"];
        }
        if ($arr["out"]){
            $out =  $arr["out"];
        }

        $outContent .= $idNo . $s .$arr["name"] . $s .  $in . $s .  $out . $s. $arr['desc'] . "\n";
        $idNo++;
    }

}
echo $outContent;

echo "file_put_contents:".$outFile;

file_put_contents($outFile,$outContent);

function getNewIdNo(){
    static $idNo;
    return $idNo++;
}






