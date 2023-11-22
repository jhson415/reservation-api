package api

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jhson415/reservation-api/db"
	"github.com/jhson415/reservation-api/types"
)

type UserHandler struct {
	ctx       context.Context
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (m *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	var (
		id  = c.Params("id")
		ctx = c.Context()
	)

	user, err := m.userStore.GetUserById(ctx, id)
	if err != nil {
		return err
	}
	return c.JSON(user)
}

func (m *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
	)
	users, err := m.userStore.GetUserList(ctx)
	if err != nil {
		return err
	}
	return c.JSON(users)
}

func (m *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.UserPostRequest
	var ctx = c.Context()
	err := c.BodyParser(&params)
	if err != nil {
		return err
	}
	errorList := types.ValidateUserRequest(params)
	if len(errorList) > 0 {
		return c.JSON(errorList)
	}

	user, err := types.CreateUserFromRequest(&params)
	if err != nil {
		return err
	}
	result, err := m.userStore.PostUser(ctx, user)
	if err != nil {
		return err
	}
	return c.JSON(result)
}
