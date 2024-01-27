package handler

import (
	errorss "clean/pkg/errors"
	services "clean/pkg/usecase/interface"
	"clean/pkg/utils/models"
	"clean/pkg/utils/response"
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type UserHandler struct {
	userUseCase services.UserUseCase
}

func NewUserHandler(usecase services.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: usecase,
	}
}
func (u *UserHandler) UserSignUp(c *gin.Context) {
	var userSignup models.SignupDetail
	if err := c.ShouldBindJSON(&userSignup); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	// checking whether the data sent by the user has all the correct constraints specified by Users struct
	err := validator.New().Struct(userSignup)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "constraints not satisfied", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	//creating a newuser signup with the given deatil passing into the bussiness logic layer

	userCreated, err := u.userUseCase.UserSignup(userSignup)
	if err != nil {
		if errors.Is(err, errorss.ErrEmailAlreadyExist) {
			errRes := response.ClientResponse(http.StatusForbidden, "email already exist", nil, err.Error())
			c.JSON(http.StatusForbidden, errRes)
			return
		}
		if errors.Is(err, errorss.ErrPhoneAlreadyExist) {
			errRes := response.ClientResponse(http.StatusForbidden, "phonenumber already exist", nil, err.Error())
			c.JSON(http.StatusForbidden, errRes)
			return
		}
		errRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong formaaaaat", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	successRes := response.ClientResponse(http.StatusCreated, "User successfully signed up", userCreated, nil)
	c.JSON(http.StatusCreated, successRes)
}

func (u *UserHandler) UserLoginWithPassword(c *gin.Context) {
	var userLoginDetail models.LoginDetail
	if err := c.ShouldBindJSON(&userLoginDetail); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	err := validator.New().Struct(userLoginDetail)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "constrains not satisfied", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	userLoggedInWithPassword, err := u.userUseCase.UserLoginWithPassword(userLoginDetail)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	successRes := response.ClientResponse(http.StatusCreated, "User successfully Logged In With password", userLoggedInWithPassword, nil)
	c.JSON(http.StatusCreated, successRes)

}
func (u *UserHandler) GetAllAddress(c *gin.Context) {

	userID, _ := c.Get("user_id")

	// userId := c.Param("user_id")
	// userID, err := strconv.Atoi(userId)

	addressInfo, err := u.userUseCase.GetAllAddress(userID.(int))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "failed to retrieve details", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "User Address", addressInfo, nil)
	c.JSON(http.StatusOK, successRes)

}
func (u *UserHandler) AddAddress(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var address models.AddressInfo

	if err := c.ShouldBindJSON(&address); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	err := validator.New().Struct(address)

	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "constraints does not match", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}
	if err := u.userUseCase.AddAddress(userID.(int), address); err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "error in adding address", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return

	}

	SuccessRes := response.ClientResponse(200, "address added successfully", nil, nil)

	c.JSON(200, SuccessRes)

}
func (u *UserHandler) UpdateAddress(c *gin.Context) {

	id := c.Param("id")
	addressId, err := strconv.Atoi(id)

	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "address id not in the right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	userID, _ := c.Get("user_id")
	user_id := userID.(int)

	var address models.AddressInfo

	if err := c.ShouldBindJSON(&address); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	updatedAddress, err := u.userUseCase.UpdateAddress(address, addressId, user_id)

	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "failed updating address", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusCreated, "address updated successfully", updatedAddress, nil)
	c.JSON(http.StatusCreated, successRes)

}
func (u *UserHandler) UserDetails(c *gin.Context) {

	userID, _ := c.Get("user_id")

	userDetails, err := u.userUseCase.UserDetails(userID.(int))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "failed to retrieve details", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "user Details", userDetails, nil)
	c.JSON(http.StatusOK, successRes)

}
func (u *UserHandler) UpdateUserDetails(c *gin.Context) {

	user_id, _ := c.Get("user_id")

	var user models.UsersProfileDetails

	if err := c.ShouldBindJSON(&user); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	updatedDetails, err := u.userUseCase.UpdateUserDetails(user, user_id.(int))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "failed update user", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Updated User Details", updatedDetails, nil)
	c.JSON(http.StatusOK, successRes)

}
func (u *UserHandler) UpdatePassword(c *gin.Context) {
	user_id, _ := c.Get("user_id")
	ctx := context.Background()
	ctx = context.WithValue(ctx, "userID", user_id.(int))
	var body models.UpdatePassword

	if err := c.ShouldBindJSON(&body); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	err := u.userUseCase.UpdatePassword(ctx, body)

	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "failed updating password", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusCreated, "Password updated successfully", nil, nil)
	c.JSON(http.StatusCreated, successRes)
}
func (u *UserHandler)AddWishList(c *gin.Context) {
	userID, _ := c.Get("user_id")
	productId := c.Param("id")
	product_id, err := strconv.Atoi(productId)

	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "product id is in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	err = u.userUseCase.AddToWishlist(product_id, userID.(int))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "failed to add product on wishlist", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return

	}

	successRes := response.ClientResponse(http.StatusOK, "SuccessFully product added to wishlist", nil, nil)
	c.JSON(http.StatusOK, successRes)
}
func(u *UserHandler) GetWishList(c *gin.Context) {

	userID, _ := c.Get("user_id")
	wishList, err := u.userUseCase.GetWishList(userID.(int))

	if err != nil {
		errorRes := response.ClientResponse(http.StatusInternalServerError, "failed to retrieve wishlist detailss", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "SuccessFully retrieved wishlist", wishList, nil)
	c.JSON(http.StatusOK, successRes)

}
func(u *UserHandler) RemoveFromWishlist(c *gin.Context) {
	userId, _ := c.Get("user_id")
	productId := c.Param("id")
	product_id, err := strconv.Atoi(productId)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "product id is in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	err = u.userUseCase.RemoveFromWishlist(product_id, userId.(int))
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "product cannot remove from the wishlist", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "SuccessFully removed product from the wishlist", nil, nil)
	c.JSON(http.StatusOK, successRes)
}
