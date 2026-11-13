package file_store

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalExists(t *testing.T) {
	dir := t.TempDir()
	store, err := NewLocalFileStore(dir)
	require.NoError(t, err)
	require.NoError(t, store.Store("present.jpg", strings.NewReader("data")))
	require.NoError(t, os.Mkdir(filepath.Join(dir, "a-directory.jpg"), 0o755))

	for name, want := range map[string]bool{
		"present.jpg":     true,
		"missing.jpg":     false,
		"a-directory.jpg": false,
	} {
		got, err := store.Exists(name)
		assert.NoError(t, err, name)
		assert.Equal(t, want, got, name)
	}
}
