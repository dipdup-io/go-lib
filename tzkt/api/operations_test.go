package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEndorsingRewards_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	api := New("https://staging.api.tzkt.io")

	bakings, err := api.GetBakings(t.Context(), map[string]string{
		"level": "9935996",
	})

	require.NoError(t, err)
	assert.Len(t, bakings, 1)
	assert.NotEmpty(t, bakings[0].Level)
}
