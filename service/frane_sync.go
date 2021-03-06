package service

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"strconv"
	"time"
	"zgoframe/protobuf/pb"
	"zgoframe/util"
)

const (
	LOCK_MODE_PESSIMISTIC = 1 //囚徒
	LOCK_MODE_OPTIMISTIC  = 2 //乐观

	//一个副本的，一条消息的，同步状态
	PLAYERS_ACK_STATUS_INIT = 1 //初始化
	PLAYERS_ACK_STATUS_WAIT = 2 //等待玩家确认
	PLAYERS_ACK_STATUS_OK   = 3 //所有玩家均已确认

	PLAYER_STATUS_ONLINE  = 1 //在线
	PLAYER_STATUS_OFFLINE = 2 //离线

	PLAYER_NO_READY  = 1
	PLAYER_HAS_READY = 2

	CTX_DONE_PRE = "ctx.done() "
)

type FrameSync struct {
	Option    FrameSyncOption
	CloseChan chan int
}

type FrameSyncOption struct {
	FPS        int32 `json:"fps"`      //frame pre second
	LockMode   int32 `json:"lockMode"` //锁模式，乐观|悲观
	MapSize    int32 `json:"mapSize"`  //地址大小，给前端初始化使用
	Store      int32 `json:"store"`
	Log        *zap.Logger
	RoomManage *RoomManager
	Netway     *util.NetWay
}

//断点调试
var debugBreakPoint int

//索引表，PlayerId=>RoomId
//var mySyncPlayerRoom map[int32]string
//同步 - 逻辑中的自增ID - 默认值
var logicFrameMsgDefaultId int32

//同步 - 逻辑中 - 操作帧的 自增ID - 默认值
var operationDefaultId int32

//var RoomSyncMetricsPool map[string]RoomSyncMetrics

func NewFrameSync(Option FrameSyncOption) *FrameSync {
	Option.Log.Info("NewSync instance")
	sync := new(FrameSync)
	sync.Option = Option

	debugBreakPoint = 0
	logicFrameMsgDefaultId = 16
	operationDefaultId = 32
	//统计
	//RoomSyncMetricsPool = make(map[string]RoomSyncMetrics)

	if sync.Option.FPS > 1000 {
		Option.Log.Error("fps > 1000 ms")
	}

	//sync.initPool()
	sync.CloseChan = make(chan int)
	return sync
}

func (sync *FrameSync) SetNetway(netway *util.NetWay) {
	sync.Option.Netway = netway
}

//func (sync *FrameSync) initPool() {
//	if sync.Option.Store == 1 {
//
//	}
//}

//进入战后，场景渲染完后，进入准确状态
func (sync *FrameSync) PlayerReady(requestPlayerReady pb.PlayerReady, conn *util.Conn) {
	//roomId := myPlayerManager.GetRoomIdByPlayerId(requestPlayerReady.PlayerId)
	roomId := conn.RoomId
	sync.Option.Log.Debug(" roomId :" + roomId)
	room, empty := sync.Option.RoomManage.GetById(roomId)
	if empty {
		sync.Option.Log.Error("playerReady getPoolElementById empty:" + roomId)
		return
	}
	room.PlayersReadyListRWLock.Lock()
	room.PlayersReadyList[requestPlayerReady.PlayerId] = PLAYER_HAS_READY
	room.PlayersReadyListRWLock.Unlock()
	playerReadyCnt := 0
	//sync.Option.Log.Info("room.PlayersReadyList:", room.PlayersReadyList)
	room.PlayersReadyListRWLock.RLock()
	for _, v := range room.PlayersReadyList {
		if v == PLAYER_HAS_READY {
			playerReadyCnt++
		}
	}
	room.PlayersReadyListRWLock.RUnlock()

	if playerReadyCnt < len(room.PlayersReadyList) {
		sync.Option.Log.Error("now ready cnt :" + strconv.Itoa(playerReadyCnt) + " ,so please wait other players...")
		return
	}
	responseStartBattle := pb.StartBattle{
		SequenceNumberStart: int32(0),
	}
	sync.boardCastFrameInRoom(room.Id, "startBattle", &responseStartBattle)
	room.UpStatus(ROOM_STATUS_EXECING)
	room.StartTime = int32(util.GetNowTimeSecondToInt())

	//RoomSyncMetricsPool[roomId] = RoomSyncMetrics{}

	sync.testFirstLogicFrame(room)
	room.ReadyCloseChan <- 1
	//开启定时器，推送逻辑帧
	go sync.logicFrameLoop(room)

}

