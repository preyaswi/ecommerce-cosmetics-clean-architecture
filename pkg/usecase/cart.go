package usecase

import (
	interfaces "clean/pkg/repository/interface"
	services "clean/pkg/usecase/interface"
	"clean/pkg/utils/models"
	"errors"
)

type cartUseCase struct {
	cartRepository    interfaces.CartRepository
	productRepository interfaces.ProductRepository
}

func NewCartUseCase(cartRepo interfaces.CartRepository, productRepo interfaces.ProductRepository) services.CartUseCase {
	return &cartUseCase{
		cartRepository: cartRepo,
	}
}
func (cr *cartUseCase) AddToCart(product_id int, user_id int) (models.CartResponse, error) {
	ok, _, err := cr.cartRepository.CheckProduct(product_id)
	//here the second return is category and we will use this later when we need to add the offer details later
	if err != nil {

		return models.CartResponse{}, err
	}

	if !ok {
		return models.CartResponse{}, errors.New("product Does not exist")
	}

	QuantityOfProductInCart, err := cr.cartRepository.QuantityOfProductInCart(user_id, product_id)
	if err != nil {

		return models.CartResponse{}, err
	}
	quantityOfProduct, err := cr.productRepository.GetQuantityFromProductID(product_id)
	if err != nil {

		return models.CartResponse{}, err
	}
	if quantityOfProduct == 0 {
		return models.CartResponse{}, errors.New("out of stock")
	}
	if quantityOfProduct == QuantityOfProductInCart {
		return models.CartResponse{}, errors.New("stock limit exceeded")
	}
	productPrice, err := cr.productRepository.GetPriceOfProductFromID(product_id)
	if err != nil {

		return models.CartResponse{}, err
	}
	if QuantityOfProductInCart == 0 {
		err := cr.cartRepository.AddItemIntoCart(user_id, product_id, 1, productPrice)
		if err != nil {

			return models.CartResponse{}, err
		}

	} else {
		currentTotal, err := cr.cartRepository.TotalPriceForProductInCart(user_id, product_id)
		if err != nil {
			return models.CartResponse{}, err
		}
		err = cr.cartRepository.UpdateCart(QuantityOfProductInCart+1, currentTotal+productPrice, user_id, product_id)
		if err != nil {
			return models.CartResponse{}, err
		}

	}
	cartDetails, err := cr.cartRepository.DisplayCart(user_id)
	if err != nil {
		return models.CartResponse{}, err
	}
	cartTotal, err := cr.cartRepository.GetTotalPrice(user_id)
	if err != nil {

		return models.CartResponse{}, err
	}
	return models.CartResponse{
		UserName:   cartTotal.UserName,
		TotalPrice: cartTotal.TotalPrice,
		Cart:       cartDetails,
	}, nil

}

func (cr *cartUseCase) RemoveFromCart(product_id int, user_id int) (models.CartResponse, error) {
	ok, err := cr.cartRepository.ProductExist(user_id, product_id)
	if err != nil {
		return models.CartResponse{}, err
	}
	if !ok {
		return models.CartResponse{}, errors.New("product doesn't exist in the cart")
	}
	var cartDetails struct {
		Quantity   int
		TotalPrice float64
	}

	cartDetails, err = cr.cartRepository.GetQuantityAndProductDetails(user_id, product_id, cartDetails)
	if err != nil {
		return models.CartResponse{}, err
	}

	cartDetails.Quantity = cartDetails.Quantity - 1

	//remove the product if quantity after deleting is 0
	if cartDetails.Quantity == 0 {
		if err := cr.cartRepository.RemoveProductFromCart(user_id, product_id); err != nil {
			return models.CartResponse{}, err
		}

	}
	if cartDetails.Quantity != 0 {

		product_price, err := cr.productRepository.GetPriceOfProductFromID(product_id)
		if err != nil {
			return models.CartResponse{}, err
		}
		cartDetails.TotalPrice = cartDetails.TotalPrice - product_price
		err = cr.cartRepository.UpdateCartDetails(cartDetails, user_id, product_id)
		if err != nil {
			return models.CartResponse{}, err
		}
	}
	updatedCart, err := cr.cartRepository.CartAfterRemovalOfProduct(user_id)
	if err != nil {
		return models.CartResponse{}, err
	}
	cartTotal, err := cr.cartRepository.GetTotalPrice(user_id)
	if err != nil {
		return models.CartResponse{}, err
	}
	return models.CartResponse{
		UserName:   cartTotal.UserName,
		TotalPrice: cartTotal.TotalPrice,
		Cart:       updatedCart,
	}, nil

}
func (cr *cartUseCase) DisplayCart(user_id int) (models.CartResponse, error) {

	cart, err := cr.cartRepository.DisplayCart(user_id)
	if err != nil {
		return models.CartResponse{}, err
	}
	cartTotal, err := cr.cartRepository.GetTotalPrice(user_id)
	if err != nil {
		return models.CartResponse{}, err
	}
	return models.CartResponse{
		UserName:   cartTotal.UserName,
		TotalPrice: cartTotal.TotalPrice,
		Cart:       cart,
	}, nil

}
func (cr *cartUseCase) EmptyCart(userID int) (models.CartResponse, error) {
	ok, err := cr.cartRepository.CartExist(userID)
	if err != nil {
		return models.CartResponse{}, err
	}
	if !ok {
		return models.CartResponse{}, errors.New("cart already empty")
	}
	if err := cr.cartRepository.EmptyCart(userID); err != nil {
		return models.CartResponse{}, err
	}

	cartTotal, err := cr.cartRepository.GetTotalPrice(userID)

	if err != nil {
		return models.CartResponse{}, err
	}

	cartResponse := models.CartResponse{
		UserName:   cartTotal.UserName,
		TotalPrice: cartTotal.TotalPrice,
		Cart:       []models.Cart{},
	}

	return cartResponse, nil

}
