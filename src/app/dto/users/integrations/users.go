package integrations

import (
	"context"
	"encoding/json"
	"errors"
	"golang-shopeekuy/src/util/helper/integrations"
	"golang-shopeekuy/src/util/repository/model/users"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/google/uuid"
)

type userRepository interface {
	RegisterUser(bReq users.User) (*uuid.UUID, error)
	GetUserDetails(bReq users.User) (*users.User, error)
}

type UserUseCase struct {
	user userRepository
}

func NewUseCase(user userRepository) *UserUseCase {
	return &UserUseCase{
		user: user,
	}
}

func (u *UserUseCase) GetUsers(bReq users.User) (*users.User, error) {
	userData, err := u.user.GetUserDetails(bReq)
	if err != nil {
		return nil, err
	}

	return userData, nil
}

func (u *UserUseCase) Register(bReq users.User) (*uuid.UUID, error) {
	result, err := u.user.RegisterUser(bReq)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (u *UserUseCase) c {
	if state != integrations.RandomString {
		return nil, errors.New("invalid state")
	}
	token, err := integrations.SSOSignUp.Exchange(context.Background(), code)
	if err != nil {
		return nil, errors.New("cannot retrieve token")
	}

	provider, err := oidc.NewProvider(context.Background(), integrations.Provider)
	if err != nil {
		return nil, errors.New("invalid token signature")
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: integrations.SSOSignUp.ClientID,
	})
	_, err = verifier.Verify(context.Background(), token.Extra("id_token").(string))
	if err != nil {
		return nil, errors.New("invalid token signature")
	}

	result, err := http.Get(integrations.UserInfoURL + token.AccessToken)
	if err != nil {
		return nil, errors.New("cannot fetch user data")
	}

	var response users.OauthUserData

	if err := json.NewDecoder(result.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil

}
