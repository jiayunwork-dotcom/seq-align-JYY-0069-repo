// Package server - CORS middleware for the alignment API server.
package server

import "net/http"

// CORSConfig holds CORS middleware configuration.
type CORSConfig struct {
	AllowOrigins []string
	AllowMethods []string
	AllowHeaders []string
	MaxAge       int // preflight cache duration in seconds
}

// DefaultCORSConfig returns a permissive CORS config for development.
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		MaxAge:       86400,
	}
}

// CORS wraps an http.Handler with CORS headers based on the given config.
func CORS(next http.Handler, cfg CORSConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			next.ServeHTTP(w, r)
			return
		}
		allowed := isOriginAllowed(origin, cfg.AllowOrigins)
		if !allowed {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", joinStrings(cfg.AllowMethods))
		w.Header().Set("Access-Control-Allow-Headers", joinStrings(cfg.AllowHeaders))

		// 处理预检请求
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// isOriginAllowed checks whether the origin is in the allowed list.
func isOriginAllowed(origin string, allowed []string) bool {
	for _, a := range allowed {
		if a == "*" || a == origin {
			return true
		}
	}
	return false
}

// joinStrings joins strings with ", " separator.
func joinStrings(ss []string) string {
	if len(ss) == 0 {
		return ""
	}
	result := ss[0]
	for _, s := range ss[1:] {
		result += ", " + s
	}
	return result
}
