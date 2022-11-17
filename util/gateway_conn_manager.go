package util

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"sync"
	"time"
	"zgoframe/protobuf/pb"
)

//管理 CONN 的容器
type ConnManager struct {
	Pool              map[int32]*Conn // map[userId]*Conn FD 连接池
	PoolRWLock        *sync.RWMutex
	CloseCheckTimeout chan int
	Option            ConnManagerOption
}

type ConnManagerOption struct {
	maxClientConnNum    int32  //客户端最大连接数
	connTimeout         int32  //一个FD的超时时间
	MsgContentMax       int32  //一条消息最大值
	DefaultContentType  int32  //每个连接的默认 内容 类型
	DefaultProtocolType int32  //每个连接的默认 协议 类型
	MsgSeparator        string //传输消息时，每条消息的间隔符，防止 粘包

	Log      *zap.Logger
	Metrics  *MyMetrics
	ProtoMap *ProtoMap //协议ID管理器
	NetWay   *NetWay
}

//消息内容协议体
type ProtocolCtrlInfo struct {
	ContentType  int32
	ProtocolType int32
}

//实例化
func NewConnManager(connManagerOption ConnManagerOption) *ConnManager {
	connManagerOption.Log.Info("NewConnManager instance:")

	connManager := new(ConnManager)

	connManager.Pool = make(map[int32]*Conn)
	connManager.CloseCheckTimeout = make(chan int) //连接超时 关闭信号
	connManager.PoolRWLock = &sync.RWMutex{}

	if connManagerOption.MsgSeparator == "" {
		//消息分隔符，主要是给TCP使用，一个字符，且最好不要用：\n，因为会与protobuf 冲突
		connManagerOption.MsgSeparator = "\f"
	}

	connManager.Option = connManagerOption

	return connManager
}
func (connManager *ConnManager) MakeError(msg string) error {
	connManager.Option.Log.Error("connManager make error:" + msg)
	return errors.New(msg)
}

//启动容器，监听 连接超时处理
func (connManager *ConnManager) CheckTimeout() {
	//defer func(ctx context.Context ) {
	//	if err := recover(); err != nil {
	//		myNetWay.RecoverGoRoutine(connManager.Start,ctx,err)
	//	}
	//}(ctx)

	connManager.Option.Log.Warn("checkConnPoolTimeout start:")
	for {
		select {
		case <-connManager.CloseCheckTimeout:
			goto end
		default:
			pool := connManager.getPoolAll()
			for _, v := range pool {
				now := int32(GetNowTimeSecondToInt())
				x := v.UpTime + connManager.Option.connTimeout
				if now > x {
					v.CloseOneConn(CLOSE_SOURCE_TIMEOUT)
				}
			}
			time.Sleep(time.Second * 1)
			//mySleepSecond(1,"checkConnPoolTimeout")
		}
	}
end:
	connManager.Option.Log.Warn(CTX_DONE_PRE + "checkConnPoolTimeout close")
}

//关闭容器，回收处理
func (connManager *ConnManager) Shutdown() {
	connManager.Option.Log.Warn("shutdown connManager")
	connManager.CloseCheckTimeout <- 1
	if len(connManager.Pool) <= 0 {
		return
	}
	pool := connManager.getPoolAll()
	for _, conn := range pool {
		conn.CloseOneConn(CLOSE_SOURCE_CONN_SHUTDOWN)
	}
}

//创建一个新的连接结构体
func (connManager *ConnManager) CreateOneConn(connFd FDAdapter, netWay *NetWay) (myConn *Conn) {
	connManager.PoolRWLock.RLock()
	defer connManager.PoolRWLock.RUnlock()

	connManager.Option.Log.Info("Create one Conn  client struct")

	now := int32(GetNowTimeSecondToInt())

	myConn = &Conn{
		Conn:           connFd,
		UserId:         0,
		AddTime:        now,
		UpTime:         now,
		Status:         CONN_STATUS_INIT, //CONN status
		ConnManager:    connManager,
		RTT:            0,
		SessionId:      "",
		ContentType:    connManager.Option.DefaultContentType,
		ProtocolType:   connManager.Option.DefaultProtocolType,
		Metrics:        ConnMetrics{},
		MsgInChan:      make(chan pb.Msg, 5000), //从底层FD中读出消息后，存储此处，等待其它协程接收
		netWay:         netWay,
		UserPlayStatus: PLAYER_STATUS_ONLINE,
		//CloseChan 		chan int
	}

	connManager.Option.Log.Info("reg conn callback CloseHandler")

	return myConn
}

