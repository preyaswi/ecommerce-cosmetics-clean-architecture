package models

type SignupDetail struct {
	Firstname    string `json:"firstname"  binding:"required"`
	Lastname     string `json:"lastname"  binding:"required"`
	Email        string `json:"email"  binding:"required" validate:"email"`
	Phone        string `json:"phone"  binding:"required"`
	Password     string `json:"password"  binding:"required"`
	ReferralCode string `json:"referral_code"`
}
type SignupDetailResponse struct {
	Id        int    `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}
type TokenUser struct {
	Users        SignupDetailResponse
	AccessToken  string
	RefreshToken string
}
type LoginDetail struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UserLoginResponse struct {
	Id        uint   `json:"id"`
	UserId    uint   `json:"user_id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
}