<?php
//golang里的常量是没法反射动态处理的，或者说处理起来很麻烦
//这里直接用PHP分析golang代码中的常量，动态生成枚举类型的常量，给后台使用，也可以当做说明文档使用

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
    public $golang_target_const_file = "util/const_handle.go";
    public $process_file_content_list = array();
    public $project_base_dir = "";
    public $parse_rs = array();

    function ConstProcess($project_base_dir){
        $this->project_base_dir = $project_base_dir;
        $this->init();
//        $this->show();
        $this->makeGolangCode();
    }

    function init(){
        $golang_const_file_path =$this->project_base_dir . "/" . $this->golang_target_const_file;
        if(!is_file($golang_const_file_path)){
            exit("golang_const_file_path err.".$this->golang_target_const_file);
        }

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
    //生成GOLANG代码，这个才是最终有用的处理
    function makeGolangCode(){
        debug("makeGolangCode:");
        $golangCode = $this->getGolangConstHandleTemplate();
        $golangEnumConstText = "";
        foreach ($this->parse_rs as $module_name=>$file_module){
            debug("parse_rs module_name:".$module_name);
            foreach ($file_module as $k2=>$v2){
//                debug("    ".$v2['desc']. " " . $v2['common_prefix']);
                $GolangEnumConstTemplate = $this->getGolangEnumConstTemplate();
                $GolangEnumConstTemplate = str_replace("#CommonPrefix#",$v2['common_prefix'],$GolangEnumConstTemplate);
                $GolangEnumConstTemplate = str_replace("#desc#",$v2['desc'],$GolangEnumConstTemplate);
                $GolangEnumConstTemplate = str_replace("#type#",'"'.$v2['type'].'"',$GolangEnumConstTemplate);
                $GolangConstItemTemplateText = "";
                foreach ($v2["const"] as $k3=>$const){
                    $GolangConstItemTemplate = $this->getGolangConstItemTemplate();
                    $GolangConstItemTemplate = str_replace("#key#",$const['key'],$GolangConstItemTemplate);
                    $GolangConstItemTemplate = str_replace("#value#",$const['value'],$GolangConstItemTemplate);
                    $GolangConstItemTemplate = str_replace("#constItemDesc#",$const['desc'],$GolangConstItemTemplate);
                    $GolangConstItemTemplateText .= $GolangConstItemTemplate . "\n\n";
                }


                $GolangEnumConstTemplate = str_replace("#ConstItem#",$GolangConstItemTemplateText,$GolangEnumConstTemplate);
                $golangEnumConstText .= $GolangEnumConstTemplate . "\n\n";
            }

        }

//        var_dump($golangEnumConstText);exit;
        $golangCode = str_replace("#EnumConst#",$golangEnumConstText,$golangCode);
        $fd = fopen($this->project_base_dir . "/" . $this->golang_target_const_file,"w");
        fwrite($fd,$golangCode);
    }

    function getGolangConstHandleTemplate(){
        $code = <<<EOF
package util

//注：此文件是由PHP动态生成，不要做任何修改，均会被覆盖

type ConstItem struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
	Desc  string      `json:"desc"`
}

type EnumConst struct {
	CommonPrefix string      `json:"common_prefix"`
	Desc         string      `json:"desc"`
	ConstList    []ConstItem `json:"const_list"`
	Type         string      `json:"type"`
}

type ConstHandle struct {
	EnumConstPool map[string]EnumConst
}

func NewConstHandle() *ConstHandle {
	constHandle := new(ConstHandle)
	constHandle.EnumConstPool = make(map[string]EnumConst)
	constHandle.Init()
	return constHandle
}

func (constHandle *ConstHandle) Init() {
	var constItemList []ConstItem
	var constItem ConstItem
	var enumConst EnumConst
    #EnumConst#
}

EOF;
        return $code;
    }
    function getGolangConstItemTemplate(){
        $code = <<<EOF
	constItem = ConstItem{
		Key:   "#key#",
		Value: #value#,
		Desc:  "#constItemDesc#",
	}
	constItemList = append(constItemList, constItem)
EOF;
    return $code;
    }
    function getGolangEnumConstTemplate(){
        $code = <<<EOF
    #ConstItem#
	enumConst = EnumConst{
		CommonPrefix: "#CommonPrefix#",
		Desc:         "#desc#",
		ConstList:    constItemList,
		Type:          #type#,
	}

	constHandle.EnumConstPool[enumConst.CommonPrefix] = enumConst
	constItemList = []ConstItem{}
EOF;
        return $code;

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
        //常量不要出现：括号关键字，会影响正则匹配
        debug("process_one:".$file_key);
        $file_content = $this->process_file_content_list[$file_key];
        preg_match_all('/\/\/@parse (.*)\nconst \(\n(.*)\)/isU',$file_content,$match);
        if(count($match) <= 0 ){
            return false;
        }

        $list = array();
        foreach ($match[1] as $k=>$const_name){
            $type = "int";
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
                if (substr($value,0,1) == '"' || substr($value,0,1) == "'"){
//                    $value = substr($value,1,strlen($value)-2);
                    $type = "string";
                }else{
                    $value = (int)$value;
                }
//                var_dump($value);
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
            $row["type"] = $type;
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

