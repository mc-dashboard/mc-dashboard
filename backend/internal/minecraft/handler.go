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
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Server start initiated",
		"result":  result,
	}); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func (h *MinecraftHandler) StopServer(w http.ResponseWriter, r *http.Request) {
	result := h.LambdaService.CallLambda("ec2-stop")

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Server stop initiated",
		"result":  result,
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
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"players": players,
		"count":   len(players),
	}); err != nil {
		log.Printf("Failed to encode players response: %v", err)
	}
}

func (h *MinecraftHandler) ExecuteCommand(w http.ResponseWriter, r *http.Request) {
	if h.RCONClient == nil {
		http.Error(w, "RCON client not initialized", http.StatusServiceUnavailable)
		return
	}

	var req struct {
		Command string `json:"command"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate command against whitelist
	commandBase := strings.Split(strings.TrimSpace(req.Command), " ")[0]
	if !allowedCommands[commandBase] && !allowedCommands[strings.Join(strings.Split(req.Command, " ")[:2], " ")] {
		http.Error(w, "Command not allowed. Only read-only commands are permitted.", http.StatusForbidden)
		return
	}

	response, err := h.RCONClient.ExecuteCommand(req.Command)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"response": response,
	}); err != nil {
		log.Printf("Failed to encode command response: %v", err)
	}
}