package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEndorsingRewards_Integration(t *testing.T) {
	api := New("https://staging.api.tzkt.io")

	bakings, err := api.GetBakings(t.Context(), map[string]string{
		"level": "9935996",
	})

	require.NoError(t, err)
	assert.Len(t, bakings, 1)
	assert.NotEmpty(t, bakings[0].Level)
}

func TestGetStaking_Integration(t *testing.T) {
	api := New("https://staging.api.tzkt.io")

	stakings, err := api.GetStaking(t.Context(), map[string]string{
		"hash": "opEK5fRFrjzyGcXS8Euh9EsxRAATKsEV94m2mejsgWHMHbxTAPo",
	})

	require.NoError(t, err)
	assert.Len(t, stakings, 1)
	assert.Equal(t, uint64(9967034), stakings[0].Level)
}
