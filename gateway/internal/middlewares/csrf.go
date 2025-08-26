package middlewares

import (
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

// "/v1/login"
// "/v1/register"
func MakeCSRFMiddleware(ignorePaths []string) func(next runtime.HandlerFunc) runtime.HandlerFunc {
	return func(next runtime.HandlerFunc) runtime.HandlerFunc {
		ignorePathMap := make(map[string]struct{}, len(ignorePaths))
		for _, path := range ignorePaths {
			ignorePathMap[path] = struct{}{}
		}

		return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete || r.Method == http.MethodPatch {
				if _, ok := ignorePathMap[r.URL.Path]; !ok {
					cookieToken, _ := r.Cookie("csrf_token")
					headerToken := r.Header.Get("X-CSRF-Token")
					if cookieToken == nil || cookieToken.Value != headerToken {
						http.Error(w, "CSRF token mismatch", http.StatusForbidden)
						return
					}
				}
			}
			next(w, r, pathParams)
		}
	}
}
