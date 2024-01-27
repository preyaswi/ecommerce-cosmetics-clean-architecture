package interfaces

import (
	"clean/pkg/domain"
	"clean/pkg/utils/models"
)

type OrderRepository interface {
	DoesCartExist(userID int) (bool, error)
	AddressExist(orderBody models.OrderIncoming) (bool, error)
	GetCouponDiscountPrice(UserID int, GrandTotal float64) (float64, error)
	UpdateCouponDetails(discount_price float64, UserID int) error
	CreateOrder(orderDetails domain.Order) error
	AddOrderItems(orderItemDetails domain.OrderItem, UserID int, ProductID uint, Quantity float64) error
	UpdateUsedOfferDetails(userID uint) error
	GetBriefOrderDetails(orderID string) (domain.OrderSuccessResponse, error)
	GetOrderDetails(userId int, page int, count int) ([]models.FullOrderDetails, error)
	UserOrderRelationship(orderID string, userID int) (int, error)
	GetProductDetailsFromOrders(orderID string) ([]models.OrderProducts, error)
	GetShipmentStatus(orderID string) (string, error)
	CancelOrders(orderID string) error
	UpdateQuantityOfProduct(orderProducts []models.OrderProducts) error
	
}
