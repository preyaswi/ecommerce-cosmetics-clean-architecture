package interfaces

import "clean/pkg/utils/models"

type CartUseCase interface {
	AddToCart(product_id int, user_id int) (models.CartResponse, error)
	RemoveFromCart(product_id int, user_id int) (models.CartResponse, error)
	DisplayCart(user_id int) (models.CartResponse, error)
	EmptyCart(userID int) (models.CartResponse, error)
	
}
