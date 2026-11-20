package agent

import (
	"testing"

	"charm.land/catwalk/v3/pkg/catwalk"
	"github.com/charmbracelet/crush/internal/config"
	"github.com/stretchr/testify/require"
)

func TestServingProvider(t *testing.T) {
	t.Parallel()

	providerCfg := func(id string, providerType catwalk.Type) config.ProviderConfig {
		return config.ProviderConfig{ID: id, Type: providerType}
	}
	model := func(modelType catwalk.Type) catwalk.Model {
		return catwalk.Model{ID: "some-model", Type: modelType}
	}

	t.Run("retained provider adapters", func(t *testing.T) {
		t.Parallel()
		for providerID, expected := range map[string]string{
			string(catwalk.InferenceProviderGemini): "google",
			"hyper":                               "hyper",
		} {
			require.Equal(t, expected, servingProvider(providerCfg(providerID, catwalk.TypeCompletions), model("")), providerID)
		}
	})

	t.Run("wire types", func(t *testing.T) {
		t.Parallel()
		for wireType, expected := range map[catwalk.Type]string{
			catwalk.TypeResponses:   "openai",
			catwalk.TypeMessages:    "",
			catwalk.TypeCompletions: "openai-compat",
		} {
			require.Equal(t, expected, servingProvider(providerCfg("custom", ""), model(wireType)), wireType)
			configuredTypeExpected := expected
			if wireType == catwalk.TypeMessages {
				configuredTypeExpected = string(wireType)
			}
			require.Equal(t, configuredTypeExpected, servingProvider(providerCfg("custom", wireType), model("")), wireType)
		}
	})

	t.Run("custom provider types are passed through", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, "litellm", servingProvider(providerCfg("litellm", "litellm"), model("")))
	})
}

func TestBuildProviderOpenCodeRouting(t *testing.T) {
	t.Parallel()

	for _, providerID := range []string{
		string(catwalk.InferenceProviderOpenCodeZen),
		string(catwalk.InferenceProviderOpenCodeGo),
	} {
		t.Run(providerID, func(t *testing.T) {
			t.Parallel()
			env := testEnv(t)
			providerCfg := config.ProviderConfig{
				ID:      providerID,
				BaseURL: "https://opencode.ai/zen/v1",
				Type:    catwalk.TypeCompletions,
				APIKey:  "$OPENCODE_API_KEY",
			}
			coord := newTestCoordinator(t, env, providerID, providerCfg)

			for _, wireType := range []catwalk.Type{catwalk.TypeCompletions, catwalk.TypeResponses} {
				provider, err := coord.buildProvider(providerCfg, config.SelectedModel{
					Model:    "some-model",
					Provider: providerID,
				}, false, wireType)
				require.NoError(t, err, wireType)
				require.NotNil(t, provider, wireType)
			}

			_, err := coord.buildProvider(providerCfg, config.SelectedModel{
				Model:    "messages-model",
				Provider: providerID,
			}, false, catwalk.TypeMessages)
			require.Error(t, err)
		})
	}
}
