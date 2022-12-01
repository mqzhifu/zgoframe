//=========
function Sync (playerId,token,GatewayConfig,DomIdPreObj,contentType,protocolType,playerInfo,actionMap,ruleConfig,mapSize,content_type_desc,protocol_type_desc) {
    this.status = 1;//1初始化 2等待准备 3运行中  4结束
    this.statusDesc = {
        1: "init",
        2: "wsLInkSuccess",
        3: "loginSuccess",
        4: "loginFailed",
        5: "matchSign",
        6: "cancelSign",
        7: "ready",
        8: "startBattle",
        9: "end",
        10: "close",
    };
    this.heartbeatLoopFunc = null;//心跳回调函数
    this.tableMax = mapSize;//地图的表格大小 5 X 5
    this.otherPlayerOffline = 0;//其它玩家调线
    this.pushLogicFrameLoopFunc = null;//定时循环 - 推送玩家操作函数
    this.playerOperationsQueue = [];//一个帧时间内，收集玩家的操作指令 队列
    this.closeFlag = 0;//关闭标识，0正常1手动关闭2后端关闭
    this.tableId = "";
    this.myClose = 0;//区分关闭事件，是自己触发，还是对方触发的
    this.playerLocation = new Object();//每个玩家的位置信息
    this.operationsInc = 0;//玩家操作指令自增ID
    this.logicFrameLoopTimeMs = 0;//毫秒，多少时间内收集一次玩家操作，推送到S端
    this.playerCommandPushLock = 0;
    this.roomId = "";
    this.sequenceNumber = 0;
    this.randSeek = 0;
    this.sessionId = "";
    this.matchSignGroupId = 0;
    this.serverPlayerState = null;//首次建立长连接后，想要玩游戏，得先登陆，登陆后要立刻从服务端摘取一下玩家的当前状态
    this.wsObj = null;//js内置ws 对象，创建WS连接成功后，会给此变量赋值
    //以下均是外面传进来的值
    this.hostUri = "ws://" + GatewayConfig.outIp + ":" + GatewayConfig.wsPort + GatewayConfig.wsUri;//ws 连接 s 端地址
    this.playerInfo = playerInfo;//玩家基础信息
    this.playerId = playerId;//玩家ID，快捷变量
    this.matchGroupPeople = ruleConfig.condition_people;//一局游戏：需要总人数
    this.offLineWaitTime = ruleConfig.off_line_wait_time;//lockStep 玩家掉线后，其它玩家等待最长时间
    this.token = token;//玩家的凭证
    this.FPS = ruleConfig.fps;//每秒多少逻辑帧
    this.ruleId = ruleConfig.id;
    this.contentTypeDesc = content_type_desc;
    this.domIdObj = DomIdPreObj;
    this.actionMap = actionMap;
    this.contentType = contentType;//protobuf|json
    this.protocolType = protocolType;//tcp|ws
    this.ClientHeartbeatTime =  GatewayConfig.client_heartbeat_time;

    //公共日志输出前缀，方便测试
    this.descPre = this.getPlayerDescById(playerId);

    this.show("玩家"+playerId + " Sync 构造函数完成.");

}
//创建ws长连接，也算是入口函数，才有后续的所有操作
Sync.prototype.create = function(){
    this.show("创建 ws 连接");
    if (this.status != 1 && this.status != 10){
        return this.show("status错误， status  !=  init or close",1)
    }
    var parent = this
    this.closeFlag = 0;//清空 关闭标识
    //根据帧数，计算出 秒数
    this.logicFrameLoopTimeMs = parseInt( 1000 / this.FPS);

    this.show("create new WebSocket"+this.hostUri ," FPS:",this.FPS , " ms:",this.logicFrameLoopTimeMs, " contentType:",this.contentType, " protocolType:",this.protocolType)
    //创建ws连接
    this.wsObj = new WebSocket(this.hostUri);
    //设置 关闭回调
    this.wsObj.onclose = function(ev){
        parent.close(ev)
    };
    //设置 接收S端消息-回调
    this.wsObj.onmessage = function(ev){
        parent.onmessage(ev)
    };
    //设置 连接建立成功回调
    this.wsObj.onopen = function(){
        parent.wsOpen();
    };
    this.wsObj.onerror = function(ev){
        alert("wsObj.onerror");
        parent.log("error:"+ev);
    };
}

//ws 接收到服务端关闭
Sync.prototype.close = function(ev){
    alert("receive server close:" +ev.code);
    this.myClearInterval();
    if (this.myClose == 1){//自己关闭的WS
        var reConnBntId = "reconn_"+this.playerId;
        var msg = "重连接";
        this.upOptBntHref(reConnBntId,msg,this.create);
    }else{
        this.closeFlag = 2;
        this.upOptBnt("服务端关闭，游戏结束，连接断开",1)
    }
};

Sync.prototype.upOptBntHref = function(domId,value,clickCallback){
    var bntContent = "<a href='javascript:void(0);' onclick='' id='"+domId+"'>"+value+"</a>";
    this.upOptBnt(bntContent, 1);
    $("#"+domId).click(clickCallback);
};