func (sync *FrameSync) StartOne(room *Room) {
	sync.Option.Log.Warn("start a new game:")

	responseClientInitRoomData := pb.EnterBattle{
		Status:         room.Status,
		AddTime:        room.AddTime,
		RoomId:         room.Id,
		SequenceNumber: -1,
		PlayerIds:      room.PlayerIds,
		RandSeek:       room.RandSeek,
		Time:           time.Now().UnixNano() / 1e6,
		UdpPort:        sync.Option.Netway.Option.UdpPort,
	}

	sync.boardCastInRoom(room.Id, "enterBattle", &responseClientInitRoomData)
	go sync.checkReadyTimeout(room)
}

//检查 所有玩家是否 都已准确，超时了
func (sync *FrameSync) checkReadyTimeout(room *Room) {
	for {
		select {
		case <-room.ReadyCloseChan:
			goto end
		default:
			now := util.GetNowTimeSecondToInt()
			if now > int(room.ReadyTimeout) {
				sync.Option.Log.Error("room ready timeout id :" + room.Id)
				requestReadyTimeout := pb.ReadyTimeout{
					RoomId: room.Id,
				}
				sync.boardCastInRoom(room.Id, "readyTimeout", &requestReadyTimeout)
				sync.roomEnd(room.Id, 0)
				goto end
			}
			time.Sleep(time.Second * 1)
		}
	}
end:
	sync.Option.Log.Warn("checkReadyTimeout loop routine close")
}

//
////在一个房间内，搜索一个用户
//func (sync *FrameSync) getPlayerByIdInRoom(playerId int32, room *Room) (myplayer *pb.Player, empty bool) {
//	for _, player := range room.PlayerList {
//		if player.Id == playerId {
//			return player, false
//		}
//	}
//	return myplayer, true
//}

//同步 玩家 操作 定时器
func (sync *FrameSync) logicFrameLoop(room *Room) {
	fpsTime := 1000 / sync.Option.FPS
	i := 0
	for {
		select {
		case <-room.CloseChan:
			goto end
		default:
			sleepMsTime := sync.logicFrameLoopReal(room, fpsTime)
			sleepMsTimeD := time.Duration(sleepMsTime)
			if sleepMsTime > 0 {
				time.Sleep(sleepMsTimeD * time.Millisecond)
			}
			i++
			//if i > 10{
			//	zlib.ExitPrint(1111)
			//}
		}
	}
end:
	sync.Option.Log.Warn("pushLogicFrame loop routine close")
}

//同上
func (sync *FrameSync) logicFrameLoopReal(room *Room, fpsTime int32) int32 {
	queue := room.PlayersOperationQueue
	end := queue.Len()
	//mylog.Debug("logicFrameLoopReal len:",end)
	if end <= 0 {
		return fpsTime
	}

	if sync.Option.LockMode == LOCK_MODE_PESSIMISTIC {
		ack := 0
		room.PlayersAckListRWLock.RLock()
		for _, v := range room.PlayersAckList {
			if v == 1 {
				ack++
			}
		}
		room.PlayersAckListRWLock.RUnlock()
		if ack < len(room.PlayersAckList) {
			sync.Option.Log.Error("还有玩家未发送操作记录,当前确认人数:" + strconv.Itoa(ack))
			return fpsTime
		}
	}

	room.SequenceNumber++

	logicFrame := pb.LogicFrame{
		Id:             0,
		RoomId:         room.Id,
		SequenceNumber: int32(room.SequenceNumber),
	}
	var operations []*pb.Operation
	i := 0
	element := queue.Front()
	for {
		if i >= end {
			break
		}
		operationsValueInterface := element.Value
		operationsValue := operationsValueInterface.(string)
		var elementOperations []pb.Operation
		err := json.Unmarshal([]byte(operationsValue), &elementOperations)
		if err != nil {
			sync.Option.Log.Error("queue json.Unmarshal err :" + err.Error())
		}
		//mylog.Debug(operationsValue,"elementOperations",elementOperations)
		for j := 0; j < len(elementOperations); j++ {
			//if elementOperations[j].Event != "empty"{
			//	mylog.Debug("elementOperations j :",elementOperations[j])
			//	debugBreakPoint = 1
			//}
			operations = append(operations, &elementOperations[j])
		}

		tmpElement := element.Next()
		queue.Remove(element)
		util.MyPrint("tmpElement:", tmpElement, " len:", queue.Len(), "i:", i)
		element = tmpElement

		i++
	}
	sync.upSyncRoomPoolElementPlayersAckStatus(room.Id, PLAYERS_ACK_STATUS_OK)

	util.MyPrint("operations:", operations)
	logicFrame.Operations = operations
	sync.boardCastFrameInRoom(room.Id, "pushLogicFrame", &logicFrame)
	return fpsTime
}

