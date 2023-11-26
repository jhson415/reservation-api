package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jhson415/reservation-api/db"
	"github.com/jhson415/reservation-api/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testDb struct {
	db.UserStore
}

func (tdb *testDb) dropDb(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Error(err)
	}
}

func setup(t *testing.T) *testDb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DB_URI))
	if err != nil {
		t.Fatal(err)
	}
	return &testDb{
		UserStore: db.NewMongoUserStore(client),
	}

}

func CreateUserRequest(t *testing.T) (*http.Request, types.UserPostParams) {
	// Create user Params
	newUser := types.UserPostParams{
		FirstName: "Jayson",
		LastName:  "Son",
		Email:     "eeee@asd.com",
		Password:  "123123123123",
	}
	b, err := json.Marshal(newUser)
	if err != nil {
		t.Error(err)
	}

	// Make a request
	req := httptest.NewRequest("POST", "/user", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	return req, newUser
}

func TestPostUser(t *testing.T) {
	// Setting up basic vars
	tdb := setup(t)
	defer tdb.dropDb(t)

	// Start Fiber server
	app := fiber.New()
	userHandler := NewUserHandler(tdb)
	app.Post("/user/", userHandler.HandlePostUser)

	req, newUser := CreateUserRequest(t)
	resp, err := app.Test(req)
	if err != nil {

		t.Error(err)
	}

	// Decode Response
	var createdUser types.User
	if err = json.NewDecoder(resp.Body).Decode(&createdUser); err != nil {
		t.Log("Possible that the user output is error, returns an array type")
		t.Error(err)
	}

	// Check Result
	if resp.StatusCode != 200 {
		t.Errorf("Status code is not 200, given code ->%d", resp.StatusCode)
	}
	if len(createdUser.EncryptedPassword) < 0 {
		t.Errorf("EncryptedPassword is not part of the result")
	}
	if createdUser.Email != newUser.Email {
		t.Errorf("Email does not match")
	}
	if createdUser.FirstName != newUser.FirstName {
		t.Errorf("FirstName does not match")
	}
	if createdUser.LastName != newUser.LastName {
		t.Errorf("LastName does not match")
	}
	t.Log(createdUser.ID)

}

func TestGetUsers(t *testing.T) {
	// Setup DB
	tdb := setup(t)
	defer tdb.dropDb(t)

	// Setup Fiber
	app := fiber.New()
	userHandler := NewUserHandler(tdb)
	app.Get("/users", userHandler.HandleGetUsers)
	app.Post("/user", userHandler.HandlePostUser)

	req, _ := CreateUserRequest(t)
	_, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	t.Log("")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	t.Log("")

	// Get ID
	var CreatedUser types.User
	if err = json.NewDecoder(resp.Body).Decode(&CreatedUser); err != nil {
		t.Error(err)
	}
	t.Log(CreatedUser.ID)

	req = httptest.NewRequest("GET", "/users", nil)
	req.Header.Add("Content-Type", "application/json")
	resp, err = app.Test(req)
	if err != nil {
		t.Error(err)
	}

	t.Log("")

	var getAllUserResult []types.User
	if err = json.NewDecoder(resp.Body).Decode(&getAllUserResult); err != nil {
		t.Error(err)
	}

	if len(getAllUserResult) < 1 {
		t.Errorf("There is no object found need to have at lest one result")
	}

	t.Log(getAllUserResult)

}

func TestGetUser(t *testing.T) {
	// Setup DB
	tdb := setup(t)
	defer tdb.dropDb(t)

	// Open Fiber
	app := fiber.New()
	userHandler := NewUserHandler(tdb)
	app.Post("/user/", userHandler.HandlePostUser)
	app.Get("/user/:id", userHandler.HandleGetUser)

	// Create user for
	req, _ := CreateUserRequest(t)
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	t.Log()

	var createdUser types.User

	if err = json.NewDecoder(resp.Body).Decode(&createdUser); err != nil {
		t.Error(err)
	}
	requestTarget := fmt.Sprintf("/user/%v", createdUser.ID.Hex())

	req = httptest.NewRequest("GET", requestTarget, nil)
	resp, err = app.Test(req)
	if err != nil {
		t.Error(err)
	}

	var SearchedUser types.User
	if err = json.NewDecoder(resp.Body).Decode(&SearchedUser); err != nil {
		t.Error(err)
	}
	if SearchedUser.Email != createdUser.Email {
		t.Errorf("The Email does not match Created Email : %v, Searched Email : %v", createdUser.Email, SearchedUser.Email)
	}
	if SearchedUser.FirstName != createdUser.FirstName {
		t.Errorf("The FirstName does not match Created FirstName : %v, Searched FirstName : %v", createdUser.FirstName, SearchedUser.FirstName)
	}
	if SearchedUser.LastName != createdUser.LastName {
		t.Errorf("The LastName does not match Created LastName : %v, Searched LastName : %v", createdUser.LastName, SearchedUser.LastName)
	}

}

func TestDeleteUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.dropDb(t)

	userHandler := NewUserHandler(tdb)
	app := fiber.New()
	app.Post("/user", userHandler.HandlePostUser)
	app.Delete("/user/:id", userHandler.HandleDeleteUser)

	req, _ := CreateUserRequest(t)
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	var createdUser types.User
	err = json.NewDecoder(resp.Body).Decode(&createdUser)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(createdUser.ID.Hex())
	// Delete User
	requestTarget := fmt.Sprintf("/user/%v", createdUser.ID.Hex())
	req = httptest.NewRequest("DELETE", requestTarget, nil)
	resp, err = app.Test(req)
	if err != nil {
		t.Error(err)
	}

	// Check Result
	var deleteResult map[string]string
	err = json.NewDecoder(resp.Body).Decode(&deleteResult)
	if err != nil {
		t.Error(err)
	}
	if deleteResult["status"] != "User Deleted" {
		t.Errorf("Delete Result is not correct")
	}

	if resp.StatusCode != 200 {
		t.Errorf("Status code is not 200, given code ->%d", resp.StatusCode)
	}

}

func TestPutUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.dropDb(t)

	userHandler := NewUserHandler(tdb)
	app := fiber.New()
	app.Post("/user", userHandler.HandlePostUser)
	app.Put("/user/:id", userHandler.HandlePutUser)

	req, _ := CreateUserRequest(t)
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	var createdUser types.User
	err = json.NewDecoder(resp.Body).Decode(&createdUser)
	if err != nil {
		t.Error(err)
	}

	// Change User Params
	createdUser.FirstName = "Changed"
	createdUser.LastName = "Changed"

	// Make a request
	b, err := json.Marshal(createdUser)
	if err != nil {
		t.Error(err)
	}
	httpReq := httptest.NewRequest("PUT", fmt.Sprintf("/user/%v", createdUser.ID.Hex()), bytes.NewReader(b))
	httpReq.Header.Add("Content-Type", "application/json")

	// Make a response
	resp, err = app.Test(httpReq)
	if err != nil {
		t.Error(err)
	}

	// Check Result
	fmt.Println(resp.StatusCode)
	var updatedUser map[string]string
	err = json.NewDecoder(resp.Body).Decode(&updatedUser)
	if err != nil {
		t.Error(err)
	}
	if updatedUser["updated"] != createdUser.ID.Hex() {
		t.Errorf("Updated User ID does not match")
	}
}
