package notifier

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckSucceedsWithValidCredentials(t *testing.T) {
	port, mailed := StartTestingServer(t, MakeAuthenticatedSessionFunc)

	assert.NoError(t, NewTestingNotifier(port, testPassword).Check())
	assert.False(t, mailed.Load(), "Check must not start sending a message")
}

func TestCheckFailsWithWrongPassword(t *testing.T) {
	port, _ := StartTestingServer(t, MakeAuthenticatedSessionFunc)

	err := NewTestingNotifier(port, "wrong").Check()
	assert.ErrorContains(t, err, "SMTP authentication")
}

func TestCheckFailsWhenServerDoesNotSupportAuth(t *testing.T) {
	port, _ := StartTestingServer(t, MakeNonAuthenticatedSessionFunc)

	err := NewTestingNotifier(port, testPassword).Check()
	assert.ErrorContains(t, err, "SMTP authentication")
}

func TestCheckFailsWhenNothingListens(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	port := l.Addr().(*net.TCPAddr).Port
	require.NoError(t, l.Close())

	err = NewTestingNotifier(port, testPassword).Check()
	assert.ErrorContains(t, err, "connecting to SMTP server")
}