//定时，接收玩家的操作记录
func (sync *FrameSync) ReceivePlayerOperation(logicFrame pb.LogicFrame, conn *util.Conn, content string) {
	//mylog.Debug(logicFrame)
	room, empty := sync.Option.RoomManage.GetById(logicFrame.RoomId)
	if empty {
		sync.Option.Log.Error("getPoolElementById is empty" + logicFrame.RoomId)
	}
	err := sync.checkReceiveOperation(room, logicFrame, conn)
	if err != nil {
		sync.Option.Log.Error("receivePlayerOperation check error:" + err.Error())
		return
	}
	if len(logicFrame.Operations) < 0 {
		sync.Option.Log.Error("len(logicFrame.Operations) < 0")
		return
	}
	//roomSyncMetrics := roomSyncMetricsPool[logicFrame.RoomId]
	//roomSyncMetrics.InputNum ++
	//roomSyncMetrics.InputSize = roomSyncMetrics.InputSize + len(content)

	logicFrameStr, _ := json.Marshal(logicFrame.Operations)
	room.PlayersOperationQueue.PushBack(string(logicFrameStr))
	room.PlayersAckListRWLock.Lock()
	room.PlayersAckList[conn.UserId] = 1
	room.PlayersAckListRWLock.Unlock()
}

//检测玩家发送的操作是否合规
func (sync *FrameSync) checkReceiveOperation(room *Room, logicFrame pb.LogicFrame, conn *util.Conn) error {
	if room.Status == ROOM_STATUS_INIT {
		return errors.New("room status err is  ROOM_STATUS_INIT  " + strconv.Itoa(int(room.Status)))
	} else if room.Status == ROOM_STATUS_END {
		return errors.New("room status err is ROOM_STATUS_END  " + strconv.Itoa(int(room.Status)))
	} else if room.Status == ROOM_STATUS_PAUSE {
		//暂时状态，囚徒模式下
		//当A掉线后，会立刻更新房间状态为:暂停，但是其它未掉线的玩家依然还会发当前帧的操作数据
		//此时，房间已进入暂停状态，如果直接拒掉该条消息，会导致A恢复后，发送当前帧数据是正常的
		//而，其它玩家因为消息被拒，导致此条消息只有A发送成功，但是迟迟等不到其它玩家再未发送消息，该帧进入死锁
		//固，这里做出改变，暂停状态下：正常玩家可以多发一帧，等待掉线玩家重新上线
		if int(logicFrame.SequenceNumber) == room.SequenceNumber {
			sync.Option.Log.Warn("logicFrame.SequenceNumber  == room.SequenceNumber")
			//只有掉线的玩家，最后这一帧的数据没有发出来，才会到这个条件里
			//但，其它正常玩家如果还是一直不停的在发 这一帧，QUEUE 就爆了
			room.PlayersAckListRWLock.RLock()
			defer room.PlayersAckListRWLock.Unlock()
			if room.PlayersAckList[conn.UserId] == 1 {
				msg := "(offline) last frame Players has ack ,don'send... "
				sync.Option.Log.Error(msg)
				return errors.New(msg)
			} else {

			}
		} else {
			c_n := strconv.Itoa(int(logicFrame.SequenceNumber))
			r_n := strconv.Itoa(int(room.SequenceNumber))
			msg := "room status is ROOM_STATUS_PAUSE ,on receive num   c_n" + c_n + " ,r_n : " + r_n
			return errors.New(msg)
		}

	} else if room.Status == ROOM_STATUS_EXECING {

	} else {
		return errors.New("room status num error.  " + strconv.Itoa(int(room.Status)))
	}

	numberMsg := "cli_sn:" + strconv.Itoa(int(logicFrame.SequenceNumber)) + ", now_sn:" + strconv.Itoa(room.SequenceNumber)
	if int(logicFrame.SequenceNumber) == room.SequenceNumber {
		sync.Option.Log.Info("checkReceiveOperation ok , " + numberMsg)
		return nil
	}

	if int(logicFrame.SequenceNumber) > room.SequenceNumber {
		return errors.New("client num > room.SequenceNumber err:" + numberMsg)
	}
	//客户端延迟较高 相对的  服务端 发送较快
	if int(logicFrame.SequenceNumber) < room.SequenceNumber {
		return errors.New("client num < room.SequenceNumber err:" + numberMsg)
	}

	return nil
}

