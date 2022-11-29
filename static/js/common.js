// 后端公共HTTP接口-头信息
var header_X_Source_Type = "11";
var header_X_Project_Id = "6";
var header_X_Access = "imzgoframe";
//长连接 通信协议 基础类型 定义
var content_type_desc = {1:"json",2:"protobuf"};
var protocol_type_desc = {1:"tcp",2:"websocket",3:"udp"};
var CONN_STATUS_DESC = {1:"初始化",2:"运行中",3:"已关闭"};
//后端相关
//域名
// var domain = "127.0.0.1:1111";
var domain = window.location.host;
//http 协议
var http_protocol = "http";
if(location.href.substring(0,5) == "https"){
    http_protocol = "https";
}
var URI_MAP = {
    "gateway_config": http_protocol + "://"+domain + "/gateway/config",
    "gateway_action_map":http_protocol + "://"+domain + "/gateway/action/map",//URI - 网关 protobuf 映射表
    "user_login":http_protocol + "://"+domain + "/base/login",//登陆
    "gateway_fd_list":http_protocol + "://"+domain + "/gateway/fd/list",//长连接用户
    "gateway_send_msg":http_protocol + "://"+domain + "/gateway/send/msg",//长连接推送消息
    "gateway_total":http_protocol + "://"+domain + "/gateway/total",//长连接 metrics
    "twin_agora_config":http_protocol + "://"+domain + "/twin/agora/config",//声网房间长连接配置
    "rule":http_protocol + "://"+domain + "/game/match/rule",
};

var UserList = {
    10 :{id:10 ,"username"  :"doctor", "password":"123456","info":{},"token":"","roomId":"","channel":"ckck"},
     9 :{id:9  ,"username":"calluser", "password":"123456","info":{},"token":"","roomId":"","channel":"ckck"},
}


function get_uri_by_name(name){
    return URI_MAP[name];
}

function UserLogin(uid,callback){
    var info = UserList[uid];
    console.log("UserLogin uid:",uid, " ",URI_MAP["user_login"] , " info:",info)

    $.ajax({
        headers: {
            "X-Source-Type": header_X_Source_Type,
            "X-Project-Id": header_X_Project_Id,
            "X-Access": header_X_Access,
        },
        type: "POST",
        data : {"username":info.username,"password":info.password},
        url:URI_MAP["user_login"],
        dataType: "json",
        async:false,
        success: function(data){
            if (data.code != 200){
                return alert("UserLogin server back err:"+data.msg);
            }

            UserList[data.data.user.id].info = data.data.user;
            UserList[data.data.user.id].token = data.data.token;
            callback(UserList[data.data.user.id]);
        }
    });
}

function formatUnixTime(us){
    if (us <= 0 ){
        return "--";
    }

    var tims = new Date(us*1000);
    var format = tims.toLocaleString()
    return format;

}