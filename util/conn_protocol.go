package util

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"regexp"
	"strconv"
	"strings"
)

type ConnProtocol struct {
	Log *zap.Logger
	ProtobufMap *ProtobufMap
}

func NewConnProtocol(protobufMap *ProtobufMap,log *zap.Logger)*ConnProtocol{
	connProtocol := new(ConnProtocol)
	connProtocol.Log = log
	connProtocol.ProtobufMap = protobufMap
	return connProtocol
}
//直接给一个FD发消息，不做任何处理
func(connProtocol *ConnProtocol)Send(conn *Conn,content string){
	msgType := connProtocol.ContentTypeCovertMegType(conn.ContentType)
	connProtocol.Log.Debug("send msgType:"+strconv.Itoa(msgType)+ " body:"+content)
	conn.Write([]byte(content),msgType)
}
//发送一条消息给一个玩家FD，
func(connProtocol *ConnProtocol)SendMsgByUid(uid int,action string , content []byte){
	conn,ok := myConnManager.getConnPoolById(uid)
	if !ok {
		connProtocol.Log.Error("conn not in pool,maybe del.")
		return
	}
	connProtocol.SendMsg(conn,action,content)
}
//发送一条消息给一个玩家，根据conn，同时将消息内容进行编码与压缩
//大部分通信都是这个方法
func(connProtocol *ConnProtocol)SendMsgCompressByConn(conn *Conn,action string , contentStruct interface{}){
	connProtocol.Log.Info("SendMsgCompressByConn  action:" + action)
	contentByte ,_ := connProtocol.CompressContent(contentStruct,conn)
	connProtocol.SendMsg(conn,action,contentByte)
}

//发送一条消息给一个玩家，根据playerId，同时将消息内容进行编码与压缩
func(connProtocol *ConnProtocol)SendMsgCompressByUid(uid int,action string , contentStruct interface{}){
	connProtocol.Log.Info("SendMsgCompressByUid playerId:" + strconv.Itoa(int(uid)) +" action:" + action)
	conn,ok := myConnManager.getConnPoolById(uid)
	if !ok {
		connProtocol.Log.Error("conn not in pool,maybe del.")
		return
	}
	contentByte ,_ := connProtocol.CompressContent(contentStruct,conn)
	connProtocol.SendMsgByUid(uid,action,contentByte)
}
//传输的内容类型转换成  协议传输内容类型
func(connProtocol *ConnProtocol)ContentTypeCovertMegType(ContentType int)int{
	if ContentType == CONTENT_TYPE_JSON || ContentType == CONTENT_TYPE_STRING {
		return websocket.TextMessage
	}else if ContentType == CONTENT_TYPE_PROTOBUF{
		return  websocket.BinaryMessage
	}else{
		return websocket.BinaryMessage
	}
}

//将 结构体 压缩成 字符串
func  (connProtocol *ConnProtocol)CompressContent(contentStruct interface{},conn *Conn)(content []byte  ,err error){
	protocolCtrlInfo := conn.GetCtrlInfoById()
	contentType := protocolCtrlInfo.ContentType

	//mylog.Debug("CompressContent contentType:",contentType)
	if contentType == CONTENT_TYPE_JSON {
		//这里有个问题：纯JSON格式与PROTOBUF格式在PB文件上 不兼容
		//严格来说是GO语言与protobuf不兼容，即：PB文件的  结构体中的 JSON-TAG
		//PROTOBUF如果想使用驼峰式变量名，即：成员变量名区分出大小写，那必须得用<下划线>分隔，编译后，下划线转换成大写字母
		//编译完成后，虽然支持了驼峰变量名，但json-tag 并不是驼峰式，却是<下划线>式
		//那么，在不想改PB文件的前提下，就得在程序中做兼容

		//所以，先将content 字符串 由下划线转成 驼峰式
		content, err = json.Marshal(JsonCamelCase{contentStruct})
		//mylog.Info("CompressContent json:",string(content),err )
	}else if  contentType == CONTENT_TYPE_PROTOBUF {
		//contentStruct := contentStruct.(proto.Message)
		//content, err = proto.Marshal(contentStruct)
	}else{
		err = errors.New(" switch err")
	}
	if err != nil{
		connProtocol.Log.Error("CompressContent err :" + err.Error())
	}
	return content,err
}

