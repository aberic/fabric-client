package response

import (
	"net/http"
	"strings"
)

var (
	// ExpNotExist 所请求事务不存在异常
	ExpNotExist = Exp("not exist", http.StatusOK)
)

// Exception 自定义异常实体
type Exception struct {
	Msg  string // 异常通用信息
	code int    // http 请求返回 code
}

// Exp 返回新自定义异常
func Exp(brief string, httpCode int) Exception {
	return Exception{
		Msg:  brief,
		code: httpCode}
}

// Fit 为内置异常补充异常信息内容前缀
func (exception *Exception) Fit(prefix string) *Exception {
	exception.Msg = strings.Join([]string{prefix, exception.Msg}, " ")
	return exception
}
