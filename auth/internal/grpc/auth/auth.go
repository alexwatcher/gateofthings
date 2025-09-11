package auth

import (
	authv1 "github.com/alexwatcher/gateofthings/protos/gen/go/auth/v1"
	sharedvalid "github.com/alexwatcher/gateofthings/shared/pkg/grpc/interceptors/valid"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func init() {

	sharedvalid.RegisterRule(&authv1.SignUpRequest{}, func(s any) error {
		rr := s.(*authv1.SignUpRequest)
		return validation.ValidateStruct(rr,
			validation.Field(&rr.Email, validation.Required, is.EmailFormat),
			validation.Field(&rr.Password, validation.Required, validation.Length(6, 128)),
		)
	})

	sharedvalid.RegisterRule(&authv1.SignInRequest{}, func(s any) error {
		lr := s.(*authv1.SignInRequest)
		return validation.ValidateStruct(lr,
			validation.Field(&lr.Email, validation.Required, is.EmailFormat),
			validation.Field(&lr.Password, validation.Required, validation.Length(6, 128)),
		)
	})
}
