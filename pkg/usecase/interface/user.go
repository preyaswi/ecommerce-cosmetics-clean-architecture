package interfaces

import "clean/pkg/utils/models"
type UserUseCase interface{
	UserSignup(user models.SignupDetail) (*models.TokenUser, error)
}