//
func (connManager *ConnManager) GetConnPoolById(userId int32) (*Conn, bool) {
	connManager.PoolRWLock.RLock()
	defer connManager.PoolRWLock.RUnlock()

	conn, ok := connManager.Pool[userId]
	return conn, ok
}

//往POOL里添加一个新的连接
func (connManager *ConnManager) addConnPool(NewConn *Conn) error {
	if NewConn.UserId <= 0 {
		connManager.Option.Log.Error("addConnPool NewConn.UserId <= 0 ")
	}
	oldConn, exist := connManager.GetConnPoolById(NewConn.UserId)
	if exist { //该UID已经创建了连接，可能是在别处登陆，直接踢掉旧的连接
		msg := strconv.Itoa(int(NewConn.UserId)) + " kickOff old pid :" + strconv.Itoa(int(oldConn.UserId))
		connManager.Option.Log.Warn(msg)
		//err := errors.New(msg)
		responseKickOff := pb.KickOff{
			Time: GetNowMillisecond(),
		}
		//给旧连接发送消息通知,这条消息不一定能发出去，因为后面马上就会接关闭FD，最好是同步执行，不要异步发送这条消息
		oldConn.SendMsgCompressByConn("SC_KickOff", &responseKickOff)
		oldConn.CloseOneConn(CLOSE_SOURCE_OVERRIDE)
	}
	connManager.Option.Log.Info("addConnPool : " + strconv.Itoa(int(NewConn.UserId)))
	connManager.PoolRWLock.Lock()
	defer connManager.PoolRWLock.Unlock()
	connManager.Pool[NewConn.UserId] = NewConn

	return nil
}

//删除一个FD
func (connManager *ConnManager) delConnPool(uid int32) {
	connManager.Option.Log.Warn("delConnPool uid :" + strconv.Itoa(int(uid)))
	connManager.PoolRWLock.Lock()
	defer connManager.PoolRWLock.Unlock()

	delete(connManager.Pool, uid)
}

//获取所有已连接用户的FD
func (connManager *ConnManager) getPoolAll() map[int32]*Conn {
	connManager.PoolRWLock.RLock()
	defer connManager.PoolRWLock.RUnlock()

	pool := make(map[int32]*Conn)
	for k, v := range connManager.Pool {
		pool[k] = v
	}
	return pool
}

func (connManager *ConnManager) GetPlayerCtrlInfoById(userId int32) ProtocolCtrlInfo {
	var contentType int32
	var protocolType int32
	if userId == 0 {
		contentType = connManager.Option.DefaultContentType
		protocolType = connManager.Option.DefaultProtocolType
	} else {
		conn, exist := connManager.GetConnPoolById(userId)
		if !exist {
			contentType = connManager.Option.DefaultContentType
			protocolType = connManager.Option.DefaultProtocolType
		} else {
			contentType = conn.ContentType
			protocolType = conn.ProtocolType
		}
	}

	protocolCtrlInfo := ProtocolCtrlInfo{
		ContentType:  contentType,
		ProtocolType: protocolType,
	}

	//connManager.Option.Log.Info("GetPlayerCtrlInfo uid : " + strconv.Itoa(int(userId)) + " contentType:" + strconv.Itoa(int(contentType)) + " , protocolType:" + strconv.Itoa(int(protocolType)))

	return protocolCtrlInfo
}

////这里是做个 容错处理，content type 跟 protocol type 不能为空，一但出现为空的情况，得给一个默认值
//func(connManager *ConnManager)GetCtrlInfo(contentType int32 ,protocolType int32)ProtocolCtrlInfo{
//	if contentType <= 0 {
//		contentType = connManager.Option.DefaultContentType
//	}
//
//	if protocolType <= 0 {
//		protocolType = connManager.Option.DefaultProtocol
//	}
//	protocolCtrlInfo := ProtocolCtrlInfo{
//		ContentType: contentType,
//		ProtocolType: protocolType,
//	}
//	return protocolCtrlInfo
//}