//更新操作按钮文字，解除 点击事件
Sync.prototype.upOptBnt = function(content,clearClick){
    $("#"+this.domIdObj.optBntId).html(content);
    if(clearClick == 1){
        $("#"+this.domIdObj.optBntId).unbind("click");
    }
};

//接收S端WS消息
Sync.prototype.onmessage = function(ev){
    // console.log("onmessage , ev :",ev );
    var msgObj = this.newMsgObj();
    var parent = this;
    //解析C端发送的数据，这一层，对于用户层的 content 数据不做处理
    //1-4字节：当前包数据总长度，~可用于：TCP粘包的情况
    //5字节：content type
    //6字节：protocol type
    //7字节 :服务Id
    //8-9字节 :函数Id
    //10-19：预留，还没想好，可以存sessionId，也可以换成UID
    //19 以后为内容体
    //结尾会添加一个字节：\f ,可用于 TCP 粘包 分隔
    var debugInfo = "onmessage ";
    if (this.contentTypeDesc[this.contentType] == 'protobuf'){
        debugInfo += " contentType: protobuf";
        var reader = new FileReader();
        reader.readAsArrayBuffer(ev.data);
        reader.onloadend = function(e) {
            // var dataBuffer = new Uint8Array(reader.result);
            //
            // msgObj.contentType = processBufferString(dataBuffer,0,1);
            // msgObj.protocolType = processBufferString(dataBuffer,1,2);
            // msgObj.actionId = processBufferString(dataBuffer,2,6);
            // msgObj.sessionId = processBufferString(dataBuffer,6,38);
            // var content = processBuffer(dataBuffer,38);
            // msgObj.action = self.getActionName(msgObj.actionId,"server")
            // //首字母转大写
            // var actionLow = msgObj.action.substring(0, 1).toUpperCase() + msgObj.action.substring(1)
            // //拼接成最终classname
            // var className =  "proto.pb.Response" + actionLow;
            // var responseProtoClass = eval(className);
            // //将进制流转换成对象
            // msgObj.content =  responseProtoClass.deserializeBinary(content).toObject();
            // self.router(msgObj);
        };
    }else if(this.contentTypeDesc[this.contentType] == "json"){
        debugInfo += " contentType: json ";
        var reader = new FileReader();
        reader.readAsArrayBuffer(ev.data);
        reader.onloadend = function(e) {
            var dataBuffer = new Uint8Array(reader.result);

            var bytes4 = processBufferRange(dataBuffer,0,4);
            msgObj.dataLength = Byte4ToInt(bytes4);
            var bytes1 = processBufferRange(dataBuffer,4,5);
            msgObj.contentType = Byte1ToInt(bytes1);
            var bytes1 = processBufferRange(dataBuffer,5,6);
            msgObj.protocolType = Byte1ToInt(bytes1);
            var bytes1 = processBufferRange(dataBuffer,6,7);
            msgObj.serviceId = Byte1ToInt(bytes1);
            var bytes2 = processBufferRange(dataBuffer,7,9);
            msgObj.funcId = Byte2ToInt(bytes2);
            var sessionBytes = processBufferRange(dataBuffer,9,19);
            msgObj.sessionId = processBufferString(sessionBytes,0)
            var content = processBufferRange(dataBuffer,19,19+msgObj.dataLength);
            content = processBufferString(content,0);
            // console.log("lenDataBuffer:",dataBuffer.length," content:",content);
            msgObj.content =  eval("("+content+")");

            msgObj.sidFid = msgObj.serviceId + "" + msgObj.funcId;

            parent.showComplex(debugInfo + "msgObj:",msgObj);
            parent.router(msgObj);
        }

    // }else if(contentTypeDesc[self.contentType] == "json"){//这种是纯JSON格式，传输的是字符流，我再想想如何处理
    //     msgObj.contentType = ev.data.substr(0,1);
    //     msgObj.protocolType = ev.data.substr(1,1);
    //     msgObj.actionId = ev.data.substr(2,4);
    //     msgObj.sessionId = ev.data.substr(6,37);
    //     msgObj.content = ev.data.substr(38);
    //     alert(msgObj.actionId);
    //     // console.log(pre + " contentType:" +contentType+ ",protocolType:" + protocolType  +" ,actionId:"+actionId +",sessionId:" +sessionId + " ,content:",content );
    //     msgObj.action = self.getActionName(msgObj.actionId,"server")
    //     msgObj.content =  eval("("+msgObj.content+")");
    //     self.router(msgObj);
    }else{
        this.show(debugInfo + "contentType err",1)
        return alert("contentType err");
    }
};

Sync.prototype.newMsgObj = function (){
    var msg = new proto.pb.Msg();
    return msg.toObject();
};

Sync.prototype.showComplex = function(str,complexType,showNoticeMsg = 0){
    console.log(this.descPre + " " + str ,complexType)
    if (showNoticeMsg){
        this.noticeMsg(str);
    }
}

Sync.prototype.show = function(str,showNoticeMsg = 0){
    console.log(this.descPre + " " + str)
    if (showNoticeMsg){
        this.noticeMsg(str);
    }
}

