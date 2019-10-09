package utils

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
)

// Response 返回值结构
type Response struct {
	Data         interface{}
	StatusCode   int
	CommandCode  int
	ErrorMessage string
}

// MakeResponse 格式化返回值
func MakeResponse(commandCode int, statusCode int, data interface{}, errorMessage string) []byte {
	response := Response{
		Data:         data,
		StatusCode:   statusCode,
		CommandCode:  commandCode,
		ErrorMessage: errorMessage,
	}
	result, err := json.Marshal(response)
	if err != nil {
		logrus.Errorln(err)
		return nil
	}
	return result
}

// SuccessResponse 成功时返回值
func SuccessResponse(commandCode int, data interface{}) []byte {
	return MakeResponse(commandCode, 0, data, "")
}

// BadRequestResponse 返回参数错误提示
func BadRequestResponse(commandCode int) []byte {
	return MakeResponse(commandCode, 1000, nil, "参数错误")
}
