package util

import (
	"gorm.io/gorm"
	"reflect"
	"strings"
	"zgoframe/model"
)

type Db struct {
	FieldTagName string
	Br string
}



type TableOption struct {
	Name string
	Comment string
	Engine string
	Charset string
	ColumnsOption map[string]TableColumnOption
}
type TableColumnOption struct {
	Comment string
	Unique string
	Index string
	Primarykey string
	Unsigned string
	Define string

	DefaultValue string
	Null string
	AutoIncrement string
}


func NewDb(gorm *gorm.DB)*Db{
	db := new (Db)
	db.FieldTagName = "db"
	db.Br = "\n"
	return db
}

func(db *Db) CreateTable(tableStruct ...interface{} ){
	for i:=0;i<len(tableStruct);i++{
		db.processOneTable(tableStruct[i])

	}
	ExitPrint(111)
}

func (db *Db)processOneTable(tableStruct interface{}){
	MyPrint("processOneTable:",tableStruct)

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
		tableOption.Name = Lcfirst( structFullName.Name() )
	}
	//MyPrint(tableOption)
	sql := db.Br+"create table " + string( CamelToSnake2([]byte(tableOption.Name)))  + "(" + db.Br
	columnsOption := make(map[string]TableColumnOption)

	TypeOfOneTableStruct := reflect.TypeOf(tableStruct)
	db.GetField(TypeOfOneTableStruct.Elem(),columnsOption)


	//type TableColumnOption struct {
	//	Unique string
	//	Index string
	//}

	sqlMid := ""
	primarykey := false
	for k,v :=  range columnsOption{
		ak := string(CamelToSnake2([]byte(k)))
		//ak := CamelToSnake2([]byte("addHeaderCddddDdddd"))
		sql +=  "    `"+ ak + "` "
		if v.Define != ""{
			sql += v.Define + " "
		}

		if v.Unsigned != ""{
			sql += " UNSIGNED "
		}
		if k != "DeletedAt"{//这个特殊，得允许null
			if v.Define == "text"{
				sql += " null "
			}else{
				//if v.Null != ""{
				sql += " not null "
				//}
			}

		}


		if v.DefaultValue != ""{
			sql += " default   "+  v.DefaultValue + " "
		}

		if v.AutoIncrement != ""{
			sql += " auto_increment "
		}

		if v.Comment != ""{
			sql += " comment '"+v.Comment+"' "
		}

		if v.Primarykey != ""{
			if primarykey {
				ExitPrint("primarykey repeat")
			}
			sql += " primary key "
			primarykey  = true
		}

		if v.Unique != ""{
			sqlMid += " Unique " + " (`"+ak+"`)  ,"
		}

		if v.Index != ""{
			sqlMid += " index " + " (`"+ak+"`) ,"
		}
		sql +=  " , " + db.Br
	}
	if sqlMid != ""{
		sqlMid = string([]byte(sqlMid)[0:len(sqlMid)-1])
		sql += sqlMid + ")" + db.Br
	}



	// charset=utf8,comment="test Table";
	Engine := "innodb"
	if tableOption.Engine != ""{
		ExitPrint(tableOption.Engine)
		Engine = tableOption.Engine
	}
	sql += " engine=" + Engine

	Charset := "utf8"
	if tableOption.Charset != ""{
		Charset = tableOption.Charset
	}
	sql += " charset="+ Charset

	comment := "''"
	if tableOption.Comment != ""{
		comment = "'" +tableOption.Comment + "'"
	}
	sql += " comment="+ comment + db.Br

	MyPrint(sql)

}

func  (db *Db)GetField(typeOfStruct reflect.Type,columnsOption map[string]TableColumnOption){
	for i:=0;i<typeOfStruct.NumField();i++{
		tableOneColumnOption := TableColumnOption{}
		structFiled := typeOfStruct.Field(i)

		//MyPrint(structFiled.Type.String())
		if structFiled.Type.Name() == "MODEL"{
			commonMODEL := model.MODEL{}
			typeOfGlobalMODEL := reflect.TypeOf(commonMODEL)
			db.GetField(typeOfGlobalMODEL,columnsOption)
			continue
		}

		structFiledTagName := structFiled.Tag.Get(db.FieldTagName)
		//MyPrint(structFiledTagName)
		if structFiledTagName != ""{
			ValueOfTableOneColumnOption := reflect.ValueOf(&tableOneColumnOption)
			//ColumnMap := make(map[string]string)
			arr := strings.Split(structFiledTagName,";")
			for _,oneKeyStr:=range arr{
				//MyPrint("oneKeyStr :",oneKeyStr)
				oneKeyArr := strings.Split(oneKeyStr,":")
				MyPrint("oneKeyArr :",oneKeyArr)
				//MyPrint(ValueOfTableOneColumnOption.Elem().FieldByName(oneKeyArr[0]).Type())
				//MyPrint(oneKeyArr[1],StrFirstToUpper(oneKeyArr[0]))
				ValueOfTableOneColumnOption.Elem().FieldByName(StrFirstToUpper(oneKeyArr[0])).SetString(oneKeyArr[1])
			}
			//MyPrint(tableOneColumnOption)
		}
		columnsOption[structFiled.Name] = tableOneColumnOption
	}
}