//======================== 以下都是接收S端消息通知 ========================
//S推送-登陆结果
Sync.prototype.SC_Login = function(LoginRes){
    if (LoginRes.code != 200) {
        this.upStatus(4);
        var msg = "登陆失败 , code: "+LoginRes.code + " , errMsg: "+LoginRes.errMsg;
        return this.show(msg);
    }

    this.show("登陆成功，开始获取玩家状态信息",1)

    this.upStatus(3);
    var playerBase = new proto.pb.PlayerBase();
    playerBase.setPlayerId(  this.playerId);
    this.sendMsg("CS_PlayerState",playerBase);

    //做一次PING 测试一下
    // var requestClientPing = new proto.pb.PingReq();
    // requestClientPing.setClientReqTime(Date.now());
    // this.sendMsg("CS_Ping",requestClientPing);

    //创建定时：心跳函数
    // this.heartbeatLoopFunc = setInterval(this.heartbeat.bind(this), this.ClientHeartbeatTime * 1000);
};
//S推送-心跳
Sync.prototype.SC_Heartbeat = function(Heartbeat){
    this.showComplex("rHeartbeat:",Heartbeat);
    if(!this.heartbeatQueue || this.heartbeatQueue.length <= 0){
        this.heartbeatQueue.push(Heartbeat);//从尾部添加一个元素
        return 1;
    }

    for(var i=0;i< this.heartbeatQueue.length ;i++){
        // console.log("=========",Heartbeat,"*********",this.heartbeatQueue[i]);
        if (this.heartbeatQueue[i].getRequestId() == Heartbeat.requestId){
            // console.log("heartbeatQueue :","ok======");
            if( this.heartbeatHistory && this.heartbeatHistory.length > 10){
                this.heartbeatHistory.shift()();//从首数删除一个元素
            }

            if(!this.heartbeatHistory || this.heartbeatHistory.length <= 0){
                this.heartbeatHistory = [Heartbeat];
            }else{
                this.heartbeatHistory.push(Heartbeat);//从尾部添加一个元素
            }

            this.heartbeatQueue = [];//只要有一次成功的响应，即把发送队列清空了
        }
    }
}
//S端响应自己的PING
Sync.prototype.SC_Pong = function(PongRes){
    this.showComplex("SC_Pong:",PongRes);
}
//S端发起了PING，C端要响应
Sync.prototype.SC_Ping = function(PingReq){
    var PongRes = new proto.pb.PongRes()
    this.sendMsg("CS_Pong:",PongRes);
}

