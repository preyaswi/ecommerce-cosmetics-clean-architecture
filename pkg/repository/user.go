package repository

import (
	"clean/pkg/domain"
	interfaces "clean/pkg/repository/interface"
	"clean/pkg/utils/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type userDatabase struct {
	DB *gorm.DB
}
func NewuserRepository(DB *gorm.DB)interfaces.UserRepository {
	return &userDatabase{DB}
	
}
func(c *userDatabase) CheckUserExistsByEmail(email string) (*domain.User, error) {
	var user domain.User
	result := c.DB.Where(&domain.User{Email: email}).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}
func(c *userDatabase) CheckUserExistsByPhone(phone string) (*domain.User, error) {
	var user domain.User
	result := c.DB.Where(&domain.User{Phone: phone}).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}
func(c *userDatabase) UserSignup(user models.SignupDetail) (models.SignupDetailResponse, error) {
	var signupDetail models.SignupDetailResponse
	err := c.DB.Raw("INSERT INTO users(firstname,lastname,email,password,phone)VALUES(?,?,?,?,?)RETURNING id,firstname,lastname,email,phone", user.Firstname, user.Lastname, user.Email, user.Password, user.Phone).Scan(&signupDetail).Error
	if err != nil {
		fmt.Println("Repository error:", err)
		return models.SignupDetailResponse{}, err
	}
	return signupDetail, nil

}
func(c *userDatabase) CreateReferralEntry(userDetails models.SignupDetailResponse, userReferral string) error {

	err := c.DB.Exec("insert into referrals (user_id,referral_code,referral_amount) values (?,?,?)", userDetails.Id, userReferral, 0).Error
	if err != nil {
		return err
	}

	return nil

}
func(c *userDatabase) GetUserIdFromReferrals(ReferralCode string) (int, error) {

	var referredUserId int
	err := c.DB.Raw("select user_id from referrals where referral_code = ?", ReferralCode).Scan(&referredUserId).Error
	if err != nil {
		return 0, nil
	}

	return referredUserId, nil
}

func(c *userDatabase) UpdateReferralAmount(referralAmount float64, referredUserId int, currentUserID int) error {

	err := c.DB.Exec("update referrals set referral_amount = ?,referred_user_id = ? where user_id = ? ", referralAmount, referredUserId, currentUserID).Error
	if err != nil {
		return err
	}

	// find the current amount in referred users referral table and add 100 with that
	err = c.DB.Exec("update referrals set referral_amount = referral_amount + ? where user_id = ? ", referralAmount, referredUserId).Error
	if err != nil {
		return err
	}

	return nil

}
