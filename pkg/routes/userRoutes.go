package routes

import (
	"clean/pkg/api/handler"
	"clean/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.RouterGroup,userHandler *handler.UserHandler,
	otpHandler *handler.OtpHandler,productHandler *handler.ProductHandler,
	cartHandler *handler.CartHandler,orderHandler *handler.OrderHandler)  {
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
	router.Use(middleware.AuthMiddleware())
	{
		address := router.Group("/address")
		{
			address.GET("", userHandler.GetAllAddress)
			address.POST("", userHandler.AddAddress)
			address.PUT("/:id", userHandler.UpdateAddress)

		}
		users := router.Group("/users")
		{

			users.GET("", userHandler.UserDetails)
			users.PUT("", userHandler.UpdateUserDetails)
			users.PUT("/update-password", userHandler.UpdatePassword)
		}
		//wishlist
		wishlist := router.Group("/wishlist")
		{

			wishlist.POST("/:id", userHandler.AddWishList)
			wishlist.GET("", userHandler.GetWishList)
			wishlist.DELETE("/:id", userHandler.RemoveFromWishlist)
		}

		//cart
		cart := router.Group("/cart")
		{
			cart.POST("/:id", cartHandler.AddToCart)
			cart.DELETE("/:id", cartHandler.RemoveFromCart)
			cart.GET("", cartHandler.DisplayCart)
			cart.DELETE("", cartHandler.EmptyCart)
		}

		//order
		order := router.Group("/order")
		{

			order.POST("", orderHandler.OrderItemsFromCart)
			order.GET("", orderHandler.GetOrderDetails)
			order.GET("/:page", orderHandler.GetOrderDetails)
			order.PUT("/:id", orderHandler.CancelOrder)
		}
				
		
	}
}