Sync.prototype.noticeMsg = function(str){
    $("#"+this.domIdObj.msgNotice).html(str);
}


//连接成功后，会执行此函数
Sync.prototype.wsOpen = function(){
    this.show("onOpen , ws connect server : Success  ",1);
    this.upStatus(2);
    //强依赖，proto 文件
    var requestLoginObj = new proto.pb.Login();
    requestLoginObj.setToken(this.token) ;
    this.sendMsg("CS_Login",requestLoginObj);
};

Sync.prototype.router = function(msgObj){
    var action = this.actionMap.server[msgObj.sidFid].func_name;
    this.showComplex("router , action:"+action,msgObj);

    eval( "this."+action+"(msgObj.content)" );
    // }else if( action == 'SC_Ping'){//获取一个当前玩家的状态，如：是否有历史未结束的游戏
    //     this.rServerPing(content);
    // }else if ( action == 'SC_RoomBaseInfo' ){
    //     this.rPushRoomInfo(content);
    // }else if ( action == 'SC_OtherPlayerOffline' ){
    //     this.rOtherPlayerOffline(content);
    // }else if( "SC_GameOver" == action){
    //     this.rGameOver(content);
    // }else if( "SC_KickOff" == action){
    //     this.rKickOff(content);
    // }else if( "SC_Pong" == action){
    //     this.rServerPong(content)
    // }else if( "SC_OtherPlayerResumeGame" == action){
    //     this.rOtherPlayerResumeGame(content)
    // }else if( "SC_RoomHistory" == action){
    //     this.rPushRoomHistory(content);
    //     // alert("接收到，玩家-房间-历史操作记录~");
    // }else{
    //     return alert("action error."+action);
    // }
}
//S推送-登陆结果
Sync.prototype.SC_Login = function(LoginRes){
    if (LoginRes.code != 200) {
        this.upStatus(4);
        var msg = "登陆失败 , code: "+LoginRes.code + " , errMsg: "+LoginRes.errMsg;
        return this.show(msg);
    }

    this.show("登陆成功，开始摘取玩家状态信息",1)

    this.upStatus(3);
    var playerBase = new proto.pb.PlayerBase();
    playerBase.setPlayerId(  this.playerId);
    this.sendMsg("CS_PlayerState",playerBase);

    // var now = Date.now();
    // var requestClientPing = new proto.pb.Ping();
    // requestClientPing.setAddTime(now);
    // this.sendMsg("CS_Ping",requestClientPing);

    // this.heartbeatLoopFunc = setInterval(this.heartbeat.bind(this), this.ClientHeartbeatTime * 1000);
};
//S推送-心跳
Sync.prototype.SC_Heartbeat = function(Heartbeat){
    this.showComplex("rHeartbeat:",Heartbeat);
}
//S推送-玩家当前的状态信息
Sync.prototype.SC_PlayerState = function(PlayerState){
    this.showComplex("ServerPlayerState ",PlayerState,1)
    this.serverPlayerState = PlayerState;
    //这里是有问题的，roomId我在外层写死了均为空，应该是动态从后端拿，且最好是短连接去取，回头优化
    if (this.serverPlayerState.roomId){
        alert("检测出，有未结束的一局游戏，开始恢复中...,先获取房间信息:rooId:"+this.serverPlayerState.roomId);
        var requestGetRoom = new proto.pb.RoomBaseInfo();
        requestGetRoom.setRoomId(playerConnInfo.roomId);
        requestGetRoom.setPlayerId(playerId);
        this.sendMsg("CS_RoomBaseInfo",requestGetRoom);
        //     var msg = {"roomId":playerConnInfo.roomId,"playerId":playerId};
    }else{
        var matchSignBntId = "matchSign_"+this.playerId;
        var hrefBody = "匹配报名";
        this.noticeMsg("连接成功，等待报名...")
        this.upOptBntHref(matchSignBntId,hrefBody,this.matchSign.bind(this));
    }
}
//S推送-匹配的结果：1. 未成团，匹配失败 2匹配成功了，等帧同步服务再发消息通知
Sync.prototype.SC_GameMatchOptResult = function(GameMatchOptResult){
    if(GameMatchOptResult.code != 200){

        var matchSignBntId = "matchSign_"+this.playerId;
        var hrefBody = "重新匹配报名";

        this.upOptBntHref(matchSignBntId,hrefBody,this.matchSign.bind(this));

        this.noticeMsg("匹配失败 , GameMatchOptResult, msg: "+GameMatchOptResult.msg + " code:"+ GameMatchOptResult.code);
        return false;
    }
    this.show("匹配成功");
    this.noticeMsg("等待S端推送初始化数据...")
    this.roomId = GameMatchOptResult.roomId;
}
//匹配成功，房间游戏基础信息S端已建立，告知C端，可进入准备状态，且初始化本地信息
Sync.prototype.SC_EnterBattle = function(EnterBattle){
    this.show("SC_EnterBattle EnterBattle:",EnterBattle , " contentType:",this.contentType)
    if(this.contentTypeDesc[this.contentType] =="protobuf"){
        console.log("rEnterBattle in protobuf ")
        EnterBattle.playerList = EnterBattle.playerListList;
    }

    this.initLocalGlobalVar(EnterBattle);

    var readySignBntId = "ready_"+this.playerId;
    var hrefBody = "准备";
    this.noticeMsg("匹配成功，已接收S端初始化数据...")
    this.upOptBntHref(readySignBntId,hrefBody,this.ready.bind(this));
}

