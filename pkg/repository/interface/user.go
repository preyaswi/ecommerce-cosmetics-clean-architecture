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
}
