// -------------------------------------------------
// Package dto
// Author: hanzhi
// Date: 2024/12/9
// -------------------------------------------------

package dto

type Code int64

const SuccessCode Code = 0

// 接口影响正常时code为0，错误时code为5位数
// 第1位,系统错误类型。
//		5:系统内部错误
//		4:业务请求错误
// 第2~3位,错误模块。
//		比如,00公共 01 用户
// 第4~5位,具体错误。
// 例子:
// 40000:请求参数错误
// 40100:用户已存在
// 50000:系统内部错误
// 50100:用户注册错误

// 业务错误 4xxxx
const (
	// 公共错误码 00

	ParamsErrCode         Code = 40000
	RecordNotFoundErrCode Code = 40001

	// 用户业务错误码 01
	UserExistsErrCode      Code = 40100
	UserTokenErrCode       Code = 40101
	UserPasswordErrCode    Code = 40102
	UserEmailExistsErrCode Code = 40103

	// 知识库业务错误码 02
	IndexExistErrCode Code = 40200
)

// 系统错误 5xxxx
const (
	InternalErrCode Code = 50000
)

var message map[Code]string

func init() {
	message = map[Code]string{}
	message[SuccessCode] = "Success"
	// 400xx错误message
	message[ParamsErrCode] = "参数错误"
	message[RecordNotFoundErrCode] = "记录不存在"

	// 401xx错误message
	message[UserExistsErrCode] = "用户已经存在"
	message[UserTokenErrCode] = "登录信息错误"
	message[UserPasswordErrCode] = "密码错误"
	message[UserEmailExistsErrCode] = "邮箱已经存在"
	// 402xx 知识库错误
	message[IndexExistErrCode] = "当前知识库已经存在"
	// 5xxxx错误message
	message[InternalErrCode] = "系统内部发生错误"

}

type BaseResponse struct {
	Code Code
	Msg  string
	Data interface{}
}

func SuccessWithData(data interface{}) BaseResponse {
	return BaseResponse{
		Code: SuccessCode,
		Msg:  message[SuccessCode],
		Data: data,
	}

}

func Success() BaseResponse {
	return BaseResponse{
		Code: SuccessCode,
		Msg:  message[SuccessCode],
		Data: "",
	}
}

func Fail(failCode Code) BaseResponse {
	return BaseResponse{
		Code: failCode,
		Msg:  message[failCode],
		Data: "",
	}
}

func FailWithMessage(failCode Code, errMessage string) BaseResponse {
	return BaseResponse{
		Code: failCode,
		Msg:  message[failCode] + ":" + errMessage,
		Data: "",
	}

}

func FailWithData(failCode Code, data interface{}) BaseResponse {
	return BaseResponse{
		Code: failCode,
		Msg:  message[failCode],
		Data: data,
	}
}
