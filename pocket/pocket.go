package pocket

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	. "ledis/common"
	"ledis/htable"
	"ledis/models"
	"reflect"
)

// Pocket 哆啦A梦的口袋，啥都能装
// Pocket 取出的数据是 DBData
// DBData 取出的数据是 TableData
// Record 是 table 中的一条记录
type Pocket htable.HashTable
type DBData htable.HashTable
type TBData []*Record
type Record htable.HashTable

var pocket *Pocket

//var Models *htable.HashTable
//var Cache *htable.HashTable

func init() {
	//Cache = htable.NewHashTable()
	//pocket = htable.NewHashTable()
	pocket = NewPocket()
}

func NewPocket() *Pocket {
	pocket := Pocket{
		Length:    htable.HashTableLength,
		BlockList: make([]*htable.Block, htable.HashTableLength),
	}
	return &pocket
}

func NewDBData() *DBData {
	data := DBData{
		Length:    htable.HashTableLength,
		BlockList: make([]*htable.Block, htable.HashTableLength),
	}
	return &data
}

func NewRecord(data map[string]interface{}) (*Record, error) {
	ht, err := htable.MapToHashTable(data)
	if err != nil {
		return nil, err
	}
	record := Record(*ht)
	return &record, nil
}

func (pocket *Pocket) GetDB(dbName string) (*DBData, error) {
	// 获取到的 DB 并没有与 pocket 强相关，因此修改 DB 需要重新给 pocket 赋值
	ht := (*htable.HashTable)(pocket)
	dbValue, err := ht.Get(dbName)
	if err != nil {
		Logger.Error(err)
		return nil, err
	}
	var dbHT DBData
	err = mapstructure.Decode(dbValue, &dbHT)
	if err != nil {
		Logger.Error(err)
		return nil, err
	}
	return &dbHT, nil
}

func (pocket *Pocket) SetDB(dbName string, data *DBData) error {
	ht := (*htable.HashTable)(pocket)
	err := ht.Set(dbName, data)
	if err != nil {
		return err
	}
	return nil
}

func (dbData *DBData) GetTBData(tbName string) (*TBData, error) {
	ht := (*htable.HashTable)(dbData)
	tbValue, err := ht.Get(tbName)
	if err != nil {
		Logger.Error(err)
		return nil, err
	}
	var tbHT []*Record
	err = mapstructure.Decode(tbValue, &tbHT)
	if err != nil {
		Logger.Error(err)
		return nil, err
	}
	r := TBData(tbHT)
	return &r, nil
}

func (dbData *DBData) SetTBData(tbName string, tbData *TBData) error {
	ht := (*htable.HashTable)(dbData)
	err := ht.Set(tbName, tbData)
	if err != nil {
		Logger.Error(err)
		return err
	}
	return nil
}

func (dbData *DBData) AddData(tbName string, record *Record) error {
	tbData, err := dbData.GetTBData(tbName)
	if err != nil {
		return err
	}
	r := ([]*Record)(*tbData)
	r = append(r, record)
	rd := TBData(r)
	err = dbData.SetTBData(tbName, &rd)
	if err != nil {
		return err
	}
	return nil
}

func (tbData *TBData) GetData(key string, value interface{}) (*Record, error) {
	data := ([]*Record)(*tbData)
	for _, d := range data {
		d2 := (*htable.HashTable)(d)
		v, err := d2.Get(key)
		if err != nil {
			return nil, err
		}
		Logger.Info("v = ", v, "\t value = ", value)
		//Logger.Info(reflect.TypeOf(v).String(), reflect.TypeOf(value).String())
		v1, _ := json.Marshal(v)
		v2, _ := json.Marshal(value)
		Logger.Info(v1, "\t", v2)
		if reflect.DeepEqual(v1, v2) {
			Logger.Info("1111")
			return d, nil
		} else {
			Logger.Info("2222")
		}
	}
	return nil, nil
}

func AddData(model *models.Model, data map[string]interface{}) error {
	record, err := NewRecord(data)
	if err != nil {
		Logger.Error(err)
		return err
	}
	dbData, err := pocket.GetDB(model.DBName)
	if err != nil {
		Logger.Error(err)
		return err
	}
	err = dbData.AddData(model.Name, record)
	if err != nil {
		Logger.Error(err)
		return err
	}
	err = pocket.SetDB(model.DBName, dbData)
	if err != nil {
		Logger.Error(err)
		return err
	}
	return nil
}

