package util

import (
	"gorm.io/gorm"
	"reflect"
	"strings"
	"zgoframe/model"
)

//根据结构体描述信息，生成一张mysql 表的创建 sql

type DbTool struct {
	FieldTagName string
	Br           string
}

type TableOption struct {
	Name    string
	Comment string
	Engine  string
	Charset string
	//ColumnsOption map[string]TableColumnOption
	ColumnsOption []TableColumnOption
}
type TableColumnOption struct {
	Name       string
	Comment    string
	Unique     string
	Index      string
	Primarykey string
	Unsigned   string
	Define     string

	DefaultValue  string
	Null          string
	AutoIncrement string
}

func NewDbTool(gorm *gorm.DB) *DbTool {
	db := new(DbTool)
	db.FieldTagName = "db"
	db.Br = "\n"
	return db
}

func (db *DbTool) CreateTable(tableStruct ...interface{}) {
	for i := 0; i < len(tableStruct); i++ {
		db.processOneTable(tableStruct[i])
		//ExitPrint(111)

	}

}

//1 2 3 4 5 6 7
//1 4 5 6 2 3 4

//1 2 3 4 5
//1 5 2 3 4

func (db *DbTool) AddField(tableStruct ...interface{}) {
	//ALTER TABLE `sms_rule` ADD `purpose` TINYINT(1) NOT NULL DEFAULT '0' COMMENT '用途,1注册2找回密码3修改密码4登陆' AFTER `third_callback_info`;
}

//公共model中的 创建时间 更新时间 是否删除 3个字段统一放到最后面
func (db *DbTool) ProcessCommonModelFiledSort(columnsOption []TableColumnOption) []TableColumnOption {
	if len(columnsOption) <= 4 { //公共字段一共就4个，小于4就没必要处理了
		return columnsOption
	}
	newColumnsOption := []TableColumnOption{}

	frontArr := db.mergeTwoArr(columnsOption[0:1], columnsOption[4:])
	newColumnsOption = db.mergeTwoArr(frontArr, columnsOption[1:4])
	//ExitPrint(newColumnsOption)
	//firstField := columnsOption[0:1]
	//frontField := tmp[1:4]
	//endField := columnsOption[4:]
	//
	//MyPrint("firstField:", firstField, "frontField:", frontField, "endField:", endField)
	//
	//frontArr := append(firstField, endField...)
	//ExitPrint(frontField)
	//newColumnsOption = append(frontArr, frontField...)

	//ExitPrint(newColumnsOption)

	return newColumnsOption
}

func (db *DbTool) mergeTwoArr(firstField []TableColumnOption, endField []TableColumnOption) []TableColumnOption {
	newColumnsOption := []TableColumnOption{}
	for _, v := range firstField {
		newColumnsOption = append(newColumnsOption, v)
	}
	//newColumnsOption = firstField
	for _, v := range endField {
		newColumnsOption = append(newColumnsOption, v)
	}
	return newColumnsOption
}

