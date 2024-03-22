package main

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
)

// simplified API keys in memory, want to be able to generate predicatble keys for testing
var APIKeys = map[string]bool{}

// function to create a hash of the API keys and store in memory
// using an in memory solution for simplicity of code review
func hashAPIKeys(keys []string) {
	for _, key := range keys {
		hash := sha256.Sum256([]byte(key))
		hashedKey := hex.EncodeToString(hash[:])
		APIKeys[hashedKey] = true
		APIKeys[key] = true
	}
}

// function to handle api key validation
func validateAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("Authorization")
		_, found := APIKeys[apiKey]
		if !found {
			logger.Println("Unauthorized request")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
