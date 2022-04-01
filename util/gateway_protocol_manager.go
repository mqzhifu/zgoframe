package util

/*
	协议管理适配器，类似 适配层，中转TCP WS UDP
*/
import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"regexp"
	"strings"
	"zgoframe/protobuf/pb"
)

type ProtocolManager struct {
	TcpServer    *TcpServer //tcp 管理器
	WsHttpServer *Websocket //ws  管理器
	Option       ProtocolManagerOption
}

type ProtocolManagerOption struct {
	Ip              string
	WsPort          string
	TcpPort         string
	UdpPort         string
	WsUri           string
	IOTimeout       int64
	MsgContentMax   int32
	Log             *zap.Logger
	OpenNewConnBack func(connFD FDAdapter) //新FD到来时，回调函数接口
}

//
func NewProtocolManager(option ProtocolManagerOption) *ProtocolManager {
	option.Log.Info("NewProtocolManager instance:")
	protocolManager := new(ProtocolManager)

	protocolManager.Option = option
	return protocolManager
}

//关闭
func (protocolManager *ProtocolManager) Shutdown() {
	protocolManager.Option.Log.Warn("protocolManager Shutdown:")

	protocolManager.TcpServer.Shutdown()
	protocolManager.WsHttpServer.Shutdown()
}

//启动监听
func (protocolManager *ProtocolManager) Start() error {
	protocolManager.Option.Log.Info("protocolManager start:")
	//创建WS协议管理器
	websocketOption := WebsocketOption{
		WsUri:           protocolManager.Option.WsUri,
		Port:            protocolManager.Option.WsPort,
		ListenerIp:      protocolManager.Option.Ip,
		OutIp:           protocolManager.Option.Ip,
		Log:             protocolManager.Option.Log,
		IOTime:          protocolManager.Option.IOTimeout,
		MsgContentMax:   protocolManager.Option.MsgContentMax,
		ProtocolManager: protocolManager,
		//OpenNewConnBack	func ( connFD FDAdapter)	//来了新连接后，回调函数
	}
	protocolManager.WsHttpServer = NewWebsocket(websocketOption)
	//这里无法立刻返回开启监听结果，回头优化
	go protocolManager.WsHttpServer.Start()
	//创建tcp协议管理器
	tcpServerOption := TcpServerOption{
		Ip:              protocolManager.Option.Ip,
		Port:            protocolManager.Option.TcpPort,
		Log:             protocolManager.Option.Log,
		IOTimeout:       protocolManager.Option.IOTimeout,
		MsgContentMax:   protocolManager.Option.MsgContentMax,
		ProtocolManager: protocolManager,
	}
	myTcpServer := NewTcpServer(tcpServerOption)
	protocolManager.TcpServer = myTcpServer
	err := myTcpServer.Start()
	return err
}
func (protocolManager *ProtocolManager) websocketHandler(connFD *websocket.Conn) {
	protocolManager.Option.Log.Info("websocketHandler: have a new client")
	imp := WebsocketConnImpNew(connFD)
	protocolManager.Option.OpenNewConnBack(imp)
}

func (protocolManager *ProtocolManager) tcpHandler(tcpConn *TcpConn) {
	imp := TcpConnImpNew(tcpConn)
	protocolManager.Option.OpenNewConnBack(imp)
}

func (protocolManager *ProtocolManager) udpHandler() {

}

//=======================================
//协议层的解包已经结束，这个时候需要将content内容进行转换成MSG结构
func (protocolManager *ProtocolManager) ParserContentMsg(msg pb.Msg, out interface{}, playerId int32) error {
	content := msg.Content
	var err error
	//protocolCtrlInfo := myPlayerManager.GetPlayerCtrlInfoById(playerId)
	//contentType := protocolCtrlInfo.ContentType
	if msg.ContentType == CONTENT_TYPE_JSON {
		unTrunVarJsonContent := CamelToSnake([]byte(content))
		err = json.Unmarshal(unTrunVarJsonContent, out)
	} else if msg.ContentType == CONTENT_TYPE_PROTOBUF {
		aaa := out.(proto.Message)
		err = proto.Unmarshal([]byte(content), aaa)
	} else {
		protocolManager.Option.Log.Error("parserContent err")
	}

	if err != nil {
		protocolManager.Option.Log.Error("parserMsgContent:" + err.Error())
		return err
	}

	protocolManager.Option.Log.Debug("protocolManager parserMsgContent:")

	return nil
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
