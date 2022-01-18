package util

import (
	"encoding/json"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"regexp"
	"strconv"
	"strings"
	"zgoframe/protobuf/pb"
)

type ProtocolManager struct {
	TcpServer 		*TcpServer
	WsHttpServer    *Websocket
	Option 			ProtocolManagerOption
	Close 			chan int
}

type ProtocolManagerOption struct {
	Ip 				string
	WsPort	 		string
	TcpPort 		string
	UdpPort			string
	WsUri 			string
	OpenNewConnBack	func ( connFD FDAdapter)
	Log 			*zap.Logger
	ProtobufMap		*ProtobufMap
}
//
func NewProtocolManager(option ProtocolManagerOption)*ProtocolManager{
	option.Log.Info("NewProtocolManager instance:")
	protocolManager := new (ProtocolManager)

	protocolManager.Option 	= option
	protocolManager.Close 	= make(chan int)

	//go protocolManager.l()
	return protocolManager
}
func  (protocolManager *ProtocolManager)Shutdown(){
	//mylog.Alert("shutdown protocolManager")
	//protocolManager.Close <- 1
	protocolManager.TcpServer.Shutdown()
	protocolManager.WsHttpServer.Shutdown()
}

//func  (protocolManager *ProtocolManager)l(){
//	<- protocolManager.Close
//	protocolManager.TcpServer.Shutdown()
//	protocolManager.WsHttpServer.Shutdown()
//}

func (protocolManager *ProtocolManager)Start( )error{
	protocolManager.Option.Log.Info("protocolManager start:")
	websocketOption := WebsocketOption{
		WsUri		:protocolManager.Option.WsUri,
		Port  		:protocolManager.Option.WsPort,
		ListenerIp 	:protocolManager.Option.Ip,
		OutIp 		:protocolManager.Option.Ip,
		Log			:	protocolManager.Option.Log,
		ProtocolManager : protocolManager,
		//OutCxt 		:	protocolManager.Option.
		//OpenNewConnBack	func ( connFD FDAdapter)	//来了新连接后，回调函数
	}
	websocket := NewWebsocket(websocketOption)

	protocolManager.WsHttpServer = websocket
	go protocolManager.WsHttpServer.Start()
	//tcp server
	tcpServerOption := TcpServerOption{
		Ip		: protocolManager.Option.Ip,
		Port	: protocolManager.Option.TcpPort,
		//OutCxt :  context.Context,
		Log 	:protocolManager.Option.Log,
		ProtocolManager :protocolManager,

	}
	myTcpServer :=  NewTcpServer(tcpServerOption)
	protocolManager.TcpServer = myTcpServer
	err := myTcpServer.Start()
	return err
}
func(protocolManager *ProtocolManager)websocketHandler( connFD *websocket.Conn) {
	protocolManager.Option.Log.Info("websocketHandler: have a new client")
	imp := WebsocketConnImpNew(connFD)
	protocolManager.Option.OpenNewConnBack(imp)
}

func(protocolManager *ProtocolManager)tcpHandler(tcpConn *TcpConn){
	imp := TcpConnImpNew(tcpConn)
	protocolManager.Option.OpenNewConnBack(imp)
}

func(protocolManager *ProtocolManager)udpHandler(){

}
//=======================================
//协议层的解包已经结束，这个时候需要将content内容进行转换成MSG结构
func  (protocolManager *ProtocolManager)parserContentMsg(msg pb.Msg ,out interface{},playerId int32)error{
	content := msg.Content
	var err error
	//protocolCtrlInfo := myPlayerManager.GetPlayerCtrlInfoById(playerId)
	//contentType := protocolCtrlInfo.ContentType
	if msg.ContentType == CONTENT_TYPE_JSON {
		unTrunVarJsonContent := CamelToSnake([]byte(content))
		err = json.Unmarshal(unTrunVarJsonContent,out)
	}else if  msg.ContentType == CONTENT_TYPE_PROTOBUF {
		aaa := out.(proto.Message)
		err = proto.Unmarshal([]byte(content),aaa)
	}else{
		protocolManager.Option.Log.Error("parserContent err")
	}

	if err != nil{
		protocolManager.Option.Log.Error("parserMsgContent:"+err.Error())
		return err
	}

	protocolManager.Option.Log.Debug("protocolManager parserMsgContent:")

	return nil
}

func  ByteTurnBytes(b byte)[]byte{
	var a []byte
	a = append(a,b)
	return a
}

