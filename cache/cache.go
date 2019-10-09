package cache

import (
	"errors"
	"github.com/mitchellh/mapstructure"
	. "ledis/common"
	"ledis/models"
	"ledis/storage"
)

//var Cache map[string]map[string][]*storage.Dict

var Models *storage.HashTable
var Cache *storage.HashTable

func init() {
	//Cache = make(map[string]map[string][]*storage.Dict)
	Cache = storage.NewHashTable()
	Models = storage.NewHashTable()
}

func FindModel(dbName string, tbName string) (*models.Model, error) {
	db, err := Models.Get(dbName)
	if err != nil {
		Logger.Error(err)
		return nil, err
	}
	var dbHt storage.HashTable
	err = mapstructure.Decode(db, &dbHt)
	if err != nil {
		Logger.Error(err)
		return nil, err
	}
	m, err := dbHt.Get(tbName)
	if err != nil {
		Logger.Error(err)
		return nil, err
	}
	var model models.Model
	err = mapstructure.Decode(m, &model)
	if err != nil {
		Logger.Error(err)
		return nil, err
	}
	return &model, nil
}

func GetDBHashTableFromCache(dbName string) (*storage.HashTable, error){
	db, err := Cache.Get(dbName)
	if err != nil {
		Logger.Error(err)
		return nil, err
	}
	var dbHT storage.HashTable
	err = mapstructure.Decode(db, &dbHT)
	if err != nil {
		Logger.Error(err)
		return nil, err
	}
	return &dbHT, nil
}

func AddModel(dbName string, tbName string, m *models.Model) error {
	Logger.Info("Begin to Add Model!")
	ht := storage.NewHashTable()
	err := ht.Set(tbName, m)
	if err != nil {
		return err
	}
	db, err := Models.Get(dbName)
	if err != nil && err.Error() == "key not found" {
		err = Models.Set(dbName, ht)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}else {
		dbHt := db.(*storage.HashTable)
		err = dbHt.Set(tbName, ht)
		if err != nil {
			return err
		}
	}
	tbHt := storage.NewHashTable()
	err = Cache.Set(dbName, tbHt)
	if err != nil {
		Logger.Error(err)
		return err
	}
	Logger.Info("Success to Add Model!")
	return nil
}

func SetCache(dbName string, tbName string, data map[string]interface{}) error {
	model, err := FindModel(dbName, tbName)
	if err != nil {
		Logger.Error(err)
		return err
	}
	tbHT, err := storage.CreateHashTableFromModel(model, data)
	if err != nil {
		Logger.Error(err)
		return err
	}
	dbHT, err := GetDBHashTableFromCache(dbName)
	if err != nil {
		Logger.Error(err)
		return err
	}
	//db, err := Cache.Get(dbName)
	//if err != nil {
	//	Logger.Error(err)
	//	return err
	//}
	//var dbHT storage.HashTable
	//err = mapstructure.Decode(db, &dbHT)
	//if err != nil {
	//	Logger.Error(err)
	//	return err
	//}
	var tbHTList []*storage.HashTable
	tbHTLs, err := dbHT.Get(tbName)
	if err != nil && err.Error() == "key not found" {

	} else if err != nil {
		Logger.Error(err)
		return err
	} else {
		tbHTList = tbHTLs.([]*storage.HashTable)
	}
	tbHTList = append(tbHTList, tbHT)
	err = dbHT.Set(tbName, tbHTList)
	if err != nil {
		Logger.Error(err)
		return err
	}
	err = Cache.Set(dbName, dbHT)
	if err != nil {
		Logger.Error(err)
		return err
	}
	Logger.Info("dbHT = ", dbHT)
	GetOneCache(dbName, tbName, "name1", "gaojian")
	return nil
}

func GetOneCache(dbName string, tbName string, key string, value interface{}) (*storage.HashTable, error) {
	//if db, found := Cache[dbName]; found {
	//	if table, found := db[tbName]; found {
	//		for _, t := range table {
	//			if t[key] == value {
	//				return t
	//			}
	//		}
	//	}
	//}
	//db, err := Cache.Get(dbName)
	//if err != nil{
	//	return nil, err
	//}
	//var dbHT storage.HashTable
	//err = mapstructure.Decode(db, &dbHT)
	//if err != nil {
	//	Logger.Error(err)
	//	return nil,err
	//}
	dbHT, err := GetDBHashTableFromCache(dbName)
	if err != nil {
		Logger.Error(err)
		return nil, err
	}
	Logger.Info("dbHT = ", dbHT)
	tbHtLs, err := dbHT.Get(tbName)
	if err != nil {
		Logger.Error(err)
		return nil, err
	}
	var tbHtList []*storage.HashTable
	//tbHtList := tbHtLs.([]*storage.HashTable)
	err = mapstructure.Decode(tbHtLs, &tbHtList)
	if err != nil {
		Logger.Error(err)
		return nil, err
	}
	for _, tbHt := range tbHtList {
		v, err := tbHt.Get(key)
		if err != nil {
			Logger.Error(err)
			continue
		}
		Logger.Info("v = ", v)
		if v == value {
			Logger.Info("return right")
			//return tbHt.Index(storage.KeyToInt(key)), nil
			return tbHt, nil
		}
	}
	return nil, errors.New("cache not found")
}

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