//游戏结束 - 结算
func (sync *FrameSync) roomEnd(roomId string, sendCloseChan int) {
	sync.Option.Log.Info("roomEnd")
	room, empty := sync.Option.RoomManage.GetById(roomId)
	if empty {
		sync.Option.Log.Error("getPoolElementById is empty" + roomId)
		return
	}
	//避免重复结束
	if room.Status == ROOM_STATUS_END {
		sync.Option.Log.Error("roomEnd status err " + roomId)
		return
	}
	room.UpStatus(ROOM_STATUS_END)
	room.EndTime = int32(util.GetNowTimeSecondToInt())
	for _, v := range room.PlayerList {
		v.UpPlayerRoomId("")
		//delete(mySyncPlayerRoom,v.Id)
	}

	if sync.Option.Store == 1 {

	}
	//给 房间FPS 协程 发送停止死循环信号
	if sendCloseChan == 1 {
		room.CloseChan <- 1
	}
}

//玩家操作后，触发C端主动发送游戏结束事件
func (sync *FrameSync) GameOver(requestGameOver pb.GameOver, conn *util.Conn) {
	responseGameOver := pb.GameOver{
		PlayerId:       requestGameOver.PlayerId,
		RoomId:         requestGameOver.RoomId,
		SequenceNumber: requestGameOver.SequenceNumber,
		Result:         requestGameOver.Result,
	}
	sync.boardCastInRoom(requestGameOver.RoomId, "gameOver", &responseGameOver)

	sync.roomEnd(requestGameOver.RoomId, 1)
}

//玩家触发了该角色死亡
func (sync *FrameSync) PlayerOver(requestGameOver pb.PlayerOver, conn *util.Conn) {
	//roomId := mySyncPlayerRoom[requestGameOver.PlayerId]
	roomId := conn.RoomId
	responseOtherPlayerOver := pb.PlayerOver{PlayerId: requestGameOver.PlayerId}
	sync.boardCastInRoom(roomId, "otherPlayerOver", &responseOtherPlayerOver)
}

//更新一个逻辑帧的确认状态
func (sync *FrameSync) upSyncRoomPoolElementPlayersAckStatus(roomId string, status int) {
	syncRoomPoolElement, _ := sync.Option.RoomManage.GetById(roomId)
	sync.Option.Log.Warn("upSyncRoomPoolElementPlayersAckStatus ,old :" + strconv.Itoa(syncRoomPoolElement.PlayersAckStatus) + "new" + strconv.Itoa(status))
	syncRoomPoolElement.PlayersAckStatus = status
}

//判定一个房间内，玩家在线的人
func (sync *FrameSync) roomOnlinePlayers(room *Room) []int32 {
	var playerOnLine []int32
	for _, v := range room.PlayerList {
		//player, empty := myPlayerManager.GetById(v.Id)
		//mylog.Debug("pinfo::",player," empty:",empty," ,pid:",v.Id)
		//if empty {
		//	continue
		//}
		//zlib.MyPrint(player.Status)
		if v.UserPlayStatus == PLAYER_STATUS_ONLINE {
			sync.Option.Log.Warn("playerOnLine append")
			playerOnLine = append(playerOnLine, v.UserId)
		}
	}
	//zlib.MyPrint(playerOnLine)
	return playerOnLine
}

