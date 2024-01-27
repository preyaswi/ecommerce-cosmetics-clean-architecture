package usecase

import (
	"clean/pkg/domain"
	"clean/pkg/helper"
	interfaces "clean/pkg/repository/interface"
	services "clean/pkg/usecase/interface"
	"clean/pkg/utils/models"
	"errors"
	"fmt"

	"github.com/jinzhu/copier"
)

type orderUseCase struct {
	orderRepository interfaces.OrderRepository
	cartRepository  interfaces.CartRepository
	userRepo interfaces.UserRepository
}

func NewOrderUseCase(orderRepo interfaces.OrderRepository, cartRepo interfaces.CartRepository,userRepo interfaces.UserRepository) services.OrderUseCase {
	return &orderUseCase{
		orderRepository: orderRepo,
		cartRepository:  cartRepo,
	}
}
func (o *orderUseCase) OrderItemsFromCart(orderFromCart models.OrderFromCart, userID int) (domain.OrderSuccessResponse, error) {
	var orderBody models.OrderIncoming
	err := copier.Copy(&orderBody, &orderFromCart)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}

	orderBody.UserID = uint(userID)
	cartExist, err := o.orderRepository.DoesCartExist(userID)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}

	if !cartExist {
		return domain.OrderSuccessResponse{}, errors.New("cart empty can't order")
	}

	addressExist, err := o.orderRepository.AddressExist(orderBody)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}

	if !addressExist {
		return domain.OrderSuccessResponse{}, errors.New("address does not exist")
	}

	// get all items a slice of carts
	cartItems, err := o.cartRepository.GetAllItemsFromCart(int(orderBody.UserID))
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	var orderDetails domain.Order
	var orderItemDetails domain.OrderItem
	// add general order details - that is to be added to orders table
	orderDetails = helper.CopyOrderDetails(orderDetails, orderBody)

	// get grand total iterating through each products in carts
	for _, c := range cartItems {
		orderDetails.GrandTotal += c.TotalPrice
	}
	discount_price, err := o.orderRepository.GetCouponDiscountPrice(int(orderBody.UserID), orderDetails.GrandTotal)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}
	err = o.orderRepository.UpdateCouponDetails(discount_price, orderDetails.UserID)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}

	orderDetails.FinalPrice = orderDetails.GrandTotal - discount_price
	if orderBody.PaymentID == 2 {
		orderDetails.PaymentStatus = "not paid"
		orderDetails.ShipmentStatus = "pending"
	}

	err = o.orderRepository.CreateOrder(orderDetails)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}

	for _, c := range cartItems {
		// for each order save details of products and associated details and use order_id as foreign key ( for each order multiple product will be there)
		orderItemDetails.OrderID = orderDetails.OrderId
		orderItemDetails.ProductID = c.ProductID
		orderItemDetails.Quantity = int(c.Quantity)
		orderItemDetails.TotalPrice = c.TotalPrice

		err := o.orderRepository.AddOrderItems(orderItemDetails, orderDetails.UserID, c.ProductID, c.Quantity)
		if err != nil {
			return domain.OrderSuccessResponse{}, err
		}

	}

	err = o.orderRepository.UpdateUsedOfferDetails(orderBody.UserID)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}

	orderSuccessResponse, err := o.orderRepository.GetBriefOrderDetails(orderDetails.OrderId)
	if err != nil {
		return domain.OrderSuccessResponse{}, err
	}

	return orderSuccessResponse, nil

}
func (o *orderUseCase) GetOrderDetails(userId int, page int, count int) ([]models.FullOrderDetails, error) {

	fullOrderDetails, err := o.orderRepository.GetOrderDetails(userId, page, count)
	if err != nil {
		return []models.FullOrderDetails{}, err
	}
	return fullOrderDetails, nil

}

func (o *orderUseCase) CancelOrders(orderId string, userId int) error {
	userTest, err := o.orderRepository.UserOrderRelationship(orderId, userId)
	if err != nil {
		return err
	}
	if userTest != userId {
		return errors.New("the order is not done by this user")
	}
	orderProductDetails, err := o.orderRepository.GetProductDetailsFromOrders(orderId)
	if err != nil {
		return err
	}
	shipmentStatus, err := o.orderRepository.GetShipmentStatus(orderId)
	if err != nil {
		return err
	}
	if shipmentStatus == "delivered" {
		return errors.New("item already delivered, cannot cancel")
	}

	if shipmentStatus == "pending" || shipmentStatus == "returned" || shipmentStatus == "return" {
		message := fmt.Sprint(shipmentStatus)
		return errors.New("the order is in" + message + ", so no point in cancelling")
	}

	if shipmentStatus == "cancelled" {
		return errors.New("the order is already cancelled, so no point in cancelling")
	}
	err = o.orderRepository.CancelOrders(orderId)
	if err != nil {
		return err
	}
	err = o.orderRepository.UpdateQuantityOfProduct(orderProductDetails)
	if err != nil {
		return err
	}
	return nil

}
func(o *orderUseCase) ExecutePurchaseCOD(userID int, addressID int) (models.Invoice, error) {
	ok, err := o.cartRepository.CartExist(userID)
	if err != nil {
		return models.Invoice{}, err
	}
	if !ok {
		return models.Invoice{}, errors.New("cart doesnt exist")
	}
	cartDetails, err := o.cartRepository.DisplayCart(userID)
	if err != nil {
		return models.Invoice{}, err
	}
	addresses, err := o.userRepo.GetAllAddress(userID)
	if err != nil {
		return models.Invoice{}, err
	}
	Invoice := models.Invoice{
		Cart:        cartDetails,
		AddressInfo: addresses,
	}
	return Invoice, nil

}