Sync.prototype.SC_ReadyTimeout = function(){
    console.log("rReadyTimeout:");

    this.upStatus(3);

    var matchSignBntId = "matchSign_"+this.playerId;
    var hrefBody = "等待准备超时，重新匹配报名";

    this.upOptBntHref(matchSignBntId,hrefBody,this.matchSign);

    this.noticeMsg("抱歉，<准备时间>已超时");
    // return alert("抱歉，<准备时间>已超时");
}

Sync.prototype.SC_StartBattle = function(StartBattle){
    this.upStatus(8);
    this.pushLogicFrameLoopFunc = setInterval(this.playerCommandPush.bind(this),this.logicFrameLoopTimeMs);
    var exceptionOffLineId = "exceptionOffLineId"+this.playerId;
    var msg = "模拟，异常掉线";
    this.noticeMsg("游戏已开始，数据同步中...")
    this.upOptBntHref(exceptionOffLineId,msg,this.closeFD);
};
//source:1正常接收服务器推送，2玩家掉线恢复使用
Sync.prototype.SC_LogicFrame = function(logicFrame,source){//接收S端逻辑帧
    this.showComplex("logicFrame:",logicFrame);
    if(this.contentTypeDesc[this.contentType] =="protobuf"){
        if (source != "rPushRoomHistory"){
            logicFrame.operations = logicFrame.operationsList;
        }
    }
    var operations = logicFrame.operations;
    this.sequenceNumber  = logicFrame.sequenceNumber;
    $("#"+this.domIdObj.seqId).html(this.sequenceNumber);//更新当前逻辑帧号，并页面显示，用于测试

    this.playerCommandPushLock = 0;//解锁，恢复玩家收集指令队列
    if (!operations || typeof (operations) == "undefined" ){
        return this.show("SC_LogicFrame operations empty")
    }
    this.show("rPushLogicFrame ,sequenceNumber:"+this.sequenceNumber+ ", operationsLen:" +  operations.length);
    for(var i=0;i<operations.length;i++){
        var playerId = operations[i].playerId;
        var str = " i=" + i + " , id: "+operations[i].id + " , event:"+operations[i].event + " , value:"+ operations[i].value + " , playerId:" + playerId;
        this.show(str);
        if ("move" == operations[i].event ){
            var LocationArr = operations[i].value.split(",");
            var LocationStart = LocationArr[0];
            var LocationEnd = LocationArr[1];

            // var lightTd = "map"+id +"_"+LocationStart + "_"+LocationEnd;
            var lightTd =this.getMapTdId(this.tableId,LocationStart,LocationEnd);
            this.show( lightTd);
            var tdObj = $("#"+lightTd);
            if(playerId == this.playerId){//地图格子上，自己用绿色标记
                tdObj.css("background", "green");
            }else{
                tdObj.css("background", "red");//地图格子上，其它人用绿色标记
            }
            var playerLocation = this.playerLocation;
            if (playerLocation[playerId] == "empty"){
                //证明是第一次移动，没有之前的数据
            }else{
                // playerLocation = getPlayerLocation(playerId);
                // alert(commands[i].playerId);
                var playerLocationArr = playerLocation[playerId].split("_");
                //非首次移动，这次移动后，要把之前所有位置清掉
                var lightTd = this.getMapTdId(this.tableId,playerLocationArr[0],playerLocationArr[1]);
                var tdObj = $("#"+lightTd);
                tdObj.css("background", "");
            }
            playerLocation[playerId] = LocationStart + "_"+LocationEnd;
        }else if(operations[i].event == "empty"){

        }
    }
};

//======================== 以上都是接收S端消息通知 ========================



Sync.prototype.ready = function(){
    this.upStatus(7);
    var requestPlayerReady = new proto.pb.PlayerReady();
    requestPlayerReady.setPlayerId(this.playerId);
    requestPlayerReady.setRoomId(this.roomId);
    this.sendMsg("CS_PlayerReady",requestPlayerReady);

    this.noticeMsg("不能取消，只能等成功或超时");

    this.upOptBntHref("","等待其它玩家准备",null);
};



Sync.prototype.initLocalGlobalVar = function(EnterBattle){
    console.log("initLocalGlobalVar:",EnterBattle)
    for(var i=0;i<EnterBattle.playerIds.length;i++){
        this.playerLocation[""+EnterBattle.playerIds[i]+""] = "empty"
    }
    // return 1;
    this.randSeek  = EnterBattle.randSeek;
    $("#"+this.domIdObj.randSeekId).html(this.randSeek);


    this.sequenceNumber  = EnterBattle.sequenceNumber;
    $("#"+this.domIdObj.seqId).html(this.sequenceNumber);

    this.roomId = EnterBattle.roomId;
    $("#"+this.domIdObj.roomId).html(this.roomId);

    var str =  " roomId:" +EnterBattle.roomId+ ", RandSeek:"+    this.randSeek +" , SequenceNumber"+    this.sequenceNumber ;
    console.log(str);
};



