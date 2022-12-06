package frame_sync

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"zgoframe/protobuf/pb"
	"zgoframe/service"
)

type FrameSync struct {
	RoomManage        *RoomManager `json:"-"` //外部指针-房间服务
	PlayerConnManager *PlayerConnManager
	Option            FrameSyncOption
	CloseChan         chan int
}

type FrameSyncOption struct {
	ProjectId             int                            `json:"project_id"`         //项目Id,给玩家推送消失的时候使用
	FPS                   int32                          `json:"fps"`                //frame pre second
	LockMode              int32                          `json:"lock_mode"`          //锁模式，乐观|悲观
	Store                 int32                          `json:"store"`              //持久化，玩家每帧的动作，暂未使用
	OffLineWaitTime       int                            `json:"off_line_wait_time"` //lockStep 玩家掉线后，其它玩家等待最长时间
	Gorm                  *gorm.DB                       `json:"-"`                  //
	RequestServiceAdapter *service.RequestServiceAdapter `json:"-"`                  //请求3方服务 适配器
	Log                   *zap.Logger                    `json:"-"`

	//MapSize               int32                          `json:"map_size"` //地址大小，给前端初始化使用
}

func NewFrameSync(Option FrameSyncOption) *FrameSync {
	Option.Log.Info("NewSync instance")
	sync := new(FrameSync)
	sync.Option = Option

	//统计
	//RoomSyncMetricsPool = make(map[string]RoomSyncMetrics)

	//帧同步 - 房间服务 - room要先实例化,math frame_sync 都强依赖room
	roomManagerOption := RoomManagerOption{
		Log:                   Option.Log,
		RequestServiceAdapter: Option.RequestServiceAdapter,
		Gorm:                  Option.Gorm,
		FrameSync:             sync,
	}
	sync.RoomManage = NewRoomManager(roomManagerOption)

	sync.PlayerConnManager = NewPlayerConnManager(Option.Log)

	//用于关闭
	sync.CloseChan = make(chan int)
	return sync
}

func (sync *FrameSync) GetPlayerBase(playerBase pb.PlayerBase) {
	conn, _ := sync.PlayerConnManager.GetById(playerBase.PlayerId)
	playerState := pb.PlayerState{
		Status:   int32(conn.Status),
		AddTime:  int32(conn.AddTime),
		RoomId:   conn.RoomId,
		PlayerId: playerBase.PlayerId,
	}
	sync.Option.RequestServiceAdapter.GatewaySendMsgByUid(playerBase.PlayerId, "SC_PlayerState", playerState)
}

func (sync *FrameSync) CreateFD(fd pb.FDCreateEvent) error {
	sync.PlayerConnManager.AddOne(fd.UserId)
	return nil
}
func (sync *FrameSync) CloseFD(FDCloseEvent pb.FDCloseEvent) error {
	user, empty := sync.PlayerConnManager.GetById(FDCloseEvent.UserId)
	if empty {
		sync.Option.Log.Error("user id not in PlayerConnManager pool")
		return errors.New("user id not in PlayerConnManager pool")
	}
	if user.RoomId != "" {
		room, empty := sync.RoomManage.GetById(user.RoomId)
		if empty {
			return errors.New("user id not in RoomManage pool")
		}
		room.Sync.CloseOne(FDCloseEvent)
	}
	sync.PlayerConnManager.DelOne(FDCloseEvent.UserId)
	return nil
}

func (sync *FrameSync) GetOption() FrameSyncOption {
	return sync.Option
}
func (sync *FrameSync) Heartbeat(hb pb.Heartbeat) error {
	return nil
}

func (sync *FrameSync) ReceivePlayerOperation(LogicFrame pb.LogicFrame) error {
	room, _ := sync.RoomManage.GetById(LogicFrame.RoomId)
	room.Sync.ReceivePlayerOperation(LogicFrame)

	return nil
}

func (sync *FrameSync) PlayerResumeGame(PlayerResumeGame pb.PlayerResumeGame) error {
	room, _ := sync.RoomManage.GetById(PlayerResumeGame.RoomId)
	room.Sync.PlayerResumeGame(PlayerResumeGame)
	return nil
}

func (sync *FrameSync) PlayerReady(PlayerReady pb.PlayerReady) error {
	room, _ := sync.RoomManage.GetById(PlayerReady.RoomId)
	room.Sync.PlayerReady(PlayerReady)
	return nil
}

func (sync *FrameSync) PlayerOver(PlayerOver pb.PlayerOver) error {
	room, _ := sync.RoomManage.GetById(PlayerOver.RoomId)
	room.Sync.PlayerOver(PlayerOver)
	return nil
}

func (sync *FrameSync) RoomHistory(ReqRoomHistory pb.ReqRoomHistory) error {
	//room,_ := sync.RoomManage.GetById(ReqRoomHistory.RoomId)
	return nil
}
