package service

import (
	"errors"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"os"
	"strconv"
	"strings"
	"zgoframe/util"
)

const (
	PERSISTENCE_TYPE_MYSQL   = 1
	PERSISTENCE_TYPE_REDIS   = 2
	PERSISTENCE_TYPE_FILE    = 3
	PERSISTENCE_TYPE_ETCD    = 4
	PERSISTENCE_TYPE_CONSULE = 5
)
func GetConstListConfigPersistenceType() map[string]int {
	list := make(map[string]int)
	list["mysql"] = PERSISTENCE_TYPE_MYSQL
	list["redis"] = PERSISTENCE_TYPE_REDIS
	list["file"] = PERSISTENCE_TYPE_FILE
	list["etcd"] = PERSISTENCE_TYPE_ETCD
	list["consul"] = PERSISTENCE_TYPE_CONSULE

	return list
}

type ConfigCenter struct {
	Option ConfigCenterOption
	pool   map[int]map[string]map[string]*viper.Viper
}

type ConfigCenterOption struct {
	envList            map[string]int
	Gorm               *gorm.DB
	Redis              *util.MyRedis
	ProjectManager     *util.ProjectManager
	PersistenceType    int
	PersistenceFileDir string
}

func NewConfigCenter(Option ConfigCenterOption) (*ConfigCenter, error) {
	configCenter := new(ConfigCenter)
	configCenter.Option = Option

	err := configCenter.Init()
	if err != nil {
		util.ExitPrint(err)
	}

	return configCenter, nil
}

func (configCenter *ConfigCenter) Init() error {
	if len(configCenter.Option.envList) <= 0 {
		return errors.New("env list  len <=0")
	}

	if configCenter.Option.PersistenceType <= 0 {
		return errors.New("PersistenceType <=0")
	}

	if configCenter.Option.PersistenceType == PERSISTENCE_TYPE_FILE {
		return configCenter.InitPersistenceFile()
	}

	return nil
}
//以模块(文件)为单位，获取该模块(文件)下的所有配置信息
func (configCenter *ConfigCenter) GetByModule(env int, projectId int, module string) (data interface{}, err error) {
	myViper ,err := configCenter.GetModuleInfo(env,projectId,module)
	if err != nil{
		return data,err
	}

	data = myViper.AllSettings()
	return data, nil

}


