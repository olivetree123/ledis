package utils

const BadRequest int = 1000
const SystemError int = 1001
const ObjectNotFound int = 1002

var Message = map[int]string{
	BadRequest:     "参数错误",
	SystemError:    "系统异常",
	ObjectNotFound: "找不到对象",
}
