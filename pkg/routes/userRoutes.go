package routes

import (
	"clean/pkg/api/handler"

	"github.com/gin-gonic/gin"
)

func UserRoutes(engine *gin.RouterGroup,userHandler *handler.UserHandler)  {
	engine.POST("/signup",userHandler.UserSignUp)
	engine.POST("/login-with-password",userHandler.UserLoginWithPassword)
}