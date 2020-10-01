package main

import (
	"fmt"
	"log"

	"localhost/medium-mongo-go-driver/config"
	"localhost/medium-mongo-go-driver/databases"
	"localhost/medium-mongo-go-driver/handlers"
	"localhost/medium-mongo-go-driver/middlewares"
	"localhost/medium-mongo-go-driver/router"
)

func main() {
	conf := config.GetConfig()
	dbClient, err := databases.NewClient(conf)

	if err != nil {
		panic(fmt.Errorf("Failed to create new database client: %s", err))
	}

	err = dbClient.Connect()

	if err != nil {
		panic(fmt.Errorf("Failed to connect to database: %s", err.Error()))
	}

	db := databases.NewDatabase(conf, dbClient)

	if err != nil {
		panic(fmt.Errorf("Failed to create new database: %s", err.Error()))
	}

	dbSession := middlewares.DatabaseSession{
		DB: db,
	}

	userDb := databases.NewUserDatabase(db)
	userHandlers := handlers.User{
		DB: userDb,
	}
	userHandlers2020 := handlers.User2020{
		DB: userDb,
	}

	log.Fatal(router.GetMainEngine(dbSession, userHandlers, userHandlers2020).Run(":3000"))

}