//==========================================================
//将 结构体 压缩成 字符串
func (connManager *ConnManager) CompressContent(contentStruct interface{}, UserId int32) (content []byte, err error) {
	//先获取该连接的通信元数据
	protocolCtrlInfo := connManager.GetPlayerCtrlInfoById(UserId)
	contentType := protocolCtrlInfo.ContentType

	if contentType == CONTENT_TYPE_JSON {
		//这里有个问题：纯JSON格式与PROTOBUF格式在PB文件上 不兼容
		//严格来说是GO语言与protobuf不兼容，即：PB文件的  结构体中的 JSON-TAG
		//PROTOBUF如果想使用驼峰式变量名，即：成员变量名区分出大小写，那必须得用<下划线>分隔，编译后，下划线转换成大写字母
		//编译完成后，虽然支持了驼峰变量名，但json-tag 并不是驼峰式，却是<下划线>式
		//那么，在不想改PB文件的前提下，就得在程序中做兼容

		//所以，先将content 字符串 由下划线转成 驼峰式
		content, err = json.Marshal(JsonCamelCase{contentStruct})
	} else if contentType == CONTENT_TYPE_PROTOBUF {
		contentStruct := contentStruct.(proto.Message)
		content, err = proto.Marshal(contentStruct)
	} else {
		err = errors.New(" contentType switch err")
	}
	if err != nil {
		return content, connManager.MakeError(err.Error())
	}

	return content, nil
}

//将 结构体 压缩成 字符串，但不依赖userId
func (connManager *ConnManager) CompressNormalContent(contentStruct interface{}, contentType int) (content []byte, err error) {
	requestClientHeartbeatStrByte := []byte{}
	if contentType == CONTENT_TYPE_PROTOBUF {
		contentStructReal := contentStruct.(proto.Message)
		requestClientHeartbeatStrByte, err = proto.Marshal(contentStructReal)
	} else {
		requestClientHeartbeatStrByte, err = json.Marshal(JsonCamelCase{contentStruct})
		//requestClientHeartbeatStrByte, err = json.Marshal(contentStruct)
	}
	return requestClientHeartbeatStrByte, err
}

//解析C端发送的数据，这一层，对于用户层的content数据不做处理
//1-4字节：当前包数据总长度，~可用于：TCP粘包的情况
//5字节：content type
//6字节：protocol type
//7字节 :服务Id
//8-9字节 :函数Id
//10-19：预留，还没想好，可以存sessionId，也可以换成UID
//19 以后为内容体
//结尾会添加一个字节：\f ,可用于 TCP 粘包 分隔
func (connManager *ConnManager) GetPackHeaderLength() int {
	return 4 + 1 + 1 + 1 + 2 + 10
}

//解析二进制流 -> msg结构体
func (connManager *ConnManager) ParserContentProtocol(content string) (message pb.Msg, err error) {
	headerLength := connManager.GetPackHeaderLength()
	if len(content) < headerLength {
		return message, errors.New("content len " + strconv.Itoa(headerLength) + "  < " + " header len" + strconv.Itoa(headerLength))
	}
	if len(content) == headerLength {
		errMsg := "content = " + strconv.Itoa(headerLength) + " ,body is empty"
		return message, errors.New(errMsg)
	}
	//数据长度
	dataLength := BytesToInt32([]byte(content[0:4]))
	if dataLength <= 0 {
		errMsg := "dataLength <= 0"
		return message, errors.New(errMsg)
	}
	//contentType + protocolType
	ctrlStream := content[4:6]
	ContentType, ProtocolType := connManager.parserProtocolCtrlInfo([]byte(ctrlStream))
	serviceId := int(content[6:7][0])
	actionId := BytesToInt32(BytesCombine([]byte{0, 0}, []byte(content[7:9])))
	//保留字
	reserved := content[9:19]
	serviceActionId, _ := strconv.Atoi(strconv.Itoa(serviceId) + strconv.Itoa(actionId))

	connManager.Option.Log.Warn(
		"contentLen:" + strconv.Itoa(len(content)) + " ,ContentType:" + strconv.Itoa(int(ContentType)) + " ,ProtocolType:" + strconv.Itoa(int(ProtocolType)) +
			" , dataLength:" + strconv.Itoa(dataLength) + " actionId:" + strconv.Itoa(actionId) + " serviceId:" + strconv.Itoa(serviceId) + " session:" + reserved)
	_, empty := connManager.Option.ProtoMap.GetServiceFuncById(serviceActionId) //这里只是提前检测一下service funcId 是否正确
	if empty {
		errMsg := "actionId ProtocolActions.GetActionName empty!!!"
		//protocolManager.Option.Log.Error(errMsg,actionId)
		return message, errors.New(errMsg)
	}
	//提取数据,ps: tcp 会自动删除末尾分隔符，而ws会有分隔符的
	data := content[19 : 19+dataLength]
	connManager.Option.Log.Debug("ParserContentProtocol content:" + string(data))
	msg := pb.Msg{
		Id:           0,
		SidFid:       int32(serviceActionId),
		FuncId:       int32(actionId),
		ServiceId:    int32(serviceId),
		DataLength:   int32(dataLength),
		Content:      data,
		ContentType:  ContentType,
		ProtocolType: ProtocolType,
		Reserved:     reserved,
	}
	//protocolManager.Option.Log.Debug("parserContentProtocol msg:",msg)
	return msg, nil
}

