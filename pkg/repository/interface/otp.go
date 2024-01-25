package interfaces

import "clean/pkg/utils/models"

type OtpRepository interface {
	FindUserByMobileNumber(phone string) bool
	UserDetailsUsingPhone(phone string) (models.SignupDetailResponse, error)
	FindUserByEmail(email string) (bool, error)
	GetUserPhoneByEmail(email string) (string, error)
}
