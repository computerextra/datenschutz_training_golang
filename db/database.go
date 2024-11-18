package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/computerextra/datenschutz_training_golang/ent"
	"github.com/computerextra/datenschutz_training_golang/ent/car"
	"github.com/computerextra/datenschutz_training_golang/ent/user"
	_ "github.com/mattn/go-sqlite3"
)

func GetDb() *ent.Client {
	// Dev Mode (Datebase in Memory)
	client, err := ent.Open(
		"sqlite3",
		"file:ent?mode=memory&cache=shared&_fk=1",
	)

	// Production Mode
	// client, err := ent.Open(
	// 	"sqlite3",
	//	"file:file.db?mode=rwc&_fk=1",
	// )

	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()

	// run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	return client
}

func CreateUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.User.Create().SetAge(30).SetName("a8m").Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed createing user %w", err)
	}
	log.Println("user was created: ", u)
	return u, nil
}

func QueryUser(ctx context.Context, client *ent.Client, username string) (*ent.User, error) {
	// 'Only' fails if no user found, or more than one user returned
	u, err := client.User.Query().Where(user.Name(username)).Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("files querying user: %w", err)
	}
	log.Println("user returned: ", u)
	return u, nil
}

func CreateCars(ctx context.Context, client *ent.Client) (*ent.User, error) {
	tesla, err := client.Car.
		Create().
		SetModel("Tesla").
		SetRegisteredAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating car: %w", err)
	}
	log.Println("car was created: ", tesla)

	ford, err := client.Car.
		Create().
		SetModel("Ford").
		SetRegisteredAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating car: %w", err)
	}
	log.Println("car was created: ", ford)

	a8m, err := client.User.
		Create().
		SetAge(30).
		SetName("a8m").
		AddCars(tesla, ford).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user with cars: %w", err)
	}
	log.Println("user with cars was created: ", a8m)
	return a8m, nil
}

func QueryCars(ctx context.Context, a8m *ent.User) error {
	cars, err := a8m.QueryCars().All(ctx)
	if err != nil {
		return fmt.Errorf("failed querying cars: %w", err)
	}
	log.Println("cars of user a8m: ", cars)

	ford, err := a8m.QueryCars().Where(car.Model("Ford")).Only(ctx)
	if err != nil {
		return fmt.Errorf("failed querying car: %w", err)
	}
	log.Println("user a8m has a ford: ", ford)
	return nil
}

func QueryCarUsers(ctx context.Context, a8m *ent.User) error {
	cars, err := a8m.QueryCars().All(ctx)
	if err != nil {
		return fmt.Errorf("failed querying cars: %w", err)
	}

	// Query the inverse edge.
	for _, c := range cars {
		owner, err := c.QueryOwner().Only(ctx)
		if err != nil {
			return fmt.Errorf("failed querying owner: %w", err)
		}
		log.Println("car is owned by: ", owner)
	}

	return nil
}
