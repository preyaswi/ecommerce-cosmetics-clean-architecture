package handler

import (
	errorss "clean/pkg/errors"
	services "clean/pkg/usecase/interface"
	"clean/pkg/utils/models"
	"clean/pkg/utils/response"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OtpHandler struct {
	otpUseCase services.OtpUseCase
}

func NewOtpHandler(usecsase services.OtpUseCase) *OtpHandler {
	return &OtpHandler{
		otpUseCase: usecsase,
	}
}
func (ot *OtpHandler) SendOTP(c *gin.Context) {
	var phone models.OTPData

	if err := c.ShouldBindJSON(&phone); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
	}
	err := ot.otpUseCase.SendOTP(phone.PhoneNumber)

	if err != nil {

		errorRes := response.ClientResponse(http.StatusInternalServerError, "Could not send OTP", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		fmt.Println()
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "OTP sent successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)

}
func(ot *OtpHandler) VerifyOTP(c *gin.Context) {

	var code models.VerifyData

	if err := c.ShouldBindJSON(&code); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	users, err := ot.otpUseCase.VerifyOTP(code)

	if err != nil {
		if errors.Is(err, errorss.ErrFailedTovalidateOtp) {
			errorRes := response.ClientResponse(http.StatusForbidden, "failed to verify OTP", nil, err.Error())
			c.JSON(http.StatusForbidden, errorRes)
			return
		}
		errorRes := response.ClientResponse(http.StatusInternalServerError, "Could not verify OTP", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errorRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Successfully verified OTP", users, nil)
	c.JSON(http.StatusOK, successRes)

}