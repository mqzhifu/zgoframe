document.write('<script src="/static/js/base64.js" type="text/javascript" charset="utf-8"></script>');
document.write('<script src="/static/js/md5.js" type="text/javascript" charset="utf-8"></script>');
document.write('<script src="/static/js/crypto-js.js" type="text/javascript" charset="utf-8"></script>');

// 后端公共HTTP接口-头信息
var header_X_Source_Type = "11";
var header_X_Project_Id = "6";
var header_X_Access = "imzgoframe";
//长连接 通信协议 基础类型 定义
var content_type_desc = {1:"json",2:"protobuf"};
var protocol_type_desc = {1:"tcp",2:"websocket",3:"udp"};
var CONN_STATUS_DESC = {1:"初始化",2:"运行中",3:"已关闭"};
//加密信息
// var DATA_ENCRYPT =  2;
// var secret = "ckZgoframe201310";
// var iv     = "ckZgoframe201310";
var DATA_ENCRYPT =  0;

//二次验证
var Second_Auth_Uname= "xiaoz"
var Second_Auth_Pwd  = "qwerASDFzxcv"
//================== 后端相关 ==================
//域名
// var domain = "127.0.0.1:1111";
var domain = window.location.host;
//http 协议
var http_protocol = "http";
if(location.href.substring(0,5) == "https"){
    http_protocol = "https";
}
//后端接口URI 的映射表
var URI_MAP = {
    "gateway_config":  "/gateway/config",
    "gateway_action_map": "/gateway/action/map",//URI - 网关 protobuf 映射表
    "user_login":  "/base/login",//登陆
    "gateway_fd_list":  "/gateway/fd/list",//长连接用户
    "gateway_send_msg":  "/gateway/send/msg",//长连接推送消息
    "gateway_total":  "/gateway/total",//长连接 metrics
    "twin_agora_config": "/twin/agora/config",//声网房间长连接配置
    "rule":  "/game/match/rule",
    "game_frame_sync_history":  "/frame/sync/room/history",

    "cicd_publish_list": "/cicd/publish/list",
    "cicd_server_service_list":  "/cicd/local/all/server/service/list",
    "cicd_service_list": "/cicd/service/list",
    "cicd_server_list":  "/cicd/server/list",
    "cicd_super_visor_list":   "/cicd/superVisor/list",
    "cicd_service_deploy":   "/cicd/service/deploy",
    "cicd_service_publish":  "/cicd/service/publish/#publishId#/2",

    "test_migu_get_para" :  "/tools/test/migu/api/para",
    "test_migu_send_back_data" :  "/tools/test/migu/api/backdata",

    "test_binary_tree_list" :  "/test/binary/tree/list/flag",
    // "test_binary_tree_list_middle" :  "/test/binary/tree/list/middle",
    // "test_binary_tree_list_end" :  "/test/binary/tree/list/end",
    "test_binary_tree_insert_one":"/test/binary/tree/insert/one/keyword",
    "test_binary_tree_each_deep":"/test/binary/tree/each/deep",


};
//用户列表，用于测试
var UserList = {
    10 :{id:10 ,"username"  :"doctor", "password":"123456","info":{},"token":"","roomId":"","channel":"ckck"},
     9 :{id:9  ,"username":"calluser", "password":"123456","info":{},"token":"","roomId":"","channel":"ckck"},
}
//根据 KEY 获取 URI 的完成URL地址
function get_uri_by_name(name){
    if(!URI_MAP.hasOwnProperty(name)){
        console.log("err : URI_MAP not has key:",name);
        return "";
    }
    return  http_protocol + "://"+domain + URI_MAP[name];
}
//给传输的数据加密
function encrypt_body_data(data){
    console.log("encrypt_body_data DATA_ENCRYPT:",DATA_ENCRYPT," data:",data);
    if(DATA_ENCRYPT <= 0 || !data){//未开启或数据为空，直接返回原数据即可
        return data;
    }
    switch (DATA_ENCRYPT){
        case 1:
            data = Base64.encode(data);
            console.log("encrypt_body_data  Base64.encode:",data);
        case 2:
            // console.log(JSON.stringify(data));
            // var data1 = CryptoJS.AES.encrypt(JSON.stringify(data), CryptoJS.enc.Utf8.parse(secret), {
            //     iv: CryptoJS.enc.Utf8.parse(iv),
            //     mode: CryptoJS.mode.CBC,
            //     padding: CryptoJS.pad.Pkcs7,
            // }).toString();
            var key = CryptoJS.enc.Utf8.parse(secret);
            var ivT = CryptoJS.enc.Utf8.parse(iv);
            let encrypted = CryptoJS.AES.encrypt(data, key, {
                iv: ivT,
                mode: CryptoJS.mode.CBC,
                padding: CryptoJS.pad.Pkcs7,
            });
            var data = encrypted.toString();
            console.log("encrypt_body_data AES CBC data:", data, " ori:", data);
        default:
            console.log("err DATA_ENCRYPT value err.")
            break;
    }
    return data;
}

