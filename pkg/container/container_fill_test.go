package container

import (
	"bytes"
	"testing"

	"github.com/4strodev/wiring_graphs/pkg/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFill(t *testing.T) {
	cont := New()

	dependantResolver := func(s testutils.MyService) *bytes.Buffer {
		return bytes.NewBuffer([]byte{})
	}

	err := cont.Dependencies(dependantResolver, testutils.NewService)
	cont.Token(map[string]any{
		"buffer": func() *bytes.Buffer {
			return bytes.NewBuffer([]byte{})
		}},
	)
	require.NoError(t, err)

	var deps testutils.MyDeps
	err = cont.Fill(&deps)
	require.NoError(t, err, "struct should be filled correctly")
	require.NoError(t, deps.CheckResolvedDependencies())
}

func TestGetToken(t *testing.T) {
	tests := []struct {
		tag   string
		token string
	}{
		{
			tag:   ",omit",
			token: "",
		},
		{
			tag:   "injected,omit",
			token: "",
		},
		{
			tag:   "injected,",
			token: "injected",
		},
		{
			tag:   "injected,nonuseful",
			token: "injected",
		},
	}

	for _, ttest := range tests {
		t.Run(ttest.tag, func(t *testing.T) {
			result := getToken(ttest.tag)
			assert.Equalf(t,
				ttest.token,
				result,
				"tag '%s' is not being processed correctly",
				ttest.tag)
		})
	}
}
