package models

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
func NewModel(name string, dbName string, fields []*Field) *Model {
	model := &Model{
		Name:   name,
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
