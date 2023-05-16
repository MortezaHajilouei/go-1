package delivery

import (
	"context"
	"micro/api/pb/base"
)

type UserController struct {
	base.UnimplementedUserServer
}

func NewUserController() UserController {
	return UserController{}
}

func (b *UserController) CreateUser(context.Context, *base.CreateUserRequest) (*base.CreateUserResponse, error) {

	return nil, nil
}
