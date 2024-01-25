package repository

import (
	interfaces "clean/pkg/repository/interface"
	"clean/pkg/utils/models"

	"gorm.io/gorm"
)

type otpRepository struct {
	DB *gorm.DB
}

func NewOtprepository(DB *gorm.DB) interfaces.OtpRepository {
	return &otpRepository{
		DB: DB,
	}
}

func (ot *otpRepository) FindUserByMobileNumber(phone string) bool {
	var count int
	if err := ot.DB.Raw("select count(*) from users where phone = ?", phone).Scan(&count).Error; err != nil {
		return false
	}

	return count > 0

}
func (ot *otpRepository) UserDetailsUsingPhone(phone string) (models.SignupDetailResponse, error) {

	var usersDetails models.SignupDetailResponse
	if err := ot.DB.Raw("select * from users where phone = ?", phone).Scan(&usersDetails).Error; err != nil {
		return models.SignupDetailResponse{}, err
	}

	return usersDetails, nil

}
func (ot *otpRepository) FindUserByEmail(email string) (bool, error) {

	var count int
	if err := ot.DB.Raw("select count(*) from users where email = ?", email).Scan(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
func (ot *otpRepository) GetUserPhoneByEmail(email string) (string, error) {
	var phone string
	if err := ot.DB.Raw("select phone from users where email = ?", email).Scan(&phone).Error; err != nil {
		return "", err
	}

	return phone, nil

}
