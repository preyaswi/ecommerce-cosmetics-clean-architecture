package interfaces

import (
	"clean/pkg/domain"
	"clean/pkg/utils/models"
)

type UserRepository interface {
	CheckUserExistsByEmail(email string) (*domain.User, error)
	CheckUserExistsByPhone(phone string) (*domain.User, error)
	UserSignup(user models.SignupDetail) (models.SignupDetailResponse, error)
	CreateReferralEntry(userDetails models.SignupDetailResponse, userReferral string) error 
	GetUserIdFromReferrals(ReferralCode string) (int, error)
	UpdateReferralAmount(referralAmount float64, referredUserId int, currentUserID int) error
	FindUserDetailsByEmail(user models.LoginDetail) (models.UserLoginResponse, error)
	GetAllAddress(userId int) (models.AddressInfoResponse, error)
	AddAddress(userId int, address models.AddressInfo) error
	UpdateAddress(address models.AddressInfo, addressID int, userID int) (models.AddressInfoResponse, error)
	UserDetails(userID int) (models.UsersProfileDetails, error)
	CheckUserAvailability(email string) bool
	UpdateUserEmail(email string, userID int) error
	UpdateUserPhone(phone string, userID int) error
	UpdateFirstName(name string, userID int) error
	UpdateLastName(name string, userID int) error
	UserPassword(userID int) (string, error)
	UpdateUserPassword(password string, userID int) error
	CheckProductExist(pid int) (bool, error)
	ProductExistInWishList(productID int, userID int) (bool, error)
	AddToWishList(userID int, productID int) error
	GetWishList(userId int) ([]models.WishListResponse, error)
	RemoveFromWishlist(userID int, productId int) error
	GetAllAddresses(userID int) ([]models.AddressInfoResponse, error)
	GetAllPaymentOption() ([]models.PaymentDetails, error)
}
