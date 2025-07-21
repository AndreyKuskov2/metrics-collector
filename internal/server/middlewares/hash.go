package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"strings"

	"github.com/AndreyKuskov2/metrics-collector/internal/server/config"
	"github.com/go-chi/render"
)

func calculateHash(data, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// CheckHashMiddleware - middleware для проверки хеша.
func CheckHashMiddleware(cfg *config.ServerConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cfg.SecretKey == "" {
				next.ServeHTTP(w, r)
				return
			}

			hash := r.Header.Get("HashSHA256")
			if hash == "" {
				render.Status(r, http.StatusBadRequest)
				render.PlainText(w, r, "")
				return
			}

			data, err := io.ReadAll(r.Body)
			if err != nil {
				render.Status(r, http.StatusBadRequest)
				render.PlainText(w, r, "")
				return
			}

			r.Body = io.NopCloser(strings.NewReader(string(data)))

			expectedHash := calculateHash(data, []byte(cfg.SecretKey))
			if expectedHash != hash {
				render.Status(r, http.StatusBadRequest)
				render.PlainText(w, r, "")
				return
			}

			next.ServeHTTP(w, r)

			responseData := []byte(w.Header().Get("Content-Type") + r.URL.Path + r.URL.RawQuery)
			responseHash := calculateHash(responseData, []byte(cfg.SecretKey))
			w.Header().Set("HashSHA256", responseHash)
		})
	}
}
