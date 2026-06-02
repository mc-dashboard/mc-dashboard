package minecraft

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// Unit Tests

func TestParseListResponse(t *testing.T) {
	tests := []struct {
		name        string
		response    string
		wantCount   int
		wantMax     int
		wantPlayers []string
	}{
		{
			name:        "no players",
			response:    "There are 0 of a max of 20 players online:",
			wantCount:   0,
			wantMax:     20,
			wantPlayers: nil,
		},
		{
			name:        "one player",
			response:    "There are 1 of a max of 20 players online: Notch",
			wantCount:   1,
			wantMax:     20,
			wantPlayers: []string{"Notch"},
		},
		{
			name:        "multiple players",
			response:    "There are 3 of a max of 20 players online: Steve, Alex, Herobrine",
			wantCount:   3,
			wantMax:     20,
			wantPlayers: []string{"Steve", "Alex", "Herobrine"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseListResponse(tt.response)

			require.Equal(t, tt.wantCount, got.PlayerCount)
			require.Equal(t, tt.wantMax, got.MaxPlayers)
			require.Len(t, got.Players, len(tt.wantPlayers))

			if len(tt.wantPlayers) > 0 {
				gotNames := make([]string, len(got.Players))
				for i, p := range got.Players {
					gotNames[i] = p.Name
				}
				require.Equal(t, tt.wantPlayers, gotNames)
			}
		})
	}
}

// Integration Tests

func TestRCONIntegration(t *testing.T) {
	host := os.Getenv("MINECRAFT_HOST")
	password := os.Getenv("MINECRAFT_RCON_PASSWORD")
	if host == "" || password == "" {
		t.Skip("MINECRAFT_HOST or MINECRAFT_RCON_PASSWORD not set")
	}

	port := os.Getenv("MINECRAFT_RCON_PORT")
	if port == "" {
		port = "25575"
	}

	client := NewRCONClient(host, port, password)
	require.NoError(t, client.Connect())
	defer func() {
		require.NoError(t, client.Disconnect(), "client.Disconnect()")
	}()

	status, err := client.GetServerStatus()
	require.NoError(t, err)
	require.True(t, status.Online)
	require.NotZero(t, status.MaxPlayers)

	t.Logf("Server: %d/%d players", status.PlayerCount, status.MaxPlayers)
}
