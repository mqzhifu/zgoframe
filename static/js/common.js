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
    "game_frame_sync_history":http_protocol + "://"+domain + "/frame/sync/room/history",

    "cicd_publish_list":http_protocol + "://"+domain +"/cicd/publish/list",
    "cicd_server_service_list": http_protocol + "://"+domain +"/cicd/local/all/server/service/list",
    "cicd_service_list": http_protocol + "://"+domain +"/cicd/service/list",
    "cicd_server_list": http_protocol + "://"+domain +"/cicd/server/list",
    "cicd_super_visor_list": http_protocol + "://"+domain +  "/cicd/superVisor/list",
    "test_migu_get_para" : http_protocol + "://"+domain + "/tools/test/migu/api/para",
    "test_migu_send_back_data" : http_protocol + "://"+domain + "/tools/test/migu/api/backdata",

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

//发起公共请求
function AjaxAdminReq(callback,urlMapKey , useToken,async,httpMethod,httpData){
    console.log("AjaxAdminReq ,urlMapKey:"+urlMapKey , "httpMethod:",httpMethod,"httpData:",httpData);
    $.ajax({
        headers: {
            "X-Second-Auth-Uname": "xiaoz",
            "X-Second-Auth-Ps":"qwerASDFzxcv",
            "X-Source-Type": header_X_Source_Type,
            "X-Project-Id": header_X_Project_Id,
            "X-Access": header_X_Access,
        },
        dataType: "json",
        contentType: "application/json;charset=utf-8",
        type: httpMethod,
        data : httpData,
        url:URI_MAP[urlMapKey],
        async:async,
        success: function(data){
            console.log(data);
            if (data.code != 200){
                return alert("AjaxReq server back err:"+data.msg);
            }
            callback(data.data);
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


//下面这些都是一些字节的处理，主要是给帧同步 protobuf 使用的，放到公共文件中了


function stringToUint8Array(str){
    var arr = [];
    for (var i = 0, j = str.length; i < j; ++i) {
        arr.push(str.charCodeAt(i));
    }

    var tmpUint8Array = new Uint8Array(arr);
    return tmpUint8Array
}

function processBufferString (dataBuffer,start,end){
    var str = "";
    for (var i = start; i < dataBuffer.length; i++) {
        if (i >= end){
            break;
        }
        str += String.fromCharCode(dataBuffer[i]);
    }
    return str;
}

function processBufferInt(dataBuffer,start,en){
    var str = "";
    for (var i = start; i < dataBuffer.length; i++) {
        if (i >= end){
            break;
        }
        str += dataBuffer[i];
    }
    return str;
}

function processBuffer(dataBuffer,start){
    //创建content ArrayBuffer和Uint8Array
    var contentArrayBuffer = new ArrayBuffer( dataBuffer.length - start );
    var contentUint8Array = new Uint8Array(contentArrayBuffer);
    var j = 0;
    for (var i = start; i < dataBuffer.length; i++) {
        contentUint8Array[j] = dataBuffer[i];
        j++;
    }
    return contentUint8Array;
}

function processBufferRange(dataBuffer,start,end){
    //创建content ArrayBuffer和Uint8Array
    var contentArrayBuffer = new ArrayBuffer( end - start );
    var contentUint8Array = new Uint8Array(contentArrayBuffer);
    var j = 0;
    for (var i = start; i < end; i++) {
        contentUint8Array[j] = dataBuffer[i];
        j++;
    }
    return contentUint8Array;
}

// function intToByte(i) {
//     var b = i & 0xFF;
//     var c = 0;
//     if (b >= 128) {
//         c = b % 128;
//         c = -1 * (128 - c);
//     } else {
//         c = b;
//     }
//     return c;
// }

function intToOneByteArr(i){
    var targets = new Uint8Array(1);
    targets[0] = i & 0xFF
    return targets;
}

function intToTwoByteArr(i){
    var targets = new Uint8Array(2);
    targets[0] = (i >> 8 & 0xFF);
    targets[1] = i & 0xFF
    return targets;
}


function intToByte4(i) {
    var targets = new Uint8Array(4);
    targets[0] = (i >> 24 & 0xFF);
    targets[1] = (i >> 16 & 0xFF);
    targets[2] = (i >> 8 & 0xFF);
    targets[3] = (i & 0xFF);
    return targets;
}

function Byte4ToInt(d) {
    var targets = new Array(4);
    targets[0] = (d[0] << 24 & 0xFF);
    targets[1] = (d[1] << 16 & 0xFF);
    targets[2] = (d[2] << 8 & 0xFF);
    targets[3] = (d[3] & 0xFF);
    return targets[0] + targets[1] +targets[2] +targets[3];
}

function Byte1ToInt(d) {
    var targets = new Array(1);
    targets[0] = (d[0] & 0xFF);
    return targets[0] ;
}

function Byte2ToInt(d) {
    var targets = new Array(2);
    targets[0] = (d[0] << 8 & 0xFF);
    targets[1] = (d[1] & 0xFF);
    // alert(targets[0]);
    // alert(targets[1]);
    return targets[0] + targets[1]  ;
}

function concatenate(...arrays) {
    let totalLen = 0;
    for (let arr of arrays)

        totalLen += arr.byteLength;

    let res = new Uint8Array(totalLen)

    let offset = 0

    for (let arr of arrays) {

        let uint8Arr = new Uint8Array(arr)

        res.set(uint8Arr, offset)

        offset += arr.byteLength

    }

    return res.buffer

}