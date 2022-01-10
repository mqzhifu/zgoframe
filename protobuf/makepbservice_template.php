<?php
class Template
{
    public $separator = "#";
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
    func  (myGrpcClient *MyGrpcClient)MountClientToConnect(serviceName string){
        var incClient interface{}
        switch serviceName {
            #switch_case#
        }
	}
EOF;
        return $str;
    }

    function GetClient()
    {
        $str = <<<EOF
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
}

class DB{
    function getConn(){
        $host = "localhost";
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


