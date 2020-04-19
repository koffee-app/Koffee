package auth

import (
	"fmt"
	"koffee/internal/config"
	"time"

	"github.com/o1egl/paseto"
)

// StreamingRustServerMessage is the JSON message stored inside the Paseto token that the Rust Streaming service  will generate
type StreamingRustServerMessage struct {
	// Duration of the new track in Seconds
	Duration string `json:"duration"`
	// ID of the new track (Path inside the Rust server)
	ID string `json:"path"`
}

// ValidateTokenFromRustService takes a paseto token and tries to parse it based on the key that generated the Rust service when saving a new mp3, these tokens specify the path that a new song was  saved in, and the duration of the song; later on we might use these for something else but at the moment it will only store that information, the transfer at the moment will be from the Rust streaming service to the client to this server CLIENT -> RSS -> CLIENT -> MONO API
func ValidateTokenFromRustService(token string) (*StreamingRustServerMessage, error) {
	key := config.StreamingServiceKey()
	var newFooter string
	var jsonToken paseto.JSONToken
	err := paseto.NewV2().Verify(token, key, &jsonToken, &newFooter)

	if err != nil {
		return nil, err
	}

	if jsonToken.Expiration.Unix() > time.Now().Unix() {
		return nil, fmt.Errorf("Token has expired")
	}

	path := jsonToken.Get("path")
	duration := jsonToken.Get("duration")
	return &StreamingRustServerMessage{ID: path, Duration: duration}, nil
}
