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

func TestRegisterLogin_Login_Success(t *testing.T) {
	ctx, st, tearDown := suite.New(t)
	defer tearDown()

	email := gofakeit.Email()
	pass := gofakeit.Password(true, true, true, true, true, passDefaultLength)

	respReg, err := st.AuthClient.Register(ctx, &authv1.RegisterRequest{
		Email:    email,
		Password: pass,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetId())

	loginTime := time.Now()

	respLogin, err := st.AuthClient.Login(ctx, &authv1.LoginRequest{
		Email:    email,
		Password: pass,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respLogin.GetToken())

	token, err := jwt.Parse(respLogin.GetToken(), func(t *jwt.Token) (interface{}, error) {
		return []byte(st.Cfg.Secret), nil
	})
	require.NoError(t, err)

	claims, ok := token.Claims.(jwt.MapClaims)
	require.True(t, ok)
	assert.Equal(t, respReg.GetId(), claims["uid"].(string))
	assert.Equal(t, email, claims["email"].(string))

	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), int64(claims["exp"].(float64)), expirationDeltaSeconds)
}

func TestRegister_DuplicateUser(t *testing.T) {
	ctx, st, tearDown := suite.New(t)
	defer tearDown()

	email := gofakeit.Email()
	pass := gofakeit.Password(true, true, true, true, true, passDefaultLength)
	pass2 := gofakeit.Password(true, true, true, true, true, passDefaultLength)

	respReg, err := st.AuthClient.Register(ctx, &authv1.RegisterRequest{
		Email:    email,
		Password: pass,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetId())

	respReg2, err := st.AuthClient.Register(ctx, &authv1.RegisterRequest{
		Email:    email,
		Password: pass2,
	})
	require.Error(t, err)
	assert.Empty(t, respReg2.GetId())
	assert.ErrorContains(t, err, "user already exists")
	assert.Equal(t, status.Code(err), codes.AlreadyExists)
}