Sync.prototype.matchSign = function(){
    this.upStatus(5)

    this.matchSignGroupId = this.makeGroupId()

    var gameMatchSign = new proto.pb.GameMatchSign();
    gameMatchSign.setRuleId(this.ruleId);
    gameMatchSign.setGroupId(this.matchSignGroupId);
    gameMatchSign.setAddition("html_test_frame_sync");

    var gameMatchSignPlayer = new proto.pb.GameMatchSignPlayer();
    gameMatchSignPlayer.setUid(this.playerId);
    gameMatchSignPlayer.clearWeightAttrMap();
    gameMatchSignPlayer.getWeightAttrMap().set("age",30);
    gameMatchSignPlayer.getWeightAttrMap().set("level",40);
    // this.showComplex("gameMatchSignPlayer:",gameMatchSignPlayer.toString())
    // return 1;
    gameMatchSign.addPlayerSets(gameMatchSignPlayer)
    // gameMatchSign.addPlayerSets(gameMatchSignPlayer);

    this.sendMsg("CS_PlayerMatchSign",gameMatchSign);

    var cancelBntId = "cancelSign_"+this.playerId;
    var hrefBody = "取消匹配报名";

    this.noticeMsg("匹配中....")
    this.upOptBntHref(cancelBntId,hrefBody,this.cancelSign.bind(this));
};
//取消报名
Sync.prototype.cancelSign = function(){
    this.upStatus(6);

    var gameMatchPlayerCancel = new proto.pb.GameMatchPlayerCancel();
    gameMatchPlayerCancel.setRuleId(this.ruleId);
    gameMatchPlayerCancel.setGroupId(this.matchSignGroupId);
    this.sendMsg("CS_PlayerMatchSignCancel",gameMatchPlayerCancel);

    var matchSignBntId = "matchSign_"+this.playerId;
    var hrefBody = "连接成功，重新匹配报名";

    this.upOptBntHref(matchSignBntId,hrefBody,this.matchSign);
};



