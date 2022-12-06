//=========
function Sync (playerId,token,GatewayConfig,DomIdPreObj,contentType,protocolType,playerInfo,actionMap,ruleConfig,mapSize,content_type_desc,protocol_type_desc) {
    this.status = 1;//1初始化 2等待准备 3运行中  4结束
    this.statusDesc = {
        1: "初始化",//
        2: "ws连接成功",//
        3: "登陆成功",//
        4: "登陆失败",//
        5: "未结束游戏，恢复中",
        6: "报名匹配中",//
        7: "取消匹配",//
        8: "准备中",//
        9: "开始游戏中",//
        10: "游戏已结束",//
        11: "连接关闭了",//
    };
    this.heartbeatLoopFunc = null;//心跳回调函数
    this.heartbeatHistory = []//记录每次心跳的数据
    this.heartbeatQueue = [];
    this.tableMax = mapSize;//地图的表格大小 5 X 5
    this.otherPlayerOffline = 0;//标识一下，有其它玩家调线
    this.pushLogicFrameLoopFunc = null;//定时循环 - 推送玩家操作函数
    this.playerOperationsQueue = [];//一个帧时间内，收集玩家的操作指令 队列
    this.closeFlag = 0;//关闭标识，0正常1手动关闭2后端关闭
    this.tableId = "";
    this.myClose = 0;//区分关闭事件，1. 是自己主动触发，2. S端(对方)触发的
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
//创建ws长连接，也算是入口函数，得选建立WS连接，才有后续的所有操作
Sync.prototype.create = function(){
    this.show("创建 ws 连接");
    if (this.status != 1 && this.status != 11){
        return this.show("status错误， status  !=  init or close",1)
    }
    var parent = this
    this.closeFlag = 0;//清空 关闭标识
    //根据帧数，计算出 秒数
    this.logicFrameLoopTimeMs = parseInt( 1000 / this.FPS);

    this.show("create new WebSocket"+this.hostUri ," FPS:",this.FPS , " ms:",this.logicFrameLoopTimeMs, " contentType:",this.contentType, " protocolType:",this.protocolType)
    //创建ws连接
    this.wsObj = new WebSocket(this.hostUri);
    //设置 连接建立成功回调
    this.wsObj.onopen = function(){
        parent.wsOpen();
    };
    //设置 关闭回调
    this.wsObj.onclose = function(ev){
        parent.close(ev)
    };
    //设置 接收S端消息-回调
    this.wsObj.onmessage = function(ev){
        parent.onmessage(ev)
    };
    //创建连接发生了错误
    this.wsObj.onerror = function(ev){
        alert("wsObj.onerror");
        parent.log("error:"+ev);
    };
}

//玩家操作 - 主动关闭
Sync.prototype.closeFD = function (){
    this.show("closeFD");
    this.myClearInterval();
    this.myClose = 1;
    this.wsObj.close();
};

//ws 接收到服务端关闭
Sync.prototype.close = function(ev){
    alert("receive server close:" +ev.code);
    this.myClearInterval();
    if (this.myClose == 1){//自己关闭的WS
        var reConnBntId = "reconn_"+this.playerId;
        var msg = "重连接";
        this.upOptBntHref(reConnBntId,msg,this.create.bind(this));
    }else{
        this.closeFlag = 2;
        this.upOptBnt("服务端关闭，游戏结束，连接断开",1)
    }
};
//连接成功后，会执行此函数
Sync.prototype.wsOpen = function(){
    this.show("onOpen , ws connect server : Success  ",1);
    this.upStatus(2);
    //强依赖，proto 文件
    var LoginObj = new proto.pb.Login();
    LoginObj.setToken(this.token) ;
    this.sendMsg("CS_Login",LoginObj);
};
//接收S端WS消息
//解析C端发送的数据，这一层，对于用户层的 content 数据不做处理
//1-4字节：当前包数据总长度，~可用于：TCP粘包的情况
//5字节：content type
//6字节：protocol type
//7字节 :服务Id
//8-9字节 :函数Id
//10-19：预留，还没想好，可以存sessionId，也可以换成UID
//19 以后为内容体
//结尾会添加一个字节：\f ,可用于 TCP 粘包 分隔
Sync.prototype.onmessage = function(ev){
    // console.log("onmessage , ev :",ev );
    var msg = new proto.pb.Msg();
    var msgObj = msg.toObject();
    var parent = this;

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
//接收S端消息后，开始进行路由，具体由哪个方法接收并处理
Sync.prototype.router = function(msgObj){
    var action = this.actionMap.server[msgObj.sidFid].func_name;
    this.showComplex("router , action:"+action,msgObj);
    //这里用了动态调用函数，减少代码量
    eval( "this."+action+"(msgObj.content)" );
}
//================== 以上是ws callback 基础函数 ====================================================


//获取当前C端到S端的RTT值，此值目前取的是过往10次成功的心跳的差值，最后再计算出一个平均值
Sync.prototype.getRTT = function(){
    if(!this.heartbeatHistory || this.heartbeatHistory.length <= 0){
        return 0;
    }

    var rttTotal = 0;
    var times = 0;
    for(var i=0;i<this.heartbeatHistory.length;i++){
        rttTotal += this.heartbeatHistory[i].server_receive_time -  this.heartbeatHistory[i].client_req_time;
        times++;
    }
    return  rttTotal / times;
}
//创建一个空的心跳包数据结构
Sync.prototype.createEmptyHeartbeatQueueElement = function (){
    var now = Date.now();
    var obj = new proto.pb.Heartbeat();
    obj.setClientReqTime(now);
    obj.setRequestId(now+ "");
    return obj;
}
//将心跳包添加到队列中
Sync.prototype.addHeartbeatQueueElement = function (){
    var obj = this.createEmptyHeartbeatQueueElement();
    var res = {"heartbeat":obj,"err":0};
    // console.log("heartbeatQueue:",this.heartbeatQueue);
    if(!this.heartbeatQueue.length){
        this.heartbeatQueue = [obj];
        return res;
    }
    //这里最多能有3次的 heartbeat
    if(this.heartbeatQueue.length > 3){
        res.err = 1;
        return res;
    }
    this.heartbeatQueue.push(obj);
    return res;
}

//心跳执行函数
Sync.prototype.heartbeat = function(){
    var obj = this.addHeartbeatQueueElement();
    if (obj.err ){//走到这里，证明3次心跳均未收到S端响应包，即：大概率是网络出了问题，可能是自己也可能是S端。
        this.closeFD()
    }else{
        // var requestClientHeartbeat = new proto.pb.Heartbeat();
        // requestClientHeartbeat.setClientReqTime(obj.client_req_time);
        // requestClientHeartbeat.setRequestId(obj.request_id);
        // this.sendMsg("CS_Heartbeat",requestClientHeartbeat);
        this.sendMsg("CS_Heartbeat",obj.heartbeat);
    }
};
//================== 以上都是心跳处理相关  ==================

//报名，主动触发
Sync.prototype.matchSign = function(){
    this.upStatus(6)

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
//取消报名，主动触发
Sync.prototype.cancelSign = function(){
    this.upStatus(7);

    var gameMatchPlayerCancel = new proto.pb.GameMatchPlayerCancel();
    gameMatchPlayerCancel.setRuleId(this.ruleId);
    gameMatchPlayerCancel.setGroupId(this.matchSignGroupId);
    this.sendMsg("CS_PlayerMatchSignCancel",gameMatchPlayerCancel);

    var matchSignBntId = "matchSign_"+this.playerId;
    var hrefBody = "连接成功，重新匹配报名";

    this.upOptBntHref(matchSignBntId,hrefBody,this.matchSign.bind(this));
};
//进入准备状态，主动触发
Sync.prototype.ready = function(){
    this.upStatus(8);
    var requestPlayerReady = new proto.pb.PlayerReady();
    requestPlayerReady.setPlayerId(this.playerId);
    requestPlayerReady.setRoomId(this.roomId);
    this.sendMsg("CS_PlayerReady",requestPlayerReady);

    this.noticeMsg("不能取消，只能等成功或超时");

    this.upOptBntHref("","等待其它玩家准备",null);
};
//两个玩家，位移碰撞了，触发了游戏结束
Sync.prototype.gameOverAndClear = function(){
    this.upStatus(10);


    var PlayerOver = new proto.pb.PlayerOver()
    PlayerOver.setRoomId(this.roomId);
    PlayerOver.setPlayerId(this.playerId);
    PlayerOver.setSequenceNumber(this.sequenceNumber);
    this.sendMsg("CS_PlayerOver", PlayerOver);

    var msg = "完犊子了，撞车了...这游戏得结束了....";
    this.noticeMsg(msg)

    this.myClearInterval();
    this.upOptBnt("游戏结束1",1)

    // return alert("完犊子了，撞车了...这游戏得结束了....");
};
//收集玩家游戏中的操作指令
Sync.prototype.playerCommandPush = function (){

    var PlayerOperations = new proto.pb.LogicFrame();
    PlayerOperations.setId(9999);//此值目前没用上
    PlayerOperations.setRoomId(this.roomId);
    PlayerOperations.setSequenceNumber(this.sequenceNumber);

    if(this.playerCommandPushLock == 1){//是否锁住了
        this.show("send msg lock...please wait server back login frame");
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
//清除定时函数
Sync.prototype.myClearInterval = function (){
    if(this.pushLogicFrameLoopFunc){
        window.clearInterval(this.pushLogicFrameLoopFunc);
        this.showComplex("myClearInterval this.pushLogicFrameLoopFunc ===================",this.pushLogicFrameLoopFunc)
        this.pushLogicFrameLoopFunc = null;
    }

    if(this.heartbeatLoopFunc){
        window.clearInterval(this.heartbeatLoopFunc);
        this.showComplex("myClearInterval this.heartbeatLoopFunc ===================",this.heartbeatLoopFunc)
        this.heartbeatLoopFunc = null;
    }
}
//游戏战局进入准备期，初始化-本地变量数据，等待所有玩家准备后，即开始游戏了
Sync.prototype.initLocalGlobalVar = function(EnterBattle){
    this.show("initLocalGlobalVar:",EnterBattle)
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
    this.show(str);
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


Sync.prototype.getMapTdId = function (tableId,i,j){
    return tableId + "_" + i +"_" + j;
}


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

//玩家报名，得有个groupId，先暂时由前端生成随机数，后期我想想如何处理
Sync.prototype.makeGroupId = function(){
    return Math.round(Math.random()*8999+1000);
}
//更新当前状态
Sync.prototype.upStatus = function(status){
    this.show("up status ,  old status:" + this.status   +  "("+ this.statusDesc[this.status]+") , new status:"  + status + "("+ this.statusDesc[status] + ")");
    this.status = status;
};

//==== debug ==========
Sync.prototype.showComplex = function(str,complexType,showNoticeMsg = 0){
    console.log(this.descPre + " " + str ,complexType)
    if (showNoticeMsg){
        this.noticeMsg(str);
    }
}

Sync.prototype.getPlayerDescById = function (id){
    return "(player_"+ id+")";
};

Sync.prototype.show = function(str,showNoticeMsg = 0){
    console.log(this.descPre + " " + str)
    if (showNoticeMsg){
        this.noticeMsg(str);
    }
}

Sync.prototype.noticeMsg = function(str){
    $("#"+this.domIdObj.msgNotice).html(str);
}
//==== debug ==========

if ("WebSocket" in window) {
    console.log("browser support the websocket.");
}else {
    // 浏览器不支持 WebSocket
    alert("您的浏览器不支持 WebSocket!");
}

