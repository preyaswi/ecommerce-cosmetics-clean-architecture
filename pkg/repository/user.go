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

func NewuserRepository(DB *gorm.DB) interfaces.UserRepository {
	return &userDatabase{
		DB: DB,
	}

}
func (c *userDatabase) CheckUserExistsByEmail(email string) (*domain.User, error) {
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
func (c *userDatabase) CheckUserExistsByPhone(phone string) (*domain.User, error) {
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
func (c *userDatabase) UserSignup(user models.SignupDetail) (models.SignupDetailResponse, error) {
	var signupDetail models.SignupDetailResponse
	err := c.DB.Raw(`
		INSERT INTO users(firstname, lastname, email, password, phone)
		VALUES(?, ?, ?, ?, ?)
		RETURNING id, firstname, lastname, email, phone
	`, user.Firstname, user.Lastname, user.Email, user.Password, user.Phone).
		Scan(&signupDetail).Error
		
	if err != nil {
		return models.SignupDetailResponse{}, err
	}
	return signupDetail, nil

}
func (c *userDatabase) CreateReferralEntry(userDetails models.SignupDetailResponse, userReferral string) error {

	err := c.DB.Exec("insert into referrals (user_id,referral_code,referral_amount) values (?,?,?)", userDetails.Id, userReferral, 0).Error
	if err != nil {
		return err
	}

	return nil

}
func (c *userDatabase) GetUserIdFromReferrals(ReferralCode string) (int, error) {

	var referredUserId int
	err := c.DB.Raw("select user_id from referrals where referral_code = ?", ReferralCode).Scan(&referredUserId).Error
	if err != nil {
		return 0, nil
	}

	return referredUserId, nil
}

func (c *userDatabase) UpdateReferralAmount(referralAmount float64, referredUserId int, currentUserID int) error {

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
func (c *userDatabase) FindUserDetailsByEmail(user models.LoginDetail) (models.UserLoginResponse, error) {
	var userdetails models.UserLoginResponse

	err := c.DB.Raw(
		`SELECT * FROM users where email = ? and blocked = false`, user.Email).Scan(&userdetails).Error

	if err != nil {
		return models.UserLoginResponse{}, errors.New("error checking user details")
	}
	return userdetails, nil

}

func (c *userDatabase) GetAllAddress(userId int) (models.AddressInfoResponse, error) {
	var addressInfoResponse models.AddressInfoResponse
	if err := c.DB.Raw("select * from addresses where user_id = ?", userId).Scan(&addressInfoResponse).Error; err != nil {
		return models.AddressInfoResponse{}, err
	}
	return addressInfoResponse, nil
}
func (c *userDatabase) AddAddress(userId int, address models.AddressInfo) error {

	if err := c.DB.Raw("insert into addresses(user_id,name,house_name,street,city,state,pin)  values(?,?,?,?,?,?,?)", userId, address.Name, address.HouseName, address.Street, address.City, address.State, address.Pin).Scan(&address).Error; err != nil {
		return err
	}
	return nil
}

func (c *userDatabase) UpdateAddress(address models.AddressInfo, addressID int, userID int) (models.AddressInfoResponse, error) {

	err := c.DB.Exec("update addresses set house_name = ?, state = ?, pin = ?, street = ?, city = ? where id = ? and user_id = ?", address.HouseName, address.State, address.Pin, address.Street, address.City, addressID, userID).Error
	if err != nil {
		return models.AddressInfoResponse{}, err
	}

	var addressResponse models.AddressInfoResponse
	err = c.DB.Raw("select * from addresses where id = ?", addressID).Scan(&addressResponse).Error
	if err != nil {
		return models.AddressInfoResponse{}, err
	}

	return addressResponse, nil

}
func (c *userDatabase) UserDetails(userID int) (models.UsersProfileDetails, error) {

	var userDetails models.UsersProfileDetails
	err := c.DB.Raw("select users.firstname,users.lastname,users.email,users.phone from users  where users.id = ?", userID).Row().Scan(&userDetails.Firstname, &userDetails.Lastname, &userDetails.Email, &userDetails.Phone)
	if err != nil {
		return models.UsersProfileDetails{}, errors.New("could not get user details")
	}
	return userDetails, nil
}
func (c *userDatabase) CheckUserAvailability(email string) bool {

	var count int
	query := fmt.Sprintf("select count(*) from users where email='%s'", email)
	if err := c.DB.Raw(query).Scan(&count).Error; err != nil {
		return false
	}
	// if count is greater than 0 that means the user already exist
	return count > 0

}
func (c *userDatabase) UpdateUserEmail(email string, userID int) error {

	err := c.DB.Exec("update users set email = ? where id = ?", email, userID).Error
	if err != nil {
		return err
	}
	return nil

}

func (c *userDatabase) UpdateUserPhone(phone string, userID int) error {

	err := c.DB.Exec("update users set phone = ? where id = ?", phone, userID).Error
	if err != nil {
		return err
	}
	return nil

}

func (c *userDatabase) UpdateFirstName(name string, userID int) error {

	err := c.DB.Exec("update users set firstname = ? where id = ?", name, userID).Error
	if err != nil {
		return err
	}
	return nil

}
func (c *userDatabase) UpdateLastName(name string, userID int) error {

	err := c.DB.Exec("update users set lastname = ? where id = ?", name, userID).Error
	if err != nil {
		return err
	}
	return nil

}
func (c *userDatabase) UserPassword(userID int) (string, error) {

	var userPassword string
	err := c.DB.Raw("select password from users where id = ?", userID).Scan(&userPassword).Error
	if err != nil {
		return "", err
	}
	return userPassword, nil

}
func (c *userDatabase) UpdateUserPassword(password string, userID int) error {
	err := c.DB.Exec("update users set password = ? where id = ?", password, userID).Error
	if err != nil {
		return err
	}
	fmt.Println("password Updated succesfully")
	return nil
}
func (c *userDatabase) CheckProductExist(pid int) (bool, error) {
	var k int
	err := c.DB.Raw("SELECT COUNT(*) FROM products WHERE id=?", pid).Scan(&k).Error
	if err != nil {
		return false, err
	}

	if k == 0 {
		return false, err
	}

	return true, err
}
func (c *userDatabase) ProductExistInWishList(productID int, userID int) (bool, error) {

	var count int
	err := c.DB.Raw("select count(*) from wish_lists where user_id = ? and product_id = ? ", userID, productID).Scan(&count).Error
	if err != nil {
		return false, errors.New("error checking user product already present")
	}

	return count > 0, nil

}
func (c *userDatabase) AddToWishList(userID int, productID int) error {
	err := c.DB.Exec("insert into wish_lists (user_id,product_id) values (?, ?)", userID, productID).Error
	if err != nil {
		return err
	}

	return nil
}

func (c *userDatabase) GetWishList(userId int) ([]models.WishListResponse, error) {
	var wishList []models.WishListResponse
	err := c.DB.Raw("select products.id as product_id,products.name as product_name,products.price as product_price from products inner join wish_lists on products.id=wish_lists.product_id where wish_lists.user_id=?", userId).Scan(&wishList).Error
	if err != nil {
		return []models.WishListResponse{}, err
	}
	return wishList, nil
}
func (c *userDatabase) RemoveFromWishlist(userID int, productId int) error {
	err := c.DB.Exec("delete from wish_lists where user_id=? and product_id =?", userID, productId).Error
	if err != nil {
		return err

	}
	return nil
}

func (c *userDatabase) GetAllAddresses(userID int) ([]models.AddressInfoResponse, error) {

	var addressResponse []models.AddressInfoResponse
	err := c.DB.Raw(`select * from addresses where user_id = $1`, userID).Scan(&addressResponse).Error
	if err != nil {
		return []models.AddressInfoResponse{}, err
	}

	return addressResponse, nil

}

func (c *userDatabase) GetAllPaymentOption() ([]models.PaymentDetails, error) {

	var paymentMethods []models.PaymentDetails
	err := c.DB.Raw("select * from payment_methods").Scan(&paymentMethods).Error
	if err != nil {
		return []models.PaymentDetails{}, err
	}

	return paymentMethods, nil

}