//func  (protocolManager *ProtocolManager)packContentMsg(content []byte,conn *Conn ,serviceId int ,actionId int )[]byte{
func  (protocolManager *ProtocolManager)PackContentMsg(msg pb.Msg)[]byte{
	dataLengthBytes := Int32ToBytes( int32( len(msg.Content) ))
	//protocolCtrlInfo := myNetWay.ConnManager.GetPlayerCtrlInfoById(conn.UserId)
	//int32 -> int -> string - > bytes
	contentTypeBytes := byte( msg.ContentType)
	protocolTypeBytes :=  byte(msg.ProtocolType)

	actionIdByte := Int32ToBytes(msg.ActionId)
	//actionIdByte := []byte(strconv.Itoa(int(msg.ActionId)))
	reserved := []byte( "reserved--")

	//serviceIdBytes := []byte(strconv.Itoa(int(msg.ServiceId)))
	serviceIdBytes := Int32ToBytes(msg.ServiceId)
	ln := "\n"
	//合并 头 + 消息内容体
	//content  := BytesCombine(dataLengthBytes,contentTypeBytes,protocolTypeBytes,serviceIdBytes,actionIdByte,reserved,[]byte(msg.Content),[]byte(ln))
	content  := BytesCombine(dataLengthBytes,ByteTurnBytes(contentTypeBytes),ByteTurnBytes(protocolTypeBytes),serviceIdBytes,actionIdByte,reserved,[]byte(msg.Content),[]byte(ln))
	return content
	////var protocolCtrlFirstByteArr []byte
	////contentTypeByte := byte(contentType)
	////protocolTypeByte := byte(player.ProtocolType)
	////contentTypeByteRight := contentTypeByte >> 5
	////protocolCtrlFirstByte := contentTypeByteRight | protocolTypeByte
	////protocolCtrlFirstByteArr = append(protocolCtrlFirstByteArr,protocolCtrlFirstByte)
	////content = zlib.BytesCombine(protocolCtrlFirstByteArr,content)
	//contentTypeStr := strconv.Itoa(int(contentType))
	//protocolTypeStr := strconv.Itoa(int(protocolType))
	//contentTypeAndprotocolType := contentTypeStr + protocolTypeStr
	//content = BytesCombine([]byte(contentTypeAndprotocolType),content)
	////myMetrics.IncNode("output_num")
	////myMetrics.PlusNode("output_size",len(content))
	////房间做统计处理
	////if action =="pushLogicFrame"{
	////	roomId := myPlayerManager.GetRoomIdByUserId(conn.UserId)
	////	roomSyncMetrics := RoomSyncMetricsPool[roomId]
	////	roomSyncMetrics.OutputNum++
	////	roomSyncMetrics.OutputSize = roomSyncMetrics.OutputSize + len(content)
	////}
	//netWay.Option.Log.Debug("final sendmsg ctrlInfo: contentType-" + string(contentTypeBytes) + " protocolType-" + string(protocolCtrlInfo.ProtocolType) + " actionId-" + string(actionIdByte))
	//netWay.Option.Log.Debug("final sendmsg content:" +string(content))
}

//解析C端发送的数据，这一层，对于用户层的content数据不做处理
//1-4字节：当前包数据总长度，~可用于：TCP粘包的情况
//5字节：content type
//6字节：protocol type
//7字节 :服务Id
//8-9字节 :函数Id
//10-19：预留，还没想好，可以存sessionId，也可以换成UID
//19 以后为内容体
//结尾会添加一个字节：\n ,可用于 TCP 粘包 分隔
func  (protocolManager *ProtocolManager)GetPackHeaderLength()int{
	return 4 + 1 + 1 + 1 + 2 + 10
}
func  (protocolManager *ProtocolManager)parserContentProtocol(content string)(message pb.Msg,err error){
	headerLength := protocolManager.GetPackHeaderLength()
	if len(content) < headerLength{
		return message,errors.New("content < "+ strconv.Itoa(headerLength))
	}
	if len(content)==headerLength{
		errMsg := "content = "+strconv.Itoa(headerLength)+" ,body is empty"
		return message,errors.New(errMsg)
	}
	//数据长度
	dataLength := BytesToInt32([]byte(content[0:4]))
	//contentType + protocolType
	ctrlStream := content[4:6]
	MyPrint("ctrlStream:",ctrlStream)
	ctrlInfo := protocolManager.parserProtocolCtrlInfo([]byte(ctrlStream))

	actionId := BytesToInt32([]byte(content[6:7]))
	serviceId := BytesToInt32([]byte(content[7:9]))
	//保留字
	reserved :=  content[9:19]

	protocolManager.Option.Log.Warn("dataLength:"+strconv.Itoa(dataLength) + " actionId:"+strconv.Itoa(actionId) +  " serviceId:"+strconv.Itoa(serviceId))
	actionMap,empty := protocolManager.Option.ProtobufMap.GetActionName(actionId)
	if empty{
		errMsg := "actionId ProtocolActions.GetActionName empty!!!"
		//protocolManager.Option.Log.Error(errMsg,actionId)
		return message,errors.New(errMsg)
	}
	//提取数据
	data := content[19:]
	msg := pb.Msg{
		ActionId: int32(actionId),
		ServiceId: int32(serviceId),
		DataLength: int32(dataLength),
		Action: actionMap.Action,
		Content:data,
		ContentType : ctrlInfo.ContentType,
		ProtocolType: ctrlInfo.ProtocolType,
		Reserved: reserved,
	}
	//protocolManager.Option.Log.Debug("parserContentProtocol msg:",msg)
	return msg,nil
}


