package main

import (
	"fmt"
	"g6/blog-api/Delivery/bootstrap"
	"g6/blog-api/Delivery/routers"
)

func main() {
	app := bootstrap.App()
	env := app.Env
	db := app.Mongo.Database(env.DBName)
	defer app.CloseDBConnection()

	router := routers.SetupRoutes(env, db)
	if err := router.Run(env.ServerAddress); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		return
	}
}
