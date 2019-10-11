package handlers

import (
	"github.com/mitchellh/mapstructure"
	. "ledis/common"
	"ledis/models"
	"ledis/pocket"
	"ledis/utils"
)

type FieldParam struct {
	Name     string
	Type     string
	RelModel string
	RelField string
}

type CreateTableParam struct {
	TbName string
	DBName string
	Fields []FieldParam
}

type DeleteTableParam struct {
	Name   string
	DBName string
}

type SetCacheParam struct {
	DBName string
	TbName string
	Data   map[string]interface{}
}

type GetCacheParam struct {
	DBName string
	TbName string
	Key    string
	Value  string
}

// CreateTableHandler 创建表
func CreateTableHandler(command Command) []byte {
	var param CreateTableParam
	err := mapstructure.Decode(command.Args, &param)
	if err != nil {
		Logger.Errorln(err)
		return ErrorResponse(command.Code, utils.BadRequest)
	}
	var fields []*models.Field
	for _, f := range param.Fields {
		field := models.NewField(f.Name, f.Type, f.RelModel, f.RelField)
		fields = append(fields, field)
	}
	md := models.NewModel(param.DBName, param.TbName, fields)
	err = pocket.AddModel(md)
	if err != nil {
		Logger.Errorln(err)
		return ErrorResponse(command.Code, utils.BadRequest)
	}
	return APIResponse(command.Code, md)
}

// DeleteTableHandler 删除表
func DeleteTableHandler(command Command) []byte {
	var param DeleteTableParam
	err := mapstructure.Decode(command.Args, &param)
	if err != nil {
		Logger.Errorln(err)
		return ErrorResponse(command.Code, utils.BadRequest)
	}
	//status := models.DeleteModel(param.DBName, param.Name)
	return APIResponse(command.Code, 0)
}

// SetCacheHandler 设置缓存
func SetCacheHandler(command Command) []byte {
	var param SetCacheParam
	err := mapstructure.Decode(command.Args, &param)
	if err != nil {
		Logger.Errorln(err)
		return ErrorResponse(command.Code, utils.BadRequest)
	}
	model, err := pocket.FindModel(param.DBName, param.TbName)
	if err != nil {
		Logger.Errorln(err)
		return ErrorResponse(command.Code, utils.BadRequest)
	}
	err = pocket.AddData(model, param.Data)
	if err != nil {
		Logger.Errorln(err)
		return ErrorResponse(command.Code, utils.BadRequest)
	}
	return APIResponse(command.Code, true)
}

// GetCacheHandler 获取缓存
func GetCacheHandler(command Command) []byte {
	var param GetCacheParam
	err := mapstructure.Decode(command.Args, &param)
	if err != nil {
		Logger.Errorln(err)
		return ErrorResponse(command.Code, utils.BadRequest)
	}
	model, err := pocket.FindModel(param.DBName, param.TbName)
	if err != nil {
		Logger.Errorln(err)
		return ErrorResponse(command.Code, utils.BadRequest)
	}
	result, err := pocket.GetData(model, param.Key, param.Value)
	if err != nil {
		Logger.Errorln(err)
		return ErrorResponse(command.Code, utils.BadRequest)
	}
	return APIResponse(command.Code, result)
}
