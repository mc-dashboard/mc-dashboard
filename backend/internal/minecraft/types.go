package minecraft

import "time"

// ServerActionResponse is returned from server control endpoints
type ServerActionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Result  any    `json:"result,omitempty"`
}

// PlayerListResponse is returned from the player list endpoint
type PlayerListResponse struct {
	Players []PlayerInfo `json:"players"`
	Count   int          `json:"count"`
}

// PlayerInfo represents a Minecraft player
type PlayerInfo struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

// ServerStatus represents the current state of the Minecraft server
type ServerStatus struct {
	Online      bool         `json:"online"`
	PlayerCount int          `json:"playerCount"`
	MaxPlayers  int          `json:"maxPlayers"`
	Players     []PlayerInfo `json:"players"`
	LastUpdate  time.Time    `json:"lastUpdate"`
}

// listResponse is an internal struct for parsing the "list" command response
type listResponse struct {
	PlayerCount int
	MaxPlayers  int
	Players     []PlayerInfo
}
