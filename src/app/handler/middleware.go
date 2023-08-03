package handler

import (
	"github.com/mjedari/vgang-project/app/configs"
	"github.com/mjedari/vgang-project/app/wiring"
	"log"
	"net"
	"net/http"
)

func RateLimiterMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !configs.Config.RateLimiter.Active {
			next.ServeHTTP(w, r)
			return
		}

		rateLimiter := wiring.Wiring.GetRateLimiter()

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "invalid ip address", http.StatusInternalServerError)
			return
		}

		if err := rateLimiter.Handle(ip); err != nil {
			http.Error(w, err.Error(), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func LoggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Our middleware logic goes here...
		log.Println("request :", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)
	}
}
