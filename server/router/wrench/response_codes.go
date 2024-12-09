// -------------------------------------------------
// Package resp_wrench
// Author: hanzhi
// Date: 2024/12/8
// -------------------------------------------------

package wrench

func FailWithMessage(code Code, errMessage string) map[string]any {
	m := map[string]any{}
	m["code"] = code
	m["msg"] = message[code] + "." + errMessage
	m["data"] = ""
	return m
}
