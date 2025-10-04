package middlewares

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func MakeAuthTokenMiddleware(ignorePaths []string) func(next runtime.HandlerFunc) runtime.HandlerFunc {
	return func(next runtime.HandlerFunc) runtime.HandlerFunc {
		ignorePathMap := make(map[string]struct{}, len(ignorePaths))
		for _, path := range ignorePaths {
			ignorePathMap[path] = struct{}{}
		}

		return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			if _, ok := ignorePathMap[r.URL.Path]; !ok {
				userId := ""
				cookies := r.CookiesNamed("token")
				if len(cookies) > 0 && len(cookies[0].Value) > 0 {
					token := cookies[0].Value

					jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
						return nil, nil
					})

					if err != nil || !jwtToken.Valid {
						// remove token
						http.SetCookie(w, &http.Cookie{
							Name:   "token",
							MaxAge: 0,
						})
						http.Error(w, "invalid token", http.StatusForbidden)
						return
					}

					claims, ok := jwtToken.Claims.(jwt.MapClaims)
					if ok {
						uid, ok := claims["uid"].(string)
						if ok {
							userId = uid
						}
					}
				}
				r.Header.Set("X-User-ID", userId)
			}
			next(w, r, pathParams)
		}
	}
}
