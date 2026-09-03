package web_test

import (
	"testing"
	"world-wide-bulb/web"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetFS(t *testing.T) {
	t.Run("returns embedded build sub-filesystem", func(t *testing.T) {
		fsys, err := web.GetFS()
		require.NoError(t, err)
		assert.NotNil(t, fsys)
	})
}
