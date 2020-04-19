package config

import "os"

var jwtKey string

// JWTConfig loads config of the jwt of the app
func JWTConfig() {
	key, ok := os.LookupEnv("JWT_KEY")
	if !ok {
		panic("JWT Key not retrieved correctly...")
	}
	jwtKey = key
}

// JWTKey returns the key for generating or parsing jwt
func JWTKey() string {
	return jwtKey
}
