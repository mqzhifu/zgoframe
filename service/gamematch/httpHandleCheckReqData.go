package gamematch

func (gamematch *Gamematch)CheckHttpSignCancelData(httpReqBusiness HttpReqBusiness)(errs error){
	if httpReqBusiness.GroupId == 0  {
		return myerr.New(457)
	}

	return nil
}

func (gamematch *Gamematch)CheckHttpSuccessDelData(httpReqBusiness HttpReqBusiness)( errs error){
	if httpReqBusiness.SuccessId == 0  {
		return myerr.New(460)
	}

	return nil
}

func (gamematch *Gamematch)CheckHttpSignData(httpReqBusiness HttpReqBusiness)(errs error){
	if httpReqBusiness.GroupId <= 0{
		return myerr.New(452)
	}
	var playerListStruct []Player
	for _,v := range httpReqBusiness.PlayerList{
		if v.Uid == 0{
			return myerr.New(456)
		}
		playerListStruct = append(playerListStruct,Player{Id:v.Uid})
	}
	return nil
}
