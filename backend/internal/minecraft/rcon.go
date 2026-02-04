package minecraft

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/james4k/rcon"
)

type RCONClient struct {
	host     string
	port     string
	password string
	conn     *rcon.RemoteConsole
	mu       sync.Mutex
	timeout  time.Duration
}

func NewRCONClient(host, port, password string) *RCONClient {
	return &RCONClient{
		host:     host,
		port:     port,
		password: password,
		timeout:  10 * time.Second,
	}
}

func (c *RCONClient) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return nil
	}

	address := fmt.Sprintf("%s:%s", c.host, c.port)

	// Connect with timeout using a goroutine
	type result struct {
		conn *rcon.RemoteConsole
		err  error
	}

	resultChan := make(chan result, 1)
	go func() {
		conn, err := rcon.Dial(address, c.password)
		resultChan <- result{conn, err}
	}()

	select {
	case res := <-resultChan:
		if res.err != nil {
			return fmt.Errorf("failed to connect to RCON: %w", res.err)
		}
		c.conn = res.conn
		log.Printf("Connected to Minecraft RCON at %s", address)
		return nil
	case <-time.After(c.timeout):
		return fmt.Errorf("RCON connection timed out after %v", c.timeout)
	}
}

func (c *RCONClient) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil
		return err
	}
	return nil
}

func (c *RCONClient) SendCommand(command string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return "", fmt.Errorf("not connected to RCON")
	}

	reqID, err := c.conn.Write(command)
	if err != nil {
		log.Printf("RCON write failed: %v", err)
		c.conn.Close()
		c.conn = nil
		return "", fmt.Errorf("failed to send command: %w", err)
	}

	response, responseID, err := c.conn.Read()
	if err != nil {
		log.Printf("RCON read failed: %v", err)
		c.conn.Close()
		c.conn = nil
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if reqID != responseID {
		return "", fmt.Errorf("response ID mismatch: want %d, got %d", reqID, responseID)
	}

	return response, nil
}

type PlayerInfo struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

type ServerStatus struct {
	Online      bool         `json:"online"`
	PlayerCount int          `json:"playerCount"`
	MaxPlayers  int          `json:"maxPlayers"`
	Players     []PlayerInfo `json:"players"`
	Version     string       `json:"version"`
	MOTD        string       `json:"motd"`
	LastUpdate  time.Time    `json:"lastUpdate"`
}

func (c *RCONClient) GetOnlinePlayers() ([]PlayerInfo, error) {
	response, err := c.SendCommand("list")
	if err != nil {
		return nil, err
	}
	return parseListResponse(response).Players, nil
}

func (c *RCONClient) GetPlayerCount() (int, int, error) {
	response, err := c.SendCommand("list")
	if err != nil {
		return 0, 0, err
	}

	// Parse response like: "There are 3 of a max of 20 players online"
	var current, max int
	_, err = fmt.Sscanf(response, "There are %d of a max of %d", &current, &max)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse player count: %w", err)
	}

	return current, max, nil
}

func (c *RCONClient) GetServerStatus() (*ServerStatus, error) {
	response, err := c.SendCommand("list")
	if err != nil {
		return &ServerStatus{
			Online:     false,
			LastUpdate: time.Now(),
		}, nil
	}

	parsed := parseListResponse(response)
	return &ServerStatus{
		Online:      true,
		PlayerCount: parsed.PlayerCount,
		MaxPlayers:  parsed.MaxPlayers,
		Players:     parsed.Players,
		LastUpdate:  time.Now(),
	}, nil
}

type listResponse struct {
	PlayerCount int
	MaxPlayers  int
	Players     []PlayerInfo
}

func parseListResponse(response string) listResponse {
	var result listResponse

	fmt.Sscanf(response, "There are %d of a max of %d", &result.PlayerCount, &result.MaxPlayers)

	if strings.Contains(response, "0 of a max") {
		return result
	}

	parts := strings.Split(response, ": ")
	if len(parts) < 2 {
		return result
	}

	for _, name := range strings.Split(parts[1], ", ") {
		name = strings.TrimSpace(name)
		if name != "" {
			result.Players = append(result.Players, PlayerInfo{Name: name})
		}
	}

	return result
}

// GetPlayerStats retrieves stats for a specific player
func (c *RCONClient) GetPlayerStats(playerName string) (string, error) {
	command := fmt.Sprintf("data get entity %s", playerName)
	return c.SendCommand(command)
}

// ExecuteCommand sends a raw command to the server
func (c *RCONClient) ExecuteCommand(command string) (string, error) {
	return c.SendCommand(command)
}
