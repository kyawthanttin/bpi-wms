package authentication

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"

	"github.com/kyawthanttin/bpi-wms/config"
	"github.com/kyawthanttin/bpi-wms/user"
	"github.com/kyawthanttin/bpi-wms/webutil"
)

type JwtToken struct {
	Token string `json:"token"`
}

func SignIn(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			webutil.RespondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
			return
		}
		var data user.User
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&data); err != nil {
			webutil.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		defer r.Body.Close()

		userInfo, err := user.GetUserByUsername(env.DB, data.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				webutil.RespondWithError(w, http.StatusNotFound, "Invalid username or password")
			} else {
				webutil.RespondWithErrorType(w, err)
			}
			return
		}

		if !userInfo.IsEnabled {
			webutil.RespondWithError(w, http.StatusUnauthorized, "The user is disabled. Contact the Administrator.")
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(data.Password))
		if err != nil {
			webutil.RespondWithError(w, http.StatusNotFound, "Invalid username or password")
			return
		}

		err = user.RecordLogin(env.DB, userInfo.Id)
		if err != nil {
			webutil.RespondWithErrorType(w, err)
		}

		// Create the jwt token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": userInfo.Username,
			"name":     userInfo.Name,
			"roles":    userInfo.Roles,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})
		tokenString, err := token.SignedString(config.JWTSigningKey)
		if err != nil {
			webutil.RespondWithError(w, http.StatusInternalServerError, "Failed to sign token")
			return
		}
		webutil.RespondWithJSON(w, http.StatusOK, JwtToken{Token: tokenString})
	})
}

func Authenticate(env *config.Env, next http.Handler) http.Handler {
	return AuthenticateWithRoles(env, next, nil)
}

func AuthenticateWithRoles(env *config.Env, next http.Handler, roles []string) http.Handler {
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
			return config.JWTSigningKey, nil
		})

		if err != nil {
			webutil.RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		if parsedToken != nil && parsedToken.Valid {
			claims := parsedToken.Claims.(jwt.MapClaims)

			userInfo, err := user.GetUserByUsername(env.DB, claims["username"].(string))
			if err != nil {
				if err == sql.ErrNoRows {
					webutil.RespondWithError(w, http.StatusNotFound, "The user may have been deleted")
				} else {
					webutil.RespondWithErrorType(w, err)
				}
				return
			}

			if !userInfo.IsEnabled {
				webutil.RespondWithError(w, http.StatusUnauthorized, "The user is disabled. Contact the Administrator.")
				return
			}

			if roles != nil {
				userRoles := strings.Split(claims["roles"].(string), ",")
				hasRole := false
				for _, role := range roles {
					for _, userRole := range userRoles {
						if role == userRole {
							hasRole = true
							break
						}
					}
				}
				if !hasRole {
					webutil.RespondWithError(w, http.StatusUnauthorized, "Insufficient Privilege")
					return
				}
			}

			// Token is valid.
			next.ServeHTTP(w, r)
			return
		}

		webutil.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	})
}