func GetData(model *models.Model, key string, value interface{}) (*Record, error) {
	dbData, err := pocket.GetDB(model.DBName)
	if err != nil {
		Logger.Error(err)
		return nil, err
	}
	tbData, err := dbData.GetTBData(model.Name)
	if err != nil {
		Logger.Error(err)
		return nil, err
	}
	Logger.Info("++++++++++")
	record, err := tbData.GetData(key, value)
	if err != nil {
		Logger.Error(err)
		return nil, err
	}
	return record, nil
}

//func GetDBHashTableFromPocket(dbName string) (*htable.HashTable, error) {
//	db, err := pocket.GetDB(dbName)
//	if err != nil {
//		Logger.Error(err)
//		return nil, err
//	}
//	return db, nil
//}

//func AddData(model *models.Model, data map[string]interface{}) error {
//	tbHT, err := model.ToHashTable(data)
//	if err != nil {
//		Logger.Error(err)
//		return err
//	}
//	db, err := pocket.GetDB(model.DBName)
//	if err != nil {
//		Logger.Error(err)
//		return err
//	}
//	var tbHTList []*htable.HashTable
//	tbHTLs, err := db.GetTB(model.Name)
//	if err != nil && err.Error() == "key not found" {
//
//	} else if err != nil {
//		Logger.Error(err)
//		return err
//	} else {
//		tbHTList = tbHTLs.([]*htable.HashTable)
//	}
//	tbHTList = append(tbHTList, tbHT)
//	err = dbHT.Set(model.Name, tbHTList)
//	if err != nil {
//		Logger.Error(err)
//		return err
//	}
//	err = pocket.Set(model.DBName, dbHT)
//	if err != nil {
//		Logger.Error(err)
//		return err
//	}
//	Logger.Info("dbHT = ", dbHT)
//	return nil
//}

//func GetOneCache(dbName string, tbName string, key string, value interface{}) (*htable.HashTable, error) {
//	//if db, found := Cache[dbName]; found {
//	//	if table, found := db[tbName]; found {
//	//		for _, t := range table {
//	//			if t[key] == value {
//	//				return t
//	//			}
//	//		}
//	//	}
//	//}
//	//db, err := Cache.Get(dbName)
//	//if err != nil{
//	//	return nil, err
//	//}
//	//var dbHT storage.HashTable
//	//err = mapstructure.Decode(db, &dbHT)
//	//if err != nil {
//	//	Logger.Error(err)
//	//	return nil,err
//	//}
//	dbHT, err := GetDBHashTableFromPocket(dbName)
//	if err != nil {
//		Logger.Error(err)
//		return nil, err
//	}
//	tbHtLs, err := dbHT.Get(tbName)
//	if err != nil {
//		Logger.Error(err)
//		return nil, err
//	}
//	var tbHtList []*htable.HashTable
//	//tbHtList := tbHtLs.([]*storage.HashTable)
//	err = mapstructure.Decode(tbHtLs, &tbHtList)
//	if err != nil {
//		Logger.Error(err)
//		return nil, err
//	}
//	for _, tbHt := range tbHtList {
//		v, err := tbHt.Get(key)
//		if err != nil {
//			Logger.Error(err)
//			continue
//		}
//		Logger.Info("v = ", v)
//		if v == value {
//			Logger.Info("return right")
//			//return tbHt.Index(storage.KeyToInt(key)), nil
//			return tbHt, nil
//		}
//	}
//	return nil, errors.New("cache not found")
//}

//func GetCache(dbName string, tbName string, key string, value interface{}) []map[string]interface{} {
//	var result []map[string]interface{}
//	if db, found := Cache[dbName]; found {
//		if table, found := db[tbName]; found {
//			for _, t := range table {
//				if t[key] == value {
//					result = append(result, t)
//				}
//			}
//		}
//	}
//	model := models.GetModel(dbName, tbName)
//	if model.Fields == nil {
//		return result
//	}
//	for _, field := range model.Fields {
//		if field.Type == models.ForeignKeyFieldType {
//			for _, r := range result {
//				d := GetOneCache(dbName, field.RelModel, field.RelField, r[field.Name])
//				r[field.Name] = d
//			}
//		} else if field.Type == models.ForeignKeyListFieldType {
//			for _, r := range result {
//				var rs []interface{}
//				values := r[field.Name].([]string)
//				for _, val := range values {
//					d := GetOneCache(dbName, field.RelModel, field.RelField, val)
//					rs = append(rs, d)
//				}
//				r[field.Name] = rs
//			}
//		}
//	}
//	return result
//}
