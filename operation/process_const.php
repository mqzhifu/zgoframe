<?php

if (count($argv) < 2){
    exit("至少有一个参数");
}

$project_base_dir = trim($argv[1]);//项目根目录
if(!is_dir($project_base_dir)){
    exit("project_base_dir 不是一个目录：".$project_base_dir);
}

debug("project_base_dir:$project_base_dir");

class ConstProcess{
    public $process_file_list = array(
        "core"=>"core/const.go",
        "service"=>"service/const.go",
        "model"=>"model/const.go",
        "util"=>"util/const.go",
    );
    public $process_file_content_list = array();
    public $project_base_dir = "";
    public $parse_rs = array();

    function ConstProcess($project_base_dir){
        $this->project_base_dir = $project_base_dir;
        $this->init();
        $this->show();
    }

    function init(){
        foreach ($this->process_file_list as $k=>$v){
            $file_path = $this->project_base_dir . "/" . $v;
            debug("check :".$file_path);
            if(!is_file($file_path)){
                exit("文件不存在：".$file_path);
            }
            $file_content = trim(file_get_contents ($file_path));
            if(!$file_content){
                exit("文件内容为空：".$file_path);
            }
            $this->process_file_content_list[$k] = $file_content;
        }
        debug("check finish,start process:");
        foreach ($this->process_file_content_list as $k=>$v){
            $this->process_one($k);
        }
    }
    function show(){
        debug("show :");
        if(count($this->parse_rs) <= 0){
            var_dump("empty,please check");exit;
        }

        foreach ($this->parse_rs as $module_name=>$file_module){
            debug("parse_rs module_name:".$module_name);
           foreach ($file_module as $k2=>$v2){
               debug("    ".$v2['desc']. " " . $v2['common_prefix']);
               foreach ($v2["const"] as $k3=>$const){
//                   foreach ($const as $k3=>$row){
                       debug("        ".$const['key']." ".$const['value']." ".$const["desc"]);
//                   }
               }
           }

        }
    }
    function process_one($file_key){
        debug("process_one:".$file_key);
        $file_content = $this->process_file_content_list[$file_key];
        preg_match_all('/\/\/@parse (.*)\nconst \(\n(.*)\)/isU',$file_content,$match);
        if(count($match) <= 0 ){
            return false;
        }

        $list = array();
        foreach ($match[1] as $k=>$const_name){
            $row = array(
                "desc"=>$const_name,
                "const"=>array(),
                "common_prefix" =>"",
            );

            $const_list_text = trim($match[2][$k]);
            $const_list_arr = explode("\n",$const_list_text);
            foreach ($const_list_arr as $k=>$line_text){
                $line_arr = explode("//", trim($line_text));
                $desc = $line_arr[1];

                $expression = explode("=",$line_arr[0]);
                $key = trim($expression[0]);
                $value = trim($expression[1]);

                $const_line = array("desc"=>$desc,"key"=>$key,"value"=>$value,"length"=>strlen($key));
                $row["const"][] = $const_line;
            }
            if(count($row["const"]) <= 1){
                continue;
            }
            //寻找公共前缀
            //先对齐，查找常量名最短的那个
            $minLength = 9999;
            foreach ($row["const"] as $k=>$v){
//                var_dump($v);
                if($v["length"] < $minLength){
                    $minLength = $v['length'];
                }
            }
            $commPrefix = "";
            for($i=0;$i<$minLength;$i++){
                $oneChar = $row["const"][0]["key"][$i];
                $allCharSum = 0;
                //累加每个字符串的第N个字符
                foreach ($row["const"] as $k=>$v){
                    $allCharSum  += ord($v["key"][$i]);
                }
                $oneCharSum = ord($oneChar) * count($row["const"]);
                if($oneCharSum != $allCharSum){
                    break;
                }
                $commPrefix .= $oneChar;
            }
            $row["common_prefix"] = $commPrefix;
            $list[] = $row;
        }

        $this->parse_rs[$file_key] = $list;
    }

}

new ConstProcess($project_base_dir);


function debug($str,$n = "\n"){
    echo $str . $n;
}

