package user

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"brianhang.me/facegraph/internal/db"
	"github.com/golang-jwt/jwt/v4"
)

var jwtKey []byte

const cookieName = "fg_user_token"

func fetchJWTKey() []byte {
	if jwtKey != nil {
		return jwtKey
	}

	jwtKey := os.Getenv("APP_USER_JWT_KEY")
	return []byte(jwtKey)
}

func createJWT(u *User) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Subject: fmt.Sprint(u.ID),
		},
	)
	return token.SignedString(fetchJWTKey())
}

func findUserFromJWT(tokenString string) (*User, error) {
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(fetchJWTKey()), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("failed to validate JWT: %v", err)
	}

	userID64, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%s is not a valid user ID: %v", claims.Subject, err)
	}
	userID := uint(userID64)

	db := db.Get()
	user := &User{}
	tx := db.First(user, userID)

	if tx.RowsAffected == 0 {
		user = nil
	}

	return user, nil
}

func SetCookie(w http.ResponseWriter, user *User) error {
	token, err := createJWT(user)
	if err != nil {
		return fmt.Errorf("failed to create JWT for user %v: %v", user, err)
	}

	cookie := http.Cookie{
		Name:     cookieName,
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   int(time.Hour * 24 * 30),
	}
	http.SetCookie(w, &cookie)

	return nil
}

func RouteWithUser(handler func(http.ResponseWriter, *http.Request, *User)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var u *User

		cookie, err := r.Cookie(cookieName)
		if err == nil {
			u, _ = findUserFromJWT(cookie.Value)
		}

		handler(w, r, u)
	}
}
