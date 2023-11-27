package api

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jhson415/reservation-api/db"
	"github.com/jhson415/reservation-api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserHandler struct {
	ctx   context.Context
	store db.Store
}

func NewUserHandler(store db.Store) *UserHandler {
	return &UserHandler{
		store: store,
	}
}

func (m *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	var (
		id  = c.Params("id")
		ctx = c.Context()
	)

	user, err := m.store.User.GetUserById(ctx, id)
	if err != nil {
		return err
	}
	return c.JSON(user)
}

func (m *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
	)
	users, err := m.store.User.GetUserList(ctx)
	if err != nil {
		return err
	}
	return c.JSON(users)
}

func (m *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.UserPostParams
	var ctx = c.Context()
	err := c.BodyParser(&params)
	if err != nil {
		return err
	}
	errorList := types.ValidateUserParams(params)
	if len(errorList) > 0 {
		return c.JSON(errorList)
	}

	user, err := types.CreateUserFromParams(&params)
	if err != nil {
		return err
	}
	result, err := m.store.User.PostUser(ctx, user)
	if err != nil {
		return err
	}
	return c.JSON(result)
}

func (m *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	var (
		params = c.Params("id")
		ctx    = c.Context()
	)

	if err := m.store.User.DeleteUser(ctx, params); err != nil {
		return err
	}
	response := map[string]string{"status": "User Deleted"}
	return c.JSON(response)
}

func (m *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	var (
		ctx    = c.Context()
		params types.UpdateUserParams
		userID = c.Params("id")
	)
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	filter := bson.M{"_id": oid}
	if err := m.store.User.PutUser(ctx, filter, params); err != nil {
		return err
	}
	return c.JSON(map[string]string{"updated": userID})
}
