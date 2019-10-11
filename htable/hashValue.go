package htable

import (
	"encoding/json"
	. "ledis/common"
)

// HashValue 哈希表 value 结构
type HashValue struct {
	Key   interface{}
	Value interface{}	// value 可以是 int/string等标准类型，也可以是 HashTable 类型
	Type  string // 类型，指明 value 是值还是 HashTable，可选值有 ValueType / HashType
}

func (hv *HashValue) IsHashType() bool {
	if hv.Type == "HashType" {
		return true
	}
	return false
}

func HashValueFromBytes(content []byte) (*HashValue, error) {
	var hashValue = HashValue{}
	err := json.Unmarshal(content, &hashValue)
	if err != nil {
		return nil, err
	}
	return &hashValue, nil
}

func HashValueToBytes(hashValue HashValue) ([]byte, error) {
	content, err := json.Marshal(hashValue)
	if err != nil {
		Logger.Info(err)
		return nil, err
	}
	return content, nil
}
