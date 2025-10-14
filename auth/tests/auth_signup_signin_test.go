package tests

import (
	"testing"
	"time"

	"github.com/alexwatcher/gateofthings/auth/tests/suite"
	authv1 "github.com/alexwatcher/gateofthings/protos/gen/go/auth/v1"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	passDefaultLength      = 10
	expirationDeltaSeconds = 2
)

func TestSignUp_SingIn_Success(t *testing.T) {
	ctx, st, tearDown := suite.New(t)
	defer tearDown()

	email := gofakeit.Email()
	pass := gofakeit.Password(true, true, true, true, true, passDefaultLength)

	respReg, err := st.AuthClient.SignUp(ctx, &authv1.SignUpRequest{
		Email:    email,
		Password: pass,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetId())

	signInTime := time.Now()

	respSignIn, err := st.AuthClient.SignIn(ctx, &authv1.SignInRequest{
		Email:    email,
		Password: pass,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respSignIn.GetToken())

	token, err := jwt.Parse(respSignIn.GetToken(), func(t *jwt.Token) (interface{}, error) {
		return []byte(st.Cfg.TokenSecret), nil
	})
	require.NoError(t, err)

	claims, ok := token.Claims.(jwt.MapClaims)
	require.True(t, ok)
	assert.Equal(t, respReg.GetId(), claims["uid"].(string))
	assert.Equal(t, email, claims["email"].(string))

	assert.InDelta(t, signInTime.Add(st.Cfg.TokenTTL).Unix(), int64(claims["exp"].(float64)), expirationDeltaSeconds)
}

func TestSignUp_DuplicateUser(t *testing.T) {
	ctx, st, tearDown := suite.New(t)
	defer tearDown()

	email := gofakeit.Email()
	pass := gofakeit.Password(true, true, true, true, true, passDefaultLength)
	pass2 := gofakeit.Password(true, true, true, true, true, passDefaultLength)

	respReg, err := st.AuthClient.SignUp(ctx, &authv1.SignUpRequest{
		Email:    email,
		Password: pass,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetId())

	respReg2, err := st.AuthClient.SignUp(ctx, &authv1.SignUpRequest{
		Email:    email,
		Password: pass2,
	})
	require.Error(t, err)
	assert.Empty(t, respReg2.GetId())
	assert.ErrorContains(t, err, "user already exists")
	assert.Equal(t, status.Code(err), codes.AlreadyExists)
}

func TestSignUp_FailedCases(t *testing.T) {
	ctx, st, tearDown := suite.New(t)
	defer tearDown()

	tests := []struct {
		name        string
		email       string
		pass        string
		expectedErr string
	}{
		{
			name:        "empty password",
			email:       gofakeit.Email(),
			pass:        "",
			expectedErr: "password: cannot be blank.",
		},
		{
			name:        "too short password",
			email:       gofakeit.Email(),
			pass:        gofakeit.Password(true, true, true, true, true, 3),
			expectedErr: "password: the length must be between",
		},
		{
			name:        "empty email",
			email:       "",
			pass:        gofakeit.Password(true, true, true, true, true, 10),
			expectedErr: "email: cannot be blank.",
		},
		{
			name:        "invalid email",
			email:       "email@mail@mail.com",
			pass:        gofakeit.Password(true, true, true, true, true, 10),
			expectedErr: "email: must be a valid email address.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.SignUp(ctx, &authv1.SignUpRequest{
				Email:    tt.email,
				Password: tt.pass,
			})
			require.Error(t, err)
			assert.ErrorContains(t, err, tt.expectedErr)
			assert.Equal(t, status.Code(err), codes.InvalidArgument)
		})
	}
}

func TestSignIn_FailedCases(t *testing.T) {
	ctx, st, tearDown := suite.New(t)
	defer tearDown()

	var validEmail = gofakeit.Email()
	var validPass = gofakeit.Password(true, true, true, true, true, 15)

	tests := []struct {
		name        string
		email       string
		pass        string
		expectedErr string
	}{
		{
			name:        "empty password",
			email:       validEmail,
			pass:        "",
			expectedErr: "password: cannot be blank.",
		},
		{
			name:        "invalid password",
			email:       validEmail,
			pass:        gofakeit.Password(true, true, true, true, true, 15),
			expectedErr: "invalid credentials.",
		},
		{
			name:        "empty email",
			email:       "",
			pass:        validPass,
			expectedErr: "email: cannot be blank.",
		},
		{
			name:        "invalid email",
			email:       "email@mail@mail.com",
			pass:        validPass,
			expectedErr: "email: must be a valid email address.",
		},
	}

	_, err := st.AuthClient.SignUp(ctx, &authv1.SignUpRequest{Email: validEmail, Password: validPass})
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.SignIn(ctx, &authv1.SignInRequest{
				Email:    tt.email,
				Password: tt.pass,
			})
			require.Error(t, err)
			assert.ErrorContains(t, err, tt.expectedErr)
			assert.Equal(t, status.Code(err), codes.InvalidArgument)
		})
	}
}
