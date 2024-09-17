package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log"
	"strings"
	"time"
)

const expiryDuration = time.Hour * 24 * 7

var secret []byte

func Init(secretKey string) {
	secret = []byte(secretKey)
}

// Creates a JWT token using HS256 algorithm
// "iat" and "exp" fields are added automatically to payload
func Create(payload map[string]any) string {
	payload["iat"] = time.Now().Unix()
	payload["exp"] = time.Now().Add(expiryDuration).Unix()
	pldBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}
	b64header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg": "HS256", "typ": "JWT"}`))
	b64payload := base64.RawURLEncoding.EncodeToString(pldBytes)
	signature := base64.RawURLEncoding.EncodeToString(
		hs256sum([]byte(b64header+"."+b64payload), secret),
	)
	return b64header + "." + b64payload + "." + signature
}

// Verifies a jwt token and also returns decoded payload if valid
func VerifyAndDecode(token string) (valid bool, payload map[string]any) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false, nil
	}
	b64header, b64payload, signature1 := parts[0], parts[1], parts[2]
	signature2 := base64.RawURLEncoding.EncodeToString(
		hs256sum([]byte(b64header+"."+b64payload), secret),
	)
	if signature1 != signature2 {
		return false, nil
	}
	decoded, _ := base64.RawURLEncoding.DecodeString(b64payload)
	json.Unmarshal(decoded, &payload)
	return true, payload
}

func hs256sum(data, key []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}
