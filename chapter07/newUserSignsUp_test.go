package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter07/db"
	"github.com/PacktPublishing/Test-Driven-Development-in-Go/chapter07/handlers"
	"github.com/cucumber/godog"
)

// contextKey is used to pass information between test steps.
type contextKey struct {
	UsersURL string
	User     db.User
}

func theBookSwapAppIsUp(ctx context.Context) (context.Context, error) {
	url, err := getTestURL()
	if err != nil {
		return ctx, fmt.Errorf("incorrect config:%v", err)
	}
	resp, err := http.Get(*url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return ctx, fmt.Errorf("bookswap not up:%v", err)
	}

	return context.WithValue(ctx, contextKey{}, contextKey{
		UsersURL: *url + "/users",
	}), nil
}

func userDetails(ctx context.Context) (context.Context, error) {
	config, ok := ctx.Value(contextKey{}).(contextKey)
	if !ok {
		return ctx, errors.New("config missing")
	}
	config.User = db.User{
		Name:     "New GoDog User",
		Address:  "1 London Road",
		PostCode: "N1",
		Country:  "United Kingdom",
	}

	return context.WithValue(ctx, contextKey{}, config), nil
}

func sentToTheUsersEndpoint(ctx context.Context) (context.Context, error) {
	config, ok := ctx.Value(contextKey{}).(contextKey)
	if !ok {
		return ctx, errors.New("config missing")
	}
	userPayload, err := json.Marshal(config.User)
	if err != nil {
		return ctx, fmt.Errorf("error marshalling user:%v", err)
	}
	r, err := http.Post(config.UsersURL, "application/json", bytes.NewBuffer(userPayload))
	if err != nil || r.StatusCode != http.StatusOK {
		return ctx, fmt.Errorf("error creating user :%v", err)
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return ctx, fmt.Errorf("error reading body:%v", err)
	}
	r.Body.Close()
	var resp handlers.Response
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return ctx, fmt.Errorf("error unmarshalling body:%v", err)
	}
	if resp.User == nil {
		return ctx, errors.New("no user in users reponse")
	}
	config.User.ID = resp.User.ID
	if *resp.User != config.User {
		return ctx, fmt.Errorf("returned user not as expected:got %v, want %v", *resp.User, config.User)
	}

	return context.WithValue(ctx, contextKey{}, config), nil
}

func aNewUserProfileIsCreated(ctx context.Context) (context.Context, error) {
	config, ok := ctx.Value(contextKey{}).(contextKey)
	if !ok {
		return ctx, errors.New("config missing")
	}
	r, err := http.Get(fmt.Sprintf("%s/%s", config.UsersURL, config.User.ID))
	if err != nil || r.StatusCode != http.StatusOK {
		return ctx, fmt.Errorf("error getting user :%v", err)
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return ctx, fmt.Errorf("error reading body:%v", err)
	}
	r.Body.Close()
	var resp handlers.Response
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return ctx, fmt.Errorf("error unmarshalling body:%v", err)
	}
	if resp.User == nil {
		return ctx, errors.New("no user in users reponse")
	}
	if *resp.User != config.User {
		return ctx, fmt.Errorf("returned user not as expected:got %v, want %v", *resp.User, config.User)
	}

	return context.WithValue(ctx, contextKey{}, config), nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^the BookSwap app is up$`, theBookSwapAppIsUp)
	ctx.Step(`^user details$`, userDetails)
	ctx.Step(`^sent to the users endpoint$`, sentToTheUsersEndpoint)
	ctx.Step(`^a new user profile is created$`, aNewUserProfileIsCreated)
}

func getTestURL() (*string, error) {
	baseURL, ok := os.LookupEnv("BOOKSWAP_BASE_URL")
	if !ok {
		return nil, errors.New("$BOOKSWAP_BASE_URL not found")
	}
	port, ok := os.LookupEnv("BOOKSWAP_PORT")
	if !ok {
		return nil, errors.New("$BOOKSWAP_PORT not found")
	}
	url := fmt.Sprintf("%s:%s", baseURL, port)

	return &url, nil
}
