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
	global.V.Util.ProjectMng, err = util.NewProjectManager(global.V.Base.Gorm)
	if err != nil {
		return err
	}
	empty := false
	global.V.Util.Project, empty = global.V.Util.ProjectMng.GetById(global.C.System.ProjectId)
	if empty {
		return errors.New("AppId not match : " + strconv.Itoa(global.C.System.ProjectId))
	}

	global.V.Base.Zap.Info(prefix + "project info ,  id : " + strconv.Itoa(global.V.Util.Project.Id) + " , name : " + global.V.Util.Project.Name)

	return nil
}
