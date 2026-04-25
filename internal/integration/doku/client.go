package doku

import (
	"codebase-app/internal/infrastructure/config"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	sandboxBaseURL    = "https://api-sandbox.doku.com"
	productionBaseURL = "https://api.doku.com"
	defaultCurrency   = "IDR"
	defaultExpiryMins = 10080
)

type Config struct {
	BaseURL      string
	ClientID     string
	SecretKey    string
	PublicKey    string
	CallbackURL  string
	Timeout      time.Duration
	HTTPClient   *http.Client
	Now          func() time.Time
	NewRequestID func() string
}

type dokuClient struct {
	baseURL      string
	clientID     string
	secretKey    string
	publicKey    string
	callbackURL  string
	httpClient   *http.Client
	now          func() time.Time
	newRequestID func() string
}

func NewClient(cfg Config) *dokuClient {
	httpClient := cfg.HTTPClient
	if httpClient == nil {
		timeout := cfg.Timeout
		if timeout <= 0 {
			timeout = 30 * time.Second
		}
		httpClient = &http.Client{Timeout: timeout}
	}

	now := cfg.Now
	if now == nil {
		now = time.Now
	}

	newRequestID := cfg.NewRequestID
	if newRequestID == nil {
		newRequestID = uuid.NewString
	}

	baseURL := strings.TrimRight(cfg.BaseURL, "/")
	if baseURL == "" {
		baseURL = sandboxBaseURL
	}

	return &dokuClient{
		baseURL:      baseURL,
		clientID:     strings.TrimSpace(cfg.ClientID),
		secretKey:    strings.TrimSpace(cfg.SecretKey),
		publicKey:    strings.TrimSpace(cfg.PublicKey),
		callbackURL:  strings.TrimSpace(cfg.CallbackURL),
		httpClient:   httpClient,
		now:          now,
		newRequestID: newRequestID,
	}
}

func NewClientFromEnv() *dokuClient {
	baseURL := strings.TrimSpace(config.Envs.Doku.BaseURL)
	if baseURL == "" {
		baseURL = sandboxBaseURL
		if strings.EqualFold(config.Envs.App.Environtment, "production") {
			baseURL = productionBaseURL
		}
	}

	return NewClient(Config{
		BaseURL:     baseURL,
		ClientID:    config.Envs.Doku.ClientID,
		SecretKey:   config.Envs.Doku.SecretKey,
		PublicKey:   config.Envs.Doku.PublicKey,
		CallbackURL: config.Envs.Doku.CallbackURL,
		Timeout:     time.Duration(config.Envs.Doku.TimeoutSeconds) * time.Second,
	})
}
