package routers

import (
	"g6/blog-api/Delivery/bootstrap"

	"g6/blog-api/Infrastructure/database/mongo"

	"github.com/gin-gonic/gin"
)

func NewUserRoutes(env *bootstrap.Env, group *gin.RouterGroup, db mongo.Database) {
	// ur := repositories.NewUserRepository(db, env.UserCollection)
	// uu := userUsecase.NewUserUsecase(ur, time.Duration(env.CtxTSeconds)*time.Second)
	// uc := controllers.NewUserController(uu)

	// group.POST("/users//", uc.CreateUser)

}
