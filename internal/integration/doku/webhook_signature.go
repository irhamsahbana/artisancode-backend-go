package doku

import (
	"codebase-app/internal/entity/coreentity"
	"crypto/hmac"
	"encoding/base64"
	"net/http"
	"strings"
)

func (c *dokuClient) VerifyWebhookSignature(headers http.Header, body []byte, targetPath string) bool {
	clientID := headers.Get("Client-Id")
	requestID := headers.Get("Request-Id")
	timestamp := headers.Get("Request-Timestamp")
	signature := headers.Get("Signature")

	if clientID == "" || requestID == "" || timestamp == "" || signature == "" {
		return false
	}

	expectedSignature := c.generateSignature(body, timestamp, requestID, targetPath)

	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

func (c *dokuClient) VerifyWebhookSignatureHeaders(headers coreentity.DOKUWebhookSignatureHeaders, body []byte, targetPath string) bool {
	httpHeaders := make(http.Header)
	httpHeaders.Set("Client-Id", headers.ClientID)
	httpHeaders.Set("Request-Id", headers.RequestID)
	httpHeaders.Set("Request-Timestamp", headers.RequestTimestamp)
	httpHeaders.Set("Signature", headers.Signature)

	return c.VerifyWebhookSignature(httpHeaders, body, targetPath)
}

func (c *dokuClient) DecodedPublicKey() string {
	if c.publicKey == "" {
		return ""
	}

	decoded, err := base64.StdEncoding.DecodeString(strings.Trim(c.publicKey, "\""))
	if err != nil {
		return c.publicKey
	}

	return string(decoded)
}
