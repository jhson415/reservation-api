package main

import (
	"flag"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jhson415/reservation-api/api"
)

func main() {
	listenAddr := flag.String("port", ":5500", "Input the port to change it")
	flag.Parse()

	app := fiber.New()
	apiV1 := app.Group("/api/v1")

	apiV1.Get("/user", api.UserHandler)
	app.Listen(*listenAddr)
	fmt.Println("Close!")
}

func handlerFoo(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"msg": "this is the landing!"})
}
