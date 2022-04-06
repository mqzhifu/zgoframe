package service

import (
	"errors"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"os"
	"strings"
	"zgoframe/util"
)

const (
	PERSITENCE_TYPE_MYSQL   = 1
	PERSITENCE_TYPE_REDIS   = 2
	PERSITENCE_TYPE_FILE    = 3
	PERSITENCE_TYPE_ETCD    = 4
	PERSITENCE_TYPE_CONSULE = 5
)

type ConfigCenter struct {
	Option ConfigCenterOption
	pool   map[string]map[string]map[string]*viper.Viper
}

type ConfigCenterOption struct {
	envList            []string
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
	return configCenter, err
}

func (configCenter *ConfigCenter) Init() error {
	if len(configCenter.Option.envList) <= 0 {
		return errors.New("env list  len <=0")
	}

	if configCenter.Option.PersistenceType <= 0 {
		return errors.New("PersistenceType <=0")
	}

	if configCenter.Option.PersistenceType == PERSITENCE_TYPE_FILE {
		return configCenter.InitPersistenceFile()
	}

	return nil
}

func (configCenter *ConfigCenter) GetByFile(env string, projectName string, fileName string) (data interface{}, err error) {
	myViper, ok := configCenter.pool[env][projectName][fileName]
	if !ok {
		return data, err
	}

	data = myViper.AllSettings()
	return data, nil

}

func (configCenter *ConfigCenter) GetByKey(env string, projectName string, fileName string, key string) (data interface{}, err error) {
	myViper, ok := configCenter.pool[env][projectName][fileName]
	if !ok {
		return data, err
	}

	data = myViper.Get(key)
	return data, nil

}

func (configCenter *ConfigCenter) SetByKey(env string, projectName string, key string, fileName string, value interface{}) (err error) {
	myViper, ok := configCenter.pool[env][projectName][fileName]
	if !ok {
		return err
	}

	myViper.Set(key, value)
	e := myViper.WriteConfig()
	return e
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
	envPool := make(map[string]map[string]map[string]*viper.Viper)
	for _, env := range configCenter.Option.envList {
		envDir := configCenter.Option.PersistenceFileDir + "/" + env
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
		//读取所有项目
		foreachDirList := util.ForeachDir(envDir)
		if len(foreachDirList) <= 0 {
			util.MyPrint("ForeachDir envDir list empty:" + envDir)
			envPool[env] = make(map[string]map[string]*viper.Viper)
			continue
		}

		projectPool := make(map[string]map[string]*viper.Viper)
		for _, projectDirInfo := range foreachDirList {
			if projectDirInfo.Cate != "dir" {
				continue
			}

			projectInfo, empty := configCenter.Option.ProjectManager.GetByName(projectDirInfo.Name)
			if empty {
				continue
			}

			projectDir := envDir + "/" + projectInfo.Name
			util.MyPrint("projectDir:" + projectDir)
			//读取一个项目下的所有配置文件列表
			foreachProjectDirList := util.ForeachDir(projectDir)
			if len(foreachProjectDirList) <= 0 {
				return nil
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
			projectPool[projectDirInfo.Name] = filePool
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
