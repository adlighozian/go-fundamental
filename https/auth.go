package https

import (
	"database/sql"
	"encoding/json"
	"go-axiata/model"
	"go-axiata/pkg/helper"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (h *handler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:

			var jwtKey = []byte("my_secret_key")
			var creds model.Credentials
			err := json.NewDecoder(r.Body).Decode(&creds)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			res, err := h.usecase.Login(r.Context(), creds.Username)
			if err != nil {
				if err == sql.ErrNoRows {
					helper.RespondJSON(w, http.StatusUnauthorized, false, "Login failed", nil, nil)
					return
				}
				helper.RespondJSON(w, http.StatusInternalServerError, false, "Login failed", nil, nil)
				return
			}

			if err := bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(creds.Password)); err != nil {
				helper.RespondJSON(w, http.StatusUnauthorized, false, "Login failed", nil, nil)
				return
			}

			expirationTime := time.Now().Add(5 * time.Hour)
			claims := &model.Claims{
				Username: res.Username,
				Role:     res.Role,
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(expirationTime),
				},
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			tokenString, err := token.SignedString(jwtKey)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:    "token",
				Value:   tokenString,
				Expires: expirationTime,
			})

			helper.RespondJSON(w, http.StatusOK, true, "Login successfully", nil, tokenString)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func (h *handler) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:

			http.SetCookie(w, &http.Cookie{
				Name:    "token",
				Value:   "",
				Expires: time.Now(),
			})
			helper.RespondJSON(w, http.StatusOK, true, "Logout successfully", nil, nil)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func (h *handler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var creds model.Credentials
			err := json.NewDecoder(r.Body).Decode(&creds)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			creds.Password = string(hashedPassword)

			err = h.usecase.Register(r.Context(), creds)
			if err != nil {
				helper.RespondJSON(w, http.StatusBadRequest, false, err.Error(), nil, nil)
			} else {
				helper.RespondJSON(w, http.StatusOK, true, "Account successfully created", nil, nil)
			}

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
