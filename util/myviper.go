package util

import (
	"errors"
	"github.com/spf13/viper"
)

func ReadConfFile(configFile string,structOut interface{})error{
	_,err := FileExist(configFile)
	if err != nil{
		return errors.New("ReadConfFile "+configFile +" file not exist:"+err.Error())
	}

	myViper := viper.New()

	myViper.SetConfigType("toml")
	myViper.SetConfigFile(configFile)
	myViper.AddConfigPath(".")
	err = myViper.ReadInConfig()
	if err != nil{
		errMsg := "ReadConfFile "+configFile +"myViper.ReadInConfig() err :" + err.Error()
		return errors.New(errMsg)
	}

	//cicdConfig := util.CicdConfig{}
	err = myViper.Unmarshal(structOut)
	if err != nil{
		errMsg := "ReadConfFile "+configFile +" myViper.Unmarshal err:"+err.Error()
		return errors.New(errMsg)
	}

	return nil
}
