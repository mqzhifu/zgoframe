package frame_sync

import "C"
import (
	"container/list"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"sync"
	"time"
	"zgoframe/model"
	"zgoframe/protobuf/pb"
	"zgoframe/service"
	"zgoframe/util"
)

//type Player struct {
//	Id     int32
//	Status int
//	RoomId string
//}

type Room struct {
	Id           string              `json:"id"`           //房间ID
	RuleId       int32               `json:"rule_id"`      //分类
	Rule         model.GameMatchRule `json:"rule"`         //分类
	AddTime      int32               `json:"addTime"`      //创建时间
	StartTime    int32               `json:"startTime"`    //开始游戏时间
	EndTime      int32               `json:"endTime"`      //游戏结束时间
	ReadyTimeout int                 `json:"readyTimeout"` //房间创建，玩家进入，需要每个玩家点击准备：超时时间
	Status       int32               `json:"status"`       //状态
	//PlayerList             []int32             `json:"playerList"`       //玩家列表
	PlayerIds              []int32         `json:"player_ids"`       //玩家列表IDS，跟上面一样，只是 方便计算
	SequenceNumber         int             `json:"sequenceNumber"`   //当前帧：顺序号
	PlayersReadyList       map[int32]int32 `json:"playersReadyList"` //存储：玩家<准备>状态的列表
	PlayersReadyListRWLock *sync.RWMutex   `json:"-"`                //准备状态的时候，是轮询，得加锁
	ReadyCloseChan         chan int        `json:"-"`                //玩家都准备后，要关闭轮询的协程，关闭信号管道
	StatusLock             *sync.Mutex     `json:"-"`                //修改状态得加锁
	RandSeek               int32           `json:"randSeek"`         //随机数种子
	PlayersAckList         map[int32]int32 `json:"playersAckList"`   //玩家确认列表
	PlayersAckStatus       int             `json:"playersAckStatus"` //玩家确认列表的状态
	PlayersAckListRWLock   *sync.RWMutex   `json:"-"`                //玩家一帧内的确认操作，需要加锁
	//接收玩家操作指令-集合
	LogicFrameWaitTime    int64      `json:"logicFrameWaitTime"`  //每一帧的等待总时长，虽然C端定时发送每帧数据，但有可能某个玩家的某帧发送的数据丢失，造成两边空等，得有个计时
	CloseChan             chan int   `json:"-"`                   //关闭信号管道
	WaitPlayerOffline     int        `json:"wait_player_offline"` //<一局游戏，某个玩家掉线，其它玩家等待它的时间>
	PlayersOperationQueue *list.List `json:"-"`                   //用于存储玩家一个逻辑帧内推送的：玩家操作指令
	//本局游戏，历史记录，玩家的所有操作
	LogicFrameHistory []*pb.RoomHistory `json:"logicFrameHistory"` //玩家的历史所有记录
	EndTotal          string            `json:"rs"`                //本房间的一局游戏，最终的比赛结果
	RoomManager       *RoomManager      `json:"-"`                 //父类
	Sync              *Sync             `json:"-"`
}

type RoomManager struct {
	Pool   map[string]*Room
	Option RoomManagerOption
}

type RoomManagerOption struct {
	RequestServiceAdapter *service.RequestServiceAdapter //请求3方服务 适配器
	Log                   *zap.Logger
	Gorm                  *gorm.DB
	FrameSync             *FrameSync
	//ReadyTimeout          int32 //房间人数满足了，等待 所有玩家确认，超时时间
	//RoomPeople            int32 //房间有多少人后，可以开始游戏了
	//MapSize               int32 `json:"mapSize"` //帧同步，地图大小，给前端初始化使用（测试使用）
	//Store                 int32 `json:"store"`   //持久化：room
}

func NewRoomManager(roomManagerOption RoomManagerOption) *RoomManager {
	roomManager := new(RoomManager)
	roomManager.Option = roomManagerOption
	roomManager.Pool = make(map[string]*Room)
	return roomManager
}

//func (roomManager *RoomManager) SetFrameSync(frameSync *FrameSync) {
//	roomManager.Option.FrameSync = frameSync
//}

