package util

import "github.com/spf13/viper"

func ReadConfFile(configFile string,structOut interface{}){
	myViper := viper.New()

	myViper.SetConfigType("toml")
	//configFile := "host" + "." +"toml"
	//util.MyPrint(configFile)
	myViper.SetConfigFile(configFile)
	myViper.AddConfigPath(".")
	err := myViper.ReadInConfig()
	if err != nil{
		MyPrint("myViper.ReadInConfig() err :",err)
		ExitPrint("err.")
	}

	//cicdConfig := util.CicdConfig{}
	err = myViper.Unmarshal(structOut)
	if err != nil{
		MyPrint(" myViper.Unmarshal err:",err)
		ExitPrint("err2")
	}
}
