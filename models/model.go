package models

import (
	"ledis/htable"
)

const ForeignKeyFieldType string = "ForeignKeyField"
const ForeignKeyListFieldType string = "ForeignKeyListField"

// Model 模型
type Model struct {
	Name   string
	DBName string
	Fields []*Field
}

// Field 域
type Field struct {
	Name     string
	Type     string
	RelModel string
	RelField string
}

// NewModel 新建 Model
func NewModel(dbName string, tbName string, fields []*Field) *Model {
	model := &Model{
		Name:   tbName,
		DBName: dbName,
		Fields: fields,
	}
	return model
}

// NewField 新建 Field
func NewField(name string, fieldType string, relModel string, relField string) *Field {
	field := &Field{
		Name:     name,
		Type:     fieldType,
		RelModel: relModel,
		RelField: relField,
	}
	return field
}

// CreateHashTableFromModel 缓存一个对象，一般来说是 map 格式。需要将对象转成 Dict 对象
func (model *Model) ToHashTable(data map[string]interface{}) (*htable.HashTable, error) {
	ht := htable.NewHashTable()
	for key, value := range data {
		tp := "ValueType" // 这个好像不对
		if model.Fields != nil {
			for _, field := range model.Fields {
				if key == field.Name {
					tp = field.Type
					break
				}
			}
		}
		err := ht.SetByType(key, value, tp)
		if err != nil {
			return nil, err
		}
	}
	return ht, nil
}
