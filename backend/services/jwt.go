package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
)

// createJWT crea un token JWT simple
func (s *AuthService) createJWT(claims TokenClaims) (string, error) {
	// Header
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}

	headerEncoded := base64.RawURLEncoding.EncodeToString(headerJSON)

	// Payload
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	claimsEncoded := base64.RawURLEncoding.EncodeToString(claimsJSON)

	// Signature
	message := headerEncoded + "." + claimsEncoded
	signature := s.sign(message)

	// Token completo
	token := message + "." + signature

	return token, nil
}

// parseJWT parsea y valida un token JWT
func (s *AuthService) parseJWT(tokenString string) (*TokenClaims, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, errors.New("token inválido")
	}

	// Verificar firma
	message := parts[0] + "." + parts[1]
	expectedSignature := s.sign(message)

	if parts[2] != expectedSignature {
		return nil, errors.New("firma del token inválida")
	}

	// Decodificar claims
	claimsJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("error decodificando claims")
	}

	var claims TokenClaims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return nil, errors.New("error parseando claims")
	}

	return &claims, nil
}

// sign genera la firma HMAC-SHA256
func (s *AuthService) sign(message string) string {
	h := hmac.New(sha256.New, s.jwtSecret)
	h.Write([]byte(message))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
