<?php
$log_base_path = "/data/www/golang/src/logs";

out( "log_base_path:$log_base_path");

include "filebeat.php";

class Op {
    public $host = "http://127.0.0.1:5555/";
    public $esHost = "59.110.167.206:9200";

    public $appTypeList = null;
    public $appList = null;
    public $app = null;
    public $projectName = "";
    public $log_base_path = "";

    //总控
    function start($log_base_path){
        $this->log_base_path = $log_base_path;
        $this->initApp();
        $this->initProjectLogDir();
        $this->initFilebeat();
    }

    function initApp(){
        if(!is_dir($this->log_base_path)){
            out("err: log_base_path not dir~ ".$this->log_base_path);
        }

        $uri = $this->host."base/appTypeList";
        $postData = array("status"=>1);
        $httpGetAppType = curl_send($uri,2,$postData,false);
        if($httpGetAppType['code'] != 200){
            out("err: httpGetAppType ~ ".$httpGetAppType['msg']);
        }

        if(count($httpGetAppType['data']) <= 0){
            out("err: httpGetAppType count <= 0 ");
        }

        $appTypeList = $httpGetAppType['data'];
        foreach ($appTypeList as $id=>$typeName){
            $appTypePath = $this->log_base_path . "/" . $typeName;
            if(!is_dir($appTypePath)){
                out("appTypePath not exist~ ".$appTypePath + " , mkdir");
                mkdir($appTypePath);
            }else{
                out("appTypePath exist~ ".$appTypePath);
            }
        }

        $scriptPath = dirname(__FILE__);
        out("scriptPath:".$scriptPath);
        $explodeScriptPath = explode("/",$scriptPath);
    //     var_dump($explodeScriptPath);
        $projectName = $explodeScriptPath[count($explodeScriptPath) - 2];//实际就是脚本 的上一级目录名
        $this->projectName = $projectName;
        out("projectName:".$projectName);


        $uri = $this->host."base/appList";
        $httpGetAppList = curl_send($uri,2,$postData,false);
        if($httpGetAppList['code'] != 200){
            out("err: httpGetAppList ~ ".$httpGetAppList['msg']);
        }
        $AppList = $httpGetAppList['data'];
        $myApp = null;
        foreach ($AppList as $id=>$app){
            if ($app['key'] == $projectName){
                $myApp = $app;
            }
        }

        if(!$myApp){
            out("projectName not exist");
        }
        $this->app = $myApp;
        $this->appTypeList = $appTypeList;
        $this->appList = $AppList;
    }

    function initProjectLogDir(){
        out("initProjectLogDir:");
        foreach ($this->appList as $k=>$v){
            $projectType = $this->appTypeList[$v['type']];
            $projectLogPath = $this->log_base_path . "/" . $projectType . "/".$v['key'];
            out("projectLogPath:".$projectLogPath);
//            if(!is_dir($projectLogPath)){
//                out("projectLogPath not exist~ ".$projectLogPath + " , mkdir");
//                mkdir($projectLogPath);
//            }else{
//                out("projectLogPath exist~ ".$projectLogPath);
//            }
        }
    }

    function initFilebeat(){
        $input_str = "";
        $index_str = "";
        foreach ($this->appTypeList as $k=>$v){
            $block = getBlock();
            $path = $this->log_base_path . "/" . $v . "/*.log" ;
            out($path);
            $block = str_replace("#paths#",$path,$block);
            $block = str_replace("#source#",$v,$block);

            $input_str .= $block;

            $index = getIndex();
            $index = str_replace("#index#",$v,$index);
            $index_str .= $index;
        }
        out($input_str);
        out($index_str);
    }
}

function getBlock(){
$str = <<<EOF
- type: log
  enabled: true
  paths:
    - #paths#

  json.keys_under_root: true
  json.add_error_key: true

  fields:
    source: #source#

EOF;
    return $str;
}

function getIndex(){
$str = <<<EOF
    - index: "ck-local-#index#-%{+yyyy.MM.dd}"
      when.equals:
        fields:
            source: "#index#"

EOF;
    return $str;
}

$class = new Op();
$class->start($log_base_path);

function out($msg,$ln = 1){
    if ($ln ){
        $msg .= "\n";
    }
    echo $msg;
}



function curl_send($url,$type = 1,$postData = null,$ssl = null ){
    $curl = curl_init();
    //设置URL
    curl_setopt($curl, CURLOPT_URL, $url);
    curl_setopt($curl, CURLOPT_USERAGENT, 'Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.1; Trident/6.0)');
    curl_setopt($curl, CURLOPT_REFERER, "https://project.shell.init.php");
    //设置获取的信息以文件流的形式返回，而不是直接输出。
    curl_setopt($curl, CURLOPT_RETURNTRANSFER, 1);
    if($type == 2){
        //设置为POST
        curl_setopt($curl, CURLOPT_POST, 1);
        //设置POST数据
        curl_setopt($curl, CURLOPT_POSTFIELDS, http_build_query($postData) );
    }

    if($ssl == 1){
        curl_setopt($curl, CURLOPT_SSL_VERIFYPEER, false);
        curl_setopt($curl, CURLOPT_SSL_VERIFYHOST, false);
    }


    $data = curl_exec($curl);
    if (curl_errno($curl)) {
        $data = array('code'=>400,'msg'=>curl_errno($curl).":".curl_error($curl));
        return data;
    }
    curl_close($curl);

    $data = json_decode($data,true);
    return $data;
}