function de_encrypt_body_data(data){
    if(DATA_ENCRYPT <= 0 || !data){//未开启或数据为空，直接返回原数据即可
        return data;
    }
    if( 1 == DATA_ENCRYPT){
        data = Base64.decode(data);
        console.log("de_encrypt_body_data  Base64.encode:",data);
    }else {
        // data = Base64.decode(data);
        // console.log("Base64.decode:",data)
        // var key = CryptoJS.enc.Utf8.parse(secret);
        // var ivT = CryptoJS.enc.Utf8.parse(iv);
        // var decrypt = CryptoJS.AES.decrypt(data, key,
        //     {
        //         iv: ivT,
        //         mode: CryptoJS.mode.ECB,
        //         padding: CryptoJS.pad.Pkcs7
        //     }
        // ).toString();
        // console.log(data);
        var key = CryptoJS.enc.Utf8.parse(secret);
        var ivT = CryptoJS.enc.Utf8.parse(iv);
        let encrypted = CryptoJS.AES.decrypt(data, key, {
            iv: ivT,
            mode: CryptoJS.mode.CBC,
            padding: CryptoJS.pad.Pkcs7,
        }).toString(CryptoJS.enc.Utf8);

        data = encrypted;
        // data = CryptoJS.enc.Utf8.stringify(decrypt).toString();
        // console.log("de_encrypt_body_data:",data);
    }

    return data;
}

//快捷 自动 登陆(根据UID)
function AutoUserLogin(uid,callback){
    var info = UserList[uid];
    console.log("UserLogin uid:",uid, " ",URI_MAP["user_login"] , " info:",info)

    var oriPostData = {"username":info.username,"password":info.password};
    request(callback,"user_login","",false,"POST",oriPostData,null);
    // var encryptData = encrypt_body_data(oriPostData);
    // // return 1;
    // $.ajax({
    //     headers: getCommonHeader(JSON.stringify(oriPostData)),
    //     type: "POST",
    //     data : encryptData,
    //     url:URI_MAP["user_login"],
    //     // dataType: "json",
    //     contentType: "application/json;charset=utf-8",
    //     async:false,
    //     success: function(data){
    //         if (data.code != 200){
    //             return alert("UserLogin server back err:"+data.msg);
    //         }
    //
    //         // UserList[data.data.user.id].info = data.data.user;
    //         // UserList[data.data.user.id].token = data.data.token;
    //         // callback(UserList[data.data.user.id]);
    //     }
    // });
}
//发起公共请求
function request(callback,urlMapKey , useToken,async,httpMethod,httpData,uriReplace){
    var httpUrl = get_uri_by_name(urlMapKey);
    if (uriReplace){
        for(let key  in uriReplace){
            httpUrl = httpUrl.replace(key,uriReplace[key]);
        }
    }
    console.log("Ajax request ,urlMapKey:"+urlMapKey , "url:"+httpUrl , " httpMethod:",httpMethod,"httpData:",httpData);

    var httpDataJsonStr = "";
    var encryptData = "";
    if(httpData){
        httpDataJsonStr = JSON.stringify(httpData);
        encryptData = encrypt_body_data(httpDataJsonStr);
    }

    $.ajax({
        headers: getCommonHeader(httpDataJsonStr,useToken),
        type: httpMethod,
        data : encryptData,
        url:httpUrl,
        // dataType: "json",
        contentType: "application/json;charset=utf-8",
        async:async,
        success: function(data){
            console.log("back data:",data);
            if (data.code != 200){
                return alert("UserLogin server back err:"+data.msg);
            }

            if(DATA_ENCRYPT > 0 ){
                data.data = eval( "(" +de_encrypt_body_data(data.data) + ")");
            }


            callback(data.data);


            // UserList[data.data.user.id].info = data.data.user;
            // UserList[data.data.user.id].token = data.data.token;
            // callback(UserList[data.data.user.id]);
        }
    });
}

// //发起公共请求
// function AjaxAdminReq(callback,urlMapKey , useToken,async,httpMethod,httpData,uriReplace){
//     console.log("AjaxAdminReq ,urlMapKey:"+urlMapKey , "httpMethod:",httpMethod,"httpData:",httpData);
//     var httpUrl = URI_MAP[urlMapKey];
//     if (uriReplace){
//         for(let key  in uriReplace){
//             httpUrl = httpUrl.replace(key,uriReplace[key]);
//         }
//     }
//
//     // if(DATA_ENCRYPT > 0){
//     //     // httpData = "ddd";
//     //     if(httpData){
//     //         var httpData = Base64.encode(httpData);
//     //         console.log("AjaxAdminReq Base64.encode:",httpData);
//     //         return 1;
//     //     }
//     // }
//
//     $.ajax({
//         headers: getCommonHeader(),
//         dataType: "json",
//         contentType: "application/json;charset=utf-8",
//         type: httpMethod,
//         data : httpData,
//         url:httpUrl,
//         async:async,
//         success: function(data){
//             console.log(data);
//             if (data.code != 200){
//                 return alert("AjaxReq server back err:"+data.msg);
//             }
//             callback(data.data);
//         }
//     });
// }

function getCommonHeader(postData,userToken){
    var now =  Math.round(new Date().getTime()/1000).toString();
    // var sign = md5( header_X_Project_Id + now + secret + postData);
    var sign = "";
    // console.log("sign:",header_X_Project_Id,now,secret,postData);
    var header = {
        "X-Second-Auth-Uname": Second_Auth_Uname,
        "X-Second-Auth-Ps":Second_Auth_Pwd,
        "X-Source-Type": header_X_Source_Type,
        "X-Project-Id": header_X_Project_Id,
        "X-Access": header_X_Access,
        "X-Token":userToken,
        "X-Client-Req-Time":now,
        "X-Sign":sign,
    }

    return header;
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