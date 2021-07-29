package main

import (
	"flag"
	"zgoframe/core"
	"zlib"
)


func main(){
	zlib.LogLevelFlag = zlib.LOG_LEVEL_DEBUG

	envList := zlib.GetEnvList()


	configType 		:= flag.String("ct", core.DEFAULT_CONFIT_TYPE, "configType")
	configFileName 	:= flag.String("cfn", core.DEFAULT_CONFIG_FILE_NAME, "configFileName")
	env 			:= flag.String("e", "must require", "env")

	flag.Parse()

	if !zlib.CheckEnvExist(*env){
		zlib.ExitPrint(  "env is err , list:",envList)
	}

	core.Init(*env,*configType,*configFileName)
}
