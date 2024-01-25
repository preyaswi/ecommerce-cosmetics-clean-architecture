package interfaces

import "clean/pkg/utils/models"

type OtpUseCase interface {
	SendOTP(phone string) error
	VerifyOTP(code models.VerifyData) (models.TokenUser, error)
}
