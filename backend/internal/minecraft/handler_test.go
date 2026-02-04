package minecraft

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// Tests for RCON Availability

func TestHandlers_RejectWhenRCONUnavailable(t *testing.T) {
	handler := &MinecraftHandler{RCONClient: nil}

	tests := []struct {
		name    string
		method  string
		path    string
		handler func(http.ResponseWriter, *http.Request)
	}{
		{"GetServerStatus", http.MethodGet, "/api/minecraft/status", handler.GetServerStatus},
		{"GetOnlinePlayers", http.MethodGet, "/api/minecraft/players", handler.GetOnlinePlayers},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			tt.handler(w, req)

			require.Equal(t, http.StatusServiceUnavailable, w.Code)
			require.Contains(t, w.Body.String(), "RCON client not initialized")
		})
	}
}

// Tests for Command Validation

func TestExecuteCommand_RejectsEmptyCommand(t *testing.T) {
	handler := &MinecraftHandler{}

	tests := []struct {
		name string
		cmd  string
	}{
		{"empty string", ""},
		{"whitespace only", "   "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := bytes.NewBufferString(`{"command": "` + tt.cmd + `"}`)
			req := httptest.NewRequest(http.MethodPost, "/api/minecraft/command", body)
			w := httptest.NewRecorder()

			handler.ExecuteCommand(w, req)

			require.Equal(t, http.StatusBadRequest, w.Code)
			require.Contains(t, w.Body.String(), "Command cannot be empty")
		})
	}
}

func TestExecuteCommand_BlocksDangerousCommands(t *testing.T) {
	handler := &MinecraftHandler{}

	tests := []struct {
		name    string
		command string
	}{
		{"blocks op", "op PlayerName"},
		{"blocks ban", "ban PlayerName"},
		{"blocks stop", "stop"},
		{"blocks kick", "kick PlayerName"},
		{"blocks unknown", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := bytes.NewBufferString(`{"command": "` + tt.command + `"}`)
			req := httptest.NewRequest(http.MethodPost, "/api/minecraft/command", body)
			w := httptest.NewRecorder()

			handler.ExecuteCommand(w, req)

			require.Equal(t, http.StatusForbidden, w.Code)
			require.Contains(t, w.Body.String(), "Command not allowed")
		})
	}
}

func TestExecuteCommand_AllowsReadOnlyCommands(t *testing.T) {
	handler := &MinecraftHandler{}

	tests := []struct {
		name    string
		command string
	}{
		{"allows list", "list"},
		{"allows seed", "seed"},
		{"allows time query", "time query daytime"},
		{"allows weather query", "weather query"},
		{"allows gamerule", "gamerule keepInventory"},
		{"allows difficulty", "difficulty"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := bytes.NewBufferString(`{"command": "` + tt.command + `"}`)
			req := httptest.NewRequest(http.MethodPost, "/api/minecraft/command", body)
			w := httptest.NewRecorder()

			handler.ExecuteCommand(w, req)

			// Commands pass validation, but fail because RCON is unavailable
			require.Equal(t, http.StatusServiceUnavailable, w.Code)
		})
	}
}

// Tests for Request Validation

func TestExecuteCommand_RejectsMalformedRequest(t *testing.T) {
	handler := &MinecraftHandler{}

	body := bytes.NewBufferString(`{invalid json}`)
	req := httptest.NewRequest(http.MethodPost, "/api/minecraft/command", body)
	w := httptest.NewRecorder()

	handler.ExecuteCommand(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), "Invalid request body")
}