//创建一个空房间
func (roomManager *RoomManager) NewEmptyRoom() *Room {
	room := new(Room)
	room.Id = CreateRoomId()
	room.Status = service.ROOM_STATUS_INIT
	room.AddTime = int32(util.GetNowTimeSecondToInt())
	room.StartTime = 0
	room.EndTime = 0
	room.SequenceNumber = -1
	room.PlayersAckList = make(map[int32]int32)
	room.PlayersAckListRWLock = &sync.RWMutex{}
	room.PlayersAckStatus = service.PLAYERS_ACK_STATUS_INIT
	room.RandSeek = int32(util.GetRandIntNum(100))
	room.PlayersOperationQueue = list.New()
	room.PlayersReadyList = make(map[int32]int32)
	room.PlayersReadyListRWLock = &sync.RWMutex{}
	room.StatusLock = &sync.Mutex{}
	room.EndTotal = ""
	room.RuleId = 0
	room.ReadyTimeout = util.GetNowTimeSecondToInt() + 10
	//readyTimeout := int32(util.GetNowTimeSecondToInt()) + roomManager.Option.ReadyTimeout

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

func (room *Room) AddPlayer(id int) {
	//p := Player{
	//	Id:     int32(id),
	//	Status: service.PLAYER_STATUS_ONLINE,
	//	RoomId: "",
	//}
	room.PlayerIds = append(room.PlayerIds, int32(id))
}

func (room *Room) GetStatus(status int32) {
	room.StatusLock.Lock()
	room.Status = status
	room.StatusLock.Unlock()
}

func (room *Room) UpStatus(status int32) {
	room.RoomManager.Option.Log.Info("room upStatus ,old :" + strconv.Itoa(int(room.Status)) + " new :" + strconv.Itoa(int(status)))
	room.StatusLock.Lock()
	room.Status = status
	room.StatusLock.Unlock()
}

//C端获取一个房间的信息
func (roomManager *RoomManager) GetRoom(requestGetRoom pb.RoomBaseInfo) error {
	roomId := requestGetRoom.RoomId
	room, _ := roomManager.GetById(roomId)

	roomBaseInfo := pb.RoomBaseInfo{
		Id:             room.Id,
		SequenceNumber: int32(room.SequenceNumber),
		AddTime:        room.AddTime,
		PlayerIds:      room.PlayerIds,
		Status:         room.Status,
		//Timeout: room.Timeout,
		RandSeek: room.RandSeek,
		RoomId:   room.Id,
	}
	roomManager.Option.RequestServiceAdapter.GatewaySendMsgByUid(requestGetRoom.SourceUid, "SC_RoomBaseInfo", &roomBaseInfo)
	//conn.SendMsgCompressByUid(requestGetRoom.SourceUid, "pushRoomInfo", &ResponsePushRoomInfo)
	return nil
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
	if len(roomManager.Pool) <= 0 {
		return
	}
	//这里只做信号关闭，即：死循环的协程，而真实的关闭由netWay.Close解决
	for _, room := range roomManager.Pool {
		if room.Status == service.ROOM_STATUS_READY {
			room.ReadyCloseChan <- 1
		} else if room.Status == service.ROOM_STATUS_EXECING {
			room.CloseChan <- 1
		}
	}
}

//给集合添加一个新的 游戏副本
//一局新游戏（副本）创建成功，告知玩家进入战场，等待 所有玩家准备确认
func (roomManager *RoomManager) AddOne(room *Room) error {
	roomManager.Option.Log.Info("roomManager addPoolElement id:" + room.Id)
	_, empty := roomManager.GetById(room.Id)
	if !empty {
		msg := "new roomId exist in mySyncRoomPool : " + room.Id
		roomManager.Option.Log.Error(msg)
		err := errors.New(msg)
		return err
	}

	var rule model.GameMatchRule
	err := roomManager.Option.Gorm.Where("id = ?", room.RuleId).First(&rule).Error
	if err != nil {
		util.MyPrint("roomManager AddOne err:" + err.Error())
		return err
	}
	room.Rule = rule
	//uids := []int32{}
	for _, pid := range room.PlayerIds {
		roomManager.Option.FrameSync.PlayerConnManager.UpRoomId(pid, room.Id)
		//v.UpPlayerRoomId(room.Id)
		roomManager.Option.Log.Debug("UpPlayerRoomId uid:" + strconv.Itoa(int(pid)) + " roomId:" + room.Id)
		//uids = append(uids, int32(pid))
	}
	//room.PlayerIds = uids
	room.UpStatus(service.ROOM_STATUS_READY)
	roomManager.Pool[room.Id] = room

	room.CloseChan = make(chan int)
	room.ReadyCloseChan = make(chan int)

	//user -> sign ->Match -> Room -> Rsync
	//帧同步服务 - 强-依赖room
	syncOption := SyncOption{
		//ProjectId:             C.System.ProjectId,
		RequestServiceAdapter: room.RoomManager.Option.RequestServiceAdapter,
		Log:                   room.RoomManager.Option.Log,
		Room:                  room,
		RoomManage:            room.RoomManager,
		FPS:                   int32(rule.Fps),
	}

	room.Sync = NewSync(syncOption)
	room.Sync.StartOne(room)

	return nil
}

//func (player *Player) UpPlayerRoomId(roomId string) {
//	player.RoomId = roomId
//}
