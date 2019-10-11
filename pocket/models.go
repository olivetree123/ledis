package pocket

import (
	"github.com/mitchellh/mapstructure"
	. "ledis/common"
	"ledis/htable"
	"ledis/models"
)

var allModel *htable.HashTable

func init() {
	allModel = htable.NewHashTable()
}

// AddModel 添加 model。是否可以简化？model 中带有 db 和 table 信息
func AddModel(model *models.Model) error {
	Logger.Info("Begin to Add Model!")
	ht := htable.NewHashTable()
	err := ht.Set(model.Name, model)
	if err != nil {
		return err
	}
	db, err := allModel.Get(model.DBName)
	if err != nil && err.Error() == "key not found" {
		err = allModel.Set(model.DBName, ht)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		dbHt := db.(*htable.HashTable)
		err = dbHt.Set(model.Name, ht)
		if err != nil {
			return err
		}
	}
	dbData := NewDBData()
	err = dbData.SetTBData(model.Name, nil)
	if err != nil {
		Logger.Error(err)
		return err
	}
	err = pocket.SetDB(model.DBName, dbData)
	if err != nil {
		Logger.Error(err)
		return err
	}
	Logger.Info("Success to Add Model!")
	return nil
}

func FindModel(dbName string, tbName string) (*models.Model, error) {
	db, err := allModel.Get(dbName)
	if err != nil {
		Logger.Error(err)
		return nil, err
	}
	var dbHt htable.HashTable
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
