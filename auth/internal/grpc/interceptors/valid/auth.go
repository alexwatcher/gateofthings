package valid

import (
	authv1 "github.com/alexwatcher/gateofthings/protos/gen/go/auth/v1"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func init() {
	rulesMap[getStructId(&authv1.RegisterRequest{})] = func(s any) error {
		rr := s.(*authv1.RegisterRequest)
		return validation.ValidateStruct(rr,
			validation.Field(&rr.Email, validation.Required, is.EmailFormat),
			validation.Field(&rr.Password, validation.Required, validation.Length(6, 128)),
		)
	}

	rulesMap[getStructId(&authv1.LoginRequest{})] = func(s any) error {
		lr := s.(*authv1.LoginRequest)
		return validation.ValidateStruct(lr,
			validation.Field(&lr.Email, validation.Required, is.EmailFormat),
			validation.Field(&lr.Password, validation.Required, validation.Length(6, 128)),
		)
	}
}
