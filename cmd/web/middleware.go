package main

import "net/http"

func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(
			"Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapi.com; font-src fonts.googleapi.com",
		)

		w.Header().Set("Referer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")

		////
		// Set to 0 (disabled) because these types of attacks are guarded against
		// with Content-Security-Policy nowadays.
		//
		w.Header().Set("X-XSS-Protection", "0")

		w.Header().Set("Server", "Go")
	})
}
