package authentication

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/kyawthanttin/bpi-wms/config"
	"github.com/kyawthanttin/bpi-wms/webutil"
)

// Set up a global string for our secret
var mySigningKey = []byte("secret")

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JwtToken struct {
	Token string `json:"token"`
}

func SignIn(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			webutil.RespondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
			return
		}
		var user User
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&user); err != nil {
			webutil.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		defer r.Body.Close()
		userInfo, err := ValidateUser(env.DB, user.Username, user.Password)
		if err != nil {
			if err == sql.ErrNoRows {
				webutil.RespondWithError(w, http.StatusNotFound, "Invalid username or password")
			} else {
				webutil.RespondWithErrorType(w, err)
			}
			return
		}

		// Create the jwt token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": userInfo.Username,
			"name":     userInfo.Name,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})
		tokenString, err := token.SignedString(mySigningKey)
		if err != nil {
			webutil.RespondWithError(w, http.StatusInternalServerError, "Failed to sign token")
			return
		}
		webutil.RespondWithJSON(w, http.StatusOK, JwtToken{Token: tokenString})
	})
}

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string

		// Get token from the Authorization header
		// format: Authorization: Bearer
		tokens, ok := r.Header["Authorization"]
		if ok && len(tokens) >= 1 {
			token = tokens[0]
			token = strings.TrimPrefix(token, "Bearer ")
		}

		// If the token is empty...
		if token == "" {
			webutil.RespondWithError(w, http.StatusUnauthorized, "Empty token")
			return
		}

		// Now parse the token
		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				msg := fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				return nil, msg
			}
			return mySigningKey, nil
		})

		if err != nil {
			webutil.RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		if parsedToken != nil && parsedToken.Valid {
			// Token is valid.
			next.ServeHTTP(w, r)
			return
		}

		webutil.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	})
}