Sync.prototype.sendMsg =  function ( action,contentObj  ){
    var id = this.getActionId(action,"client");
    if (!id){
        this.show("sendMsg err:id empty action:"+action);
        return false;
    }
    // console.log( contentObj.toObject())
    var content = null;

    //解析C端发送的数据，这一层，对于用户层的content数据不做处理
    //1-4字节：当前包数据总长度，~可用于：TCP粘包的情况
    //5字节：content type
    //6字节：protocol type
    //7字节 :服务Id
    //8-9字节 :函数Id
    //10-19：预留，还没想好，可以存sessionId，也可以换成UID
    //19 以后为内容体
    //结尾会添加一个字节：\f ,可用于 TCP 粘包 分隔

    var serviceId = id.toString().substring(0,2);
    var funcId = id.toString().substring(2);
    var session = "1234567890";
    this.show( " <sendMsg> action: "+ action  + " fullId: " + id + " serviceId: " + serviceId + " funcId: " + funcId   );
    if (this.contentTypeDesc[this.contentType] == "json"){
        var debugLog = " contentType: json ";
        content = contentObj.toObject();
        this.showComplex("debug contentObj.toObject():",content)

        //js 编译完的 proto 类文件，正常使用二进制的protobuf传输是OK的，但是要直接当成json使用，它有几个问题：
        //1. 所有的数组类型，它自动给 key 加了list
        //2. 所有的map类型，它并不是一个对象，而是用一个数组结构来存储，下面就是把数组转回成map 对象格式
        if (action == "CS_PlayerMatchSign" && content.playerSetsList.length > 0 ){
            for(var i=0;i<content.playerSetsList.length;i++){
                if(content.playerSetsList[i].weightAttrMap && content.playerSetsList[i].weightAttrMap.length > 0){
                    var newWeightAttrMap = {};
                    for(var j=0;j<content.playerSetsList[i].weightAttrMap.length;j++){
                        newWeightAttrMap[ content.playerSetsList[i].weightAttrMap[j][0]] =  content.playerSetsList[i].weightAttrMap[j][1]
                    }
                    content.playerSetsList[i].weightAttrMap = newWeightAttrMap;
                    console.log(newWeightAttrMap);

                }
            }
        }

        content = JSON.stringify(content);
        //这里有个坑，注意下. JS编译proto文件后，会把 数组类型 自动加上：List 关键字，而 map 类型会 多加一个 map ,这里得去掉
        content = content.replace("List","");
        content = content.replace("Map","");



        this.showComplex("debug JSON.stringify:",content)

        var contentLenByte = intToByte4( content.length);
        var contentTypeByte = intToOneByteArr(contentType);
        var protocolTypeByte = intToOneByteArr(protocolType);
        var serviceIdByte = intToOneByteArr(parseInt(serviceId));
        var funcIdByte = intToTwoByteArr(parseInt(funcId));
        var sessionByte = stringToUint8Array(session);
        var contentByte = stringToUint8Array(content);

        // console.log("bbyte:",contentLenByte,contentTypeByte,protocolTypeByte,serviceIdByte,funcIdByte);
        // concatenate(contentLenByte,contentTypeByte[3]);
        // content =  concatenate(contentLenByte,contentTypeByte[3],protocolTypeByte[3],serviceIdByte[3],funcIdByte[3],sessionByte,contentByte);
        var endStr = new Uint8Array(1);
        endStr[0] = "\f";
        content =  concatenate(contentLenByte,contentTypeByte,protocolTypeByte,serviceIdByte,funcIdByte,sessionByte,contentByte,endStr)  ;

        // content = contentLenByte +"" + contentTypeByte + "" + protocolTypeByte + "" +serviceIdByte + "" + funcIdByte +"" +sessionByte + contentByte+ "\f";

        // var finalContent = contentLenByte + contentTypeByte +  "" + protocolTypeByte + serviceIdByte + "" +  funcIdByte + sessionByte + contentByte  + "\f";


        // var protocolCtrl = contentType +  "" + protocolType + id;
        // if (action == "login" ){
        //     content = contentLen + protocolCtrl + content;
        // }else{
        //     content = contentLen + protocolCtrl + self.sessionId +  content;
        // }
        //下面未使用
        // var contentTypeByte = contentType << 5;
        // console.log("aLeftL:",aLeft);
        // var firstbyte = contentTypeByte | protocolType;
        // console.log("firstbyte:",firstbyte);
        // var myArrayBuffer = new Uint8Array(1)
        // var myIntToByte =  intToByte(firstbyte);
        // myArrayBuffer[0] = firstbyte;
        // var b_s = String.fromCharCode.apply(null, new Uint8Array(myArrayBuffer));
        // console.log("myArrayBuffer:",myArrayBuffer,",b_s:",b_s);
        // var emptyByte = intToByte(0);
        // content =  b_s +emptyByte + content;
        this.showComplex("<sendMsg final>" + debugLog + " " , content);

        this.wsObj.send(content);
    }else if ( this.contentTypeDesc[this.contentType]  == "protobuf"){
        debugLog += " contentType: protobuf ";
        // content = contentObj.serializeBinary();
        // var protocolCtrl = contentType +  "" + protocolType + id;
        // if (action != "login" ){
        //     protocolCtrl = protocolCtrl + self.sessionId  ;
        // }
        // var idBinary = stringToUint8Array(protocolCtrl);
        // var mergedArray = new Uint8Array(idBinary.length + content.length);
        // mergedArray.set(idBinary);
        // mergedArray.set(content, idBinary.length);
        //
        // self.wsObj.send(mergedArray);
    }
};
//玩家报名，得有个groupId，先暂时由前端生成随机数，后期我想想如何处理
Sync.prototype.makeGroupId = function(){
    return Math.round(Math.random()*8999+1000);
}
//更新当前状态
Sync.prototype.upStatus = function(status){
    this.show("up status ,  old status:" + this.status   +  "("+ this.statusDesc[this.status]+") , new status:"  + status + "("+ this.statusDesc[status] + ")");
    this.status = status;
};

//心跳执行函数
Sync.prototype.heartbeat = function(){
    var now = Date.now();
    var requestClientHeartbeat = new proto.pb.Heartbeat();
    requestClientHeartbeat.setTime(now);
    this.sendMsg("CS_Heartbeat",requestClientHeartbeat);
};

Sync.prototype.rKickOff = function(ev){
    return alert("您被踢下线喽~");
};

Sync.prototype.getPlayerDescById = function (id){
    return "(player_"+ id+")";
};
Sync.prototype.getMapTdId = function (tableId,i,j){
    return tableId + "_" + i +"_" + j;
}

Sync.prototype.getActionId = function (action,category){
    var data = this.actionMap[category];
    for(let key  in data){
        if (data[key].func_name == action){
            return data[key].id;
        }
    }
    this.show("getActionId is empty, action:"+action + " category:"+category)
    return "";
};

Sync.prototype.getActionName = function (actionId,category){
    var data = this.actionMap[category];
    return data[actionId].func_name;
};


Sync.prototype.myClearInterval = function (){
    var a = clearInterval(this.pushLogicFrameLoopFunc);
    console.log("=================================",a,this.pushLogicFrameLoopFunc);
    // a = clearInterval(this.heartbeatLoopFunc);
    // console.log(a)
}
//玩家操作 - 主动关闭
Sync.prototype.closeFD = function (){
    this.show("closeFD");
    this.myClearInterval();
    this.myClose = 1;
    this.wsObj.close();
};

Sync.prototype.getMap = function (tableId) {
    this.tableId = tableId;
    var tableObj = $("#" + tableId);
    var matrix = new Array();
    var matrixSize = this.tableMax;
    var inc = 0;
    for (var i = 0; i < matrixSize; i++) {
        matrix[i] = new Array();
        var trTemp = $("<tr></tr>");
        for (var j = 0; j < matrixSize; j++) {
            // var tdId = tableId + "_" + i +"_" + j;
            var tdId = this.getMapTdId(tableId,i,j);
            matrix[i][j] = inc;
            trTemp.append("<td id='"+tdId+"'>"+ i +","+j +"</td>");
            inc++;1
        }
        // alert(trTemp);
        trTemp.appendTo(tableObj);
    }
};

