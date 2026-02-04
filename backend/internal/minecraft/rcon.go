package minecraft

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/james4k/rcon"
)

// RCONClient manages the connection to a Minecraft server's RCON interface
type RCONClient struct {
	host     string
	port     string
	password string
	conn     *rcon.RemoteConsole
	mu       sync.Mutex
	timeout  time.Duration
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

type listResponse struct {
	PlayerCount int
	MaxPlayers  int
	Players     []PlayerInfo
}

// Constructor

func NewRCONClient(host, port, password string) *RCONClient {
	return &RCONClient{
		host:     host,
		port:     port,
		password: password,
		timeout:  10 * time.Second,
	}
}

// Connection Lifecycle

func (c *RCONClient) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return nil
	}

	address := fmt.Sprintf("%s:%s", c.host, c.port)

	type dialResult struct {
		conn *rcon.RemoteConsole
		err  error
	}
	ch := make(chan dialResult, 1)
	go func() {
		conn, err := rcon.Dial(address, c.password)
		ch <- dialResult{conn, err}
	}()

	select {
	case r := <-ch:
		if r.err != nil {
			return fmt.Errorf("failed to connect to RCON: %w", r.err)
		}
		c.conn = r.conn
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

// Core Communication

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

// High-Level Queries

func (c *RCONClient) GetOnlinePlayers() ([]PlayerInfo, error) {
	response, err := c.SendCommand("list")
	if err != nil {
		return nil, err
	}
	return parseListResponse(response).Players, nil
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

// Internal Helpers

func parseListResponse(response string) listResponse {
	var result listResponse
	fmt.Sscanf(response, "There are %d of a max of %d", &result.PlayerCount, &result.MaxPlayers)

	parts := strings.Split(response, ": ")
	if len(parts) < 2 {
		return result
	}

	for name := range strings.SplitSeq(parts[1], ", ") {
		name = strings.TrimSpace(name)
		if name != "" {
			result.Players = append(result.Players, PlayerInfo{Name: name})
		}
	}
	return result
}
