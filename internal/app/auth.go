package app

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func IssueToken(user User) string {
	raw := fmt.Sprintf("%d:%t:%d", user.ID, user.IsAdmin, time.Now().Unix())
	return base64.StdEncoding.EncodeToString([]byte(raw))
}

func ParseToken(token string) (AuthInfo, error) {
	if token == "" {
		return AuthInfo{}, errors.New("missing token")
	}

	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return AuthInfo{}, errors.New("invalid token")
	}

	parts := strings.Split(string(decoded), ":")
	if len(parts) < 2 {
		return AuthInfo{}, errors.New("invalid token format")
	}

	id, err := strconv.Atoi(parts[0])
	if err != nil {
		return AuthInfo{}, errors.New("invalid user id")
	}

	isAdmin, _ := strconv.ParseBool(parts[1])

	return AuthInfo{UserID: id, IsAdmin: isAdmin, Token: token}, nil
}

func AuthFromRequest(r *http.Request) (AuthInfo, error) {
	if header := r.Header.Get("Authorization"); header != "" {
		parts := strings.SplitN(header, " ", 2)
		if len(parts) == 2 {
			return ParseToken(strings.TrimSpace(parts[1]))
		}
	}

	if token := r.Header.Get("X-Auth-Token"); token != "" {
		return ParseToken(token)
	}

	if token := r.URL.Query().Get("token"); token != "" {
		return ParseToken(token)
	}

	return AuthInfo{}, errors.New("missing token")
}
