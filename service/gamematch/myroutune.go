package gamematch

import (
	"go.uber.org/zap"
	"strconv"
	"time"
)

//关闭一条rule的所有守护协和
func  (gamematch *Gamematch)closeOneRuleDemonRoutine( ruleId int )int{
	mylog.Info("closeOneRuleDemonRoutine ,ruleId :  " + strconv.Itoa(ruleId)  + " start" )
	//1.先关闭httpd 入口
	_,ok := gamematch.HttpdRuleState[ruleId]
	if !ok{
		mylog.Error("closeOneRuleDemonRoutine ,gamematch.HttpdRuleState[ruleId]  is null "+strconv.Itoa(ruleId))
		return -1
	}
	mylog.Info("close httpd ,up state...")
	gamematch.HttpdRuleState[ruleId] = HTTPD_RULE_STATE_CLOSE
	//判断ruleId 是否已正确加载到rule池中
	_,ok = gamematch.RuleConfig.GetById(ruleId)
	if !ok{
		mylog.Error("closeOneRuleDemonRoutine ,ruleId not in db  ... "+strconv.Itoa(ruleId))
		return -1
	}
	//获取该rule下的所有管道
	gamematch.signalChanRWLock.RLock()
	_,ok = gamematch.signalChan[ruleId]
	gamematch.signalChanRWLock.RUnlock()
	if !ok{
		mylog.Error("closeOneRuleDemonRoutine ,gamematch.signalChan[ruleId]  is null "+strconv.Itoa(ruleId))
		return -2
	}
	//循环，给该条rule下的每个协程发送管道信号
	for title,mychann := range gamematch.signalChan[ruleId]{
		//mylog.Warning(SIGNAL_SEND_DESC  , ruleId,title )
		//mychann <- SIGN_GOROUTINE_EXEC_EXIT
		gamematch.signSend(mychann,SIGNAL_GOROUTINE_EXEC_EXIT,title)
	}

	closeRoutineCnt := 0
	for{
		if closeRoutineCnt == len(gamematch.signalChan[ruleId]){
			break
		}
		for title,chanLink := range gamematch.signalChan[ruleId]{
			select {
			case sign := <- chanLink:
				mylog.Warn("SIGNAL send" +strconv.Itoa(sign)+"退出完成 , ruleId: " + strconv.Itoa(ruleId) + title)
				closeRoutineCnt++
			//default:
			//	mylog.Info("closeDemonRoutine  waiting for goRoutine callback signal...")
			}
		}
	}
	//协程均已经结束，再回收管道
	for title,mychann := range gamematch.signalChan[ruleId]{
		mylog.Warn("close mychann ,"+strconv.Itoa(ruleId)+title)
		close(mychann)
	}
	gamematch.signalChanRWLock.Lock()
	//从map中移除该管道变量
	delete(gamematch.signalChan,ruleId)
	gamematch.signalChanRWLock.Unlock()
	mylog.Warn( "delete gamematch.signalChan map key :"+strconv.Itoa(ruleId))

	mylog.Info("closeOneRuleDemonRoutine ,ruleId :  " + strconv.Itoa(ruleId)+" finish." )
	return closeRoutineCnt
}
func  (gamematch *Gamematch)closeDemonRoutine(  )int{
	mylog.Info("closeDemonRoutine : start " )

	//要先关闭入口，也就是HTTPD
	httpdChan := gamematch.signalChan[0]["httpd"]
	//select {
	//	case httpdChan <- SIGN_GOROUTINE_EXEC_EXIT:
	//		mylog.Warning(SIGNAL_SEND_DESC ,SIGN_GOROUTINE_EXEC_EXIT,0,"httpd")
	//	case sign := <- httpdChan:
	//		mylog.Warning(SIGNAL_RECE_DESC ,sign,0,"httpd")
	//}
	gamematch.signSend(httpdChan,SIGNAL_GOROUTINE_EXEC_EXIT,"httpd")
	gamematch.signReceive(httpdChan,"httpd")

	mylog.Warn("gamematch.signalChan len:"+strconv.Itoa(len(gamematch.signalChan)))
	for ruleId,_ := range gamematch.signalChan{
		if ruleId == 0{//0是特殊管理，仅给HTTPD使用
			continue
		}
		gamematch.closeOneRuleDemonRoutine(ruleId)
		//for title,chanLink := range set{
		//	mylog.Warning(SIGNAL_SEND_DESC  , ruleId,title )
		//	chanLink <- SIGN_GOROUTINE_EXEC_EXIT
		//}
	}
	mylog.Info("closeDemonRoutine : all routine quit success~~~")
	return 1

	//chanHasFinished := make(map[int] map[string]  int)
	//chanHasFinishedCnt := 0
	////signalChanLen := len(gamematch.signalChan)
	//countSignalChanNumber := gamematch.countSignalChan(0)
	//for{
	//	if chanHasFinishedCnt == countSignalChanNumber{
	//		break
	//	}
	//
	//	for ruleId,set := range gamematch.signalChan{
	//		if ruleId == 0{//0是特殊通道，留给httpd
	//			continue
	//		}
	//		for title,chanLink := range set{
	//			_,ok :=  chanHasFinished[ruleId][title]
	//			if ok{
	//				continue
	//			}
	//
	//			select {
	//				case sign := <- chanLink:
	//					mylog.Warning(SIGNAL_RECE_DESC ,sign,ruleId,title)
	//					_,ok := chanHasFinished[ruleId]
	//					if !ok {
	//						tmp := make(map[string] int)
	//						chanHasFinished[ruleId] = tmp
	//					}
	//					chanHasFinished[ruleId][title] = 1
	//					chanHasFinishedCnt++
	//				default:
	//					mylog.Info("closeDemonRoutine  waiting for goRoutine callback signal...")
	//			}
	//		}
	//	}
	//	//	_,ok :=  chanHasFinished[i]
	//	//	if ok{
	//	//		continue
	//	//	}
	//	//	select {
	//	//		case sign := <- gamematch.signalChan[i]:
	//	//			mylog.Warning(SIGNAL_RECE_DESC + strconv.Itoa(sign))
	//	//			chanHasFinished[i] = sign
	//	//		default:
	//	//			mylog.Info("closeDemonRoutine  waiting for goRoutine callback signal...")
	//	//	}
	//	//}
	//	mySleepSecond(1,"closeDemonRoutine")
	//}

}
//接收一个管道的数据-会阻塞
func (gamematch *Gamematch)signReceive(mychan chan int,keyword string)int{
	mylog.Warn("SIGNAL receive blocking.... ")
	sign := <- mychan
	mylog.Warn("SIGNAL receive: "+strconv.Itoa(sign) + getSignalDesc(sign) + keyword)
	return sign
}
func (gamematch *Gamematch)signSend(mychan chan int,data int,keyword string){
	mylog.Warn("SIGNAL send blocking.... " + keyword)
	mychan <- data
	mylog.Warn("SIGNAL send: " + strconv.Itoa(data) + getSignalDesc(data) +keyword)
}
//两个入口在调用：HTTPD MYDEMON
func (gamematch *Gamematch)NewSignalChan(ruleId int,title string)chan int{
	//mylog.Debug("getNewSignalChan: ",ruleId,title)
	signalChann := make(chan int)
	gamematch.signalChanRWLock.RLock()
	_,ok := gamematch.signalChan[ruleId]
	gamematch.signalChanRWLock.RUnlock()


	gamematch.signalChanRWLock.Lock()
	if !ok {
		tmp := make(map[string]chan int)
		gamematch.signalChan[ruleId] = tmp
	}
	gamematch.signalChan[ruleId][title] = signalChann

	gamematch.signalChanRWLock.Unlock()
	return gamematch.signalChan[ruleId][title]
}
//开启一个守护协程，这里只是统一管理
//handler:每秒回调一次这个函数
func (gamematch *Gamematch)StartOneGoroutineDemon( ruleId int ,title string,demonLog *zap.Logger, handler func( )){
	msg := "StartOneGoroutineDemon :   rId  "+strconv.Itoa(ruleId)+" t: " +title
	mylog.Warn(msg)
	//demonLog.Warning(msg)
	signalChan := gamematch.NewSignalChan(ruleId,title)
	//rule,_ := gamematch.RuleConfig.GetById(ruleId)
	//signalChan <- SIGNAL_GOROUTINE_EXEC_ING
	//inc := 0
	for{
		select {
		case signal := <-signalChan:
			mylog.Warn("SIGNAL receive: "+strconv.Itoa(signal) + getSignalDesc(signal) + title)
			mylog.Warn("SIGNAL send: "+strconv.Itoa(SIGNAL_GOROUTINE_EXIT_FINISH) + getSignalDesc(SIGNAL_GOROUTINE_EXIT_FINISH) + title)
			signalChan <- SIGNAL_GOROUTINE_EXIT_FINISH
			goto forEnd
		default:
			handler()
			time.Sleep(time.Second * 1)
		}
	}
forEnd:
	//demonLog.Notice("MyDemon end : ",title)
	mylog.Warn("MyDemon end : "+title)
}


