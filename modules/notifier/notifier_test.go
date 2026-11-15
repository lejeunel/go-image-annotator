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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testUsername = "user"
	testPassword = "secret"
)

// session accepts messages without authentication, and records whether a
// message was started.
type session struct {
	mailed *atomic.Bool
}

func (s session) Mail(string, *smtp.MailOptions) error {
	s.mailed.Store(true)
	return nil
}

func (s session) Rcpt(string, *smtp.RcptOptions) error { return nil }

func (s session) Data(
	r io.Reader,
) error {
	_, err := io.Copy(io.Discard, r)
	return err
}
func (s session) Reset()        {}
func (s session) Logout() error { return nil }

// authSession also advertises the AUTH extension, and accepts PLAIN
// authentication with testUsername and testPassword.
type authSession struct {
	session
}

func (s authSession) AuthMechanisms() []string {
	return []string{sasl.Plain}
}

func (s authSession) Auth(mech string) (sasl.Server, error) {
	return sasl.NewPlainServer(func(identity, username, password string) error {
		if username != testUsername || password != testPassword {
			return smtp.ErrAuthFailed
		}
		return nil
	}), nil
}

// startTestServer serves SMTP on a random local port, without TLS, and returns
// that port. Also return a boolean that indicate whether an email was sent.
// Plain authentication is allowed over 127.0.0.1 by net/smtp.
func startTestServer(t *testing.T, newSession func(*atomic.Bool) smtp.Session) (int, *atomic.Bool) {
	t.Helper()
	mailed := &atomic.Bool{}
	s := smtp.NewServer(smtp.BackendFunc(func(*smtp.Conn) (smtp.Session, error) {
		return newSession(mailed), nil
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

func newTestNotifier(port int, password string) SMTPPasswordResetNotifier {
	return NewSMTPPasswordResetNotifier(*slog.New(slog.DiscardHandler),
		testUsername, password, "127.0.0.1", port)
}

func withAuth(mailed *atomic.Bool) smtp.Session    { return authSession{session{mailed}} }
func withoutAuth(mailed *atomic.Bool) smtp.Session { return session{mailed} }

func TestCheckSucceedsWithValidCredentials(t *testing.T) {
	port, mailed := startTestServer(t, withAuth)

	assert.NoError(t, newTestNotifier(port, testPassword).Check())
	assert.False(t, mailed.Load(), "Check must not start sending a message")
}

func TestCheckFailsWithWrongPassword(t *testing.T) {
	port, _ := startTestServer(t, withAuth)

	err := newTestNotifier(port, "wrong").Check()
	assert.ErrorContains(t, err, "SMTP authentication")
}

func TestCheckFailsWhenServerDoesNotSupportAuth(t *testing.T) {
	port, _ := startTestServer(t, withoutAuth)

	err := newTestNotifier(port, testPassword).Check()
	assert.ErrorContains(t, err, "SMTP authentication")
}

func TestCheckFailsWhenNothingListens(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	port := l.Addr().(*net.TCPAddr).Port
	require.NoError(t, l.Close())

	err = newTestNotifier(port, testPassword).Check()
	assert.ErrorContains(t, err, "connecting to SMTP server")
}
