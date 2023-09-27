package errors

import (
	"strconv"
	"strings"
)

type _error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
}

// 1:系統相關錯誤 2:訂單相關錯誤

var (
	ErrAuth            = &_error{Code: 1001, Message: "驗證失敗", Details: ""}
	ErrPerm            = &_error{Code: 1002, Message: "權限不足", Details: ""}
	ErrValidation      = &_error{Code: 1003, Message: "參數驗證失敗", Details: ""}
	ErrOp              = &_error{Code: 1004, Message: "操作失敗", Details: ""}
	ErrEmptyDevice     = &_error{Code: 1005, Message: "參數錯誤", Details: ""}
	ErrEmptyOrderNo    = &_error{Code: 1006, Message: "參數錯誤", Details: ""}
	ErrNotFound        = &_error{Code: 1007, Message: "找不到資源", Details: ""}
	ErrWrongTimeFormat = &_error{Code: 1008, Message: "時間格式錯誤", Details: ""}

	ErrStockNotEnough = &_error{Code: 2001, Message: "庫存不足", Details: ""}
	ErrEmptyCart      = &_error{Code: 2002, Message: "購物車不可為空", Details: ""}
	ErrOrderDelivered = &_error{Code: 2003, Message: "取消失敗，訂單已出貨", Details: ""}
)

func NewErrValidation(err error) *_error {
	return &_error{
		Code:    3,
		Message: "參數驗證失敗",
		Details: err.Error(),
	}
}

func (e *_error) Error() string {
	var b strings.Builder
	_, _ = b.WriteRune('[')
	_, _ = b.WriteString(strconv.Itoa(e.Code))
	_, _ = b.WriteRune(']')
	_, _ = b.WriteString(e.Message)
	if e.Details != "" {
		_, _ = b.WriteRune('，')
		_, _ = b.WriteString(e.Details)
	}
	return b.String()
}