//字节 转换 协议控制信息
func (connManager *ConnManager) parserProtocolCtrlInfo(stream []byte) (int32, int32) {
	//firstByteHighThreeBit := (firstByte >> 5 ) & 7
	//firstByteLowThreeBit := ((firstByte << 5 ) >> 5 )  & 7
	//protocolCtrlInfo := connManager.GetCtrlInfo(int32(stream[0]),int32(stream[1]))
	//mylog.Debug("parserProtocolCtrlInfo ContentType:",protocolCtrlInfo.ContentType,",ProtocolType:",protocolCtrlInfo.ProtocolType)
	return int32(stream[0]), int32(stream[1])
}

//将消息 压缩成二进制
//func  (protocolManager *ProtocolManager)packContentMsg(content []byte,conn *Conn ,serviceId int ,actionId int )[]byte{
func (connManager *ConnManager) PackContentMsg(msg pb.Msg) []byte {
	dataLengthBytes := Int32ToBytes(int32(len(msg.Content)))
	contentTypeBytes := byte(msg.ContentType)
	protocolTypeBytes := byte(msg.ProtocolType)
	//actionIdByte := Int32ToBytes(msg.ActionId)
	//actionIdByte = actionIdByte[2:4]
	funcId, _ := strconv.Atoi(strconv.Itoa(int(msg.FuncId))[2:])
	actionIdByte := Int32ToBytes(int32(funcId))[2:]
	reserved := []byte("reserved--")
	serviceIdBytes := Int32ToBytes(msg.ServiceId)[3]
	ln := connManager.Option.MsgSeparator
	//connManager.Option.Log.Info("PackContentMsg dataLengthBytes:" + strconv.Itoa(len(msg.Content)))
	//合并: 头 + 消息内容体 + 分隔符
	content := BytesCombine(dataLengthBytes, ByteTurnBytes(contentTypeBytes), ByteTurnBytes(protocolTypeBytes), ByteTurnBytes(serviceIdBytes), actionIdByte, reserved, []byte(msg.Content), []byte(ln))
	return content
}

//==============================
type ConnMetrics struct {
	ReceiveMsgTimes int   `json:"receive_msg_times"`
	ReceiveMsgSize  int64 `json:"receive_msg_size"`
	SendMsgTimes    int   `json:"send_msg_times"`
	SendMsgSize     int64 `json:"send_msg_size"`
}

//一个连接
type Conn struct {
	AddTime      int32        `json:"add_time"`      //添加时间
	UpTime       int32        `json:"up_time"`       //最后更新时间
	UserId       int32        `json:"user_id"`       //FD 绑定 uid
	Status       int          `json:"status"`        //状态
	RTT          int64        `json:"rtt"`           //延迟
	SessionId    string       `json:"session_id"`    //会话ID
	ContentType  int32        `json:"content_type"`  //传输数据的内容类型	此值由第一次通信时确定，直到断开连接前，不会变更
	ProtocolType int32        `json:"protocol_type"` //传输数据的类型		此值由第一次通信时确定，直到断开连接前，不会变更
	Metrics      ConnMetrics  `json:"metrics"`       //测量/统计，公共metrics是放汇总数据的，这个是针对一个用户的
	Conn         FDAdapter    `json:"-"`             //TCP/WS 真实的 socket FD
	CloseChan    chan int     `json:"-"`             //减半信号管道
	netWay       *NetWay      `json:"-"`             //父父类,快捷指针
	ConnManager  *ConnManager `json:"-"`             //父类，快捷指针
	MsgInChan    chan pb.Msg  `json:"-"`             //上层真实sock FD 转发消息进来
	//帧同步在用，后期优化掉
	RoomId         string `json:"-"` //使用网关模式后，room_id 下放到各微服务上，网关不做记录
	UserPlayStatus int    `json:"-"` //status 冲突 暂放弃,不能删 frame_sync game_match 还在使用，后期优化吧
	//RTTCancelChan chan int
	//UdpConn 		bool
}

