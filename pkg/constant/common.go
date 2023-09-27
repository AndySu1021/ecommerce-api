package constant

// YesNo 是否
type YesNo int8

const (
	All YesNo = -1
	No  YesNo = 0
	Yes YesNo = 1
)

type Pagination struct {
	Page     int32 `json:"page" form:"page" binding:"required,gte=1"`
	PageSize int32 `json:"page_size" form:"page_size" binding:"required,gte=1"`
	Total    int64 `json:"total" form:"total"`
}

type HeaderKey string

const (
	HeaderKeyHost        HeaderKey = "X-Host"
	HeaderKeyMemberToken HeaderKey = "X-Token"
	HeaderKeyAdminToken  HeaderKey = "X-Admin-Token"
	HeaderKeyTraceID     HeaderKey = "X-Trace-ID"
)

func (h HeaderKey) String() string {
	return string(h)
}

type ContextKey string

const (
	ContextKeyHost       ContextKey = "host"
	ContextKeyToken      ContextKey = "token"
	ContextKeyAdminToken ContextKey = "admin_token"
	ContextKeyTraceID    ContextKey = "trace_id"
)

type Currency int8

const (
	CurrencyTWD Currency = iota + 1
)
