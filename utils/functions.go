package utils

import (
	"encoding/json"
	"hash/fnv"
	. "ledis/common"
)

// HashInt 根据 hash 算法获取数值
func HashInt(s interface{}) int {
	sBytes, err := json.Marshal(s)
	if err != nil {
		Logger.Error(err)
		return -1
	}
	h := fnv.New32a()
	_, err = h.Write(sBytes)
	if err != nil {
		Logger.Error(err)
		return -1
	}
	return int(h.Sum32())
}