//最后更新时间
func (conn *Conn) UpLastTime() {
	conn.UpTime = int32(GetNowTimeSecondToInt())
}

//直接从FD中读取一条原始消息(未做解析)
func (conn *Conn) Read() (content string, err error) {
	// 设置消息的最大长度 - 暂无
	//conn.Conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(mynetWay.Option.IOTimeout)))
	_, dataByte, err := conn.Conn.ReadMessage()
	//messageType, dataByte, err := conn.Conn.ReadMessage()
	//_ , dataByte  , err  := conn.Conn.ReadMessage()
	if err != nil {
		//myMetrics.fastLog("total.input.err.num",METRICS_OPT_INC,0)
		conn.ConnManager.Option.Log.Error("conn.Conn.ReadMessage err: " + err.Error())
		return content, err
	}
	conn.Metrics.ReceiveMsgTimes++
	conn.Metrics.ReceiveMsgSize = conn.Metrics.ReceiveMsgSize + int64(len(dataByte))
	conn.ConnManager.Option.Metrics.CounterInc("total_input_num")
	//conn.ConnManager.Option.Metrics.GaugeAdd("total.input.size",float64(StringToFloat(strconv.Itoa(len(dataByte)))))
	conn.ConnManager.Option.Metrics.GaugeAdd("total_input_size", float64(len(dataByte)))

	//pid := strconv.Itoa(int(conn.UserId))
	//myMetrics.fastLog("player.fd.num."+pid,METRICS_OPT_INC,0)
	//myMetrics.fastLog("player.fd.size."+pid,METRICS_OPT_PLUS,len(content))

	//conn.ConnManager.Option.Log.Info("conn.ReadMessage messageType:" + strconv.Itoa(messageType) + " len :" + strconv.Itoa(len(dataByte)) + " data:" + string(dataByte))
	content = string(dataByte)
	return content, nil
}
func (conn *Conn) UpPlayerRoomId(roomId string) {
	conn.RoomId = roomId
}
func (conn *Conn) IOLoop() {
	conn.ConnManager.Option.Log.Info("conn IOLoop:")
	conn.ConnManager.Option.Log.Info("set conn status :" + strconv.Itoa(CONN_STATUS_EXECING) + " make close chan")
	conn.Status = CONN_STATUS_EXECING
	conn.CloseChan = make(chan int)
	ctx, cancel := context.WithCancel(context.Background())
	go conn.ReadLoop(ctx)       //读取消息，拆包，然后投入到队列中
	go conn.ProcessMsgLoop(ctx) //从队列中取出已拆包的值，做下一步处理：router
	//进入阻塞，监听外部<取消>操作
	<-conn.CloseChan
	cancel()
	conn.ConnManager.Option.Log.Warn("IOLoop receive chan quit~~~")
	conn.Conn.Close()
	//取消上面两个协程

}

////一个协程挂了，再给拉起来
//func  (conn *Conn) RecoverReadLoop(ctx context.Context){
//	conn.ConnManager.Option.Log.Warn("recover ReadLoop:")
//	go conn.ReadLoop(ctx)
//}
//死循环，从底层已读取出的消息中，再读取消息
func (conn *Conn) ReadLoop(ctx context.Context) {
	//defer func(ctx context.Context) {
	//	if err := recover(); err != nil {
	//		conn.ConnManager.Option.Log.Panic("conn.ReadLoop panic defer :")
	//		conn.RecoverReadLoop(ctx)
	//	}
	//}(ctx)
	for {
		select {
		case <-ctx.Done():
			conn.ConnManager.Option.Log.Warn("connReadLoop receive signal: ctx.Done.")
			goto end
		default:
			//从ws 读取 数据
			content, err := conn.Read()
			if err != nil {
				IsUnexpectedCloseError := websocket.IsUnexpectedCloseError(err)
				conn.ConnManager.Option.Log.Warn("connReadLoop connRead err:" + err.Error() + "IsUnexpectedCloseError:")
				if IsUnexpectedCloseError {
					conn.CloseOneConn(CLOSE_SOURCE_CLIENT_WS_FD_GONE)
					goto end
				} else {
					//对端，如果直接关闭网络，或者崩溃之类的，类库捕捉不到这个事件
					errMsg := strings.ToLower(err.Error())
					fdNetworkAbnormal := strings.Contains(errMsg, "connection reset by peer")
					if fdNetworkAbnormal {
						//这里可能有重复关闭的情况，如：C端监听到网络断了，然后立刻重连，等于kick off 会关闭一次，这里又监测到一次
						conn.CloseOneConn(CLOSE_SOURCE_CONN_RESET_BY_PEER)
						goto end
					}

					continue
				}
			}

			if content == "" {
				continue
			}
			//最后更新时间
			conn.UpLastTime()
			if len(content) > int(conn.ConnManager.Option.MsgContentMax) {
				errMsg := "msg content len > max content " + strconv.Itoa(int(conn.ConnManager.Option.MsgContentMax)) + " " + strconv.Itoa(len(content))
				conn.ConnManager.Option.Log.Error(errMsg)
				return
			}
			//解析消息内容
			msg, err := conn.ConnManager.ParserContentProtocol(content)
			if err != nil {
				conn.ConnManager.Option.Log.Warn("parserContent err :" + err.Error())
				continue
			}
			msg.SourceUid = conn.UserId
			//写入队列，等待其它协程处理，继续死循环
			conn.MsgInChan <- msg
		}
	}
end:
	conn.ConnManager.Option.Log.Warn("connReadLoop receive signal: done.")
}

