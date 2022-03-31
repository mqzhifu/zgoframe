package service

import (
	"container/list"
	"errors"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"time"
	"zgoframe/protobuf/pb"
	"zgoframe/util"
)

const (
	ROOM_STATUS_INIT    = 1 //新房间，刚刚初始化，等待其它操作
	ROOM_STATUS_EXECING = 2 //已开始游戏
	ROOM_STATUS_END     = 3 //已结束
	ROOM_STATUS_READY   = 4 //准备中
	ROOM_STATUS_PAUSE   = 5 //有玩家掉线，暂停中
)

type Room struct {
	Id                     string          `json:"id"`           //房间ID
	AddTime                int32           `json:"addTime"`      //创建时间
	StartTime              int32           `json:"startTime"`    //开始游戏时间
	EndTime                int32           `json:"endTime"`      //游戏结束时间
	ReadyTimeout           int32           `json:"readyTimeout"` //玩家准备超时时间
	Status                 int32           `json:"status"`       //状态
	PlayerList             []*util.Conn    `json:"playerList"`   //玩家列表
	PlayerIds              []int32         `json:"player_ids"`
	SequenceNumber         int             `json:"sequenceNumber"`   //当前帧频逻辑帧顺序号
	PlayersReadyList       map[int32]int32 `json:"playersReadyList"` //玩家准备列表
	RandSeek               int32           `json:"randSeek"`         //随机数种子
	PlayersAckList         map[int32]int32 `json:"playersAckList"`   //玩家确认列表
	PlayersAckStatus       int             `json:"playersAckStatus"` //玩家确认列表的状态
	PlayersAckListRWLock   *sync.RWMutex   `json:"-"`                //玩家一帧内的确认操作，需要加锁
	PlayersReadyListRWLock *sync.RWMutex   `json:"-"`
	//接收玩家操作指令-集合
	PlayersOperationQueue      *list.List `json:"-"` //用于存储玩家一个逻辑帧内推送的：玩家操作指令
	CloseChan                  chan int   `json:"-"` //关闭信号管道
	ReadyCloseChan             chan int   `json:"-"` //<玩家准备>协程的关闭信号管道
	WaitPlayerOfflineCloseChan chan int   `json:"-"` //<一局游戏，某个玩家掉线，其它玩家等待它的时间>
	//本局游戏，历史记录，玩家的所有操作
	LogicFrameHistory []*pb.RoomHistory `json:"logicFrameHistory"` //玩家的历史所有记录
	RoomManager       *RoomManager      //父类
}

type RoomManager struct {
	Pool   map[string]*Room
	Option RoomManagerOption
}

type RoomManagerOption struct {
	Log          *zap.Logger
	FrameSync    *FrameSync
	ReadyTimeout int32
	RoomPeople   int32
}

func NewRoomManager(roomManagerOption RoomManagerOption) *RoomManager {
	roomManager := new(RoomManager)
	roomManager.Option = roomManagerOption
	return roomManager
}

func (roomManager *RoomManager) SetFrameSync(frameSync *FrameSync) {
	roomManager.Option.FrameSync = frameSync
}

func (roomManager *RoomManager) NewRoom() *Room {
	room := new(Room)
	room.Id = CreateRoomId()
	room.Status = ROOM_STATUS_INIT
	room.AddTime = int32(util.GetNowTimeSecondToInt())
	room.StartTime = 0
	room.EndTime = 0
	room.SequenceNumber = 0
	room.PlayersAckList = make(map[int32]int32)
	room.PlayersAckListRWLock = &sync.RWMutex{}
	room.PlayersAckStatus = PLAYERS_ACK_STATUS_INIT
	room.RandSeek = int32(util.GetRandIntNum(100))
	room.PlayersOperationQueue = list.New()
	room.PlayersReadyList = make(map[int32]int32)
	room.PlayersReadyListRWLock = &sync.RWMutex{}

	readyTimeout := int32(util.GetNowTimeSecondToInt()) + roomManager.Option.ReadyTimeout
	room.ReadyTimeout = readyTimeout

	room.RoomManager = roomManager
	//myMetrics.fastLog("total.RoomNum", 2, 0)
	//roomManager.Pool[room.Id] = room

	return room
}

func CreateRoomId() string {
	tt := time.Now().UnixNano() / 1e6
	string := strconv.FormatInt(tt, 10)
	return string
}

func (room *Room) AddPlayer(conn *util.Conn) {
	room.PlayerList = append(room.PlayerList, conn)
}

func (room *Room) UpStatus(status int32) {
	room.RoomManager.Option.Log.Info("room upStatus ,old :" + strconv.Itoa(int(room.Status)) + " new :" + strconv.Itoa(int(status)))
	room.Status = status
}

//C端获取一个房间的信息
func (roomManager *RoomManager) GetRoom(requestGetRoom pb.RoomBaseInfo, conn *util.Conn) {
	roomId := requestGetRoom.RoomId
	room, _ := roomManager.GetById(roomId)

	ResponsePushRoomInfo := pb.RoomBaseInfo{
		Id:             room.Id,
		SequenceNumber: int32(room.SequenceNumber),
		AddTime:        room.AddTime,
		PlayerIds:      room.PlayerIds,
		Status:         room.Status,
		//Timeout: room.Timeout,
		RandSeek: room.RandSeek,
		RoomId:   room.Id,
	}
	conn.SendMsgCompressByUid(conn.UserId, "pushRoomInfo", &ResponsePushRoomInfo)
}

//根据ROOID  有池子里找到该roomInfo
func (roomManager *RoomManager) GetById(roomId string) (room *Room, empty bool) {
	room, exist := roomManager.Pool[roomId]
	if !exist {
		roomManager.Option.Log.Error("getPoolElementById is empty," + roomId)
		return room, true
	}
	return room, false
}

func (roomManager *RoomManager) Shutdown() {
	roomManager.Option.Log.Warn("shutdown mySync")
	//sync.CloseChan <- 1
	if len(roomManager.Pool) <= 0 {
		return
	}
	//这里只做信号关闭，即：死循环的协程，而真实的关闭由netWay.Close解决
	for _, room := range roomManager.Pool {
		if room.Status == ROOM_STATUS_READY {
			room.ReadyCloseChan <- 1
		} else if room.Status == ROOM_STATUS_EXECING {
			room.CloseChan <- 1
		}
	}
}

//给集合添加一个新的 游戏副本
//一局新游戏（副本）创建成功，告知玩家进入战场，等待 所有玩家准备确认
func (roomManager *RoomManager) AddPoolElement(room *Room) error {
	roomManager.Option.Log.Info("addPoolElement")
	_, exist := roomManager.GetById(room.Id)
	if exist {
		msg := "mySyncRoomPool has exist : " + room.Id
		roomManager.Option.Log.Error(msg)
		err := errors.New(msg)
		return err
	}

	uids := []int32{}
	for _, v := range room.PlayerList {
		uids = append(uids, v.UserId)
	}
	room.PlayerIds = uids
	room.UpStatus(ROOM_STATUS_READY)
	roomManager.Pool[room.Id] = room

	room.CloseChan = make(chan int)
	room.ReadyCloseChan = make(chan int)

	roomManager.Option.FrameSync.StartOne(room)

	return nil
}
