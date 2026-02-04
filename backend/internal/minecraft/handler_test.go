package minecraft

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetServerStatus_NilClient(t *testing.T) {
	handler := &MinecraftHandler{
		RCONClient: nil,
	}

	req := httptest.NewRequest(http.MethodGet, "/api/minecraft/status", nil)
	w := httptest.NewRecorder()

	handler.GetServerStatus(w, req)

	require.Equal(t, http.StatusServiceUnavailable, w.Code)
	require.Contains(t, w.Body.String(), "RCON client not initialized")
}

func TestGetOnlinePlayers_NilClient(t *testing.T) {
	handler := &MinecraftHandler{
		RCONClient: nil,
	}

	req := httptest.NewRequest(http.MethodGet, "/api/minecraft/players", nil)
	w := httptest.NewRecorder()

	handler.GetOnlinePlayers(w, req)

	require.Equal(t, http.StatusServiceUnavailable, w.Code)
	require.Contains(t, w.Body.String(), "RCON client not initialized")
}

func TestExecuteCommand_EmptyCommand(t *testing.T) {
	handler := &MinecraftHandler{
		RCONClient: nil,
	}

	body := bytes.NewBufferString(`{"command": ""}`)
	req := httptest.NewRequest(http.MethodPost, "/api/minecraft/command", body)
	w := httptest.NewRecorder()

	handler.ExecuteCommand(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), "Command cannot be empty")
}

func TestExecuteCommand_WhitespaceOnly(t *testing.T) {
	handler := &MinecraftHandler{
		RCONClient: nil,
	}

	body := bytes.NewBufferString(`{"command": "   "}`)
	req := httptest.NewRequest(http.MethodPost, "/api/minecraft/command", body)
	w := httptest.NewRecorder()

	handler.ExecuteCommand(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), "Command cannot be empty")
}

func TestExecuteCommand_DisallowedCommand(t *testing.T) {
	handler := &MinecraftHandler{
		RCONClient: nil,
	}

	tests := []struct {
		name    string
		command string
	}{
		{"op command blocked", "op PlayerName"},
		{"ban command blocked", "ban PlayerName"},
		{"stop command blocked", "stop"},
		{"kick command blocked", "kick PlayerName"},
		{"unknown command blocked", "unknown"},
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

func TestExecuteCommand_AllowedCommandPassesValidation(t *testing.T) {
	handler := &MinecraftHandler{
		RCONClient: nil,
	}

	tests := []struct {
		name    string
		command string
	}{
		{"list command allowed", "list"},
		{"seed command allowed", "seed"},
		{"time query allowed", "time query daytime"},
		{"weather query allowed", "weather query"},
		{"gamerule allowed", "gamerule keepInventory"},
		{"difficulty allowed", "difficulty"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := bytes.NewBufferString(`{"command": "` + tt.command + `"}`)
			req := httptest.NewRequest(http.MethodPost, "/api/minecraft/command", body)
			w := httptest.NewRecorder()

			handler.ExecuteCommand(w, req)

			// Should pass whitelist validation (not 403), failing at RCON client check (503)
			require.Equal(t, http.StatusServiceUnavailable, w.Code)
		})
	}
}

func TestExecuteCommand_InvalidJSON(t *testing.T) {
	handler := &MinecraftHandler{
		RCONClient: nil,
	}

	body := bytes.NewBufferString(`{invalid json}`)
	req := httptest.NewRequest(http.MethodPost, "/api/minecraft/command", body)
	w := httptest.NewRecorder()

	handler.ExecuteCommand(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), "Invalid request body")
}
