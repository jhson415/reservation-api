package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jhson415/reservation-api/types"
)

func UserHandler(c *fiber.Ctx) error {
	user := types.User{
		FirstName: "Jayson",
		LastName:  "Son",
	}
	return c.JSON(user)
}
