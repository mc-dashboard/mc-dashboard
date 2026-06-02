package minecraft

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/rohanvsuri/minecraft-dashboard/internal/lambda"
)

// allowedCommands defines the whitelist of safe, read-only commands
var allowedCommands = map[string]bool{
	"list":          true,
	"seed":          true,
	"time query":    true,
	"weather query": true,
	"gamerule":      true,
	"difficulty":    true,
}

// MinecraftHandler handles HTTP requests for Minecraft server operations
type MinecraftHandler struct {
	lambdaService *lambda.FunctionWrapper
	rconClient    *RCONClient
}

func NewMinecraftHandler(lambdaService *lambda.FunctionWrapper, rconClient *RCONClient) *MinecraftHandler {
	return &MinecraftHandler{
		lambdaService: lambdaService,
		rconClient:    rconClient,
	}
}

func (h *MinecraftHandler) StartServer(w http.ResponseWriter, r *http.Request) {
	result := h.lambdaService.CallLambda("ec2-start")

	h.writeJSON(w, ServerActionResponse{
		Success: true,
		Message: "Server start initiated",
		Result:  result,
	})
}

func (h *MinecraftHandler) StopServer(w http.ResponseWriter, r *http.Request) {
	result := h.lambdaService.CallLambda("ec2-stop")

	h.writeJSON(w, ServerActionResponse{
		Success: true,
		Message: "Server stop initiated",
		Result:  result,
	})
}

func (h *MinecraftHandler) GetServerStatus(w http.ResponseWriter, r *http.Request) {
	if h.rconClient == nil {
		http.Error(w, "RCON client not initialized", http.StatusServiceUnavailable)
		return
	}

	status, err := h.rconClient.GetServerStatus()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, status)
}

func (h *MinecraftHandler) GetOnlinePlayers(w http.ResponseWriter, r *http.Request) {
	if h.rconClient == nil {
		http.Error(w, "RCON client not initialized", http.StatusServiceUnavailable)
		return
	}

	players, err := h.rconClient.GetOnlinePlayers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, PlayerListResponse{
		Players: players,
		Count:   len(players),
	})
}

func (h *MinecraftHandler) ExecuteCommand(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Command string `json:"command"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	trimmedCmd := strings.TrimSpace(req.Command)
	if trimmedCmd == "" {
		http.Error(w, "Command cannot be empty", http.StatusBadRequest)
		return
	}

	if !isCommandAllowed(trimmedCmd) {
		http.Error(w, "Command not allowed. Only read-only commands are permitted.", http.StatusForbidden)
		return
	}

	if h.rconClient == nil {
		http.Error(w, "RCON client not initialized", http.StatusServiceUnavailable)
		return
	}

	response, err := h.rconClient.SendCommand(req.Command)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, map[string]string{"response": response})
}

func (h *MinecraftHandler) writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

// isCommandAllowed uses a two-level check because some Minecraft commands
// require specific subcommands (like "time query") while others are standalone.
func isCommandAllowed(cmd string) bool {
	parts := strings.Split(cmd, " ")
	if allowedCommands[parts[0]] {
		return true
	}
	if len(parts) >= 2 && allowedCommands[parts[0]+" "+parts[1]] {
		return true
	}
	return false
}
