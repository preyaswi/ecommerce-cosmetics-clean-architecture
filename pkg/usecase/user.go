package usecase

import (
	"clean/pkg/config"
	errorss "clean/pkg/errors"
	"clean/pkg/helper"
	interfaces "clean/pkg/repository/interface"
	"clean/pkg/utils/models"
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
)

type userUseCase struct {
	userRepo interfaces.UserRepository
	cartRepo interfaces.CartRepository
	cfg      config.Config
}


func NewUserUseCase(repo interfaces.UserRepository, cartRepo interfaces.CartRepository, cfg config.Config) *userUseCase {
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

func (u *userUseCase) GetAllAddress(userId int) (models.AddressInfoResponse, error) {
	addressInfo, err := u.userRepo.GetAllAddress(userId)
	if err != nil {
		return models.AddressInfoResponse{}, err
	}
	return addressInfo, nil

}
func (u *userUseCase) AddAddress(userId int, address models.AddressInfo) error {
	if err := u.userRepo.AddAddress(userId, address); err != nil {
		return err
	}

	return nil

}
func (u *userUseCase) UpdateAddress(address models.AddressInfo, addressID int, userID int) (models.AddressInfoResponse, error) {

	return u.userRepo.UpdateAddress(address, addressID, userID)

}
func (u *userUseCase) UserDetails(userID int) (models.UsersProfileDetails, error) {
	return u.userRepo.UserDetails(userID)
}

func (u *userUseCase) UpdateUserDetails(userDetails models.UsersProfileDetails, userID int) (models.UsersProfileDetails, error) {
	if !IsEmailValid(userDetails.Email) {
		return models.UsersProfileDetails{}, errors.New("invalid email format")
	}

	if !IsValidPhoneNumber(userDetails.Phone) {
		return models.UsersProfileDetails{}, errors.New("invalid phone number format")
	}
	userExist := u.userRepo.CheckUserAvailability(userDetails.Email)
	// update with email that does not already exist
	if userExist {
		return models.UsersProfileDetails{}, errors.New("user already exist, choose different email")
	}
	userExistByPhone, err := u.userRepo.CheckUserExistsByPhone(userDetails.Phone)

	if err != nil {
		return models.UsersProfileDetails{}, errors.New("error with server")
	}
	if userExistByPhone != nil {
		return models.UsersProfileDetails{}, errors.New("user with this phone is already exists")
	}
	// which all field are not empty (which are provided from the front end should be updated)
	if userDetails.Email != "" {
		u.userRepo.UpdateUserEmail(userDetails.Email, userID)
	}

	if userDetails.Firstname != "" {
		u.userRepo.UpdateFirstName(userDetails.Firstname, userID)
	}
	if userDetails.Lastname != "" {
		u.userRepo.UpdateLastName(userDetails.Lastname, userID)
	}

	if userDetails.Phone != "" {
		u.userRepo.UpdateUserPhone(userDetails.Phone, userID)
	}

	return u.userRepo.UserDetails(userID)

}

func (u *userUseCase) UpdatePassword(ctx context.Context, body models.UpdatePassword) error {
	var userID int
	var ok bool
	if userID, ok = ctx.Value("userID").(int); !ok {
		return errors.New("error retrieving user details")
	}
	fmt.Println("user id is", userID)
	userPassword, err := u.userRepo.UserPassword(userID)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(body.OldPassword))
	if err != nil {
		return errors.New("password incorrect")
	}
	if body.NewPassword != body.ConfirmNewPassword {
		return errors.New("password not matching")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), 10)
	if err != nil {
		return err
	}
	if err := u.userRepo.UpdateUserPassword(string(hashedPassword), userID); err != nil {
		return err
	}
	return nil
}
func (u *userUseCase) AddToWishlist(product_id int, user_id int) error {

	productExist, err := u.userRepo.CheckProductExist(product_id)
	if err != nil {
		return err
	}

	if !productExist {
		return errors.New("product does not exist")
	}

	productExistInWishList, err := u.userRepo.ProductExistInWishList(product_id, user_id)
	if err != nil {
		return err
	}
	if productExistInWishList {
		return errors.New("product already exist in wishlist")
	}

	err = u.userRepo.AddToWishList(user_id, product_id)
	if err != nil {
		return err
	}

	return nil
}
func (u *userUseCase) GetWishList(userID int) ([]models.WishListResponse, error) {

	wishList, err := u.userRepo.GetWishList(userID)
	if err != nil {
		return []models.WishListResponse{}, err
	}

	return wishList, err
}
func (u *userUseCase) RemoveFromWishlist(productId int, userID int) error {
	productExistInWishlist, err := u.userRepo.ProductExistInWishList(productId, userID)
	if err != nil {
		return err
	}
	if !productExistInWishlist {
		return errors.New("error deleting product doesnot exist in the wishlist")
	}

	err = u.userRepo.RemoveFromWishlist(userID, productId)
	if err != nil {
		return err
	}
	return nil
}
