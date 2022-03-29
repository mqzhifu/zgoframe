package initialize

import (
	"errors"
	"strconv"
	"zgoframe/core/global"
	"zgoframe/util"
)

func InitProject(prefix string) (err error) {
	if global.C.System.ProjectId <= 0 {
		return errors.New("ProjectId empty")
	}
	global.V.ProjectMng, err = util.NewProjectManager(global.V.Gorm)
	if err != nil {
		return err
	}
	empty := false
	global.V.Project, empty = global.V.ProjectMng.GetById(global.C.System.ProjectId)
	if empty {
		return errors.New("AppId not match : " + strconv.Itoa(global.C.System.ProjectId))
	}

	global.V.Zap.Info(prefix + "project info ,  id : " + strconv.Itoa(global.V.Project.Id) + " , name : " + global.V.Project.Name)

	return nil
}
