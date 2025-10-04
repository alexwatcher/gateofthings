package interceptors

import (
	"context"
	"net/http"

	authv1 "github.com/alexwatcher/gateofthings/protos/gen/go/auth/v1"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

func SetSignInCookies(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
	if r, ok := resp.(*authv1.SignInResponse); ok {
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    r.Token,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
		})

		csrfToken := uuid.NewString()
		w.Header().Set("X-CSRF-Token", csrfToken)
		http.SetCookie(w, &http.Cookie{
			Name:     "csrf_token",
			Value:    csrfToken,
			Path:     "/",
			HttpOnly: false,
			Secure:   true,
		})
	}
	return nil
}