//S推送-玩家当前的状态信息
Sync.prototype.SC_PlayerState = function(PlayerState){
    this.showComplex("ServerPlayerState ",PlayerState,1)
    this.serverPlayerState = PlayerState;
    //这里是有问题的，roomId我在外层写死了均为空，应该是动态从后端拿，且最好是短连接去取，回头优化
    if (this.serverPlayerState.roomId){
        this.upStatus(5);
        this.show("检测出，有未结束的一局游戏，开始恢复中...,先获取房间信息:rooId:"+this.serverPlayerState.roomId,1)
        var RoomBaseInfo = new proto.pb.RoomBaseInfo();
        RoomBaseInfo.setRoomId(PlayerState.roomId);
        this.sendMsg("CS_RoomBaseInfo",RoomBaseInfo);

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

    this.upOptBntHref(matchSignBntId,hrefBody,this.matchSign.bind(this));

    this.noticeMsg("抱歉，<准备时间>已超时");
    // return alert("抱歉，<准备时间>已超时");
}

Sync.prototype.SC_StartBattle = function(StartBattle){
    this.StartBattle(StartBattle,1);
};
//source:1正常开始2断线追帧完毕后，开始游戏
Sync.prototype.StartBattle  = function(StartBattle,source){
    this.upStatus(9);
    this.pushLogicFrameLoopFunc = setInterval(this.playerCommandPush.bind(this),this.logicFrameLoopTimeMs);
    this.showComplex("this.pushLogicFrameLoopFunc===================",this.pushLogicFrameLoopFunc);
    var exceptionOffLineId = "exceptionOffLineId"+this.playerId;
    var msg = "模拟，异常掉线";
    this.noticeMsg("游戏已开始，数据同步中...")
    this.upOptBntHref(exceptionOffLineId,msg,this.closeFD.bind(this));
}
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

//S推送-本局游戏的房间总信息，一般到这里，证明有玩家掉线了，重新连接，需要重新播放，追帧
Sync.prototype.SC_RoomBaseInfo = function(RoomBaseInfo ){
    if(this.contentTypeDesc[this.contentType] =="protobuf"){
        logicFrame.playerList = logicFrame.playerListList;
    }
    this.initLocalGlobalVar(RoomBaseInfo);

    var ReqRoomHistory = new proto.pb.ReqRoomHistory();
    ReqRoomHistory.setRoomId(this.roomId);
    ReqRoomHistory.setSequenceNumberstart(0);
    ReqRoomHistory.setSequenceNumberend(-1);
    ReqRoomHistory.setPlayerId(this.playerId);
    //这里拉回来的数量量略有点大，不好处理，改成短连接吧
    // this.sendMsg("CS_RoomHistory",ReqRoomHistory);
    // content = JSON.stringify(content);
    var postData = ReqRoomHistory.toObject();
    // console.log("SC_RoomBaseInfo postData: ",postData);
    var parent = this;
    $.ajax({
        dataType: "json",
        url: URI_MAP["game_frame_sync_history"],
        type: 'POST',
        headers: {
            "X-Token":this.token,
        },
        sync:false,
        data:JSON.stringify(postData) ,
        success: function (res) {
            var msgObj = {"sets":res.data};
            parent.SC_RoomHistory(msgObj);
            // console.log("game_frame_sync_history res: ",res);
            // for(var i=0;i<res.data.length;i++){
            //     var row = res.data[i].content;
            //     console.log(row);
            // }
            // var data = eval( "(" + res.data + ")" );
            // console.log("game_frame_sync_history data:",data);
        }
    });

    this.noticeMsg("游戏恢复中，已拉取回房间基础信息，开始拉取游戏历史逻辑帧...")

}
//接收本局游戏，其它游戏帧，开始录播追帧
Sync.prototype.SC_RoomHistory = function(RoomHistorySets ){
    if(this.contentTypeDesc[this.contentType] =="protobuf"){
        RoomHistorySets.sets = RoomHistorySets.setsList;
    }
    var list = RoomHistorySets.sets;
    for(var i=0;i<list.length;i++){
        // if (  list[i].event == "pushLogicFrame"){
            var data = eval( "(" + list[i].content + ")" );
            this.SC_LogicFrame(data,2);
        // }
    }
    this.noticeMsg("游戏恢复成功，通知其它玩家，并重新开始继续游戏...")
    //播放完成，告知其它玩家，我上线了，可以重新开始游戏
    var requestPlayerResumeGame = new proto.pb.PlayerResumeGame();
    requestPlayerResumeGame.setRoomId(this.roomId);
    requestPlayerResumeGame.setSequenceNumber(this.sequenceNumber);
    requestPlayerResumeGame.setPlayerId(this.playerId);
    this.sendMsg("CS_PlayerResumeGame",requestPlayerResumeGame);
    this.StartBattle(RoomHistorySets,2);

}

Sync.prototype.SC_OtherPlayerOver = function(PlayerOver){
    this.showComplex("OtherPlayerOver: ",PlayerOver);
    this.noticeMsg("有玩家死亡了 , player_id:"+ PlayerOver.playerId);
}

Sync.prototype.SC_GameOver = function(GameOver){
    this.myClearInterval()
    this.upOptBnt("游戏结束2",1)
    this.noticeMsg("某玩家触发了撞车，所以游戏要结束了",1)
    // return alert("have player game end...");
};

Sync.prototype.SC_KickOff = function(KickOff){
    this.showComplex("您被踢下线喽",KickOff,1)
    this.closeFD()
};

Sync.prototype.SC_OtherPlayerResumeGame = function(PlayerResumeGame){
    this.otherPlayerOffline = 0;//其它玩家恢复了
    if(PlayerResumeGame.playerId != this.playerId){
        var tdId = self.tableId + "_" + self.playerLocation[content.playerId];
        var tdObj = $("#"+tdId)
        tdObj.css("background", "red");
    }
}

Sync.prototype.SC_OtherPlayerOffline = function(OtherPlayerOffline){
    this.otherPlayerOffline = OtherPlayerOffline.playerId;//标识出，其它玩家掉线了，后面有移动操作，直接就给拒绝了，不让继续操作了
    this.show("其它玩家掉线了："+OtherPlayerOffline.playerId +"，略等："+this.offLineWaitTime +"秒",1);

    var tdId = this.tableId + "_" + this.playerLocation[OtherPlayerOffline.playerId];
    var tdObj = $("#"+tdId)
    tdObj.css("background", "#A9A9A9");
    // var lightTd =self.getMapTdId(self.tableId,LocationStart,LocationEnd);
    // console.log(pre+"  "+lightTd);
    // var tdObj = $("#"+lightTd);
    // if(commands[i].playerId == playerId){
    //     tdObj.css("background", "green");
    // }else{
    //     tdObj.css("background", "red");
    // }
};