type ProtocolCtrlInfo struct {
	ContentType int32
	ProtocolType int32
}
func (protocolManager *ProtocolManager)parserProtocolCtrlInfo(stream []byte)ProtocolCtrlInfo{
	//firstByte := stream[0:1][0]
	//mylog.Debug("firstByte:",firstByte)
	//firstByteHighThreeBit := (firstByte >> 5 ) & 7
	//firstByteLowThreeBit := ((firstByte << 5 ) >> 5 )  & 7
	firstByteHighThreeBit , _:= strconv.Atoi(string(stream[0:1]))
	firstByteLowThreeBit , _:= strconv.Atoi(string(stream[1:2]))
	protocolCtrlInfo := ProtocolCtrlInfo{
		ContentType : int32(firstByteHighThreeBit),
		ProtocolType : int32(firstByteLowThreeBit),
	}
	//mylog.Debug("parserProtocolCtrlInfo ContentType:",protocolCtrlInfo.ContentType,",ProtocolType:",protocolCtrlInfo.ProtocolType)
	return protocolCtrlInfo
}
//将 结构体 压缩成 字符串
func  (protocolManager *ProtocolManager)CompressContent(contentStruct interface{},UserId int32)(content []byte  ,err error){
	//先获取该连接的通信元数据
	protocolCtrlInfo := myNetWay.ConnManager.GetPlayerCtrlInfoById(UserId)
	contentType 	:= protocolCtrlInfo.ContentType

	if contentType == CONTENT_TYPE_JSON {
		//这里有个问题：纯JSON格式与PROTOBUF格式在PB文件上 不兼容
		//严格来说是GO语言与protobuf不兼容，即：PB文件的  结构体中的 JSON-TAG
		//PROTOBUF如果想使用驼峰式变量名，即：成员变量名区分出大小写，那必须得用<下划线>分隔，编译后，下划线转换成大写字母
		//编译完成后，虽然支持了驼峰变量名，但json-tag 并不是驼峰式，却是<下划线>式
		//那么，在不想改PB文件的前提下，就得在程序中做兼容

		//所以，先将content 字符串 由下划线转成 驼峰式
		content, err = json.Marshal(JsonCamelCase{contentStruct})
	}else if  contentType == CONTENT_TYPE_PROTOBUF {
		contentStruct := contentStruct.(proto.Message)
		content, err = proto.Marshal(contentStruct)
	}else{
		err = errors.New(" contentType switch err")
	}
	if err != nil{
		protocolManager.Option.Log.Error("CompressContent err :"+err.Error())
	}
	return content,err
}

type JsonCamelCase struct {
	Value interface{}
}
//下划线 转 驼峰命
func Case2Camel(name string) string {
	//将 下划线 转 空格
	name = strings.Replace(name, "_", " ", -1)
	//将 字符串的 每个 单词 的首字母转大写
	name = strings.Title(name)
	//最后再将空格删掉
	return strings.Replace(name, " ", "", -1)
}

func (c JsonCamelCase) MarshalJSON() ([]byte, error) {
	var keyMatchRegex = regexp.MustCompile(`\"(\w+)\":`)
	marshalled, err := json.Marshal(c.Value)
	converted := keyMatchRegex.ReplaceAllFunc(
		marshalled,
		func(match []byte) []byte {
			matchStr := string(match)
			key := matchStr[1 : len(matchStr)-2]
			resKey := Lcfirst(Case2Camel(key))
			return []byte(`"` + resKey + `":`)
		},
	)
	return converted, err
}