//func  (conn *Conn) RecoverProcessMsgLoop(ctx context.Context){
//	conn.ConnManager.Option.Log.Warn("recover ReadLoop:")
//	go conn.ProcessMsgLoop(ctx)
//}

//关闭一个已登陆成功的FD,之所以放在最外层，是方便统一管理
func (conn *Conn) CloseOneConn(source int) {
	conn.ConnManager.Option.Log.Info("Conn close ,source : " + strconv.Itoa(source) + " , " + strconv.Itoa(int(conn.UserId)))
	if conn.Status == CONN_STATUS_CLOSE {
		conn.ConnManager.Option.Log.Error("CloseOneConn error :Conn.Status == CLOSE")
		return
	}

	if conn.Status == CONN_STATUS_CLOSE_ING {
		conn.ConnManager.Option.Log.Error("CloseOneConn error :Conn.Status == CONN_STATUS_CLOSE_ING")
		return
	}
	conn.Status = CONN_STATUS_CLOSE_ING

	connCloseEvent := pb.FDCloseEvent{
		UserId:       conn.UserId,
		Source:       int32(source),
		ContentType:  conn.ContentType,
		ProtocolType: conn.ProtocolType,
	}
	//先通知外层，FD即将要关闭了，外层网关再广播给所有微服务
	requestClientCloseStrByte, _ := conn.ConnManager.CompressNormalContent(connCloseEvent, int(connCloseEvent.ContentType))
	msg, _, _ := conn.ConnManager.MakeMsgByActionName(connCloseEvent.UserId, "FdClose", requestClientCloseStrByte)
	conn.ConnManager.Option.NetWay.Router(msg, conn)
	//通知同步服务，先做构造处理
	//mySync.CloseOne(conn)//这里可能还要再发消息

	//状态更新为已关闭，防止重复关闭
	conn.Status = CONN_STATUS_CLOSE
	//把后台守护  协程 先关了，不再收消息了
	conn.CloseChan <- 1
	//netWay.Players.delById(Conn.PlayerId)//这个不能删除，用于玩家掉线恢复的记录
	//先把玩家的在线状态给变更下，不然 mySync.close 里面获取房间在线人数，会有问题
	//myPlayerManager.upPlayerStatus(conn.UserId, PLAYER_STATUS_OFFLINE)
	if source != CLOSE_SOURCE_CLIENT {
		//客户端主动关闭，本层属于被动通知，底层已经知道了连接断开了，不用再关闭FD了
		err := conn.Conn.Close()
		if err != nil {
			conn.ConnManager.Option.Log.Error("Conn.Conn.Close err:" + err.Error())
		}
	}
	conn.ConnManager.delConnPool(conn.UserId)
	//处理掉-已报名的玩家
	//myMatch.realDelOnePlayer(conn.PlayerId)
	//mySleepSecond(2,"CloseOneConn")
	//myMetrics.fastLog("total.fd.num",METRICS_OPT_DIM,0)
	//myMetrics.fastLog("history.fd.destroy",METRICS_OPT_INC,0)
	//netWay.Metrics.CounterDec("total.fd.num")
	if source == CLOSE_SOURCE_CLIENT {
		conn.ConnManager.Option.Metrics.CounterInc("server_close_fd")
	} else {
		conn.ConnManager.Option.Metrics.CounterInc("client_close_fd")
	}

}

