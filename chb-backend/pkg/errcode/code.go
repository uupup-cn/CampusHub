package errcode

import "net/http"

type ErrorCode struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	HTTP    int    `json:"-"`
}

func (e ErrorCode) Error() string {
	return e.Message
}

var Success = ErrorCode{Code: 0, Message: "success", HTTP: http.StatusOK}

var ErrParamMissing = ErrorCode{Code: 1001, Message: "参数缺失", HTTP: http.StatusBadRequest}
var ErrParamInvalid = ErrorCode{Code: 1002, Message: "参数格式错误", HTTP: http.StatusBadRequest}
var ErrNotFound = ErrorCode{Code: 1003, Message: "资源不存在", HTTP: http.StatusNotFound}

var ErrUnauthorized = ErrorCode{Code: 2001, Message: "未登录", HTTP: http.StatusUnauthorized}
var ErrPermissionDenied = ErrorCode{Code: 2002, Message: "权限不足", HTTP: http.StatusForbidden}
var ErrTokenInvalid = ErrorCode{Code: 2003, Message: "Token 无效", HTTP: http.StatusUnauthorized}
var ErrTokenExpired = ErrorCode{Code: 2004, Message: "Token 过期", HTTP: http.StatusUnauthorized}
var ErrTrustLevelLow = ErrorCode{Code: 2005, Message: "等级不足", HTTP: http.StatusForbidden}
var ErrScopeInsufficient = ErrorCode{Code: 2006, Message: "Scope 不足", HTTP: http.StatusForbidden}

var ErrBalanceInsufficient = ErrorCode{Code: 3001, Message: "余额不足", HTTP: http.StatusBadRequest}
var ErrPoolInsufficient = ErrorCode{Code: 3002, Message: "积分池余额不足", HTTP: http.StatusInternalServerError}
var ErrDailyCapReached = ErrorCode{Code: 3003, Message: "超出日上限", HTTP: http.StatusBadRequest}
var ErrCooldown = ErrorCode{Code: 3004, Message: "冷却中", HTTP: http.StatusTooManyRequests}
var ErrDuplicateRequest = ErrorCode{Code: 3005, Message: "重复请求", HTTP: http.StatusOK}
var ErrAccountFrozen = ErrorCode{Code: 3006, Message: "账户已冻结", HTTP: http.StatusForbidden}

var ErrItemSoldOut = ErrorCode{Code: 4001, Message: "商品已售罄", HTTP: http.StatusBadRequest}
var ErrItemPending = ErrorCode{Code: 4002, Message: "商品待审核", HTTP: http.StatusBadRequest}
var ErrMerchantAppExists = ErrorCode{Code: 4003, Message: "入驻申请已存在", HTTP: http.StatusBadRequest}
var ErrNotMerchant = ErrorCode{Code: 4004, Message: "非商家", HTTP: http.StatusForbidden}

var ErrInternal = ErrorCode{Code: 5001, Message: "系统内部错误", HTTP: http.StatusInternalServerError}
var ErrDatabase = ErrorCode{Code: 5002, Message: "数据库错误", HTTP: http.StatusInternalServerError}
var ErrServiceUnavailable = ErrorCode{Code: 5003, Message: "服务不可用", HTTP: http.StatusServiceUnavailable}
