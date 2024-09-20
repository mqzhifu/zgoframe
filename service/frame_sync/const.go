package frame_sync

// INIT -> READY -> EXECING -> PAUSE -> END
// @parse 帧同步-房间状态
const (
	ROOM_STATUS_INIT    = 1 //新房间，刚刚初始化，等待其它操作
	ROOM_STATUS_EXECING = 2 //已开始游戏
	ROOM_STATUS_END     = 3 //已结束
	ROOM_STATUS_READY   = 4 //准备中
	ROOM_STATUS_PAUSE   = 5 //有玩家掉线，暂停中
)

// @parse 帧同步-锁模式
const (
	LOCK_MODE_PESSIMISTIC = 1 //囚徒
	LOCK_MODE_OPTIMISTIC  = 2 //乐观
)

// @parse 帧同步，一个副本的，一条消息的，同步状态
const (
	PLAYERS_ACK_STATUS_INIT = 1 //初始化
	PLAYERS_ACK_STATUS_WAIT = 2 //等待玩家确认
	PLAYERS_ACK_STATUS_OK   = 3 //所有玩家均已确认
)

// (帧同步 游戏匹配 好像都在用)
// @parse 玩家状态
const (
	PLAYER_STATUS_ONLINE  = 1 //在线
	PLAYER_STATUS_OFFLINE = 2 //离线
)

// @parse 帧同步-房间内，玩家准备状态
const (
	PLAYER_NO_READY  = 1 //玩家未准备
	PLAYER_HAS_READY = 2 //玩家已准备
)
