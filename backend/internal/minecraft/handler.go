package minecraft

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/rohanvsuri/minecraft-dashboard/internal/lambda"
)

var allowedCommands = map[string]bool{
	"list":          true,
	"seed":          true,
	"time query":    true,
	"weather query": true,
	"gamerule":      true,
	"difficulty":    true,
}

// Response types for API endpoints
type ServerActionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Result  any    `json:"result,omitempty"`
}

type PlayerListResponse struct {
	Players []PlayerInfo `json:"players"`
	Count   int          `json:"count"`
}

type CommandResponse struct {
	Response string `json:"response"`
}

type MinecraftHandler struct {
	LambdaService *lambda.FunctionWrapper
	RCONClient    *RCONClient
}

func NewMinecraftHandler(lambdaService *lambda.FunctionWrapper, rconClient *RCONClient) *MinecraftHandler {
	return &MinecraftHandler{
		LambdaService: lambdaService,
		RCONClient:    rconClient,
	}
}

func (h *MinecraftHandler) StartServer(w http.ResponseWriter, r *http.Request) {
	result := h.LambdaService.CallLambda("ec2-start")

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ServerActionResponse{
		Success: true,
		Message: "Server start initiated",
		Result:  result,
	}); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func (h *MinecraftHandler) StopServer(w http.ResponseWriter, r *http.Request) {
	result := h.LambdaService.CallLambda("ec2-stop")

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ServerActionResponse{
		Success: true,
		Message: "Server stop initiated",
		Result:  result,
	}); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func (h *MinecraftHandler) GetServerStatus(w http.ResponseWriter, r *http.Request) {
	if h.RCONClient == nil {
		http.Error(w, "RCON client not initialized", http.StatusServiceUnavailable)
		return
	}

	status, err := h.RCONClient.GetServerStatus()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		log.Printf("Failed to encode status response: %v", err)
	}
}

func (h *MinecraftHandler) GetOnlinePlayers(w http.ResponseWriter, r *http.Request) {
	if h.RCONClient == nil {
		http.Error(w, "RCON client not initialized", http.StatusServiceUnavailable)
		return
	}

	players, err := h.RCONClient.GetOnlinePlayers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(PlayerListResponse{
		Players: players,
		Count:   len(players),
	}); err != nil {
		log.Printf("Failed to encode players response: %v", err)
	}
}

func (h *MinecraftHandler) ExecuteCommand(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Command string `json:"command"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate command against whitelist
	trimmedCmd := strings.TrimSpace(req.Command)
	if trimmedCmd == "" {
		http.Error(w, "Command cannot be empty", http.StatusBadRequest)
		return
	}

	if !isCommandAllowed(trimmedCmd) {
		http.Error(w, "Command not allowed. Only read-only commands are permitted.", http.StatusForbidden)
		return
	}

	if h.RCONClient == nil {
		http.Error(w, "RCON client not initialized", http.StatusServiceUnavailable)
		return
	}

	response, err := h.RCONClient.ExecuteCommand(req.Command)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(CommandResponse{
		Response: response,
	}); err != nil {
		log.Printf("Failed to encode command response: %v", err)
	}
}

// isCommandAllowed checks if a command matches the whitelist.
// Supports both single-word commands (e.g., "list") and two-word prefixes (e.g., "time query").
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
