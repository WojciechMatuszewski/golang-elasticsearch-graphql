package testing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFullPath(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		fPath := "/serverless.yml"
		fullPath, err := GetFullPath(fPath)

		assert.NoError(t, err)
		assert.Contains(t, fullPath, fPath)
	})

}
