//=========
function Sync (playerId,token,data,DomIdPreObj,contentType,protocolType){
    var self = this;
    this.wsObj = null;//js内置ws 对象
    //ws 连接 s 端地址
    this.hostUri = "ws://"+data.outIp + ":"+ data.wsPort + data.wsUri;
    this.statusDesc = {
        1:"init",
        2:"wsLInkSuccess",
        3:"loginSuccess",
        4:"loginFailed",
        5:"matchSign",
        6:"cancelSign",
        7:"ready",
        8:"startBattle",
        9:"end",
        10:"close",
    };
    this.status = 1;//1初始化 2等待准备 3运行中  4结束
    this.playerId = playerId;//玩家ID
    this.matchGroupPeople = data.roomPeople;//一个副本的人数
    this.heartbeatLoopFunc = null;//心跳回调函数
    this.offLineWaitTime = data.offLineWaitTime;//lockStep 玩家掉线后，其它玩家等待最长时间
    //以上都是S端返回的一些配置值

    this.token = token;//玩家的凭证
    this.tableMax = data.mapSize;//地址的表格大小
    this.otherPlayerOffline = 0;//其它玩家调线
    this.pushLogicFrameLoopFunc = null;//定时循环 - 推送玩家操作函数
    this.playerOperationsQueue = [];//一个帧时间内，收集玩家的操作指令 队列
    this.closeFlag = 0;//关闭标识，0正常1手动关闭2后端关闭

    this.tableId = "";
    this.domIdObj = DomIdPreObj ;
    this.playerLocation = new Object();//每个玩家的位置信息
    this.operationsInc = 0;//玩家操作指令自增ID
    this.logicFrameLoopTimeMs = 0;//毫秒，多少时间内收集一次玩家操作，推送到S端
    this.FPS = data.fps;//每秒多少逻辑帧
    this.playerCommandPushLock = 0;
    //下面是帧同步初始化信息，是由S端供给
    this.roomId = "";
    this.actionMap = data.actionMap;
    this.sequenceNumber = 0;
    this.randSeek = 0;
    this.sessionId = "";
    this.contentType = contentType;//protobuf|json
    this.protocolType = protocolType;//tcp|ws
    //入口函数，必须得先建立连接后，都有后续的所有操作
    this.create  = function(){
        // console.log("entrance : create ws , this status :",self.status);
        if (self.status != 1 && self.status != 10){
            return alert(" status !=  init or close");
        }
        self.closeFlag = 0;//清空 关闭标识
        //根据帧数，计算出 秒数
        self.logicFrameLoopTimeMs = parseInt( 1000 / self.FPS);

        console.log("create new WebSocket"+self.hostUri ," FPS:",self.FPS , " ms:",self.logicFrameLoopTimeMs, " contentType:",self.contentType, " protocolType:",self.protocolType)
        //创建ws连接
        self.wsObj = new WebSocket(self.hostUri);
        //设置 关闭回调
        self.wsObj.onclose = function(ev){
            self.onclose(ev)
        };
        //设置 接收S端消息-回调
        self.wsObj.onmessage = function(ev){
            self.onmessage(ev)
        };
        //设置 连接建立成功回调
        self.wsObj.onopen = function(){
            self.wsOpen();
        };
        self.wsObj.onerror = function(ev){
            alert("wsObj.onerror");
            console.log("error:"+ev);
        };
    };
    //连接成功后，会执行此函数
    this.wsOpen = function(){
        console.log("onOpen : ws connect server : Success  ");
        this.upStatus(2);
        //强依赖，proto 文件
        // var requestLoginObj = new proto.pb.RequestLogin();
        var requestLoginObj = new proto.pb.Login();
        requestLoginObj.setToken(self.token) ;
        this.sendMsg("login",requestLoginObj);
    };
    //更新当前状态
    this.upStatus = function(status){
        console.log("up status:"," old status:",this.status, " new status:" ,status,this.statusDesc[status]);
        this.status = status;
    };
    this.sendMsg =  function ( action,contentObj  ){
        var id = self.getActionId(action,"client");
        console.log( self.descPre , " <sendMsg>" ,  " actionName: "+action , " actionName:" , id  , " content:" ,contentObj.toObject());
        var content = null;
        if (contentTypeDesc[self.contentType] == "json"){
            content = contentObj.toObject();
            content = JSON.stringify(content);
            if(action == "playerOperations"){
                console.log(content);
                content = content.replace("operationsList","operations");
                console.log(content);
            }
            var protocolCtrl = contentType +  "" + protocolType + id;
            if (action == "login" ){
                content = protocolCtrl + content;
            }else{
                content = protocolCtrl + self.sessionId +  content;
            }
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
            console.log("<sendMsg final>",content);

            self.wsObj.send(content);
        }else if ( contentTypeDesc[self.contentType]  == "protobuf"){
            content = contentObj.serializeBinary();
            var protocolCtrl = contentType +  "" + protocolType + id;
            if (action != "login" ){
                protocolCtrl = protocolCtrl + self.sessionId  ;
            }
            var idBinary = stringToUint8Array(protocolCtrl);
            var mergedArray = new Uint8Array(idBinary.length + content.length);
            mergedArray.set(idBinary);
            mergedArray.set(content, idBinary.length);

            self.wsObj.send(mergedArray);
        }
    };
    //更新操作按钮文字，解除 点击事件
    this.upOptBnt = function(content,clearClick){
        $("#"+self.domIdObj.optBntId).html(content);
        if(clearClick == 1){
            $("#"+self.domIdObj.optBntId).unbind("click");
        }
    };
    //玩家操作 - 主动关闭
    this.closeFD = function (){
        console.log("closeFD");
        // window.clearInterval(self.heartbeatLoopFunc);
        clearInterval(self.pushLogicFrameLoopFunc);
        self.myClose = 1;
        self.wsObj.close();
    };
    //ws 接收到服务端关闭
    this.onclose = function(ev){
        alert("receive server close:" +ev.code);
        clearInterval(self.pushLogicFrameLoopFunc);
        self.upStatus(10);
        // window.clearInterval(self.heartbeatLoopFunc);
        if (self.myClose == 1){
            var reConnBntId = "reconn_"+self.playerId;
            var msg = "重连接";
            self.upOptBntHref(reConnBntId,msg,self.create);
        }else{
            self.closeFlag = 2;
            self.upOptBnt("服务端关闭，游戏结束，连接断开",1)
        }
    };
    //心跳执行函数
    this.heartbeat = function(){
        var now = Date.now();
        var requestClientHeartbeat = new proto.myproto.RequestClientHeartbeat();
        requestClientHeartbeat.setTime(now);
        this.sendMsg("clientHeartbeat",requestClientHeartbeat);
    };
    //接收S端WS消息
    this.onmessage = function(ev){
        var pre = self.descPre;
        console.log("onmessage:"+ pre + " " +ev.data);
        var msgObj = self.newMsgObj();

        if (contentTypeDesc[self.contentType] == 'protobuf'){
            var reader = new FileReader();
            reader.readAsArrayBuffer(ev.data);
            reader.onloadend = function(e) {
                var dataBuffer = new Uint8Array(reader.result);

                msgObj.contentType = processBufferString(dataBuffer,0,1);
                msgObj.protocolType = processBufferString(dataBuffer,1,2);
                msgObj.actionId = processBufferString(dataBuffer,2,6);
                msgObj.sessionId = processBufferString(dataBuffer,6,38);
                var content = processBuffer(dataBuffer,38);
                msgObj.action = self.getActionName(msgObj.actionId,"server")
                //首字母转大写
                var actionLow = msgObj.action.substring(0, 1).toUpperCase() + msgObj.action.substring(1)
                //拼接成最终classname
                var className =  "proto.myproto.Response" + actionLow;
                var responseProtoClass = eval(className);
                //将进制流转换成对象
                msgObj.content =  responseProtoClass.deserializeBinary(content).toObject();
                self.router(msgObj);
            };
        }else if(contentTypeDesc[self.contentType] == "json"){
            msgObj.contentType = ev.data.substr(0,1);
            msgObj.protocolType = ev.data.substr(1,1);
            msgObj.actionId = ev.data.substr(2,4);
            msgObj.sessionId = ev.data.substr(6,37);
            msgObj.content = ev.data.substr(38);
            // console.log(pre + " contentType:" +contentType+ ",protocolType:" + protocolType  +" ,actionId:"+actionId +",sessionId:" +sessionId + " ,content:",content );
            msgObj.action = self.getActionName(msgObj.actionId,"server")
            msgObj.content =  eval("("+msgObj.content+")");
            self.router(msgObj);
        }else{
            return alert("contentType err");
        }
    };
    this.parserContent = function(){

    };
    this.router = function(msgObj){
        console.log("router:",msgObj)
        var action  = msgObj.action;
        var content = msgObj.content;

        // var actionUp = msgObj.action.substring(0, 1).toUpperCase() + msgObj.action.substring(1)
        // var selfFuncName =  "r" + actionUp;
        // console.log("<router> ",selfFuncName,msgObj);
        // eval("self."+selfFuncName+"("+msgObj.content+")");
        // self.call("tttt","bbb")
        // return 1;
        // console.log("router:",action,content)
        if ( action == 'loginRes' ) {
            self.rLoginRes(content);
        }else if( action == 'serverPing'){//获取一个当前玩家的状态，如：是否有历史未结束的游戏
            self.rServerPing(content);
        }else if ( action == 'startBattle' ){
            self.rStartBattle(content);
        }else if ( action == 'pushRoomInfo' ){
            self.rPushRoomInfo(content);
        }else if ( action == 'otherPlayerOffline' ){
            self.rOtherPlayerOffline(content);
        }else if ( action == 'enterBattle' ){
            self.rEnterBattle(content);
        }else if( "gameOver" == action){
            self.rGameOver(content);
        }else if( "kickOff" == action){
            self.rKickOff(content);
        }else if( "pushLogicFrame" == action){
            self.rPushLogicFrame(content,"router")
        }else if( "readyTimeout" == action){
            self.rReadyTimeout(content)
        }else if( "serverPong" == action){
            self.rServerPong(content)
        }else if( "otherPlayerResumeGame" == action){
            self.rOtherPlayerResumeGame(content)
        }else if( "pushRoomHistory" == action){
            self.rPushRoomHistory(content);
            // alert("接收到，玩家-房间-历史操作记录~");
        }else{
            return alert("action error.");
        }
    };
    //=================== 以下都是 接收S端的处理函数========================================
    this.rOtherPlayerResumeGame = function(content){
        if(content.playerId != self.playerId){
            var tdId = self.tableId + "_" + self.playerLocation[content.playerId];
            var tdObj = $("#"+tdId)
            tdObj.css("background", "red");
        }
    };
    this.rReadyTimeout= function(logicFrame){
        console.log("rReadyTimeout:",logicFrame);

        this.upStatus(3);

        var matchSignBntId = "matchSign_"+self.playerId;
        var hrefBody = "连接成功，匹配报名";

        self.upOptBntHref(matchSignBntId,hrefBody,self.matchSign);

        return alert("抱歉，<准备时间>已超时");
    };
    this.rPushRoomHistory = function(logicFrame){
        console.log("rPushRoomHistory:");
        if(contentTypeDesc[self.contentType] =="protobuf"){
            logicFrame.list = logicFrame.listList;
        }
        var list = logicFrame.list;
        for(var i=0;i<list.length;i++){
            // console.log( "rPushRoomHistory:" + logicFrame[i].Action);
            if (  list[i].action == "pushLogicFrame"){
                var data = eval( "(" + list[i].content + ")" );
                // console.log("rPushRoomHistory data:",data);
                self.rPushLogicFrame(data,"rPushRoomHistory");
                }
        }
        var requestPlayerResumeGame = new proto.myproto.RequestPlayerResumeGame();
        requestPlayerResumeGame.setRoomId(self.roomId);
        requestPlayerResumeGame.setSequenceNumber(self.sequenceNumber);
        requestPlayerResumeGame.setPlayerId(self.playerId);
        self.sendMsg("playerResumeGame",requestPlayerResumeGame);
    };
    this.upOptBntHref = function(domId,value,clickCallback){
        var bntContent = "<a href='javascript:void(0);' onclick='' id='"+domId+"'>"+value+"</a>";
        self.upOptBnt(bntContent, 1);
        $("#"+domId).click(clickCallback);
    };
    this.rLoginRes = function(logicFrame){
        if (logicFrame.code != 200) {
            this.upStatus(4);
            return alert("loginRes failed!!!"+logicFrame.code + " , "+logicFrame.errMsg);
        }

        var playerConnInfo = logicFrame.player;
        self.sessionId = playerConnInfo.sessionId;

        var now = Date.now();
        var requestClientPing = new proto.myproto.RequestClientPing();
        requestClientPing.setAddTime(now);
        this.sendMsg("clientPing",requestClientPing);
        this.upStatus(3);

        if (playerConnInfo.roomId){
            alert("检测出，有未结束的一局游戏，开始恢复中...,先获取房间信息:rooId:"+playerConnInfo.roomId);
            var requestGetRoom = new proto.myproto.RequestGetRoom();
            requestGetRoom.setRoomId(playerConnInfo.roomId);
            requestGetRoom.setPlayerId(playerId);
            self.sendMsg("getRoom",requestGetRoom);
            //     var msg = {"roomId":playerConnInfo.roomId,"playerId":playerId};
        }else{
            var matchSignBntId = "matchSign_"+self.playerId;
            var hrefBody = "连接成功，匹配报名";

            self.upOptBntHref(matchSignBntId,hrefBody,self.matchSign);
        }
        // self.heartbeatLoopFunc = setInterval(self.heartbeat, 5000);
    };
    this.rServerPong = function(logicFrame){
        console.log("rServerPong:",logicFrame)
    };
    this.rServerPing = function(logicFrame){
        var now = Date.now();

        var requestClientPong = new proto.myproto.RequestClientPong();
        requestClientPong.setClientReceiveTime(now);
        requestClientPong.setAddTime(logicFrame.addTime);
        requestClientPong.setRttTimeout(logicFrame.rttTimeout);
        requestClientPong.setRttTimes(logicFrame.rttTimes);
        this.sendMsg("clientPong",requestClientPong);
        //     logicFrame.clientReceiveTime =  now
    };
    this.rStartBattle = function(logicFrame){
        self.upStatus(8);
        self.pushLogicFrameLoopFunc = setInterval(self.playerCommandPush,self.logicFrameLoopTimeMs);
        var exceptionOffLineId = "exceptionOffLineId"+self.playerId;
        var msg = "异常掉线";
        self.upOptBntHref(exceptionOffLineId,msg,self.closeFD);
    };
    this.rPushRoomInfo = function(logicFrame){
        if(contentTypeDesc[self.contentType] =="protobuf"){
            logicFrame.playerList = logicFrame.playerListList;
        }
        self.initLocalGlobalVar(logicFrame);
        var requestRoomHistory = new proto.myproto.RequestRoomHistory();
        requestRoomHistory.setRoomId(self.roomId);
        requestRoomHistory.setSequenceNumberstart(0);
        requestRoomHistory.setSequenceNumberend(-1);
        requestRoomHistory.setPlayerId(self.playerId);
        // var history ={"roomId":self.roomId,"sequenceNumber":0,"playerId":self.playerId };
        // self.sendMsg("getRoomHistory",history);
        self.sendMsg("roomHistory",requestRoomHistory);

        var readySignBntId = "ready_"+self.playerId;
        var hrefBody = "匹配成功，准备";

        self.upOptBntHref(readySignBntId,hrefBody,self.ready);
    };
    this.rPushLogicFrame = function(logicFrame,source){//接收S端逻辑帧
        console.log("logicFrame:",logicFrame);
        var pre = self.descPre;
        if(contentTypeDesc[self.contentType] =="protobuf"){
            if (source != "rPushRoomHistory"){
                logicFrame.operations = logicFrame.operationsList;
            }
        }
        var operations = logicFrame.operations;
        self.sequenceNumber  = logicFrame.sequenceNumber;
        $("#"+self.domIdObj.seqId).html(self.sequenceNumber);

        self.playerCommandPushLock = 0;

        console.log("rPushLogicFrame ,sequenceNumber:"+self.sequenceNumber+ ", operationsLen:" +  operations.length);
        for(var i=0;i<operations.length;i++){
            playerId= operations[i].playerId;
            var str = pre + " i=i , id: "+operations[i].id + " , event:"+operations[i].event + " , value:"+ operations[i].value + " , playerId:" + playerId;
            console.log(str);
            if (operations[i].event == "move"){
                var LocationArr = operations[i].value.split(",");
                var LocationStart = LocationArr[0];
                var LocationEnd = LocationArr[1];

                // var lightTd = "map"+id +"_"+LocationStart + "_"+LocationEnd;
                var lightTd =self.getMapTdId(self.tableId,LocationStart,LocationEnd);
                console.log(pre+"  "+lightTd);
                var tdObj = $("#"+lightTd);
                if(playerId == self.playerId){
                    tdObj.css("background", "green");
                }else{
                    tdObj.css("background", "red");
                }
                var playerLocation = self.playerLocation;
                if (playerLocation[playerId] == "empty"){
                    //证明是第一次移动，没有之前的数据
                }else{
                    // playerLocation = getPlayerLocation(playerId);
                    // alert(commands[i].playerId);
                    var playerLocationArr = playerLocation[playerId].split("_");
                    //非首次移动，这次移动后，要把之前所有位置清掉
                    var lightTd = self.getMapTdId(self.tableId,playerLocationArr[0],playerLocationArr[1]);
                    var tdObj = $("#"+lightTd);
                    tdObj.css("background", "");
                }
                playerLocation[playerId] = LocationStart + "_"+LocationEnd;
            }else if(operations[i].event == "empty"){

            }
        }
        // self.sendPlayerLogicFrameAck( self.sequenceNumber)
    };

    this.rOtherPlayerOffline = function(logicFrame){
        //房间内有其它玩家掉线了
        console.log("in test:",logicFrame.playerId,logicFrame)
        self.otherPlayerOffline = logicFrame.playerId;
        alert("其它玩家掉线了："+logicFrame.playerId +"，略等："+self.offLineWaitTime +"秒");

        var tdId = self.tableId + "_" + self.playerLocation[logicFrame.playerId];
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
    this.rEnterBattle = function(logicFrame){
        if(contentTypeDesc[self.contentType] =="protobuf"){
            logicFrame.playerList = logicFrame.playerListList;
        }

        self.initLocalGlobalVar(logicFrame);

        var readySignBntId = "ready_"+self.playerId;
        var hrefBody = "匹配成功，准备";

        self.upOptBntHref(readySignBntId,hrefBody,self.ready);
    };
    this.rGameOver = function(ev){
        clearInterval(self.pushLogicFrameLoopFunc);
        self.upOptBnt("游戏结束2",1)
        return alert("have player game end...");
    };

    this.rKickOff = function(ev){
        return alert("您被踢下线喽~");
    };

    //=================== 以上都是 接收S端的处理函数========================================
    this.initLocalGlobalVar = function(logicFrame){
        var pre = self.descPre;
        console.log("initLocalGlobalVar:",logicFrame)
        for(var i=0;i<logicFrame.playerList.length;i++){
            self.playerLocation[""+logicFrame.playerList[i].id+""] = "empty"
        }
        // return 1;
        self.randSeek  = logicFrame.randSeek;
        $("#"+self.domIdObj.randSeekId).html(self.randSeek);


        self.sequenceNumber  = logicFrame.sequenceNumber;
        $("#"+self.domIdObj.seqId).html(self.sequenceNumber);

        self.roomId = logicFrame.roomId;
        $("#"+self.domIdObj.roomId).html(self.roomId);

        var str =  pre+", roomId:" +logicFrame.roomId+ ", RandSeek:"+    self.randSeek +" , SequenceNumber"+    self.sequenceNumber ;
        console.log(str);
    };

    this.getMap = function (tableId) {
        // var tableDivPre = "map";
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
                inc++;
            }
            // alert(trTemp);
            trTemp.appendTo(tableObj);
        }
    };
    this.ready = function(){
        self.upStatus(7);
        var requestPlayerReady = new proto.myproto.RequestPlayerReady();
        requestPlayerReady.setPlayerId(self.playerId);
        requestPlayerReady.setRoomId(self.roomId);
        self.sendMsg("playerReady",requestPlayerReady);
        self.upOptBntHref("","等待其它玩家准备",this.ready);
    };
    this.cancelSign = function(){
        self.upStatus(6);

        var requestPlayerMatchSignCancel = new proto.myproto.RequestPlayerMatchSignCancel();
        requestPlayerMatchSignCancel.setPlayerId(self.playerId);
        self.sendMsg("playerMatchSignCancel",requestPlayerMatchSignCancel);

        var matchSignBntId = "matchSign_"+self.playerId;
        var hrefBody = "连接成功，匹配报名";

        self.upOptBntHref(matchSignBntId,hrefBody,self.matchSign);
    };
    this.matchSign = function(){
        self.upStatus(5)

        var requestPlayerMatchSign = new proto.myproto.RequestPlayerMatchSign();
        requestPlayerMatchSign.setPlayerId(self.playerId);
        self.sendMsg("playerMatchSign",requestPlayerMatchSign);

        var cancelBntId = "cancelSign_"+self.playerId;
        var hrefBody = "取消匹配报名";

        self.upOptBntHref(cancelBntId,hrefBody,self.cancelSign);
    };
    this.move = function ( dirObj ){

        if (self.otherPlayerOffline){
            return alert("其它玩家掉线了，请等待一下...");
        }

        if (self.closeFlag > 0 ){
            return alert("WS FD 已关闭...");
        }

        if (self.status != 8){
            return alert("status err , != startBattle ， 游戏还未开始，请等待一下...");
        }

        var dir = dirObj.data;
        var playerLocation = self.playerLocation;
        var nowLocationStr = playerLocation[self.playerId]
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
            if(nowLocationLine == self.tableMax - 1 ){
                return alert("nowLocationLine == "+ self.tableMax);
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
            if(nowLocationColumn ==  self.tableMax - 1 ){
                return alert("nowLocationColumn == "+ self.tableMax);
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
                 return self.gameOverAndClear()
            }
        }

        console.log("dir:"+dir+"oldLocation"+nowLocationStr+" , newLocation:"+newLocation);
        self.playerOperationsQueue.push({"id":self.operationsInc,"event":"move","value":newLocation,"playerId":self.playerId});
        self.operationsInc++;
        var playerLocationArr = playerLocation[self.playerId].split("_");
        var lightTd = self.getMapTdId(self.tableId,playerLocationArr[0],playerLocationArr[1]);
        var tdObj = $("#"+lightTd);
        tdObj.css("background", "");
    }
    this.playerCommandPush = function (){
        var requestPlayerOperations = new proto.myproto.RequestPlayerOperations();
        requestPlayerOperations.setId(3);
        requestPlayerOperations.setRoomId(self.roomId);
        requestPlayerOperations.setSequenceNumber(self.sequenceNumber);

        if(self.playerCommandPushLock == 1){
            console.log("send msg lock...please wait server back login frame");
            return
        }
        if (self.playerOperationsQueue.length > 0){
            // {"id":self.operationsInc,"event":"move","value":newLocation,"playerId":self.playerId}
            var operations = new Array(self.playerOperationsQueue.length)
            for(var i=0;i<self.playerOperationsQueue.length;i++){
                var operation = new proto.myproto.Operation();
                operation.setId(self.playerOperationsQueue[i].id);
                operation.setEvent(self.playerOperationsQueue[i].event);
                operation.setValue(self.playerOperationsQueue[i].value);
                operation.setPlayerId(self.playerId),
                // operations.push(operation);
                operations[i] = operation;
            }
            requestPlayerOperations.setOperationsList(operations);
            // loginFrame.operations = self.playerOperationsQueue;
            self.playerOperationsQueue = [];//将当前队列里的，当前帧的数据，清空
        }else{
            // window.clearInterval(self.pushLogicFrameLoopFunc);

            var operations = new Array(1);
            var operation = new proto.myproto.Operation();
            operation.setId(1);
            operation.setEvent("empty");
            operation.setValue("-1");
            operation.setPlayerId(self.playerId),
            // operations.push(operation);
            operations[0] = operation;
            // console.log(operations)
            requestPlayerOperations.setOperationsList(operations);
            // var emptyCommand = [{"id":1,"event":"empty","value":"-1","playerId":self.playerId}];
            // loginFrame.operations = emptyCommand;
        }
        self.sendMsg("playerOperations",requestPlayerOperations);
        self.playerCommandPushLock = 1;
        // self.sendMsg("playerOperations",loginFrame);
    };

    //两个玩家，位移碰撞了，触发了游戏结束
    this.gameOverAndClear = function(){
        this.upStatus(9);

        var requestGameOver = new proto.myproto.RequestGameOver()
        requestGameOver.setRoomId(self.roomId);
        requestGameOver.setSequenceNumber(self.sequenceNumber);
        requestGameOver.setResult("ccccccWin");
        this.sendMsg("gameOver", requestGameOver);

        clearInterval(self.pushLogicFrameLoopFunc);
        self.upOptBnt("游戏结束1",1)

        return alert("完犊子了，撞车了...这游戏得结束了....");
    };

    this.getPlayerDescById = function (id){
        return "(player_"+ id+")";
    };
    this.getMapTdId = function (tableId,i,j){
        return tableId + "_" + i +"_" + j;
    }

    this.getActionId = function (action,category){
        var data = self.actionMap[category];
        for(let key  in data){
            if (data[key].action == action){
                return data[key].id;
            }
        }
        alert(action + ": no match");
        return "";
    };
    this.getActionName = function (actionId,category){
        var data = self.actionMap[category];
        return data[actionId].action;
    };

    this.descPre = this.getPlayerDescById(playerId);
    this.newMsgObj = function (){
        var msg = new proto.myproto.Msg();
        return msg.toObject();
    };
};


if ("WebSocket" in window) {
    console.log("browser has websocket netway.");
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