//玩家断开连接后
func (sync *FrameSync) CloseOne(conn *util.Conn) {
	sync.Option.Log.Warn("sync.close one")
	//根据连接中的playerId，在用户缓存池中，查找该连接是否有未结束的游戏房间ID
	//roomId := myPlayerManager.GetRoomIdByPlayerId(conn.UserId)
	roomId := conn.RoomId
	if roomId == "" {
		//这里会先执行roomEnd，然后清空了player roomId 所有获取不到
		sync.Option.Log.Warn("roomid = empty " + strconv.Itoa(int(conn.UserId)))
		return
	}
	//根据roomId 查找房间信息
	room, empty := sync.Option.RoomManage.GetById(roomId)
	if empty {
		sync.Option.Log.Warn("room == empty , " + roomId)
		return
	}
	sync.Option.Log.Info("room.Status:" + strconv.Itoa(int(room.Status)))
	if room.Status == ROOM_STATUS_EXECING || room.Status == ROOM_STATUS_PAUSE {
		//判断下所有玩家是否均下线了
		playerOnLine := sync.roomOnlinePlayers(room)
		//mylog.Debug("playerOnLine:",playerOnLine, "len :",len(playerOnLine))
		playerOnLineCount := len(playerOnLine)
		//playerOnLineCount-- //这里因为，已有一个玩家关闭中，但是还未处理
		sync.Option.Log.Info("has check roomEnd , playerOnLineCount : " + strconv.Itoa(playerOnLineCount))
		if playerOnLineCount <= 1 { //这里这个判断有点不好处理，按说应该是<=0，也就是netway.close 应该先关闭了在线状态，但是如果全关了，后面可能要发消息就不行了
			sync.roomEnd(roomId, 1)
		} else {
			if room.Status == ROOM_STATUS_EXECING {
				room.UpStatus(ROOM_STATUS_PAUSE)
				responseOtherPlayerOffline := pb.OtherPlayerOffline{
					PlayerId: conn.UserId,
				}
				sync.boardCastInRoom(roomId, "otherPlayerOffline", &responseOtherPlayerOffline)
			}
		}
	} else {
		sync.Option.Log.Error("room.Status exception~~~")
		//能走到这个条件，肯定是发生过异常
		if room.Status == ROOM_STATUS_INIT {
			//本该room进入ready状态，但异常了
			sync.roomEnd(roomId, 0)
		} else if room.Status == ROOM_STATUS_END {
			//roomEnd 结算方法没有执行完毕，没有清空player的room id
			for _, v := range room.PlayerList {
				v.UpPlayerRoomId("")
			}
		} else if room.Status == ROOM_STATUS_READY {
			//<房间准备超时>守护协程  发生异常，未捕获到此房间已超时
			sync.roomEnd(room.Id, 0)
		}
	}
}

//单纯的给一个房间里的人发消息，不考虑是否有顺序号的情况
func (sync *FrameSync) boardCastInRoom(roomId string, action string, contentStruct interface{}) {
	room, empty := sync.Option.RoomManage.GetById(roomId)
	if empty {
		sync.Option.Log.Warn("syncRoomPoolElement is empty!!!")
	}
	for _, player := range room.PlayerList {
		if player.UserPlayStatus == PLAYER_STATUS_OFFLINE {
			sync.Option.Log.Error("player offline")
			continue
		}
		player.SendMsgCompressByUid(player.UserId, action, contentStruct)
	}
	//content ,_:= json.Marshal(contentStruct)
	content, _ := json.Marshal(util.JsonCamelCase{contentStruct})
	sync.addOneRoomHistory(room, action, string(content))
}

