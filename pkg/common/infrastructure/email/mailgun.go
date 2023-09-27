package email

import (
	"context"
	"ecommerce-api/internal/config"
	"ecommerce-api/internal/logger"
	"ecommerce-api/pkg/common/domain/vo"
	"encoding/json"
	"fmt"
	"github.com/mailgun/mailgun-go/v4"
	"time"
)

type SendResult struct {
	ID      string
	Message string
}

type Mailgun struct {
	Sender string
	Client *mailgun.MailgunImpl
}

func NewMailgun(cfg config.MailgunConfig) *Mailgun {
	return &Mailgun{
		Sender: cfg.Sender,
		Client: mailgun.NewMailgun(cfg.Domain, cfg.ApiKey),
	}
}

func (m *Mailgun) Send(ctx context.Context, params vo.EmailSendParams) error {
	sender := fmt.Sprintf("%s %s", params.Sender, m.Sender)

	payload := m.Client.NewMessage(
		sender,
		params.Subject,
		"",
		params.Recipient,
	)

	payload.SetHtml(params.Content)

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	message, id, err := m.Client.Send(ctxTimeout, payload)
	if err != nil {
		return err
	}

	request, _ := json.Marshal(params)
	body, _ := json.Marshal(SendResult{
		ID:      id,
		Message: message,
	})

	logger.Logger.Infow("third log",
		"vendor", "Mailgun",
		"type", "send_email",
		"request", string(request),
		"response", string(body),
	)

	return nil
}
