package handlers

import (
	"ledis/utils"
)

func APIResponse(commandCode int, data interface{}) []byte {
	return utils.SuccessResponse(commandCode, data)
}

func BadRequestResponse(commandCode int) []byte {
	return utils.BadRequestResponse(commandCode)
}

func ErrorResponse(commandCode int, statusCode int) []byte {
	return utils.MakeResponse(commandCode, statusCode, nil, utils.Message[statusCode])
}
