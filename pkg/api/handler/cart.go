package handler

import (
	services "clean/pkg/usecase/interface"
	"clean/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	cartUseCase services.CartUseCase
}

func NewCartHandler(usecase services.CartUseCase) *CartHandler {
	return &CartHandler{
		cartUseCase: usecase,
	}
}
func(cr *CartHandler) AddToCart(c *gin.Context) {
	id := c.Param("id")

	product_id, err := strconv.Atoi(id)
	if err != nil {
		errResponse := response.ClientResponse(http.StatusBadGateway, "Prodcut id is given in the wrong format", nil, err.Error())
		c.JSON(http.StatusBadGateway, errResponse)
		return
	}
	user_ID, _ := c.Get("user_id")

	cartResponse, err := cr.cartUseCase.AddToCart(product_id, user_ID.(int))
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadGateway, "could not add product to the cart", nil, err.Error())
		c.JSON(http.StatusBadGateway, errRes)
		return
	}
	successRes := response.ClientResponse(200, "Added porduct Successfully to the cart", cartResponse, nil)
	c.JSON(200, successRes)

}
func(cr *CartHandler)  RemoveFromCart(c *gin.Context) {
	id := c.Param("id")
	product_id, err := strconv.Atoi(id)

	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "product not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	userID, _ := c.Get("user_id")

	// user_ID := c.Query("user_id")
	// user_id, _ := strconv.Atoi(user_ID)

	updatedCart, err := cr.cartUseCase.RemoveFromCart(product_id, userID.(int))
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadGateway, "cannot remove product from the cart", nil, err.Error())
		c.JSON(http.StatusBadGateway, errRes)
		return
	}
	succesRes := response.ClientResponse(200, "product removed successfully", updatedCart, nil)
	c.JSON(200, succesRes)

}
func (cr *CartHandler) DisplayCart(c *gin.Context) {

	userID, _ := c.Get("user_id")
	// user_ID := c.Query("user_id")
	// user_id, _ := strconv.Atoi(user_ID)
	cart, err := cr.cartUseCase.DisplayCart(userID.(int))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "cannot display cart", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Cart items displayed successfully", cart, nil)
	c.JSON(http.StatusOK, successRes)

}
func(cr *CartHandler)  EmptyCart(c *gin.Context) {
	userID, _ := c.Get("user_id")
	// user_ID := c.Query("user_id")
	// user_id, _ := strconv.Atoi(user_ID)
	cart, err := cr.cartUseCase.EmptyCart(userID.(int))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "cannot empty the cart", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Cart emptied successfully", cart, nil)
	c.JSON(http.StatusOK, successRes)

}
