package notifier

import (
	"crypto/tls"
	"fmt"
	"log/slog"
	"net"
	"net/smtp"
	"strconv"
	"time"

	g "github.com/lejeunel/go-image-annotator/globals"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

// smtpCheckTimeout bounds the whole exchange made by Check.
const smtpCheckTimeout = 10 * time.Second

type Notifier interface {
	Notify(n Notification)
	Check() error
}

type Notification struct {
	Email string
	URL   string
}

type VoidNotifier struct {
	slog.Logger
}

func (n VoidNotifier) Notify(notification Notification) {
	n.Logger.Info("notifying password reset token", "notification", notification)
}

func (n VoidNotifier) Check() error {
	return nil
}

type SMTPPasswordResetNotifier struct {
	slog.Logger
	Username string
	Password string
	Host     string
	Port     int
}

func NewSMTPPasswordResetNotifier(
	l slog.Logger,
	u, p, h string,
	port int,
) SMTPPasswordResetNotifier {
	return SMTPPasswordResetNotifier{l, u, p, h, port}
}

func (n SMTPPasswordResetNotifier) Notify(notification Notification) {
	auth := smtp.PlainAuth(
		"",
		n.Username,
		n.Password,
		n.Host,
	)

	msg := fmt.Appendf(nil, "To: %v\r\n"+
		"Subject: [%v] Password reset\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/plain; charset=UTF-8\r\n"+
		"\r\n"+
		"You have requested a password reset to %v.\r\n"+
		"Click the following link to reset your password.\r\n"+
		"%v\r\n",
		notification.Email,
		g.AppName,
		g.AppName,
		notification.URL)

	err := smtp.SendMail(
		fmt.Sprintf("%v:%v", n.Host, n.Port),
		auth,
		n.Username,
		[]string{notification.Email},
		msg,
	)
	if err != nil {
		n.Logger.Error(fmt.Errorf("sending password reset email: %w", err).Error())
	}
	n.Logger.Info("notified password reset token", "recipient", notification.Email)
}

// Check attempts to connects and authenticates to the SMTP server.
func (n SMTPPasswordResetNotifier) Check() error {
	addr := net.JoinHostPort(n.Host, strconv.Itoa(n.Port))
	conn, err := net.DialTimeout("tcp", addr, smtpCheckTimeout)
	if err != nil {
		return fmt.Errorf("connecting to SMTP server %v: %w", addr, err)
	}
	// net/smtp has no timeouts of its own: bound the whole exchange.
	if err := conn.SetDeadline(time.Now().Add(smtpCheckTimeout)); err != nil {
		conn.Close()
		return fmt.Errorf("setting SMTP connection deadline: %w", err)
	}

	c, err := smtp.NewClient(conn, n.Host)
	if err != nil {
		conn.Close()
		return fmt.Errorf("SMTP handshake with %v: %w", addr, err)
	}
	defer c.Close()

	if err := c.Hello("localhost"); err != nil {
		return fmt.Errorf("SMTP hello: %w", err)
	}
	if ok, _ := c.Extension("STARTTLS"); ok {
		if err := c.StartTLS(&tls.Config{ServerName: n.Host}); err != nil {
			return fmt.Errorf("SMTP STARTTLS: %w", err)
		}
	}
	auth := smtp.PlainAuth("", n.Username, n.Password, n.Host)
	if err := c.Auth(auth); err != nil {
		return fmt.Errorf("SMTP authentication: %w", err)
	}
	return c.Quit()
}

func NewNotifierWithVoidDefault(
	host *string,
	port int,
	username *string,
	password *string,
	l slog.Logger,
) (Notifier, error) {
	if host == nil {
		return VoidNotifier{}, fmt.Errorf(
			"validating SMTP host: found nul value: %w", e.ErrValidation)
	}

	if username == nil {
		return VoidNotifier{}, fmt.Errorf(
			"validating SMTP username: %w", e.ErrValidation)
	}
	if password == nil {
		return VoidNotifier{}, fmt.Errorf(
			"validating SMTP password: %w", e.ErrValidation)
	}
	n := NewSMTPPasswordResetNotifier(l, *username, *password, *host, port)
	return &n, nil
}