//从：FD里读取的消息（缓存队列），拿出来，做分发路由，处理
func (conn *Conn) ProcessMsgLoop(ctx context.Context) {
	//defer func(ctx context.Context) {
	//	if err := recover(); err != nil {
	//		conn.ConnManager.Option.Log.Panic("conn.ProcessMsgLoop panic defer :")
	//		conn.RecoverProcessMsgLoop(ctx)
	//	}
	//}(ctx)

	for {
		ctxHasDone := 0
		select {
		case <-ctx.Done():
			ctxHasDone = 1
		case msg := <-conn.MsgInChan:
			conn.ConnManager.Option.Log.Info("ProcessMsgLoop receive msg , SidFid:" + strconv.Itoa(int(msg.SidFid)))
			conn.ConnManager.Option.NetWay.Router(msg, conn)
		}
		if ctxHasDone == 1 {
			goto end
		}
	}
end:
	conn.ConnManager.Option.Log.Warn("ProcessMsgLoop receive signal: done.")
}

//监听到某个FD被关闭后，回调函数
func (conn *Conn) CloseHandler(code int, text string) error {
	conn.CloseOneConn(CLOSE_SOURCE_CLIENT)
	return nil
}

//===================================================================
//发送一条消息给一个玩家，根据conn，同时将消息内容进行编码与压缩
//大部分通信都是这个方法
func (conn *Conn) SendMsgCompressByConn(actionName string, contentStruct interface{}) error {
	if conn.UserId <= 0 {
		return conn.ConnManager.MakeError("conn.UserId <= 0")
	}
	conn.ConnManager.Option.Log.Info("SendMsgCompressByConn  actionName:" + actionName)
	//conn.UserId=0 时，由函数内部做兼容，主要是用来取content type ,protocol type
	contentByte, err := conn.ConnManager.CompressContent(contentStruct, conn.UserId)
	if err != nil {
		return err
	}
	return conn.SendMsg(actionName, contentByte)
}

//发送一条消息给一个玩家，根据UserId，同时将消息内容进行编码与压缩
func (conn *Conn) SendMsgCompressByUid(UserId int32, action string, contentStruct interface{}) error {
	if UserId <= 0 {
		return conn.ConnManager.MakeError("conn UserId <= 0")
	}

	conn.ConnManager.Option.Log.Info("SendMsgCompressByUid UserId:" + strconv.Itoa(int(UserId)) + " action:" + action)
	contentByte, err := conn.ConnManager.CompressContent(contentStruct, UserId)
	if err != nil {
		return err
	}
	return conn.SendMsgByUid(UserId, action, contentByte)
}

//发送一条消息给一个玩家,根据UserId,且不做压缩处理
func (conn *Conn) SendMsgByUid(UserId int32, action string, content []byte) error {
	conn, ok := conn.ConnManager.GetConnPoolById(UserId)
	if !ok {
		return conn.ConnManager.MakeError("conn not in pool,maybe del.")
	}
	return conn.SendMsg(action, content)
}

//发送一条消息给一个玩家,根据UserId,且不做压缩处理
//SendMsg的别名方法，不确定还有什么用，可能其它程序还在用吧
func (conn *Conn) SendMsgByConn(action string, content []byte) {
	conn.SendMsg(action, content)
}

func (connManager *ConnManager) MakeMsgByActionName(userId int32, action string, content []byte) (msg pb.Msg, protocolCtrlInfo ProtocolCtrlInfo, err error) {
	//获取协议号结构体
	actionMap, empty := connManager.Option.ProtoMap.GetServiceFuncByFuncName(action)
	if empty {
		return msg, protocolCtrlInfo, connManager.MakeError("GetActionId is  empty:" + action)
	}
	protocolCtrlInfo = connManager.GetPlayerCtrlInfoById(userId)

	//SidFid, _ := strconv.Atoi(strconv.Itoa(actionMap.ServiceId) + strconv.Itoa(actionMap.Id))
	msg = pb.Msg{
		//SidFid:       int32(SidFid),
		DataLength:   int32(len(content)),
		Content:      string(content),
		SidFid:       int32(actionMap.Id),
		ServiceId:    int32(actionMap.ServiceId),
		FuncId:       int32(actionMap.Id),
		ContentType:  protocolCtrlInfo.ContentType,
		ProtocolType: protocolCtrlInfo.ProtocolType,
	}

	return msg, protocolCtrlInfo, nil
}

