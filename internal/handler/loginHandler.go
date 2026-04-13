package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"gin_demo/models"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var Users = map[string]string{
	"admin": "123456",
}

var jwtKey = []byte("my_secret_key")

// ===== LOGIN =====

// @Summary Login
// @Description login user
// @Accept json
// @Produce json
// @Param user body User true "User Info"
// @Success 200 {object} map[string]string
// @Router /login [post]
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var u User
	json.NewDecoder(r.Body).Decode(&u)

	pass, ok := Users[u.Username]
	if !ok || pass != u.Password {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(1 * time.Hour)

	claims := &models.Claims{
		Username: u.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString(jwtKey)

	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenStr,
	})
}
