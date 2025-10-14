package options

import (
	"context"
	"net/http"

	"github.com/alexwatcher/gateofthings/gateway/internal/consts"
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
		w.Header().Set(consts.HttpCsrfTokenHeader, csrfToken)
		http.SetCookie(w, &http.Cookie{
			Name:     "csrf_token",
			Value:    csrfToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
		})
	}
	return nil
}

// RemoveSignInToken is a ForwardResponseRewriter
func RemoveSignInToken(ctx context.Context, response proto.Message) (any, error) {
	if _, ok := response.(*authv1.SignInResponse); ok {
		return struct{}{}, nil
	}
	return response, nil
}
