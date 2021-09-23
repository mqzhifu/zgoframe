package util

import (
	"gorm.io/gorm"
	"strconv"
)

type DbQueryListPara struct {
	Where	string
	Fields 	string
	Order 	string
	Group   string
	Limit 	int
	Offset	int
}

type Db struct {
	Orm *gorm.DB
}

func NewDb(gorm *gorm.DB)*Db{
	db := new (Db)
	db.Orm = gorm
	return db
}
//gorm.ErrRecordNotFound
//查询
//根据主键ID查找一条记录
func (db *Db) GetRowById(model interface{} , Id int)(interface{}, error){
	//这里没必要用first|last ，因为它会排序，既然ID是主键不可能出现重复，不可能有多条记录，也就不需要排序
	//first last take 会自动给sql 加上：limit 1
	txDb := db.Orm.Take( model,Id)
	return model, txDb.Error
}
//
func  (db *Db)GetRow(model interface{} , where string)(interface{}, error){
	err := db.Orm.Where(where).First(model).Error
	return model, err
}

//根据主键ID查找一组记录
func (db *Db)GetListByIds(modelList interface{},Ids []int)(interface{}, error){
	txDb := db.Orm.Find( modelList,Ids)
	return modelList, txDb.Error
}

func  (db *Db)  GetList(modelList interface{},para DbQueryListPara )(interface{}, error){
	query := db.Orm

	if para.Where != ""{
		MyPrint("im where")
		query = query.Where(para.Where)
	}

	if para.Fields != ""{
		query = query.Select(para.Fields)
	}

	if para.Group != ""{
		query = query.Order(para.Group)
	}

	if para.Order != ""{
		query = query.Order(para.Order)
	}

	if para.Offset >= 0 {
		query = query.Offset(para.Offset)
	}

	if para.Limit  >= 0{
		query = query.Limit(para.Limit )
	}

	err := query.Find(modelList).Error
	return modelList, err
}
//
func  (db *Db)Count(model interface{} ,  where string)(int64, error){
	//if fields == ""{
	//	fields = " count ( id ) as cnt "
	//}
	var cnt int64
	//err := db.Orm.Where(where).Select(fields).Take(model).Error
	err := db.Orm.Model(model).Where(where).Count(&cnt).Error
	return cnt, err
}


//save是全保存，即使字段中有为 空串 or 0 ，且不需要加where ，会根据结构体里自带的ID字段匹配
//另外，问题：没有主键的时候是新增，有主键的时候是更新...很乱！
//func (db *Db) UpSave(model interface{})(interface{}, error){
//	return nil,nil
//}

//修改
func (db *Db) UpRowById(fields interface{} , Id int)(int64 , error){
	txDb := db.Orm.Updates(fields).Where(" id = "+ strconv.Itoa(Id))
	return txDb.RowsAffected, txDb.Error
}
//根据主键ID查找一组记录
func (db *Db)UpListByIds(fields interface{},Ids string)(int64, error){
	txDb := db.Orm.Updates(fields).Where(" id in ( " + Ids + " ) ")
	return txDb.RowsAffected, txDb.Error
}
//
func  (db *Db) UpRow(fields interface{} , where string)(int64, error){
	txDb := db.Orm.Updates(fields).Where(where).Limit(1)
	return txDb.RowsAffected, txDb.Error
}

func  (db *Db) UpList(fields interface{} , where string)(int64, error){
	txDb := db.Orm.Updates(fields).Where(where)
	return txDb.RowsAffected, txDb.Error
}
//删除
func  (db *Db) DeleteList(model interface{} , where string)(int64, error){
	txDb := db.Orm.Where(where).Delete(model)
	return txDb.RowsAffected, txDb.Error
}
//根据主键ID查找一条记录
func (db *Db) DeleteRowById(model interface{} , Id int)(int64, error){
	txDb := db.Orm.Delete( model,Id)
	return txDb.RowsAffected, txDb.Error
}
//根据主键ID查找一组记录
func (db *Db)DeleteListByIds(modelList interface{},Ids []int)(int64, error){
	txDb := db.Orm.Delete( modelList,Ids)
	return txDb.RowsAffected, txDb.Error
}
//删除一条记录
func  (db *Db) DeleteRow(model interface{} , where string)(int64, error){
	txDb := db.Orm.Where(where).Delete(model)
	return txDb.RowsAffected , txDb.Error
}
//插入一条新记录
func  (db *Db) InsertRow(fields interface{})(int64, error){
	txDb := db.Orm.Create(fields)
	return txDb.RowsAffected , txDb.Error
}
