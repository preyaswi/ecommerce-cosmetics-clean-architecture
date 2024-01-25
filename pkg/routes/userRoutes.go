package routes

import (
	"clean/pkg/api/handler"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.RouterGroup,userHandler *handler.UserHandler,otpHandler *handler.OtpHandler,productHandler *handler.ProductHandler)  {
	router.POST("/signup",userHandler.UserSignUp)
	router.POST("/login-with-password",userHandler.UserLoginWithPassword)

	router.POST("/send-otp",otpHandler.SendOTP)
	router.POST("verify-otp",otpHandler.VerifyOTP)
	
	products := router.Group("/products")
	{
		products.GET("", productHandler.ShowAllProducts)
		products.GET("/page/:page", productHandler.ShowAllProducts) //TO ARRANGE PAGE WITH COUNT
		products.GET("/:id", productHandler.ShowIndividualProducts)
		products.POST("/filter", productHandler.FilterCategory)

	}
}
