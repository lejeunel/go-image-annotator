package auth

import (
	"fmt"
	g "github.com/lejeunel/go-image-annotator/globals"
	"log/slog"
	"net/smtp"
	"os"
	"strconv"
)

type Notifier interface {
	Notify(n Notification)
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

type SMTPPasswordResetNotifier struct {
	slog.Logger
	Username string
	Password string
	Host     string
	Port     int
}

func NewSMTPPasswordResetNotifier(l slog.Logger, u, p, h string, port int) SMTPPasswordResetNotifier {
	return SMTPPasswordResetNotifier{l, u, p, h, port}
}

func (n SMTPPasswordResetNotifier) Notify(notification Notification) {
	auth := smtp.PlainAuth(
		"",
		n.Username,
		n.Password,
		n.Host,
	)

	msg := []byte(
		fmt.Sprintf("To: %v\r\n"+
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
			notification.URL))

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

func MakeNotifierFromEnv(l slog.Logger) Notifier {
	username := os.Getenv("GOIA_SMTPUSERNAME")
	password := os.Getenv("GOIA_SMTPPASSWORD")
	host := os.Getenv("GOIA_SMTPHOST")
	if host == "" {
		l.Warn("GOIA_SMTPHOST environment variable is not set. Users will not receive notification!")
		return VoidNotifier{l}
	}
	port, err := strconv.Atoi(os.Getenv("GOIA_SMTPPORT"))
	if err != nil {
		panic(fmt.Errorf("extracting SMTP port from environment variable GOIA_SMTPPORT: %w", err))
	}
	if username == "" {
		panic("when using SMTP notification, username must be defined, but found none")
	}
	return NewSMTPPasswordResetNotifier(l, username, password, host, port)
}
