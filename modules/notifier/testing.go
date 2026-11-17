package notifier

import (
	"errors"
	"io"
	"log/slog"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/stretchr/testify/require"
)

const (
	testUsername = "user"
	testPassword = "secret"
)

// TestingSession accepts messages without authentication, and records whether a
// message was started.
type TestingSession struct {
	mailed *atomic.Bool
}

func (s TestingSession) Mail(string, *smtp.MailOptions) error {
	s.mailed.Store(true)
	return nil
}

func (s TestingSession) Rcpt(string, *smtp.RcptOptions) error { return nil }

func (s TestingSession) Data(
	r io.Reader,
) error {
	_, err := io.Copy(io.Discard, r)
	return err
}
func (s TestingSession) Reset()        {}
func (s TestingSession) Logout() error { return nil }

// AuthenticatedTestingSession also advertises the AUTH extension, and accepts PLAIN
// authentication with testUsername and testPassword.
type AuthenticatedTestingSession struct {
	TestingSession
}

func (s AuthenticatedTestingSession) AuthMechanisms() []string {
	return []string{sasl.Plain}
}

func (s AuthenticatedTestingSession) Auth(mech string) (sasl.Server, error) {
	return sasl.NewPlainServer(func(identity, username, password string) error {
		if username != testUsername || password != testPassword {
			return smtp.ErrAuthFailed
		}
		return nil
	}), nil
}

// StartTestingServer serves SMTP on a random local port, without TLS, and returns
// that port. Also return a boolean that indicate whether an email was sent.
// Plain authentication is allowed over 127.0.0.1 by net/smtp.
func StartTestingServer(
	t *testing.T,
	makeSessionFunc func(*atomic.Bool) smtp.Session,
) (int, *atomic.Bool) {
	t.Helper()
	mailed := &atomic.Bool{}
	s := smtp.NewServer(smtp.BackendFunc(func(*smtp.Conn) (smtp.Session, error) {
		return makeSessionFunc(mailed), nil
	}))
	s.Domain = "127.0.0.1"
	s.AllowInsecureAuth = true
	s.ReadTimeout = 5 * time.Second
	s.WriteTimeout = 5 * time.Second

	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	go func() {
		if err := s.Serve(l); err != nil && !errors.Is(err, smtp.ErrServerClosed) {
			t.Errorf("SMTP server: %v", err)
		}
	}()
	t.Cleanup(func() { s.Close() })
	return l.Addr().(*net.TCPAddr).Port, mailed
}

func NewTestingNotifier(port int, password string) SMTPPasswordResetNotifier {
	return NewSMTPPasswordResetNotifier(*slog.New(slog.DiscardHandler),
		testUsername, password, "127.0.0.1", port)
}

func MakeAuthenticatedSessionFunc(mailed *atomic.Bool) smtp.Session {
	return AuthenticatedTestingSession{TestingSession{mailed}}
}

func MakeNonAuthenticatedSessionFunc(
	mailed *atomic.Bool,
) smtp.Session {
	return TestingSession{mailed}
}
