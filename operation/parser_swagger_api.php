<?php
class ParserSwaggerApi{
    public $swaggerFilePath = "";
    public $outDir = "";
    public $template = null;
    public $yamlObject = null;
    public $tmpLoopForeachObj = null;
    public function __construct($swaggerFilePath,$outDir,$templateType)
    {
        $this->swaggerFilePath = $swaggerFilePath;
        $this->outDir = $outDir;
        $this->template = new Template($templateType);
    }
    //开始分析
    function Start(){
        $this->out("swaggerFilePath:".$this->swaggerFilePath);
        //读取 yaml 文件，并转化成一个对象
        $content = file_get_contents($this->swaggerFilePath);
        $yaml = yaml_parse($content);
        $this->yamlObject = $yaml;

        //1. 处理 函数
        $funcTemplateCodeByTags = [];//保存所有函数的处理结果集
        foreach ($yaml["paths"] as $path=>$pathOne){
            $this->out("path:".$path);
            foreach ($pathOne as $method=>$row){
                $method = strtoupper($method);//转大写
                $this->out("method:".$method);
                $summary = $this->checkUnset( $row,"summary");
                $description = $this->checkUnset( $row,"description");
                $parameters = $this->checkUnset( $row,"parameters");
                $responses  = $this->checkUnset( $row,"responses");
                $produces = $this->checkUnset( $row,"responses");
                $tags = $this->checkUnset($row,"tags");
                if(!$parameters){
                    $this->out("warning ,parameters empty.");
                    continue;
                }
                //处理请求参数
                $funcTemplate = $this->ParserParameters($parameters,$path,$method,$description);
                $funcTemplateCodeByTags[$tags[0]][] = $funcTemplate;
            }
        }
        //2. 处理 class ，将函数代码 替换到类中
        foreach ($funcTemplateCodeByTags as $tagName=>$functions){
            $TemplateClass = $this->template->GetClass();
            $TemplateClass = str_replace("#class#",$tagName,$TemplateClass);
            $functionsCode = "";
            foreach ($functions as $k=>$code){
                $functionsCode .= $code;
            }
            $TemplateClass = str_replace("#functions#",$functionsCode,$TemplateClass);
            $fileName = $tagName . ".js";
            $path = $this->outDir. "/".$fileName;
            $fd = fopen($path,"w+");
            fwrite($fd,$TemplateClass);
        }
    }

    function ParserParameters($parameters,$path,$method,$description){
        $inType = "";
        $paraList = [];
        foreach ($parameters as $k=>$paraOne){
            $paraDescription = $this->checkUnset( $paraOne,"description");
            $name = $this->checkUnset( $paraOne,"name");
            $required = $this->checkUnset( $paraOne,"required");
            $valueType = $this->checkUnset( $paraOne,"type");
            $in = $this->checkUnset( $paraOne,"in");

            switch ($in){
                case "header":
                    $inType = "header";
                    break;
                case "body":
                    $schema = $paraOne["schema"]['$ref'];
                    $this->out("schema:$schema");
                    $obj = $this->GetDefinitionsObj($schema);
//                    var_dump($obj);exit;
                    $realObj = [];
                    $parserNewObj = $this->LoopForeachObj($obj,1);
                    $inType = "body";
//                    $paraList[] = $obj;
                    $paraList = $parserNewObj;//这里没有用数组存，因为一但 body 里传 JSON，只能有一个大的 HTTP BODY参数
                    break;
                case "path":
                    $paraList[] = $paraOne;
                    $inType = "path";
                    break;
                case "formData":
                    $inType = "formData";
                    $paraList[] = $paraOne;
                    break;
                default:
                    $this->out("err ,Parameters <in> value err:"+$in);
                    break;
            }
//            var_dump($paraOne);exit;
//            if(row.parameters[para].in == "header"){
        }
        if(count($paraList) <= 0){
            return "";
        }
        $FunctionTemplateStr = $this->MakeOneFunctionTemplate($path,$method,$description,$inType,$paraList);
        return $FunctionTemplateStr;
    }
    function MakeOneFunctionTemplate($path,$method,$description,$type,$paraList){
        $functionName = $this->ParserFuncNameByPath($path);//解析函数名
        $functionDataBody = "";
        switch ($type){
            case "formData":
            case  "path":
                $obj = [];
                foreach ($paraList as $k=>$v){
                    $obj[$v['name']] = "";
                }
                $functionDataBody = json_encode($obj);
                break;
            case  "body":
                $functionDataBody = json_encode($paraList);
                break;
            default:
                $this->out("MakeOneFunctionTemplate type value err:".$type);
                break;
        }

        $functionTemplate = $this->template->GetFunction();
        $functionTemplate = str_replace("#uri#",$path,$functionTemplate);
        $functionTemplate = str_replace("#method#",$method,$functionTemplate);
        $functionTemplate = str_replace("#body#",$functionDataBody,$functionTemplate);
        $functionTemplate = str_replace("#desc#",$description,$functionTemplate);
        $functionTemplate = str_replace("#Funcname#",$functionName,$functionTemplate);

        return $functionTemplate;
    }
    function ParserFuncNameByPath($path){
        $functionName = "";
        $pathArr  = explode("/",$path);
        foreach ($pathArr as $k=>$v){
            if(!$v){
                continue;
            }

            if(substr($v,0,1) == "/"|| substr($v,0,1) == "{"){
                continue;
            }

            //首字母转大写
            $functionName .= ucfirst($v);
        }

        return $functionName;
    }
    function GetDefinitionsObj($schema){
        $className = str_replace("#/definitions/","",$schema);
        $obj = $this->yamlObject["definitions"][$className];
//        if($className == "request.HeaderRequest"){
//            var_dump($obj);exit;
//        }
        return $obj;
    }
    function checkUnset($arr,$key){
        if(isset($arr[$key])){
            return $arr[$key];
        }
        return "";

    }
    function out($str,$ln = 1){
        if($ln){
            $str .= "\n";
        }
        echo $str;
    }

