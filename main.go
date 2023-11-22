// TODO Create update user func

package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jhson415/reservation-api/api"
	"github.com/jhson415/reservation-api/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	const dbUri = "mongodb://localhost:27017"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbUri))
	if err != nil {
		log.Fatal(err)
	}

	var config = fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.JSON(map[string]string{"error": err.Error()})
		},
	}

	listenAddr := flag.String("port", ":5500", "Input the port to change it")
	flag.Parse()

	app := fiber.New(config)

	userHandler := api.NewUserHandler(db.NewMongoUserStore(client))

	apiV1 := app.Group("/api/v1")
	apiV1.Put("/user/:id", userHandler.HandlePutUser)
	apiV1.Post("/user/", userHandler.HandlePostUser)
	apiV1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiV1.Get("/user/", userHandler.HandleGetUsers)
	apiV1.Get("/user/:id", userHandler.HandleGetUser)
	app.Listen(*listenAddr)
	fmt.Println("Close!")
}

func handlerFoo(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"msg": "this is the landing!"})
}
