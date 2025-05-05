package main

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
)

var pluginName = "validator"

// HandlerRegisterer is the symbol the plugin loader will try to load. It must implement the Registerer interface
var HandlerRegisterer = registerer(pluginName)

type registerer string

func (r registerer) RegisterHandlers(f func(
	name string,
	handler func(context.Context, map[string]interface{}, http.Handler) (http.Handler, error),
)) {
	f(string(r), r.registerHandlers)
}

func (r registerer) registerHandlers(_ context.Context, extra map[string]interface{}, h http.Handler) (http.Handler, error) {
	config, ok := extra[pluginName].(map[string]interface{})
	if !ok {
		return h, errors.New("configuration not found")
	}

	path, _ := config["path"].(string)
	logger.Debug(fmt.Sprintf("The plugin is now hijacking the path %s", path))

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != path {
			h.ServeHTTP(w, req)
			return
		}

		fmt.Fprintf(w, "Hello, %q", html.EscapeString(req.URL.Path))

		// Load the public key
		publicKey, err := loadPublicKey("/etc/krakend/certs/demo/public_key.pem")
		if err != nil {
			logger.Error("failed to load public key:", err)
			http.Error(w, "failed to load public key", http.StatusInternalServerError)
			return
		}

		// Decode the base64 signature from the request body
		body, err := io.ReadAll(req.Body)
		if err != nil {
			logger.Error("failed to read request body:", err)
			http.Error(w, "failed to read request body", http.StatusBadRequest)
			return
		}

		// Assume the request body contains the base64-encoded signature
		signature, err := base64.StdEncoding.DecodeString(string(body))
		if err != nil {
			logger.Error("invalid base64 signature format:", err)
			http.Error(w, "invalid signature format", http.StatusBadRequest)
			return
		}

		// Example data that should be verified
		dataToEncode := []byte("example data")

		// Validate the signature
		if !validateSignature(publicKey, dataToEncode, signature) {
			logger.Error("certificate validation failed")
			http.Error(w, "certificate validation failed", http.StatusUnauthorized)
			return
		}

		logger.Debug("request:", html.EscapeString(req.URL.Path))
	}), nil
}

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

func validateSignature(publicKey *rsa.PublicKey, data []byte, signature []byte) bool {
	hashed := sha256.Sum256(data)
	err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signature)
	if err != nil {
		fmt.Printf("Signature verification failed: %v\n", err)
		return false
	}
	return true
}

func main() {}

// This logger is replaced by the RegisterLogger method to load the one from KrakenD
var logger Logger = noopLogger{}

func (registerer) RegisterLogger(v interface{}) {
	l, ok := v.(Logger)
	if !ok {
		return
	}
	logger = l
	logger.Debug(fmt.Sprintf("[PLUGIN: %s] Logger loaded", HandlerRegisterer))
}

type Logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Warning(v ...interface{})
	Error(v ...interface{})
	Critical(v ...interface{})
	Fatal(v ...interface{})
}

// Empty logger implementation
type noopLogger struct{}

func (n noopLogger) Debug(_ ...interface{})    {}
func (n noopLogger) Info(_ ...interface{})     {}
func (n noopLogger) Warning(_ ...interface{})  {}
func (n noopLogger) Error(_ ...interface{})    {}
func (n noopLogger) Critical(_ ...interface{}) {}
func (n noopLogger) Fatal(_ ...interface{})    {}
