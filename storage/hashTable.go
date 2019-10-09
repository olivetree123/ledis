package storage

import (
	"encoding/json"
	"errors"
	"github.com/mitchellh/mapstructure"
	. "ledis/common"
	"ledis/models"
	"ledis/utils"
)

//type HashTable []*Block // 哈希表

type HashTable struct {
	Content []*Block
}

var HashTableLength int // 哈希表长度

type HashValue struct {
	Key   interface{}
	Value interface{}
	Type  string // 类型，指明 value 是值还是另一个 hashValue，可选值有 ValueType / HashType
}

type HashTypeContent struct {
	Index       int // HashTable 下标
	BlockOffset int
}

func init() {
	HashTableLength = 100
	//HashTable = make([]*Block, HashTableLength)
}

// KeyToInt 将 Key 转换为数组下标
func KeyToInt(key interface{}) int {
	return utils.HashInt(key) % HashTableLength
}

func NewHashTable() *HashTable {
	//ht := make(HashTable, HashTableLength)
	ht := HashTable{
		Content: make([]*Block, HashTableLength),
	}
	return &ht
}

func (ht *HashTable) Index(keyInt int) *Block {
	return ht.Content[keyInt]
	//bl := []*Block(*ht)
	//if keyInt >= len(bl) {
	//	return nil
	//}
	//return bl[keyInt]
}

func (ht *HashTable) Set(key string, value interface{}) error {
	return ht.setByType(key, value, "ValueType")
}

func (ht *HashTable) setByType(key string, value interface{}, tp string) error {
	hashValue := HashValue{
		Key:   key,
		Value: value,
		Type:  tp,
	}
	content, err := json.Marshal(hashValue)
	if err != nil {
		Logger.Info(err)
		return err
	}
	block := Store(content)
	if block == nil {
		return errors.New("failed to store")
	}
	keyInt := KeyToInt(key)
	Logger.Info("keyInt = ", keyInt, key)
	//val := ht.Index(keyInt)
	//if val == nil {
	//	Logger.Info("set data ++++")
	//	ht.setData(keyInt, block)
	//} else {
	//	for val.Next != nil {
	//		val = val.Next
	//	}
	//	val.Next = block
	//}
	ht.setData(keyInt, block)
	return nil
}

func (ht *HashTable) setData(keyInt int, block *Block) {
	ht.Content[keyInt] = block
}

func (ht *HashTable) Get(key string) (interface{}, error) {
	val := ht.Index(KeyToInt(key))
	for val != nil {
		content := Content[val.Offset : val.Offset+val.Length]
		var hashValue = HashValue{}
		err := json.Unmarshal(content, &hashValue)
		if err != nil {
			return nil, err
		}
		if hashValue.Key == key {
			value := hashValue.Value
			for hashValue.Type == "HashType" {
				hc := HashTypeContent{}
				err := mapstructure.Decode(value, &hc)
				if err != nil {
					Logger.Error(err)
					return nil, err
				}
				b := ht.Index(hc.Index)
				for b.Offset != hc.BlockOffset {
					if b.Next == nil {
						Logger.Error("Error, not found.")
						return nil, err
					}
					b = b.Next
				}
				//Content[b.Offset : b.Offset+b.Length]
				// 递归？
			}
			return hashValue.Value, nil
		}
		val = val.Next
	}
	return nil, errors.New("key not found")
}

func (ht *HashTable) Print() {
	if ht.Content == nil {
		Logger.Error("ht.Content is nil")
		return
	}
	for _, block := range ht.Content {
		if block == nil {
			continue
		}
		Logger.Info("block = ", block)
		content := GetStore(block.Offset, block.Length)
		var hv HashValue
		err:= json.Unmarshal(content, &hv)
		if err != nil {
			Logger.Error(err)
			return
		}
		Logger.Info("key = ", hv.Key, "	value = ", hv.Value)
	}
}

// CreateHashTableFromModel 缓存一个对象，一般来说是 map 格式。需要将对象转成 Dict 对象
func CreateHashTableFromModel(model *models.Model, data map[string]interface{}) (*HashTable, error) {
	ht := NewHashTable()
	Logger.Info("model = ", model)
	for key, value := range data {
		tp := "ValueType"
		if model.Fields != nil {
			for _, field := range model.Fields {
				if key == field.Name {
					tp = field.Type
					break
				}
			}
		}
		err := ht.setByType(key, value, tp)
		if err != nil {
			return nil, err
		}
	}
	return ht, nil
}

// AddHash 向 hashTable 添加数据，一般来说是是 map 中的某一个 k/v
//func AddHash(key interface{}, value interface{}, tp string) (int, error) {
//	hashValue := HashValue{
//		Key:   key,
//		Value: value,
//		Type:  tp,
//	}
//	content, err := json.Marshal(hashValue)
//	if err != nil {
//		Logger.Info(err)
//		return 0, err
//	}
//	block := Store(content)
//	keyInt := KeyToInt(key)
//	val := HashTable[keyInt]
//	if val == nil {
//		HashTable[keyInt] = block
//	} else {
//		for val.Next != nil {
//			val = val.Next
//		}
//		val.Next = block
//	}
//	return keyInt, nil
//}

// GetHash 获取数据
//func GetHash(key string) interface{} {
//	keyInt := KeyToInt(key)
//	val := HashTable[keyInt]
//	for val != nil {
//		content := Content[val.Offset : val.Offset+val.Length]
//		var hashValue = HashValue{}
//		err := json.Unmarshal(content, &hashValue)
//		if err != nil {
//			Logger.Error(err)
//			return nil
//		}
//		if hashValue.Key == key {
//			value := hashValue.Value
//			for hashValue.Type == "HashType" {
//				hc := HashTypeContent{}
//				err := mapstructure.Decode(value, &hc)
//				if err != nil {
//					Logger.Error(err)
//					return nil
//				}
//				b := HashTable[hc.Index]
//				for b.Offset != hc.BlockOffset {
//					if b.Next == nil {
//						Logger.Error("Error, not found.")
//						return nil
//					}
//					b = b.Next
//				}
//				//Content[b.Offset : b.Offset+b.Length]
//				// 递归？
//			}
//			return hashValue.Value
//		}
//		val = val.Next
//	}
//	return nil
//}
