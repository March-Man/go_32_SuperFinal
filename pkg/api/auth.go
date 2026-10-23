package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("my_secret_key")

// Claims represents the JWT token payload.
type Claims struct {
	Hash string `json:"hash"`
	jwt.RegisteredClaims
}

// signInHandler handles the POST /api/signin request to authenticate a user.
func signInHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, map[string]string{"error": "ошибка десериализации"}, http.StatusBadRequest)
		return
	}

	expectedPass := os.Getenv("TODO_PASSWORD")
	if expectedPass == "" {
		writeJSON(w, map[string]string{"token": ""}, http.StatusOK)
		return
	}

	if req.Password != expectedPass {
		writeJSON(w, map[string]string{"error": "Неверный пароль"}, http.StatusUnauthorized)
		return
	}

	hash := sha256.Sum256([]byte(expectedPass))
	hashStr := hex.EncodeToString(hash[:])

	claims := Claims{
		Hash: hashStr,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Printf("ошибка создания токена: %v", err)
		writeJSON(w, map[string]string{"error": "Внутренняя ошибка"}, http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"token": tokenString}, http.StatusOK)
}

// auth is a middleware that validates JWT authentication for protected routes.
func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) == 0 {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseWithClaims(cookie.Value, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("неверный метод подписи")
			}
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		hash := sha256.Sum256([]byte(pass))
		hashStr := hex.EncodeToString(hash[:])
		if claims.Hash != hashStr {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	})
}
