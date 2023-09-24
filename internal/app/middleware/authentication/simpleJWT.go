package authentication

import (
	"crypto/rand"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const authorizationURL = "/api/user/urls"

type userStorage interface {
	GetFreeID() (int, error)
	SaveNewUser() (int, error)
}

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

type TokenController struct {
	users     userStorage
	secretKey []byte
}

func NewTokenController(u userStorage) *TokenController {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		log.Fatalln(err)
	}
	return &TokenController{
		users:     u,
		secretKey: key,
	}
}

func (t *TokenController) CheckToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("Authorization")
		userID := t.getUserID(cookie.String())

		if userID == -1 || err != nil {

			if r.URL.Path == authorizationURL && r.Method == http.MethodGet {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			_, err := t.users.SaveNewUser()
			if err != nil {
				log.Fatalln(err)
			}

			id, err := t.users.GetFreeID()
			if err != nil {
				log.Fatalln(err)
			}

			newToken := t.createToken(id)
			newCookie := http.Cookie{
				Name:  "Authorization",
				Value: newToken,
			}
			http.SetCookie(w, &newCookie)

			r.Header.Add("userID", strconv.Itoa(id))
		}

		r.Header.Add("userID", strconv.Itoa(userID))

		next.ServeHTTP(w, r)
	})
}

func (t *TokenController) createToken(id int) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 3)),
		},
		UserID: id,
	})

	tokenString, err := token.SignedString(t.secretKey)
	if err != nil {
		log.Fatalln(err)
	}

	return tokenString
}

func (t *TokenController) getUserID(tokenString string) int {
	if tokenString == "" {
		return -1
	}
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(token *jwt.Token) (interface{}, error) {
			return t.secretKey, nil
		})
	if err != nil {
		return -1
	}

	if !token.Valid {
		return -1
	}
	return claims.UserID
}
