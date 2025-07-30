package bootstrap

import (
	"g6/blog-api/Configs"
	"g6/blog-api/Infrastructure/database"

	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	Env   *Configs.Env
	Mongo *mongo.Client
}

func App() Application {
	env, err := Configs.NewEnv(".env")
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}
	app := &Application{}
	app.Env = env
	app.Mongo = database.NewMongoDatabase(app.Env)
	return *app
}

func (app *Application) CloseDBConnection() {
	database.CloseMongoDBConnection(app.Mongo)
}