Sync.prototype.playerCommandPush = function (){

    var PlayerOperations = new proto.pb.LogicFrame();
    PlayerOperations.setId(9999);//此值目前没用上
    PlayerOperations.setRoomId(this.roomId);
    PlayerOperations.setSequenceNumber(this.sequenceNumber);

    if(this.playerCommandPushLock == 1){//是否锁住了
        console.log("send msg lock...please wait server back login frame");
        return
    }
    if (this.playerOperationsQueue.length > 0){
        // {"id":self.operationsInc,"event":"move","value":newLocation,"playerId":self.playerId}
        var operations = new Array(this.playerOperationsQueue.length)
        for(var i=0;i<this.playerOperationsQueue.length;i++){
            var operation = new proto.pb.Operation();
            operation.setId(this.playerOperationsQueue[i].id);
            operation.setEvent(this.playerOperationsQueue[i].event);
            operation.setValue(this.playerOperationsQueue[i].value);
            operation.setPlayerId(this.playerId);
            // operations.push(operation);
            operations[i] = operation;
        }
        PlayerOperations.setOperationsList(operations);
        // loginFrame.operations = self.playerOperationsQueue;
        this.playerOperationsQueue = [];//将当前：队列里的 玩家操作数据，清空
    }else{//当前玩家在此帧没有产生操作数据
        var operations = new Array(1);//创建一个长度为1的数组
        var operation = new proto.pb.Operation();
        operation.setId(1);
        operation.setEvent("empty");//empty和-1 代表，当前帧，该玩家没有任何操作
        operation.setValue("-1");
        operation.setPlayerId(this.playerId);
        // operations.push(operation);
        operations[0] = operation;
        // console.log(operations)
        PlayerOperations.setOperationsList(operations);
        // var emptyCommand = [{"id":1,"event":"empty","value":"-1","playerId":self.playerId}];
        // loginFrame.operations = emptyCommand;
    }
    this.sendMsg("CS_PlayerOperations",PlayerOperations);//发送一帧玩家操作数据
    this.playerCommandPushLock = 1;//发完包，直接锁住，等待S端返回信息，才解锁
};
//玩家移动
Sync.prototype.move = function ( dirObj ){
    this.show("in move");
    if (this.otherPlayerOffline){
        return alert("其它玩家掉线了，请等待一下...");
    }

    if (this.closeFlag > 0 ){
        return alert("WS FD 已关闭...");
    }

    if (this.status != 8){
        return alert("status err , != startBattle ， 游戏还未开始，请等待一下...");
    }

    var dir = dirObj.data;
    var playerLocation = this.playerLocation;
    var nowLocationStr = playerLocation[this.playerId]
    var nowLocationArr = nowLocationStr.split("_");
    var nowLocationLine =  nowLocationArr[0];
    var nowLocationColumn = nowLocationArr[1];

    nowLocationLine = Number(nowLocationLine)
    nowLocationColumn = Number(nowLocationColumn)
    var newLocation = "";
    if(dir == "up"){
        if(nowLocationLine == 0 ){
            return alert("nowLocationLine == 0");
        }
        var newLocationLine =  nowLocationLine - 1;
        newLocation = newLocationLine + "," + nowLocationColumn;
    }else if(dir == "down"){
        if(nowLocationLine == this.tableMax - 1 ){
            return alert("nowLocationLine == "+ this.tableMax);
        }
        var newLocationLine =  nowLocationLine + 1;
        newLocation = newLocationLine + "," + nowLocationColumn;
    }else if(dir == "left"){
        if(nowLocationColumn == 0 ){
            return alert("nowLocationColumn == 0");
        }
        var newLocationColumn =  nowLocationColumn - 1;
        newLocation = nowLocationLine + "," + newLocationColumn;
    }else if(dir == "right"){
        if(nowLocationColumn ==  this.tableMax - 1 ){
            return alert("nowLocationColumn == "+ this.tableMax);
        }
        var newLocationColumn =  nowLocationColumn + 1;
        newLocation = nowLocationLine + "," + newLocationColumn;
    }else {
        return alert("move dir error."+dir);
    }

    var localNewLocation = newLocation.replace(',','_');
    for(let key  in playerLocation){
        // alert(playerLocation[key]);
        if(playerLocation[key] == localNewLocation){
             return this.gameOverAndClear()
        }
    }

    this.show("dir:"+dir+"oldLocation"+nowLocationStr+" , newLocation:"+newLocation);
    this.playerOperationsQueue.push({"id":this.operationsInc,"event":"move","value":newLocation,"playerId":this.playerId});
    this.operationsInc++;
    var playerLocationArr = playerLocation[this.playerId].split("_");
    var lightTd = this.getMapTdId(this.tableId,playerLocationArr[0],playerLocationArr[1]);
    var tdObj = $("#"+lightTd);
    tdObj.css("background", "");
}

