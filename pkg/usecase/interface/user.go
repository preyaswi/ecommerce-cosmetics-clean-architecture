package interfaces

import (
	"clean/pkg/utils/models"
	"context"
)

type UserUseCase interface {
	UserSignup(user models.SignupDetail) (*models.TokenUser, error)
	UserLoginWithPassword(user models.LoginDetail) (*models.TokenUser, error)
	GetAllAddress(userId int) (models.AddressInfoResponse, error)
	AddAddress(userId int, address models.AddressInfo) error
	UpdateAddress(address models.AddressInfo, addressID int, userID int) (models.AddressInfoResponse, error)
	UserDetails(userID int) (models.UsersProfileDetails, error)
	UpdateUserDetails(userDetails models.UsersProfileDetails, userID int) (models.UsersProfileDetails, error)
	UpdatePassword(ctx context.Context, body models.UpdatePassword) error
	AddToWishlist(product_id int, user_id int) error
	GetWishList(userID int) ([]models.WishListResponse, error)
	RemoveFromWishlist(productId int, userID int) error
}
