package email

import (
	"context"
	"ecommerce-api/pkg/common/domain/vo"
)

type Service interface {
	Send(ctx context.Context, params vo.EmailSendParams) error
}