func (connManager *ConnManager) MakeMsgByActionId(userId int32, actionId int, content []byte) (msg pb.Msg, protocolCtrlInfo ProtocolCtrlInfo, err error) {
	//获取协议号结构体
	actionMap, empty := connManager.Option.ProtoMap.GetServiceFuncById(actionId)
	if empty {
		return msg, protocolCtrlInfo, connManager.MakeError("GetActionId is  empty:" + strconv.Itoa(actionId))
	}
	protocolCtrlInfo = connManager.GetPlayerCtrlInfoById(userId)

	//SidFid, _ := strconv.Atoi(strconv.Itoa(actionMap.ServiceId) + strconv.Itoa(actionMap.Id))
	msg = pb.Msg{
		//SidFid:       int32(SidFid),
		DataLength:   int32(len(content)),
		Content:      string(content),
		SidFid:       int32(actionMap.Id),
		ServiceId:    int32(actionMap.ServiceId),
		FuncId:       int32(actionMap.Id),
		ContentType:  protocolCtrlInfo.ContentType,
		ProtocolType: protocolCtrlInfo.ProtocolType,
	}

	return msg, protocolCtrlInfo, nil
}

//这里才是，最终给一个User发送一条消息
//此方法轻易不要调用，因为没有做容错，比如：conn 不在pool中，直接调用 ProtoMap/uid=0 等，会把程序带崩溃了
func (conn *Conn) SendMsg(action string, content []byte) error {
	if conn.Status == CONN_STATUS_CLOSE {
		conn.ConnManager.Option.Log.Error("Conn status =CONN_STATUS_CLOSE.")
		return conn.ConnManager.MakeError("Conn status =CONN_STATUS_CLOSE.")
	}

	msg, protocolCtrlInfo, err := conn.ConnManager.MakeMsgByActionName(conn.UserId, action, content)
	if err != nil {
		return err
	}
	debugMsg := "SendMsg , ContentType:" + strconv.Itoa(int(protocolCtrlInfo.ContentType)) + " ProtocolType: " + strconv.Itoa(int(protocolCtrlInfo.ProtocolType)) + " , actionId:" + strconv.Itoa(int(msg.Id)) + " , userId:" + strconv.Itoa(int(conn.UserId)) + " , actionName:" + action
	conn.ConnManager.Option.Log.Info(debugMsg)

	contentBytes := conn.ConnManager.PackContentMsg(msg)

	//这里先注释掉，发现WS协议，传输内容必须统一：要么全是字符，要么就是二进制，而我的协议中 头消息是二进制的，内容如果是json那就是字符，貌似WS不行
	//if protocolCtrlInfo.ContentType == CONTENT_TYPE_PROTOBUF {
	return conn.Write(contentBytes, websocket.BinaryMessage)
	//} else {
	//	conn.Write(contentBytes, websocket.TextMessage)
	//}
}

//这里才是，真正的往sock FD息(严格说：跟用户没关联了) 里写内容，也就是发送消
func (conn *Conn) Write(content []byte, messageType int) error {
	//defer func() {
	//	if err := recover(); err != nil {
	//		conn.ConnManager.Option.Log.Error("conn.Conn.WriteMessage failed:")
	//		myNetWay.CloseOneConn(conn,CLOSE_SOURCE_SEND_MESSAGE)
	//	}
	//}()

	conn.Metrics.SendMsgTimes++
	conn.Metrics.SendMsgSize = conn.Metrics.ReceiveMsgSize + int64(len(content))

	//myMetrics.fastLog("total.output.num",METRICS_OPT_INC,0)
	//myMetrics.fastLog("total.output.size",METRICS_OPT_PLUS,len(content))
	conn.ConnManager.Option.Metrics.CounterInc("total_output_num")
	//conn.ConnManager.Option.Metrics.GaugeAdd("total.output.size",float64(StringToFloat(strconv.Itoa(len(content)))))
	conn.ConnManager.Option.Metrics.GaugeAdd("total_output_size", float64(len(content)))

	//pid := strconv.Itoa(int(conn.UserId))
	//myMetrics.fastLog("player.fd.num."+pid,METRICS_OPT_INC,0)
	//myMetrics.fastLog("player.fd.size."+pid,METRICS_OPT_PLUS,len(content))

	return conn.Conn.WriteMessage(messageType, content)
}