//以以模块(文件)+里面具体的key 为单位，获取配置信息
func (configCenter *ConfigCenter) GetByKey(env int,projectId int, module string, key string) (data interface{}, err error) {
	myViper ,err := configCenter.GetModuleInfo(env,projectId,module)
	if err != nil{
		return data,err
	}

	data = myViper.Get(key)
	return data, nil

}
//以模块(文件)+里面具体的key 为单位，设置置信息(如果存在，覆盖)
func (configCenter *ConfigCenter) SetByKey(env int, projectId int,  module string, key string, value interface{}) (err error) {
	myViper ,err := configCenter.GetModuleInfo(env,projectId,module)
	if err != nil{
		return err
	}

	myViper.Set(key, value)
	e := myViper.WriteConfig()
	return e
}
func (configCenter *ConfigCenter)GetModuleInfo(env int , projectId int,module string)(myViper *viper.Viper, err error){
	util.MyPrint("GetModuleInfo ,  env:"+strconv.Itoa(env)  + " projectId:"+strconv.Itoa(projectId) + " module:"+ module)
	project , empty := configCenter.Option.ProjectManager.GetById(projectId)
	if empty{
		return myViper,errors.New("projectId is empty")
	}

	myViper, ok := configCenter.pool[env][project.Name][module]
	if !ok {
		return myViper, errors.New("module is empty")
	}

	return myViper,nil
}
func (configCenter *ConfigCenter) InitPersistenceFile() error {
	prefix := "InitPersistenceFile "
	if configCenter.Option.PersistenceFileDir == "" {
		return errors.New(prefix + "PersistenceFileDir == '' ")
	}

	_, err := util.PathExists(configCenter.Option.PersistenceFileDir)
	if err != nil {
		return errors.New(prefix + err.Error())
	}
	//             map[env][projectName][fileName]myViper
	envPool := make(map[int]map[string]map[string]*viper.Viper)
	for _, env := range configCenter.Option.envList {
		//envDir := configCenter.Option.PersistenceFileDir + "/" + strconv.Itoa(env)
		envDir := configCenter.Option.PersistenceFileDir + "/" + util.GetConstListEnvStr()[env]

		util.MyPrint("envDir:" + envDir)
		_, err = util.PathExists(envDir)
		if err != nil { //
			if !os.IsNotExist(err) {
				return errors.New(prefix + err.Error())
			}
		}

		if os.IsNotExist(err) {
			util.MyPrint(prefix + " mkdir:" + envDir)
			err = os.Mkdir(envDir, os.ModePerm)
			if err != nil {
				return errors.New(prefix + err.Error())
			}
			envPool[env] = make(map[string]map[string]*viper.Viper)
			continue
		}
		////读取所有项目
		//foreachDirList := util.ForeachDir(envDir)
		//if len(foreachDirList) <= 0 {
		//	util.MyPrint("ForeachDir envDir list empty:" + envDir)
		//	envPool[env] = make(map[string]map[string]*viper.Viper)
		//	continue
		//}

		projectPool := make(map[string]map[string]*viper.Viper)
		//util.MyPrint(configCenter.Option.ProjectManager.Pool)
		//util.ExitPrint(configCenter.Option.ProjectManager.Pool)
		for _, projectInfo := range configCenter.Option.ProjectManager.Pool {
			//util.MyPrint(projectInfo)
			//if projectDirInfo.Cate != "dir" {
			//	continue
			//}
			//
			//projectInfo, empty := configCenter.Option.ProjectManager.GetByName(projectDirInfo.Name)
			//if empty {
			//	continue
			//}

			projectDir := envDir + "/" + projectInfo.Name
			util.MyPrint("projectDir:" + projectDir)
			_, err = util.PathExists(projectDir)
			if err != nil { //
				if !os.IsNotExist(err) {
					//util.ExitPrint("err 1:",err)
					return errors.New(prefix + err.Error())
				}else{
					err = os.Mkdir(projectDir, os.ModePerm)
					if err != nil {
						//util.ExitPrint("err 2:",err)
						return errors.New(prefix + err.Error())
					}
				}
			}
			//读取一个项目下的所有配置文件列表
			foreachProjectDirList := util.ForeachDir(projectDir)
			if len(foreachProjectDirList) <= 0 {
				continue
			}
			filePool := make(map[string]*viper.Viper)
			for _, fileInfo := range foreachProjectDirList {
				if fileInfo.Cate != "file" {
					continue
				}

				fileNameArr := strings.Split(fileInfo.Name, ".")
				if len(fileNameArr) < 2 {
					continue
				}

				extName := fileNameArr[len(fileNameArr)-1]
				if extName != "toml" {
					continue
				}

				util.MyPrint("file:" + fileInfo.Name)

				fileDir := projectDir + "/" + fileInfo.Name
				myViper, err := configCenter.ViperRead(fileDir)
				if err != nil {
					return err
				}
				//info := myViper.Get("mysql")
				//util.ExitPrint(err, "tt:", info)
				filePool[fileInfo.Name] = myViper
			}
			projectPool[projectInfo.Name] = filePool
		}
		envPool[env] = projectPool
	}
	configCenter.pool = envPool
	return nil
}

func (configCenter *ConfigCenter) ViperRead(configFile string) (*viper.Viper, error) {
	myViper := viper.New()

	myViper.SetConfigType("toml")
	myViper.SetConfigFile(configFile)
	myViper.AddConfigPath(".")

	err := myViper.ReadInConfig()
	if err != nil {
		errMsg := "ReadConfFile " + configFile + "myViper.ReadInConfig() err :" + err.Error()
		return myViper, errors.New(errMsg)
	}
	return myViper, nil
}
