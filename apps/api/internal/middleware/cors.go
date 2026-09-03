package middleware

import "net/http"

// CORS allows credentialed cross-origin requests from a single configured
// origin (the Next.js dev server runs on a different port than the API
// locally). credentials: "include" on the frontend means the response can
// never use a wildcard origin — it must echo back the exact allowed origin
// and set Access-Control-Allow-Credentials, or the browser discards any
// Set-Cookie from the response. This must wrap outside the mux: a preflight
// OPTIONS request won't match any registered "METHOD /path" pattern, so it
// has to be answered here before it ever reaches the mux.
func CORS(allowedOrigin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" && origin == allowedOrigin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Vary", "Origin")
			}

			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
