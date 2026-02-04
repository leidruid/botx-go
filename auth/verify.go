package auth

import (
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrMissingAuthorization = errors.New("authorization header missing")
	ErrInvalidToken         = errors.New("invalid token")
	ErrInvalidAudience      = errors.New("invalid audience")
	ErrInvalidIssuer        = errors.New("invalid issuer")
)

// VerifyRequest verifies BotX request JWT token.
// It checks aud (must contain a single bot_id) and issuer (iss).
func VerifyRequest(authHeader string, secretKey string, expectedIssuer string, expectedAudience string, trustedIssuers map[string]struct{}) error {
	if authHeader == "" {
		return ErrMissingAuthorization
	}

	parts := strings.Split(authHeader, " ")
	tokenStr := parts[len(parts)-1]

	// Decode without verification to extract claims.
	parser := jwt.NewParser(jwt.WithoutClaimsValidation())
	claims := jwt.MapClaims{}
	_, _, err := parser.ParseUnverified(tokenStr, claims)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidToken, err.Error())
	}

	audRaw, ok := claims["aud"]
	if !ok {
		return ErrInvalidAudience
	}

	// Expect aud to be array with single bot_id.
	var audList []string
	switch v := audRaw.(type) {
	case []any:
		for _, item := range v {
			if s, ok := item.(string); ok {
				audList = append(audList, s)
			}
		}
	case []string:
		audList = v
	case string:
		audList = []string{v}
	}
	if len(audList) != 1 {
		return ErrInvalidAudience
	}
	if expectedAudience != "" && audList[0] != expectedAudience {
		return ErrInvalidAudience
	}

	iss, _ := claims["iss"].(string)
	if iss == "" {
		return ErrInvalidIssuer
	}

	// Verify signature. Audience/issuer are checked manually to match BotX rules.
	_, err = jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Method.Alg())
		}
		return []byte(secretKey), nil
	}, jwt.WithoutClaimsValidation())
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidToken, err.Error())
	}

	if iss != expectedIssuer {
		if trustedIssuers == nil {
			return ErrInvalidIssuer
		}
		if _, ok := trustedIssuers[iss]; !ok {
			return ErrInvalidIssuer
		}
	}

	return nil
}
