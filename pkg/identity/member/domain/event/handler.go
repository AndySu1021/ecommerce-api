package event

import (
	"bytes"
	"context"
	"ecommerce-api/pkg/common/domain/email"
	"ecommerce-api/pkg/common/domain/vo"
	"fmt"
	"html/template"
)

type Handler struct {
	emailSvc email.Service
}

func NewHandler(emailSvc email.Service) *Handler {
	return &Handler{emailSvc: emailSvc}
}

func (h *Handler) Handle(event interface{}) error {
	switch op := event.(type) {
	case MemberForgetPassword:
		// TODO: template path should fetch from config
		t, err := template.ParseFiles("pkg/identity/member/templates/forget_password.tmpl")
		if err != nil {
			return err
		}

		var b bytes.Buffer
		if err = t.Execute(&b, op); err != nil {
			return err
		}

		if err = h.emailSvc.Send(context.Background(), vo.EmailSendParams{
			Sender:    op.MerchantName,
			Recipient: op.Email,
			Subject:   fmt.Sprintf("[%s]會員密碼確認信", op.MerchantName),
			Content:   b.String(),
		}); err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("invalid event")
}