//两个玩家，位移碰撞了，触发了游戏结束
Sync.prototype.gameOverAndClear = function(){
    this.upStatus(9);

    var requestGameOver = new proto.pb.GameOver()
    requestGameOver.setRoomId(this.roomId);
    requestGameOver.setSequenceNumber(this.sequenceNumber);
    requestGameOver.setResult("ccccccWin");
    this.sendMsg("CS_GameOver", requestGameOver);


    var msg = "完犊子了，撞车了...这游戏得结束了....";
    this.noticeMsg(msg)

    this.myClearInterval();
    this.upOptBnt("游戏结束1",1)

    // return alert("完犊子了，撞车了...这游戏得结束了....");
};



    // //=================== 以下都是 接收S端的处理函数========================================
    // this.rOtherPlayerResumeGame = function(content){
    //     if(content.playerId != self.playerId){
    //         var tdId = self.tableId + "_" + self.playerLocation[content.playerId];
    //         var tdObj = $("#"+tdId)
    //         tdObj.css("background", "red");
    //     }
    // };

    // this.rPushRoomHistory = function(logicFrame){
    //     console.log("rPushRoomHistory:");
    //     if(self.contentTypeDesc[self.contentType] =="protobuf"){
    //         logicFrame.list = logicFrame.listList;
    //     }
    //     var list = logicFrame.list;
    //     for(var i=0;i<list.length;i++){
    //         // console.log( "rPushRoomHistory:" + logicFrame[i].Action);
    //         if (  list[i].action == "pushLogicFrame"){
    //             var data = eval( "(" + list[i].content + ")" );
    //             // console.log("rPushRoomHistory data:",data);
    //             self.rPushLogicFrame(data,"rPushRoomHistory");
    //             }
    //     }
    //     var requestPlayerResumeGame = new proto.pb.PlayerResumeGame();
    //     requestPlayerResumeGame.setRoomId(self.roomId);
    //     requestPlayerResumeGame.setSequenceNumber(self.sequenceNumber);
    //     requestPlayerResumeGame.setPlayerId(self.playerId);
    //     self.sendMsg("CS_PlayerResumeGame",requestPlayerResumeGame);
    // };


    // this.rServerPong = function(logicFrame){
    //     console.log("rServerPong:",logicFrame)
    // };
    // this.rServerPing = function(logicFrame){
    //     var now = Date.now();
    //
    //     var requestClientPong = new proto.pb.Pong();
    //     requestClientPong.setClientReceiveTime(now);
    //     requestClientPong.setAddTime(logicFrame.addTime);
    //     requestClientPong.setRttTimeout(logicFrame.rttTimeout);
    //     requestClientPong.setRttTimes(logicFrame.rttTimes);
    //     this.sendMsg("CS_Pong",requestClientPong);
    //     //     logicFrame.clientReceiveTime =  now
    // };
    // this.rPushRoomInfo = function(logicFrame){
    //     if(self.contentTypeDesc[self.contentType] =="protobuf"){
    //         logicFrame.playerList = logicFrame.playerListList;
    //     }
    //     self.initLocalGlobalVar(logicFrame);
    //     var requestRoomHistory = new proto.pb.RoomHistory();
    //     requestRoomHistory.setRoomId(self.roomId);
    //     requestRoomHistory.setSequenceNumberstart(0);
    //     requestRoomHistory.setSequenceNumberend(-1);
    //     requestRoomHistory.setPlayerId(self.playerId);
    //     // var history ={"roomId":self.roomId,"sequenceNumber":0,"playerId":self.playerId };
    //     // self.sendMsg("getRoomHistory",history);
    //     self.sendMsg("CS_RoomHistory",requestRoomHistory);
    //
    //     var readySignBntId = "ready_"+self.playerId;
    //     var hrefBody = "匹配成功，准备";
    //
    //     self.upOptBntHref(readySignBntId,hrefBody,self.ready);
    // };

    //
    // this.rOtherPlayerOffline = function(logicFrame){
    //     //房间内有其它玩家掉线了
    //     console.log("in test:",logicFrame.playerId,logicFrame)
    //     self.otherPlayerOffline = logicFrame.playerId;
    //     alert("其它玩家掉线了："+logicFrame.playerId +"，略等："+self.offLineWaitTime +"秒");
    //
    //     var tdId = self.tableId + "_" + self.playerLocation[logicFrame.playerId];
    //     var tdObj = $("#"+tdId)
    //     tdObj.css("background", "#A9A9A9");
    //     // var lightTd =self.getMapTdId(self.tableId,LocationStart,LocationEnd);
    //     // console.log(pre+"  "+lightTd);
    //     // var tdObj = $("#"+lightTd);
    //     // if(commands[i].playerId == playerId){
    //     //     tdObj.css("background", "green");
    //     // }else{
    //     //     tdObj.css("background", "red");
    //     // }
    // };
    // this.rGameOver = function(ev){
    //     clearInterval(self.pushLogicFrameLoopFunc);
    //     self.upOptBnt("游戏结束2",1)
    //     return alert("have player game end...");
    // };
    //





if ("WebSocket" in window) {
    console.log("browser support the websocket.");
}else {
    // 浏览器不支持 WebSocket
    alert("您的浏览器不支持 WebSocket!");
}

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