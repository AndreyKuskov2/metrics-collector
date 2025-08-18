package middlewares

import (
	"net"
	"net/http"

	"github.com/AndreyKuskov2/metrics-collector/internal/server/config"
	"github.com/go-chi/render"
)

func CheckTrustedSubnetMiddleware(cfg *config.ServerConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cfg.TrustedSubnet == "" {
				next.ServeHTTP(w, r)
			} else {
				ipStr := r.Header.Get("X-Real-IP")
				if ipStr == "" {
					render.Status(r, http.StatusForbidden)
					render.PlainText(w, r, "cannot get x-real-ip header")
					return
				}

				_, subnet, err := net.ParseCIDR(cfg.TrustedSubnet)
				if err != nil {
					render.Status(r, http.StatusInternalServerError)
					render.PlainText(w, r, "cannot parse trusted subnet")
					return
				}

				agentIP := net.ParseIP(ipStr)
				if agentIP == nil || !subnet.Contains(agentIP) {
					render.Status(r, http.StatusForbidden)
					render.PlainText(w, r, "agent address is not in trusted subnet")
					return
				}

				next.ServeHTTP(w, r)
			}
		})
	}
}
