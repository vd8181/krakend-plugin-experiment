package main

import (
	"bytes"
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
	"log"
	"net/http"
	"os"
	"time"

	// KrakenD HTTP-server plugin interface
	"github.com/devopsfaith/krakend/plugin/http-server"
)

// PluginConfig holds the fields configurable via krakend.json
type PluginConfig struct {
	SignatureKey  string `json:"signature_key"`
	DataKey       string `json:"data_key"`
	PublicKeyPath string `json:"public_key_path"`
}

// global config instance
var config PluginConfig

// init sets up logging for the plugin
func init() {
	logDir := "/tmp/krakend/logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Printf("[CertificateValidationPlugin] cannot create log dir: %v", err)
		return
	}
	f, err := os.OpenFile(logDir+"/plugin.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("[CertificateValidationPlugin] cannot open log file: %v", err)
		return
	}
	log.SetOutput(f)
	log.Println("[CertificateValidationPlugin] logger initialized")
}

// HandlerRegisterer is the symbol KrakenD looks for in the plugin .so
type HandlerRegisterer struct{}

// RegisterHandlers registers our plugin under the given name
func (h *HandlerRegisterer) RegisterHandlers(
	register func(name string, handler func(cfg map[string]interface{}) (http.Handler, error)),
) {
	register("CertificateValidationPlugin", NewCertificateValidationHandler)
}

// NewCertificateValidationHandler constructs the HTTP handler from the plugin config
func NewCertificateValidationHandler(cfg map[string]interface{}) (http.Handler, error) {
	// parse and store config
	conf := parseConfig(cfg)
	config = conf
	log.Printf("[CertificateValidationPlugin] loaded config: %+v", conf)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// only intercept POST /validatecertificate
		if r.Method != http.MethodPost || r.URL.Path != "/validatecertificate" {
			http.DefaultServeMux.ServeHTTP(w, r)
			return
		}
		log.Println("[CertificateValidationPlugin] intercepting request")

		// read and restore body
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("error reading body: %v", err)
			writeError(w, "Error reading request", http.StatusInternalServerError)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// parse JSON
		var body map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &body); err != nil {
			log.Printf("invalid JSON: %v", err)
			writeError(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		// extract data and signature
		rawData, ok1 := body[config.DataKey].(map[string]interface{})
		sig, ok2 := body[config.SignatureKey].(string)
		if !ok1 || !ok2 {
			log.Printf("missing %q or %q", config.DataKey, config.SignatureKey)
			writeError(w, "Missing data or signature", http.StatusBadRequest)
			return
		}

		// marshal nested data back to JSON
		dataBytes, err := json.Marshal(rawData)
		if err != nil {
			log.Printf("error marshaling data: %v", err)
			writeError(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// validate signature
		valid, err := validateSignature(dataBytes, sig, config.PublicKeyPath)
		if err != nil {
			log.Printf("validation error: %v", err)
			writeError(w, "Certificate validation error", http.StatusInternalServerError)
			return
		}
		if !valid {
			log.Println("signature not valid")
			writeError(w, "Invalid certificate signature", http.StatusUnauthorized)
			return
		}

		log.Println("signature valid; forwarding to backend")
		http.DefaultServeMux.ServeHTTP(w, r)
	}), nil
}

// parseConfig applies defaults and reads user-provided fields
func parseConfig(cfg map[string]interface{}) PluginConfig {
	conf := PluginConfig{
		SignatureKey:  "signature",
		DataKey:       "data",
		PublicKeyPath: "/tmp/krakend/plugins/certs/public_key.pem",
	}
	if pluginCfg, ok := cfg["CertificateValidationPlugin"].(map[string]interface{}); ok {
		if v, ok := pluginCfg["signature_key"].(string); ok {
			conf.SignatureKey = v
		}
		if v, ok := pluginCfg["data_key"].(string); ok {
			conf.DataKey = v
		}
		if v, ok := pluginCfg["public_key_path"].(string); ok {
			conf.PublicKeyPath = v
		}
	}
	return conf
}

// writeError sends a structured JSON error response
func writeError(w http.ResponseWriter, msg string, code int) {
	resp := map[string]interface{}{
		"error":     msg,
		"status":    "ERROR",
		"timestamp": time.Now().UTC().Format("2006-01-02T15:04:05Z"),
	}
	b, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(b)
}

// validateSignature checks the base64 signature against the data using RSA/SHA256
func validateSignature(data []byte, signatureB64, pubKeyPath string) (bool, error) {
	// load public key
	pub, err := loadPublicKey(pubKeyPath)
	if err != nil {
		return false, fmt.Errorf("load public key: %w", err)
	}

	// decode signature
	sigBytes, err := base64.StdEncoding.DecodeString(signatureB64)
	if err != nil {
		return false, fmt.Errorf("decode signature: %w", err)
	}

	// hash data
	h := sha256.New()
	h.Write(data)
	digest := h.Sum(nil)

	// verify
	if err := rsa.VerifyPKCS1v15(pub, crypto.SHA256, digest, sigBytes); err != nil {
		return false, nil // invalid signature
	}
	return true, nil
}

// loadPublicKey reads an RSA public key from PEM
func loadPublicKey(path string) (*rsa.PublicKey, error) {
	pemBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(pemBytes)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("invalid PEM public key")
	}
	if key, err := x509.ParsePKIXPublicKey(block.Bytes); err != nil {
		return nil, err
	} else if rsaKey, ok := key.(*rsa.PublicKey); !ok {
		return nil, errors.New("not an RSA public key")
	} else {
		return rsaKey, nil
	}
}
