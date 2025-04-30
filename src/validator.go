package main

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/luraproject/lura/v2/proxy"
)

// Plugin exports the required symbol for KrakenD.
var Plugin = func() interface{} {
	return certificateValidator
}

// Config holds the plugin configuration.
type Config struct {
	PublicKeyPath string `json:"public_key_path"`
}

// RequestData represents the data structure in the request.
type RequestData struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	DOB   string `json:"dob"`
}

// Payload represents the complete request payload.
type Payload struct {
	Data      RequestData `json:"data"`
	Signature string      `json:"signature"`
}

// HTTPError represents a custom HTTP error with a status code.
type HTTPError struct {
	StatusCode int
	Message    string
}

func (e *HTTPError) Error() string {
	return e.Message
}

// global public key path, set via the PUBLIC_KEY_PATH environment variable.
// If not set, defaults to "/etc/krakend/certs/demo/public_key.pem".
var publicKeyPath string

func init() {
	publicKeyPath = os.Getenv("PUBLIC_KEY_PATH")
	if publicKeyPath == "" {
		publicKeyPath = "/etc/krakend/certs/demo/public_key.pem"
	}
}

// certificateValidator is the middleware that validates the certificate.
func certificateValidator(next proxy.Proxy) proxy.Proxy {
	return func(ctx context.Context, req *proxy.Request) (*proxy.Response, error) {
		// Load the public key using the global publicKeyPath.
		publicKey, err := loadPublicKey(publicKeyPath)
		if err != nil {
			return nil, &HTTPError{StatusCode: 500, Message: "failed to load public key"}
		}

		// Read and decode the request body.
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, &HTTPError{StatusCode: 400, Message: "failed to read request body"}
		}

		var payload Payload
		if err := json.Unmarshal(body, &payload); err != nil {
			return nil, &HTTPError{StatusCode: 400, Message: "invalid JSON payload"}
		}

		// Decode the signature.
		signature, err := base64.StdEncoding.DecodeString(payload.Signature)
		if err != nil {
			return nil, &HTTPError{StatusCode: 400, Message: "invalid signature format"}
		}

		// Serialize the 'data' field for validation.
		dataToEncode, err := json.Marshal(payload.Data)
		if err != nil {
			return nil, &HTTPError{StatusCode: 500, Message: "failed to serialize data for validation"}
		}

		// Validate the signature.
		if !validateSignature(publicKey, dataToEncode, signature) {
			return nil, &HTTPError{StatusCode: 401, Message: "certificate validation failed"}
		}

		// Continue processing if validation is successful.
		return next(ctx, req)
	}
}

// loadPublicKey loads the RSA public key from a PEM file.
func loadPublicKey(filePath string) (*rsa.PublicKey, error) {
	publicKeyPEM, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file: %v", err)
	}

	block, _ := pem.Decode(publicKeyPEM)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	rsaPublicKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}

	return rsaPublicKey, nil
}

// validateSignature validates the RSA signature using the public key.
func validateSignature(publicKey *rsa.PublicKey, data []byte, signature []byte) bool {
	hashed := sha256.Sum256(data)
	err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signature)
	if err != nil {
		fmt.Printf("Signature verification failed: %v\n", err)
		return false
	}
	return true
}

func main() {
	// This function is required to compile the plugin.
}