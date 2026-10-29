package cli

import (
	"errors"
	"strings"

	dootask "github.com/dootask/tools/server/go"
)

// 错误码固定清单（--json 输出的 code 字段）。新增须同步 doo 规范文档。
const (
	CodeInvalidArgs      = "invalid_args"      // 参数/命令不合法
	CodeAuthRequired     = "auth_required"     // 未登录或 token 失效
	CodePermissionDenied = "permission_denied" // 无权限
	CodeNotFound         = "not_found"         // 对象不存在或不可见
	CodeNeedsConfirm     = "needs_confirm"     // 危险操作需确认（见 confirm 字段）
	CodeAPIError         = "api_error"         // 主程序业务错误（message 为原因）
	CodeServerError      = "server_error"      // 主程序 HTTP 非 200
	CodeNetworkError     = "network_error"     // 网络不可达
	CodeTimeout          = "timeout"           // 请求超时
	CodeCancelled        = "cancelled"         // 交互确认时用户取消
	CodeError            = "error"             // 其他本地错误
)

// Error 带错误码的错误。
type Error struct {
	Code    string
	Message string
	Confirm *ConfirmInfo // 仅 needs_confirm
}

func (e *Error) Error() string { return e.Message }

// Errorf 构造带码错误。
func Errorf(code, msg string) *Error { return &Error{Code: code, Message: msg} }

// Classify 把任意错误归类为错误码。
func Classify(err error) *Error {
	var e *Error
	if errors.As(err, &e) {
		return e
	}
	if errors.Is(err, ErrNoAuth) {
		return Errorf(CodeAuthRequired, err.Error())
	}
	var apiErr *dootask.APIError
	if errors.As(err, &apiErr) {
		return Errorf(apiCode(apiErr.Ret, apiErr.Msg), err.Error())
	}
	var httpErr *dootask.HTTPError
	if errors.As(err, &httpErr) {
		switch httpErr.StatusCode {
		case 401:
			return Errorf(CodeAuthRequired, err.Error())
		case 403:
			return Errorf(CodePermissionDenied, err.Error())
		case 404:
			return Errorf(CodeNotFound, err.Error())
		}
		return Errorf(CodeServerError, err.Error())
	}
	var tErr *dootask.TransportError
	if errors.As(err, &tErr) {
		if tErr.Timeout() {
			return Errorf(CodeTimeout, err.Error())
		}
		return Errorf(CodeNetworkError, err.Error())
	}
	return Errorf(CodeError, err.Error())
}

// apiCode 按主程序 ret 码与文案归类业务错误。
// ret：-1 未登录/身份失效；-4001 项目、-4002 任务、-4003 对话 不存在或不在成员内。
func apiCode(ret int, msg string) string {
	switch {
	case ret == -1:
		return CodeAuthRequired
	case ret <= -4001 && ret >= -4003:
		return CodeNotFound
	case strings.Contains(msg, "权限") || strings.Contains(msg, "无权"):
		return CodePermissionDenied
	case strings.Contains(msg, "不存在"):
		return CodeNotFound
	}
	return CodeAPIError
}

// exitCodes 进程退出码：0 成功；1 一般错误；2 参数错误；3 未登录；4 需确认。
func exitCodeOf(code string) int {
	switch code {
	case CodeInvalidArgs:
		return 2
	case CodeAuthRequired:
		return 3
	case CodeNeedsConfirm:
		return 4
	}
	return 1
}
