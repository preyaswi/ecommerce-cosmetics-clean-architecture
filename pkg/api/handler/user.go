package handler

import (
	errorss "clean/pkg/errors"
	services "clean/pkg/usecase/interface"
	"clean/pkg/utils/models"
	"clean/pkg/utils/response"
	"errors"
	"net/http"

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
