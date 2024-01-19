package usecase

import (
	"clean/pkg/config"
	errorss "clean/pkg/errors"
	"clean/pkg/helper"
	interfaces "clean/pkg/repository/interface"
	"clean/pkg/utils/models"
	"errors"
	"regexp"
	"strconv"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
)

type userUseCase struct {
	userRepo interfaces.UserRepository
	cfg      config.Config
}

func NewUserUseCase(repo interfaces.UserRepository, cfg config.Config) *userUseCase {
	return &userUseCase{
		userRepo: repo,
		cfg:      cfg,
	}
}

func IsEmailValid(email string) bool {

	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(emailRegex, email)
	if match {
		return true
	} else {
		return false
	}
}
func IsValidPhoneNumber(phoneNumber string) bool {

	phoneRegex := `^[789]\d{9}$`
	match, _ := regexp.MatchString(phoneRegex, phoneNumber)
	if match {
		return true
	} else {
		return false
	}
}
func (u *userUseCase) UserSignup(user models.SignupDetail) (*models.TokenUser, error) {

	if !IsEmailValid(user.Email) {
		return &models.TokenUser{}, errors.New("invalid email format")
	}

	if !IsValidPhoneNumber(user.Phone) {
		return &models.TokenUser{}, errors.New("invalid phone number format")
	}
	//check whether the user already exsist by looking the email and the phone number provided
	email, err := u.userRepo.CheckUserExistsByEmail(user.Email)

	if err != nil {
		return &models.TokenUser{}, errors.New("error with server")
	}
	if email != nil {
		return &models.TokenUser{}, errorss.ErrEmailAlreadyExist
	}

	phone, err := u.userRepo.CheckUserExistsByPhone(user.Phone)

	if err != nil {
		return &models.TokenUser{}, errors.New("error with server")
	}
	if phone != nil {
		return &models.TokenUser{}, errors.New("user with this phone is already exists")
	}

	//if the signing up is a new user then hashing the password
	hashedPassword, err := helper.PasswordHashing(user.Password)
	if err != nil {
		return &models.TokenUser{}, errors.New("error in hashing password")
	}

	user.Password = hashedPassword
	//after hashing adding the user detail into the database and taking the added user detail to the userdata
	userData, err := u.userRepo.UserSignup(user)
	if err != nil {
		return &models.TokenUser{}, errors.New("could not add the user ")
	}

	// create referral code for the user and send in details of referred id of user if it exist
	id := uuid.New().ID()
	str := strconv.Itoa(int(id))
	userReferral := str[:8]
	err = u.userRepo.CreateReferralEntry(userData, userReferral)
	if err != nil {
		return &models.TokenUser{}, err
	}

	if user.ReferralCode != "" {
		// first check whether if a user with that referralCode exist
		referredUserId, err := u.userRepo.GetUserIdFromReferrals(user.ReferralCode)
		if err != nil {
			return &models.TokenUser{}, err
		}

		if referredUserId != 0 {
			referralAmount := 100
			err := u.userRepo.UpdateReferralAmount(float64(referralAmount), referredUserId, userData.Id)
			if err != nil {
				return &models.TokenUser{}, err
			}

		}
	}

	//creating a jwt token for the new user with the detail that has been stored in the database
	accessToken, err := helper.GenerateAccessToken(userData)
	if err != nil {
		return &models.TokenUser{}, errors.New("counldnt create access token due to error")
	}

	refreshToken, err := helper.GenerateRefreshToken(userData)
	if err != nil {
		return &models.TokenUser{}, errors.New("counldnt create refresh token due to error")
	}

	return &models.TokenUser{
		Users:        userData,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}
func (u *userUseCase) UserLoginWithPassword(user models.LoginDetail) (*models.TokenUser, error) {
	email, err := u.userRepo.CheckUserExistsByEmail(user.Email)

	if err != nil {
		return &models.TokenUser{}, errors.New("error with server")
	}
	if email == nil {
		return &models.TokenUser{}, errors.New("email  does not exsist")
	}

	userDetails, err := u.userRepo.FindUserDetailsByEmail(user)
	if err != nil {
		return &models.TokenUser{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userDetails.Password), []byte(user.Password))

	if err != nil {
		return &models.TokenUser{}, errors.New("password not matching")
	}
	var user_details models.SignupDetailResponse
	err = copier.Copy(&user_details, &userDetails)
	if err != nil {
		return &models.TokenUser{}, err
	}
	accessToken, err := helper.GenerateAccessToken(user_details)
	if err != nil {
		return &models.TokenUser{}, errors.New("could not create accesstoken due to internal error")
	}
	refreshToken, err := helper.GenerateRefreshToken(user_details)
	if err != nil {
		return &models.TokenUser{}, errors.New("counldnt create refresh token due to error")
	}

	return &models.TokenUser{
		Users:        user_details,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}
