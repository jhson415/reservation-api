package api

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jhson415/reservation-api/db"
)

type HotelHandler struct {
	ctx   context.Context
	store db.Store
}

func NewHotelHandler(store db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

func (m *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
	)

	hotel, err := m.store.Hotel.GetHotelList(ctx)
	if err != nil {
		return err
	}
	return c.JSON(hotel)
}

func (m *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	param := c.Params("id")
	hotel, err := m.store.Hotel.GetHotelById(c.Context(), param)
	if err != nil {
		return err
	}
	return c.JSON(hotel)

}

func (m *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	param := c.Params("id")
	rooms, err := m.store.Room.GetRoomListByHotelId(c.Context(), param)
	if err != nil {
		log.Fatalln("Error while calling db.GetRoomListByHotelId", err)
		return err
	}
	return c.JSON(rooms)
}
