package model

// @description 站内信
type ProjectPushMsg struct {
	MODEL
	Type            int    `json:"type" db:"define:tinyint(1);comment:1-文字消息 2-截图 3-图片 4-模型文件 5-视频;defaultValue:0"` //
	SourceId        int    `json:"source_id" db:"define:int;comment:发送者的UID;defaultValue:0"`                        //
	SourceProjectId int    `json:"source_project_id" db:"define:int;comment:发送者的项目ID;defaultValue:0"`               //
	Content         string `json:"content" db:"define:varchar(255);comment:内容;defaultValue:''"`                     //
	TargetProjectId int    `json:"target_project_id" db:"define:int;comment:接收者的项目ID;defaultValue:0"`               //
	TargetUids      string `json:"target_uids" db:"define:varchar(255);comment:接收者的UID;defaultValue:''"`            //
	Date            int    `json:"date" db:"define:int;comment:年月日,用于统计;defaultValue:0"`                            //
}

func (projectPushMsg *ProjectPushMsg) TableOptions() map[string]string {
	m := make(map[string]string)
	m["comment"] = "长连接-推荐消息-日志表"

	return m
}