func (db *DbTool) processOneTable(tableStruct interface{}) {
	MyPrint("processOneTable:", tableStruct)

	ValueOfTableStruct := reflect.ValueOf(tableStruct)
	//查找方法
	method := ValueOfTableStruct.MethodByName("TableOptions")
	//动态执行方法
	rs := method.Call([]reflect.Value{})
	//执行方法结果：一个表的基础信息
	tableOptionString := rs[0].Interface().(map[string]string)
	tableOption := TableOption{
		Comment: tableOptionString["comment"],
	}
	//MyPrint(tableOption)
	if tableOption.Name == "" {
		structFullName := ValueOfTableStruct.Elem().Type()
		tableOption.Name = Lcfirst(structFullName.Name())
		MyPrint("tableOption.Name:", tableOption.Name)
	} else {
		MyPrint("err: tableOption.Name empty")
	}
	//MyPrint(tableOption)
	sql := db.Br + "create table " + string(CamelToSnake2([]byte(tableOption.Name))) + "(" + db.Br
	//columnsOption := []TableColumnOption{}

	TypeOfOneTableStruct := reflect.TypeOf(tableStruct)
	columnsOption := db.GetField(TypeOfOneTableStruct.Elem())

	columnsOption = db.ProcessCommonModelFiledSort(columnsOption)

	sqlMid := ""
	primarykey := false
	for _, v := range columnsOption {
		//ak := string(CamelToSnake2([]byte(k)))
		ak := string(CamelToSnake2([]byte(v.Name)))
		sql += "    `" + ak + "` "
		if v.Define != "" {
			sql += v.Define + " "
		}

		if v.Unsigned != "" {
			sql += " UNSIGNED "
		}
		if v.Name != "DeletedAt" { //这个特殊，得允许null
			if v.Define == "text" {
				sql += " null "
			} else {
				//if v.Null != ""{
				sql += " not null "
				//}
			}

		}

		if v.DefaultValue != "" {
			sql += " default   " + v.DefaultValue + " "
		}

		if v.AutoIncrement != "" {
			sql += " auto_increment "
		}

		if v.Comment != "" {
			sql += " comment '" + v.Comment + "' "
		}

		if v.Primarykey != "" {
			if primarykey {
				ExitPrint("primarykey repeat")
			}
			sql += " primary key "
			primarykey = true
		}

		if v.Unique != "" {
			sqlMid += " Unique " + " (`" + ak + "`)  ,"
		}

		if v.Index != "" {
			sqlMid += " index " + " (`" + ak + "`) ,"
		}
		sql += " , " + db.Br
	}
	if sqlMid != "" {
		sqlMid = string([]byte(sqlMid)[0 : len(sqlMid)-1])
		sql += sqlMid + ")" + db.Br
	}

	// charset=utf8,comment="test Table";
	Engine := "innodb"
	if tableOption.Engine != "" {
		ExitPrint(tableOption.Engine)
		Engine = tableOption.Engine
	}
	sql += " engine=" + Engine

	Charset := "utf8"
	if tableOption.Charset != "" {
		Charset = tableOption.Charset
	}
	sql += " charset=" + Charset

	comment := "''"
	if tableOption.Comment != "" {
		comment = "'" + tableOption.Comment + "'"
	}
	sql += " comment=" + comment + db.Br

	MyPrint(sql)
}

func (db *DbTool) GetField(typeOfStruct reflect.Type) []TableColumnOption {
	columnsOption := []TableColumnOption{}
	//func (db *DbTool) GetField(typeOfStruct reflect.Type, columnsOption map[string]TableColumnOption) {
	//MyPrint("typeOfStruct:", typeOfStruct)
	for i := 0; i < typeOfStruct.NumField(); i++ {
		tableOneColumnOption := TableColumnOption{}
		structFiled := typeOfStruct.Field(i)

		//MyPrint(structFiled.Type.String())
		if structFiled.Type.Name() == "MODEL" {
			commonMODEL := model.MODEL{}
			typeOfGlobalMODEL := reflect.TypeOf(commonMODEL)
			columnsOption = db.GetField(typeOfGlobalMODEL)
			continue
		}

		structFiledTagName := structFiled.Tag.Get(db.FieldTagName)
		//MyPrint(structFiledTagName)
		if structFiledTagName != "" {
			ValueOfTableOneColumnOption := reflect.ValueOf(&tableOneColumnOption)
			//ColumnMap := make(map[string]string)
			arr := strings.Split(structFiledTagName, ";")
			for _, oneKeyStr := range arr {
				//MyPrint("oneKeyStr :",oneKeyStr)
				oneKeyArr := strings.Split(oneKeyStr, ":")
				//MyPrint("oneKeyArr :",oneKeyArr)
				//MyPrint(ValueOfTableOneColumnOption.Elem().FieldByName(oneKeyArr[0]).Type())
				//MyPrint(oneKeyArr[1],StrFirstToUpper(oneKeyArr[0]))
				ValueOfTableOneColumnOption.Elem().FieldByName(StrFirstToUpper(oneKeyArr[0])).SetString(oneKeyArr[1])
			}
			//MyPrint(tableOneColumnOption)
		}
		tableOneColumnOption.Name = structFiled.Name
		//columnsOption[structFiled.Name] = tableOneColumnOption
		columnsOption = append(columnsOption, tableOneColumnOption)
	}
	return columnsOption
	//ExitPrint(1)
}
