package service

import (
	"context"
	"go.uber.org/zap"
	"strconv"
	"time"
	"zgoframe/protobuf/pb"
	"zgoframe/util"
)

type Match struct {
	Option MatchOption
	Close  chan int //关闭 匹配程序 死循环
	//RecoverTimes int
}

type MatchOption struct {
	Log         *zap.Logger
	RoomManager *RoomManager
}

type PlayerSign struct {
	AddTime  int32
	PlayerId int32
	Conn     *util.Conn
}

var signPlayerPool []PlayerSign

func NewMatch(matchOption MatchOption) *Match {
	matchOption.Log.Info("NewMatch instance")
	match := new(Match)
	match.Option = matchOption
	match.Close = make(chan int)
	return match
}

func (match *Match) getOneSignPlayerById(playerId int32) (playerSign PlayerSign, empty bool) {
	for _, v := range signPlayerPool {
		if v.PlayerId == playerId {
			return v, false
		}
	}
	return playerSign, true
}

func (match *Match) addOnePlayer(requestPlayerMatchSign pb.PlayerMatchSign, conn *util.Conn) {
	//playerId := conn.UserId
	//player, empty := myPlayerManager.GetById(requestPlayerMatchSign.PlayerId)
	if conn.UserId <= 0 {
		msg := "conn.UserId <= 0"
		match.matchSignErrAndSend(msg, conn)
		return
	}

	if conn.UserPlayStatus != PLAYER_STATUS_ONLINE {
		msg := "playerMatchSign status != online , status=" + strconv.Itoa(conn.UserPlayStatus)
		match.matchSignErrAndSend(msg, conn)
		return
	}

	if conn.RoomId != "" {
		msg := "playerMatchSign player.RoomId != '' , roomId=" + conn.RoomId
		match.matchSignErrAndSend(msg, conn)
		return
	}
	//检查，是否重复报名
	_, empty := match.getOneSignPlayerById(conn.UserId)
	if !empty {
		msg := "match sign addOnePlayer : player has exist" + strconv.Itoa(int(conn.UserId))
		match.matchSignErrAndSend(msg, conn)
		return
	}
	newPlayerSign := PlayerSign{PlayerId: conn.UserId, AddTime: int32(util.GetNowTimeSecondToInt()), Conn: conn}
	signPlayerPool = append(signPlayerPool, newPlayerSign)
	return
}

//玩家报名时，可能因为BUG，造成一些系统级的错误，如：丢失玩家状态等
//出现这种S端异常的情况，除了报错还要通知一下C端
func (match *Match) matchSignErrAndSend(msg string, conn *util.Conn) {
	//mylog.Error(msg)
	playerMatchSignFailed := pb.PlayerMatchSignFailed{
		PlayerId: conn.UserId,
		Msg:      msg,
	}
	conn.SendMsgCompressByUid(conn.UserId, "playerMatchSignFailed", &playerMatchSignFailed)
}

func (match *Match) delOnePlayer(requestCancelSign pb.PlayerMatchSignCancel, conn *util.Conn) {
	//playerId := requestCancelSign.PlayerId
	match.Option.Log.Info("cancel : delOnePlayer " + strconv.Itoa(int(requestCancelSign.PlayerId)))
	for k, v := range signPlayerPool {
		if v.PlayerId == requestCancelSign.PlayerId {
			if len(signPlayerPool) == 1 {
				signPlayerPool = []PlayerSign{}
			} else {
				signPlayerPool = append(signPlayerPool[:k], signPlayerPool[k+1:]...)
			}
			return
		}
	}
	match.Option.Log.Warn("no match playerId" + strconv.Itoa(int(conn.UserId)))
}

func (match *Match) Shutdown() {
	match.Close <- 1
}
func (match *Match) Start(ctx context.Context) {
	for {
		select {
		case <-match.Close:
			goto end
		default:
			//计算：当前池子里的人数，是否满足可以开启一局游戏
			if int32(len(signPlayerPool)) < match.Option.RoomManager.Option.RoomPeople {
				//不满足即睡眠等待500毫秒
				time.Sleep(time.Millisecond * 500)
				break
			}
			//创建一个新的房间，用于装载玩家及同步数据等
			newRoom := match.Option.RoomManager.NewRoom()
			for i := 0; i < len(signPlayerPool); i++ {
				if signPlayerPool[i].Conn.UserPlayStatus != PLAYER_STATUS_ONLINE {
					match.Option.Log.Error("matching Players.status != online , " + strconv.Itoa(int(signPlayerPool[i].PlayerId)))
				}
				if signPlayerPool[i].Conn.RoomId != "" {
					match.Option.Log.Error("matching Players.roomId != '' , " + strconv.Itoa(int(signPlayerPool[i].PlayerId)))
				}
				newRoom.AddPlayer(signPlayerPool[i].Conn)
			}
			//删除上面匹配成功的玩家
			signPlayerPool = append(signPlayerPool[match.Option.RoomManager.Option.RoomPeople:])
			util.MyPrint("newRoom:", newRoom)
			match.Option.Log.Info("create a room :" + newRoom.Id)
			match.Option.RoomManager.AddPoolElement(newRoom)

			time.Sleep(time.Millisecond * 100)
		}
	}
end:
	match.Option.Log.Warn(CTX_DONE_PRE + "matchingPlayerCreateRoom close")
}
