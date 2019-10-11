package htable

import (
	"encoding/json"
	"errors"
	. "ledis/common"
	"reflect"
)

// 哈希表
/*
0:—–>&{20 hello }
1:
2:—–>&{67 hello }
3:—–>&{68 hello }
4:
5:—–>&{35 hello }
6:—–>&{16 hello 0xc0420024e0}—–>&{76 hello }
*/

// HashTableLength 哈希表默认数组长度
const HashTableLength int = 100 // 哈希表默认长度

// HashTable 哈希表
// 理论上讲哈希表的 value 应该是 HashValue，即哈希表应该是 HashValue 的数组；
// 但是因为要控制内存空间的大小，所以必须将HashValue 存储在 block 中，所以哈希表的 value 是 block 的指针。
type HashTable struct {
	Length    int
	BlockList []*Block
}

//type HashTypeContent struct {
//	Index       int // HashTable 下标
//	BlockOffset int
//}

func NewHashTable() *HashTable {
	//ht := make(HashTable, HashTableLength)
	ht := HashTable{
		Length:    HashTableLength,
		BlockList: make([]*Block, HashTableLength),
	}
	return &ht
}

// Index 通过下标获取数据
func (ht *HashTable) Index(keyInt int) *Block {
	return ht.BlockList[keyInt]
}

func (ht *HashTable) Set(key string, value interface{}) error {
	tp := reflect.TypeOf(value)
	if tp.String() == "*storage.HashTable" {
		return ht.SetByType(key, value, "HashType")
	}
	return ht.SetByType(key, value, "ValueType")
}

func (ht *HashTable) SetByType(key string, value interface{}, tp string) error {
	if tp == "HashType" {
		Logger.Info("YES, HashType!")
	}
	hashValue := HashValue{
		Key:   key,
		Value: value,
		Type:  tp,
	}
	content, err := HashValueToBytes(hashValue)
	if err != nil {
		Logger.Error(err)
		return err
	}
	block := Store(content)
	if block == nil {
		return errors.New("failed to store")
	}

	keyInt := KeyToInt(key)
	//Logger.Info("key = ", key, "\t keyInt = ", keyInt, "\t value = ", value)
	val := ht.Index(keyInt)
	if val == nil {
		ht.setData(keyInt, block)
		return nil
	}
	hv, err := HashValueFromBytes(val.GetContent())
	if err != nil {
		Logger.Error(err)
		return err
	}
	if hv != nil && hv.Key == key {
		block.Next = val.Next
		ht.setData(keyInt, block)
	} else {
		for val.Next != nil {
			val = val.Next
		}
		val.Next = block
	}
	return nil
}

func (ht *HashTable) setData(keyInt int, block *Block) {
	ht.BlockList[keyInt] = block
}

func (ht *HashTable) Get(key string) (interface{}, error) {
	block := ht.Index(KeyToInt(key))
	for block != nil {
		hashValue, err := HashValueFromBytes(block.GetContent())
		if err != nil {
			return nil, err
		}
		if hashValue.Key == key {
			Logger.Info("找到了, key = ", key)
			Logger.Info(hashValue.Value, "\t", hashValue.Type)
			return hashValue.Value, nil
		}
		block = block.Next
	}
	return nil, errors.New("key not found")

}

// Keys 列出 HashTable 所有的 key
func (ht *HashTable) Keys() ([]interface{}, error) {
	var keys []interface{}
	for i := 0; i < ht.Length; i++ {
		if ht.BlockList[i] == nil {
			continue
		}
		block := ht.BlockList[i]
		for block != nil {
			content := block.GetContent()
			hv, err := HashValueFromBytes(content)
			if err != nil {
				return nil, err
			}
			keys = append(keys, hv.Key)
			block = block.Next
		}
	}
	return keys, nil
}

func (ht *HashTable) String() string {
	m, err := ht.ToMap()
	if err != nil {
		Logger.Error(err)
		return ""
	}
	b, err := json.Marshal(m)
	if err != nil {
		Logger.Error(err)
		return ""
	}
	return string(b)
}

// ToMap 将 HashTable 转成 map
func (ht *HashTable) ToMap() (map[string]interface{}, error) {
	result := make(map[string]interface{})
	if ht.BlockList == nil {
		return result, nil
	}
	for _, block := range ht.BlockList {
		for block != nil {
			content := block.GetContent()
			hv, err := HashValueFromBytes(content)
			if err != nil {
				Logger.Error(err)
				return nil, err
			}
			if hv.IsHashType() {
				ht2, err := MapToHashTable(hv.Value.(map[string]interface{}))
				if err != nil {
					Logger.Error(err)
					return nil, err
				}
				result[hv.Key.(string)], err = ht2.ToMap()
				if err != nil {
					Logger.Error(err)
					return nil, err
				}
			} else {
				result[hv.Key.(string)] = hv.Value
			}
			block = block.Next
		}
	}
	return result, nil
}

// MapToHashTable 将 map 对象装成 HashTable
func MapToHashTable(v map[string]interface{}) (*HashTable, error) {
	ht := NewHashTable()
	for key, value := range v {
		err := ht.Set(key, value)
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
