package frame_sync

import (
	"errors"
	"go.uber.org/zap"
	"strconv"
	"zgoframe/util"
)

const (
	STATUS_NORMAL  = 1
	STATUS_PLAYING = 2
)

type PlayerConn struct {
	Id      int32
	AddTime int
	Status  int
	RoomId  string
}

type PlayerConnManager struct {
	Log  *zap.Logger
	Pool map[int32]*PlayerConn
}

func NewPlayerConnManager(log *zap.Logger) *PlayerConnManager {
	playerConnManager := new(PlayerConnManager)
	playerConnManager.Pool = make(map[int32]*PlayerConn)
	playerConnManager.Log = log
	return playerConnManager
}

func (playerConnManager PlayerConnManager) GetById(id int32) (*PlayerConn, bool) {
	playerConn, exist := playerConnManager.Pool[id]
	return playerConn, exist
}

func (playerConnManager PlayerConnManager) AddOne(id int32) error {
	playerConnManager.Log.Info("add one id:" + strconv.Itoa(int(id)))
	checkPlayerConn, exist := playerConnManager.GetById(id)
	if exist {
		if checkPlayerConn.RoomId != "" {
			return errors.New("player has roomId")
		}

		if checkPlayerConn.Status != STATUS_NORMAL {
			return errors.New("player stats != STATUS_NORMAL")
		}
	}
	checkPlayerConn = &PlayerConn{
		AddTime: util.GetNowTimeSecondToInt(),
		Id:      id,
		RoomId:  "",
		Status:  STATUS_NORMAL,
	}
	playerConnManager.Pool[id] = checkPlayerConn
	return nil
}

func (playerConnManager PlayerConnManager) DelOne(id int32) error {
	_, exist := playerConnManager.GetById(id)
	if exist {
		return errors.New("id not exist")
	}
	delete(playerConnManager.Pool, id) //只删除游戏信息，房间信息不要动
	return nil
}

func (playerConnManager PlayerConnManager) UpStatus(id int32, status int) {
	playerConn, exist := playerConnManager.Pool[id]
	if !exist {
		playerConnManager.Log.Error("UpStatus id not exist" + strconv.Itoa(int(id)))
		return
	}
	playerConn.Status = status
}

func (playerConnManager PlayerConnManager) UpRoomId(id int32, roomId string) {

	playerConn, exist := playerConnManager.Pool[id]
	if !exist {
		playerConnManager.Log.Error("UpRoomId id not exist" + strconv.Itoa(int(id)))
		return
	}
	playerConn.RoomId = roomId
}