    function LoopForeachObj($obj,$root){
        $list =[];
        $type = $this->checkUnset($obj,"type");
        if($root){
            if( $type == "object") {
                if(!$this->checkUnset($obj,"properties")){//空结构体
                    return [];
                }
                $list = $this->LoopForeachObj($obj["properties"],0);
            }else{
                var_dump("根对象如果不是 object 接口定义就有问题的");exit;
            }
        }else{
            //没有 type 属性，又不是一个数组
            if(!$type && count($obj) <= 0){
                var_dump($obj);
                var_dump('xxxx');
                exit;
            }else{
//                $row = [];
                foreach ($obj as $k=>$v){
                    //单个对象
                    if(is_array($v) && $this->checkUnset($v,'$ref')){
                        $obj = $this->GetDefinitionsObj($v['$ref']);
                        $list[$k] = $this->LoopForeachObj($obj,1);
                        continue;
                    }
                    //数组里是一个对象
                    if(is_array($v) && $this->checkUnset($v ,"items") && $this->checkUnset($v["items"] , '$ref')){
                        $obj = $this->GetDefinitionsObj($v["items"]['$ref']);
                        $list[$k][] = $this->LoopForeachObj($obj,1);
                        continue;
                    }
                    switch ($v['type']){
                        case "boolean":
                            $list[$k] = false;
                            break;
                        case "string":
                            $list[$k] = "";
                            break;
                        case "integer":
                            $list[$k] = 0;
                            break;
                        case "object":
                            $this->out("in object:");
                            $additionalProperties = $this->checkUnset($v,'additionalProperties');
                            if($additionalProperties){//这是一个自定义的MAP类型的结构
                                $list[$k] = array("a"=>"bbbb");
                            }else{
                                //这里没写完，有个BUG
//                                $row[$k] = $this->LoopForeachObj($obj[$k],0);
                            }
                            break;
                        case "array":
                            $this->out("in array , k:$k");
                            if($obj[$k]["items"]["type"] == "string"){
                                $array = ["aaaa","bbbbb"];
                            }elseif($obj[$k]["items"]["type"] == "integer"){
                                $array = [0,1];
                            }else{
                                var_dump("00998877");
                                var_dump($obj);exit;
                            }
                            var_dump($obj[$k]);
                            $list[$k] = $array;
                            break;
                        default:
                            var_dump($v);exit;
                            $this->out("err err err:".$type);
                            var_dump($obj);exit;

                    }
                }
            }
        }
        return $list;



//        var_dump($type);
        if( $type == "object") {
            $row[] = $this->LoopForeachObj($obj["properties"]);
        }elseif( $type == "array"){
            if(1){

            }
            $row[] = $this->LoopForeachObj($obj["properties"]);
        }else{

            return $row;
        }
        return $row;
//        foreach ($obj["properties"] as $objKey=>$objValue){
//            $objRow = [];
//            if($objValue['type'] == "string"){
//                $objRow[$objKey] = "";
//            }elseif($objValue['type'] == "integer"){
//                $objRow[$objKey] = 0;
//            }elseif($objValue['type'] == "object"){
//                var_dump($obj);
//                exit("aaaa");
//            }elseif($objValue['type'] == "array"){
//                var_dump($objValue);
//            }else{
//                var_dump($objValue);exit;
//                exit("bbbbb");
//            }
//        }


    }
}

class Template{
    public $_type = 0;
    public function __construct($_type)
    {
        $this->_type = $_type;
    }
    function GetClass(){
        $code = <<<EOF
class #class#{
    #functions#
}
EOF;
        return $code;
    }
    function GetFunction(){
        $code = <<<EOF
//#desc#
    #Funcname#(obj,callback){
        let uri = "#uri#";
        let method = "#method#";
        //let loginData = #body#;
        this.HttpRequest.request(this.CommonCallback.bind(this),uri,this.token,false,method,obj,"");
    }
    
EOF;
        return $code;
    }
}


$path = "/data/golang/zgoframe/docs/swagger.yaml";
$outDir = "/data/php/zhongyuhuacai/storage";
$ParserSwaggerApiClass = new ParserSwaggerApi($path,$outDir,"js");
$ParserSwaggerApiClass->Start();