//给一个副本里的所有玩家广播数据，且该数据必须得有C端ACK
func (sync *FrameSync) boardCastFrameInRoom(roomId string, action string, contentStruct interface{}) {
	sync.Option.Log.Warn("boardCastFrameInRoom:" + roomId + action)
	syncRoomPoolElement, empty := sync.Option.RoomManage.GetById(roomId)
	if empty {
		sync.Option.Log.Panic("syncRoomPoolElement is empty!!!")
	}
	if sync.Option.LockMode == LOCK_MODE_PESSIMISTIC {
		if syncRoomPoolElement.PlayersAckStatus == PLAYERS_ACK_STATUS_WAIT {
			util.MyPrint(syncRoomPoolElement.PlayersAckList)
			sync.Option.Log.Error("syncRoomPoolElement PlayersAckStatus = " + strconv.Itoa(PLAYERS_ACK_STATUS_WAIT))
			return
		}
	}
	PlayersAckList := make(map[int32]int32)
	for _, player := range syncRoomPoolElement.PlayerList {
		PlayersAckList[player.UserId] = 0
		if player.Status == PLAYER_STATUS_OFFLINE {
			sync.Option.Log.Error("player offline")
			continue
		}
		player.SendMsgCompressByUid(player.UserId, action, contentStruct)

	}

	if sync.Option.LockMode == LOCK_MODE_PESSIMISTIC {
		syncRoomPoolElement.PlayersAckList = PlayersAckList
		sync.upSyncRoomPoolElementPlayersAckStatus(roomId, PLAYERS_ACK_STATUS_WAIT)
	}
	//content,_ := json.Marshal(contentStruct)
	content, _ := json.Marshal(util.JsonCamelCase{contentStruct})
	sync.addOneRoomHistory(syncRoomPoolElement, action, string(content))

	//if debugBreakPoint == 1{
	//	zlib.MyPrint(contentStruct)
	//	zlib.ExitPrint(3333)
	//}
}
func (sync *FrameSync) addOneRoomHistory(room *Room, action, content string) {
	logicFrameHistory := pb.RoomHistory{
		Action:  action,
		Content: content,
	}
	//该局副本的所有玩家操作日志，用于断线重连-补放/重播
	room.LogicFrameHistory = append(room.LogicFrameHistory, &logicFrameHistory)
}

//一个房间的玩家的所有操作记录，一般用于C端断线重连时，恢复
func (sync *FrameSync) RoomHistory(requestRoomHistory pb.ReqRoomHistory, conn *util.Conn) {
	roomId := requestRoomHistory.RoomId
	room, _ := sync.Option.RoomManage.GetById(roomId)
	responsePushRoomHistory := pb.RoomHistoryList{}
	responsePushRoomHistory.List = room.LogicFrameHistory
	conn.SendMsgCompressByUid(conn.UserId, "pushRoomHistory", &responsePushRoomHistory)
}

//玩家掉线了，重新连接后，恢复游戏了~这个时候，要通知另外的玩家
func (sync *FrameSync) PlayerResumeGame(requestPlayerResumeGame pb.PlayerResumeGame, conn *util.Conn) {
	room, empty := sync.Option.RoomManage.GetById(requestPlayerResumeGame.RoomId)
	if empty {
		sync.Option.Log.Error("playerResumeGame get room empty")
		return
	}
	var restartGame = 0
	var playerIds []int32
	if room.Status == ROOM_STATUS_PAUSE {
		playerOnlineNum := sync.roomOnlinePlayers(room)
		if len(playerOnlineNum) == len(room.PlayerList) {
			room.UpStatus(ROOM_STATUS_EXECING)
			restartGame = 1
			for _, v := range room.PlayerList {
				playerIds = append(playerIds, v.UserId)
			}
		}
	}

	responseOtherPlayerResumeGame := pb.PlayerResumeGame{
		PlayerId:       requestPlayerResumeGame.PlayerId,
		SequenceNumber: requestPlayerResumeGame.SequenceNumber,
		RoomId:         requestPlayerResumeGame.RoomId,
	}
	sync.boardCastInRoom(room.Id, "otherPlayerResumeGame", &responseOtherPlayerResumeGame)
	if restartGame == 1 {
		responseRestartGame := pb.RestartGame{
			RoomId:    requestPlayerResumeGame.RoomId,
			PlayerIds: playerIds,
		}
		sync.boardCastInRoom(room.Id, "restartGame", &responseRestartGame)
	}

}

func (sync *FrameSync) testFirstLogicFrame(room *Room) {
	//初始结束后，这里方便测试，再补一帧，所有玩家的随机位置
	if room.PlayerList[0].UserId < 999 {
		var operations []*pb.Operation
		for _, player := range room.PlayerList {
			location := strconv.Itoa(util.GetRandInt32Num(sync.Option.MapSize)) + "," + strconv.Itoa(util.GetRandInt32Num(sync.Option.MapSize))
			operation := pb.Operation{
				Id:       logicFrameMsgDefaultId,
				Event:    "move",
				Value:    location,
				PlayerId: player.UserId,
			}
			operations = append(operations, &operation)
		}
		logicFrameMsg := pb.LogicFrame{
			Id:             operationDefaultId,
			RoomId:         room.Id,
			SequenceNumber: int32(room.SequenceNumber),
			Operations:     operations,
		}
		sync.boardCastInRoom(room.Id, "pushLogicFrame", &logicFrameMsg)
	}
}
