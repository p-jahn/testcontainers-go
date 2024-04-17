package testcontainers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/internal/core"
)

func TestProviderTypeGetProviderAutodetect(t *testing.T) {
	const dockerSocket = "unix://$XDG_RUNTIME_DIR/docker.sock"
	const podmanSocket = "unix://$XDG_RUNTIME_DIR/podman/podman.sock"
	defaultHostExtractor := hostExtractorFn
	t.Cleanup(func() {
		hostExtractorFn = defaultHostExtractor
	})

	tests := []struct {
		name         string
		providerType ProviderType
		inferredHost string
		want         string
	}{
		{
			name:         "default provider without podman.socket",
			providerType: ProviderDefault,
			inferredHost: dockerSocket,
			want:         Bridge,
		},
		{
			name:         "default provider with podman.socket",
			providerType: ProviderDefault,
			inferredHost: podmanSocket,
			want:         Podman,
		},
		{
			name:         "docker provider without podman.socket",
			providerType: ProviderDocker,
			inferredHost: dockerSocket,
			want:         Bridge,
		},
		{
			// Explicitly setting Docker provider should not be overridden by auto-detect
			name:         "docker provider with podman.socket",
			providerType: ProviderDocker,
			inferredHost: podmanSocket,
			want:         Bridge,
		},
		{
			name:         "Podman provider without podman.socket",
			providerType: ProviderPodman,
			inferredHost: dockerSocket,
			want:         Podman,
		},
		{
			name:         "Podman provider with podman.socket",
			providerType: ProviderPodman,
			inferredHost: podmanSocket,
			want:         Podman,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.providerType == ProviderPodman && core.IsWindows() {
				t.Skip("Podman provider is not implemented for Windows")
			}

			hostExtractorFn = func(_ context.Context) string {
				return tt.inferredHost
			}

			got, err := tt.providerType.GetProvider()
			require.NoError(t, err)

			provider, ok := got.(*DockerProvider)
			require.True(t, ok, "ProviderType.GetProvider() = %T, want %T", got, &DockerProvider{})
			assert.Equal(t, tt.want, provider.defaultBridgeNetworkName)
		})
	}
}
