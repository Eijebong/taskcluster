package util

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

// Min returns minimum of two ints
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

var (
	jwtRe = regexp.MustCompile(`^Bearer ([\w-\.]+)$`)
)

// MakeWsURL converts http:// to ws://
func MakeWsURL(url string) string {
	if u, trimmed := strings.CutPrefix(url, "http"); trimmed {
		return "ws" + u
	}
	return url
}

func ExtractJWT(authHeader string) string {
	c := jwtRe.FindStringSubmatch(authHeader)
	if len(c) != 2 {
		return ""
	}
	return c[1]
}

func GetTokenExp(tokenString string) time.Time {
	token, err := jwt.Parse(tokenString, nil)
	if err.(*jwt.ValidationError).Errors == jwt.ValidationErrorMalformed {
		return time.Time{}
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return time.Time{}
	}

	exp, ok := claims["exp"]
	if !ok {
		return time.Time{}
	}
	switch exp := exp.(type) {
	case float64:
		return time.Unix(int64(exp), 0)
	case json.Number:
		v, _ := exp.Int64()
		return time.Unix(v, 0)
	}
	return time.Time{}

}

// verify token is valid, and also exp and nbf on token
func IsTokenUsable(tokenString string) bool {
	token, err := jwt.Parse(tokenString, nil)
	if err.(*jwt.ValidationError).Errors == jwt.ValidationErrorMalformed {
		return false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	now := time.Now().Unix()
	return claims.VerifyExpiresAt(now, true) && claims.VerifyNotBefore(now, true)
}
