package service

import (
	"github.com/gin-gonic/gin"
)

func Receive(c *gin.Context){

}

func ProcessOne(data []byte){
	//logSlaveMsg := OneMsg{}
	//err := json.Unmarshal(data,&logSlaveMsg)
	//if err != nil{
	//	zlib.MyPrint("json.Unmarshal err:",err.Error())
	//	//ResponseStatusCode(w,512,"json.Unmarshal err:"+err.Error())
	//	return
	//}
	//msg := zlib.Msg{
	//	AppId: logSlaveMsg.AppId,
	//	ModuleId: logSlaveMsg.ModuleId,
	//	Content: logSlaveMsg.Content,
	//}
	//logSlave.OutLog.SlaveOut(logSlaveMsg.LeaveId,msg)
}
