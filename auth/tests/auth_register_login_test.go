package tests

import (
	"testing"

	"github.com/alexwatcher/gateofthings/auth/tests/suite"
	authv1 "github.com/alexwatcher/gateofthings/protos/gen/go/auth/v1"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	passDefaultLength = 10
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
	assert.NotEmpty(t, respReg.Id)

	respLogin, err := st.AuthClient.Login(ctx, &authv1.LoginRequest{
		Email:    email,
		Password: pass,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respLogin.Token)
}
