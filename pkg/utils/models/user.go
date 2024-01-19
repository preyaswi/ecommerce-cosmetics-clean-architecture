package models

type SignupDetail struct {
	Firstname    string `json:"firstname"  validate:"required"`
	Lastname     string `json:"lastname"  validate:"required"`
	Email        string `json:"email"  validate:"required"`
	Phone        string `json:"phone"  validate:"required"`
	Password     string `json:"password"  validate:"required"`
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