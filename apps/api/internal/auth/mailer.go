package auth

import (
	"context"
	"log/slog"
)

// Mailer sends transactional auth emails. It is an interface so a real SMTP / provider
// implementation can be wired in later without touching the service; that wiring is
// infra/human-gated (external network) and out of scope here.
type Mailer interface {
	// SendVerificationEmail delivers an email-verification link to the address.
	SendVerificationEmail(ctx context.Context, to string, verifyURL string) error
}

// LogMailer is the default development Mailer: instead of sending email (no provider is
// configured), it logs the verification link so a developer can complete the flow
// locally. It never fails.
type LogMailer struct {
	logger *slog.Logger
}

func NewLogMailer(logger *slog.Logger) *LogMailer {
	return &LogMailer{logger: logger}
}

func (m *LogMailer) SendVerificationEmail(_ context.Context, to string, verifyURL string) error {
	m.logger.Info("email verification link (dev log mailer; no email sent)",
		"to", to, "verify_url", verifyURL)
	return nil
}
