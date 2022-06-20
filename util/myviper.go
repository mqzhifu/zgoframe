package util

import (
	"errors"
	"github.com/spf13/viper"
	"strings"
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

func ReadConfFileAutoExt(configFile string,structOut interface{})error{
	_,err := FileExist(configFile)
	if err != nil{
		return errors.New("ReadConfFile "+configFile +" file not exist:"+err.Error())
	}

	myViper := viper.New()

	configFileArr := strings.Split(configFile,".")
	if len(configFileArr) < 2{
		return errors.New("配置文件名中：必须得有扩展名")
	}
	configFileExtName := configFileArr[len(configFileArr)-1]
	if configFileExtName == "toml"{
		myViper.SetConfigType("toml")
	}else if configFileExtName == "yaml"{
		myViper.SetConfigType("yaml")
	}else{
		return errors.New("配置文件扩展名仅支持：toml yaml")
	}

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
