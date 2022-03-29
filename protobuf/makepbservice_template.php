<?php
class Template
{
    public $separator = "#";

    function fastCallFileHeader()
    {
        $str = <<<EOF
package util

import (
    "encoding/json"
    "errors"
    "google.golang.org/grpc"
    "zgoframe/protobuf/pb"
    "context"
)

EOF;
        return $str;
    }


    function serviceImplementHeader()
    {
        $str = <<<EOF
package #serviceImplementPackage#

import (
	"context"
	"zgoframe/protobuf/#package#"
)

EOF;
        return $str;
    }

    function serviceImplementStruct()
    {
        $go_class_str = <<<EOF
type #service_class_name# struct{}

#funcs#
EOF;
        return $go_class_str;
    }

//一个GRPC服务的，具体实现函数体
    function serviceImplementFunc()
    {
        $go_func = <<<EOF
func (#class_name# *#class_type#)#func_name#(ctx context.Context,#para_in_name# *#para_in_type#) (*#return_type#,error){
    #return_name# := &#return_type#{}
    return #return_name#,nil
}
EOF;
        return $go_func;
    }

    function MountClientSwitchCase()
    {
        $str = <<<EOF
    case "#service_name#":
        incClient = #package_name#.New#service_name#Client(myGrpcClient.ClientConn)
EOF;
        return $str;
    }

    function MountClientSwitch()
    {
        $str = <<<EOF
//根据服务名获取一个GRPC-CLIENT 连接(c端使用)
func  (myGrpcClient *MyGrpcClient)GetGrpcClientByServiceName(serviceName string,clientConn *grpc.ClientConn)(interface{},error){
    var incClient interface{}
    switch serviceName {
#switch_case#
    default:
        return incClient,errors.New("service name router failed.")
    }
    return incClient,nil
}
EOF;
        return $str;
    }

    function GetClient()
    {
        $str = <<<EOF
//获取一个服务的grpc client : #service_name#
func (grpcManager *GrpcManager)Get#service_name#Client(name string,balanceFactor string)(#package_name#.#service_name#Client,error){
    client, err := grpcManager.GetClientByLoadBalance(name,balanceFactor)
    if err != nil{
        return nil,err
    }

    return client.(pb.#service_name#Client),nil
}
EOF;
        return $str;
    }

    function CallGrpcCase(){
        $str = <<<EOF
    case "#service_name#":
        resData , err = grpcManager.CallServiceFunc#service_name#(funcName,balanceFactor,requestData)
EOF;
        return $str;
    }

    function CallGrpc(){
        $str = <<<EOF
//动态调用一个GRPC-SERVER 的一个方法(c端使用)
func (grpcManager *GrpcManager) CallGrpc(serviceName string,funcName string,balanceFactor string,requestData []byte)( resData interface{},err error){
	switch serviceName {
#case#
    default:
            return requestData,errors.New("service name router failed.")
	}

	return resData,err
}
EOF;
        return $str;
    }


    function CallServiceFuncCase(){
        $str = <<<EOF
    case "#func_name#":
        request := #package_name#.#request#{}
        err := json.Unmarshal(postData,&request)
        if err != nil{
            return data,err
        }
        data ,err = grpcClient.#func_name#(ctx,&request)
EOF;
        return $str;
    }

    function CallServiceFunc(){
        $str = <<<EOF
//动态调用服务的函数 : #service_name#
func (grpcManager *GrpcManager) CallServiceFunc#service_name#(funcName string,balanceFactor string,postData []byte)( data interface{},err error){
    //获取GRPC一个连接
    grpcClient,err := grpcManager.Get#service_name#Client("#service_name#",balanceFactor)
    if err != nil{
        return data,err
    }

    ctx := context.Background()
    switch funcName {
#case#
    default:
            return data,errors.New("func name router failed.")
    }
    return data,err
}
EOF;
        return $str;
    }
}

class DB{
    function getConn(){
        $host = "8.142.177.235";
        $username = "root";
        $ps = "mqzhifu";
        $db = "test";
        $conn = mysqli_connect($host,$username,$ps,$db);
//        var_dump(mysqli_connect_error());exit;
//        var_dump($conn);exit;
        return $conn;
    }
    function GetProjectList(){
        $conn = $this->getConn();
        $sql = "select * from project";
        $resultQuery = mysqli_query($conn,$sql);
        if (!$resultQuery){
            return array();
        }

        $result = array();
        while($row = mysqli_fetch_assoc($resultQuery)){
            $result[]   =   $row;
        }
        return $result;
    }
}

function CurlGetProjectListInfo(){
//     curl -X POST "http://127.0.0.1:1111/tools/project/list" -H "accept: application/json" -H "X-Source-Type: 11" -H "X-Project-Id: 6" -H "X-Access: imzgoframe"
    $url = "http://127.0.0.1:1111/tools/project/list";
    $header = array("X-Source-Type:11","X-Project-Id:6","X-Access:imzgoframe");
// var_dump($header);

    $ch = curl_init($url);
    curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, FALSE);
    curl_setopt($ch, CURLOPT_SSL_VERIFYHOST, FALSE);
    curl_setopt($ch, CURLOPT_POST, true);
    curl_setopt($ch, CURLOPT_HTTPHEADER, $header);
//     curl_setopt($ch, CURLOPT_POSTFIELDS, $payload);
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);
//     curl_setopt($ch, CURLINFO_HEADER_OUT, true);

    $output = curl_exec($ch);
    curl_close($ch);

    $jsonData  = json_decode($output);
    return $jsonData["data"];


}
