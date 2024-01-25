package usecase

import (
	"clean/pkg/config"
	"clean/pkg/helper"
	interfaces "clean/pkg/repository/interface"
	"clean/pkg/utils/models"
	"errors"

	"github.com/jinzhu/copier"
)

type otpUseCase struct {
	cfg           config.Config
	otpRepository interfaces.OtpRepository
}

func NewOtpUseCase(cfg config.Config, repo interfaces.OtpRepository) *otpUseCase {
	return &otpUseCase{
		cfg:           cfg,
		otpRepository: repo,
	}
}

func (ot *otpUseCase) SendOTP(phone string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	ok := ot.otpRepository.FindUserByMobileNumber(phone)
	if !ok {
		return errors.New("the user does not exist")
	}
	helper.TwilioSetup(cfg.ACCOUNTSID, cfg.AUTHTOKEN)

	_, err = helper.TwilioSendOtp(phone, cfg.SERVICESSID)

	if err != nil {
		return errors.New("error occured while generating otp")
	}

	return nil
}
func (ot *otpUseCase) VerifyOTP(code models.VerifyData) (models.TokenUser, error) {

	cfg, err := config.LoadConfig()
	if err != nil {
		return models.TokenUser{}, err
	}
	helper.TwilioSetup(cfg.ACCOUNTSID, cfg.AUTHTOKEN)
	err = helper.TwilioVerifyOTP(cfg.SERVICESSID, code.Code, code.User.PhoneNumber)
	if err != nil {
		return models.TokenUser{}, errors.New("error while verifying")
	}
	userDetails, err := ot.otpRepository.UserDetailsUsingPhone(code.User.PhoneNumber)
	if err != nil {
		return models.TokenUser{}, err
	}
	accessToken, err := helper.GenerateAccessToken(userDetails)
	if err != nil {
		return models.TokenUser{}, errors.New("could not create token due to some internal error")
	}
	refreshToken, err := helper.GenerateRefreshToken(userDetails)
	if err != nil {
		return models.TokenUser{}, errors.New("could not create token due to some internal error")
	}
	var user models.SignupDetailResponse
	err = copier.Copy(&user, &userDetails)
	if err != nil {
		return models.TokenUser{}, err
	}
	return models.TokenUser{
		Users:        user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}
