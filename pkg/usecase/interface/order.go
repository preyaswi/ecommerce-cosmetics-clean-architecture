package interfaces

import (
	"clean/pkg/domain"
	"clean/pkg/utils/models"
)

type OrderUseCase interface {
	OrderItemsFromCart(orderFromCart models.OrderFromCart, userID int) (domain.OrderSuccessResponse, error)
	GetOrderDetails(userId int, page int, count int) ([]models.FullOrderDetails, error)
	CancelOrders(orderId string, userId int) error
	ExecutePurchaseCOD(userID int, addressID int) (models.Invoice, error)
}