func(connProtocol *ConnProtocol)SendMsg(conn *Conn,action string,content []byte){
	//获取协议号结构体
	actionMapT,empty := connProtocol.ProtobufMap.GetActionId(action)
	//connProtocol.Log.Info("SendMsg",actionMapT.Id,conn.PlayerId,action)
	if empty{
		connProtocol.Log.Error("GetActionId empty:" + action)
		return
	}
	protocolCtrlInfo := conn.GetCtrlInfoById()
	contentType := protocolCtrlInfo.ContentType
	protocolType := protocolCtrlInfo.ProtocolType
	//player ,_ := myPlayerManager.GetById(conn.PlayerId)
	//SessionIdBtye := []byte(player.SessionId)
	content  = BytesCombine([]byte(conn.SessionId),content)
	//协议号
	strId := strconv.Itoa(int(actionMapT.Id))
	//合并 协议号 + 消息内容体
	content = BytesCombine([]byte(strId),content)
	if conn.Status == CONN_STATUS_CLOSE {
		connProtocol.Log.Error("Conn status =CONN_STATUS_CLOSE.")
		return
	}

	//var protocolCtrlFirstByteArr []byte
	//contentTypeByte := byte(contentType)
	//protocolTypeByte := byte(player.ProtocolType)
	//contentTypeByteRight := contentTypeByte >> 5
	//protocolCtrlFirstByte := contentTypeByteRight | protocolTypeByte
	//protocolCtrlFirstByteArr = append(protocolCtrlFirstByteArr,protocolCtrlFirstByte)
	//content = zlib.BytesCombine(protocolCtrlFirstByteArr,content)
	contentTypeStr := strconv.Itoa(int(contentType))
	protocolTypeStr := strconv.Itoa(int(protocolType))
	contentTypeAndprotocolType := contentTypeStr + protocolTypeStr
	content = BytesCombine([]byte(contentTypeAndprotocolType),content)
	//myMetrics.IncNode("output_num")
	//myMetrics.PlusNode("output_size",len(content))
	//房间做统计处理
	//if action =="pushLogicFrame"{
	//	roomId := myPlayerManager.GetRoomIdByPlayerId(conn.PlayerId)
	//	roomSyncMetrics := RoomSyncMetricsPool[roomId]
	//	roomSyncMetrics.OutputNum++
	//	roomSyncMetrics.OutputSize = roomSyncMetrics.OutputSize + len(content)
	//}
	connProtocol.Log.Debug("final sendmsg ctrlInfo: contentType-" + contentTypeStr + " protocolType-" + protocolTypeStr + " pid-" + strId)
	connProtocol.Log.Debug("final sendmsg content:" + string(content))
	//if contentType == CONTENT_TYPE_PROTOBUF {
	//	conn.Write(content,websocket.BinaryMessage)
		//netWay.myWriteMessage(Conn,websocket.BinaryMessage,content)
	//}else{
	//	conn.Write(content,websocket.TextMessage)
		//netWay.myWriteMessage(Conn,websocket.TextMessage,content)
	//}
	connProtocol.Send(conn,string(content))
}


//=======================================
//协议层的解包已经结束，这个时候需要将content内容进行转换成MSG结构
func  (connProtocol *ConnProtocol)ParserContentMsg(msg ConnMsg ,out interface{})error{
	content := msg.Content
	var err error
	//protocolCtrlInfo := myPlayerManager.GetPlayerCtrlInfoById(playerId)
	//contentType := protocolCtrlInfo.ContentType
	if msg.ContentType == CONTENT_TYPE_JSON {
		unTrunVarJsonContent := CamelToSnake([]byte(content))
		err = json.Unmarshal(unTrunVarJsonContent,out)
	}else if  msg.ContentType == CONTENT_TYPE_PROTOBUF {
		//aaa := out.(proto.Message)
		//err = proto.Unmarshal([]byte(content),aaa)
	}else{
		connProtocol.Log.Error("parserContent err")
	}

	if err != nil{
		connProtocol.Log.Error("parserMsgContent:" + err.Error())
		return err
	}

	connProtocol.Log.Debug("protocolManager parserMsgContent:")

	return nil
}
func (connProtocol *ConnProtocol)ParserProtocolCtrlInfo(stream []byte)ProtocolCtrlInfo{
	//firstByte := stream[0:1][0]
	//mylog.Debug("firstByte:",firstByte)
	//firstByteHighThreeBit := (firstByte >> 5 ) & 7
	//firstByteLowThreeBit := ((firstByte << 5 ) >> 5 )  & 7
	firstByteHighThreeBit , _:= strconv.Atoi(string(stream[0:1]))
	firstByteLowThreeBit , _:= strconv.Atoi(string(stream[1:2]))
	protocolCtrlInfo := ProtocolCtrlInfo{
		ContentType : firstByteHighThreeBit,
		ProtocolType : firstByteLowThreeBit,
	}
	//mylog.Debug("parserProtocolCtrlInfo ContentType:",protocolCtrlInfo.ContentType,",ProtocolType:",protocolCtrlInfo.ProtocolType)
	return protocolCtrlInfo
}

//解析C端发送的数据，这一层，对于用户层的content数据不做处理
//前2个字节控制流，3-6为协议号，7-38为sessionId
func (connProtocol *ConnProtocol)ParserContentProtocol(content string)(message ConnMsg,err error){
	protocolSum := 6
	if len(content) < protocolSum {
		return message,errors.New("content < "+ strconv.Itoa(protocolSum))
	}

	if len(content) == protocolSum {
		errMsg := "content = "+strconv.Itoa(protocolSum)+" ,body is empty"
		return message,errors.New(errMsg)
	}
	ctrlStream := content[0:2]
	ctrlInfo := connProtocol.ParserProtocolCtrlInfo([]byte(ctrlStream))
	actionIdStr := content[2:6]
	actionId,_ := strconv.Atoi(actionIdStr)
	actionName,empty := connProtocol.ProtobufMap.GetActionName(actionId)
	if empty{
		errMsg := "actionId ProtocolActions.GetActionName empty!!!"
		connProtocol.Log.Error(errMsg + strconv.Itoa(actionId))
		return message,errors.New("actionId ProtocolActions.GetActionName empty!!!")
	}

	//mylog.Info("parserContent actionid:",actionId, ",actionName:",actionName.Action)

	sessionId := ""
	userData := ""
	if actionName.Action != "login"{
		sessionId = content[6:38]
		userData = content[38:]
	}else{
		userData = content[6:]
	}

	msg := ConnMsg{
		Action: actionName.Action,
		Content:userData,
		ContentType : ctrlInfo.ContentType,
		ProtocolType: ctrlInfo.ProtocolType,
		SessionId: sessionId,
	}
	//mylog.Debug("parserContentProtocol msg:",msg)
	return msg,nil
}

//==================================================
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

