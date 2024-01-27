package repository

import (
	"clean/pkg/domain"
	"clean/pkg/helper"
	interfaces "clean/pkg/repository/interface"
	"clean/pkg/utils/models"
	"fmt"

	"gorm.io/gorm"
)

type orderRepository struct {
	DB *gorm.DB
}

func NewOrderRepository(DB *gorm.DB) interfaces.OrderRepository {
	return &orderRepository{
		DB: DB,
	}
}
func (o *orderRepository) DoesCartExist(userID int) (bool, error) {

	var exist bool
	err := o.DB.Raw("select exists(select 1 from carts where user_id = ?)", userID).Scan(&exist).Error
	if err != nil {
		return false, err
	}

	return exist, nil
}

func (o *orderRepository) AddressExist(orderBody models.OrderIncoming) (bool, error) {

	var count int
	if err := o.DB.Raw("select count(*) from addresses where user_id = ? and id = ?", orderBody.UserID, orderBody.AddressID).Scan(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil

}
func (o *orderRepository) GetCouponDiscountPrice(UserID int, GrandTotal float64) (float64, error) {

	discountPrice, err := helper.GetCouponDiscountPrice(UserID, GrandTotal, o.DB)
	if err != nil {
		return 0.0, err
	}

	return discountPrice, nil

}

func (o *orderRepository) UpdateCouponDetails(discount_price float64, UserID int) error {

	if discount_price != 0.0 {
		err := o.DB.Exec("update used_coupons set used = true where user_id = ?", UserID).Error
		if err != nil {
			return err
		}
	}
	return nil
}
func (o *orderRepository) CreateOrder(orderDetails domain.Order) error {

	err := o.DB.Create(&orderDetails).Error
	if err != nil {
		return err
	}
	return nil

}
func (o *orderRepository) AddOrderItems(orderItemDetails domain.OrderItem, UserID int, ProductID uint, Quantity float64) error {

	// after creating the order delete all cart items and also update the quantity of the product
	err := o.DB.Omit("id").Create(&orderItemDetails).Error
	if err != nil {
		return err
	}

	err = o.DB.Exec("delete from carts where user_id = ? and product_id = ?", UserID, ProductID).Error
	if err != nil {
		return err
	}

	err = o.DB.Exec("update products set quantity = quantity - ? where id = ?", Quantity, ProductID).Error
	if err != nil {
		return err
	}

	return nil

}
func (o *orderRepository) UpdateUsedOfferDetails(userID uint) error {

	o.DB.Exec("update category_offer_useds set used = true where user_id = ?", userID)
	o.DB.Exec("update product_offer_useds set used = true where user_id = ?", userID)

	return nil
}

func (o *orderRepository) GetBriefOrderDetails(orderID string) (domain.OrderSuccessResponse, error) {

	var orderSuccessResponse domain.OrderSuccessResponse
	o.DB.Raw("select order_id,shipment_status from orders where order_id = ?", orderID).Scan(&orderSuccessResponse)
	return orderSuccessResponse, nil

}
func(o *orderRepository) GetOrderDetails(userId int, page int, count int) ([]models.FullOrderDetails, error) {
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * count
	var orderDetails []models.OrderDetails
	fmt.Println("userid is", userId, "page is ", page, "count is ", count, "offset is", offset)
	o.DB.Raw("select order_id,final_price,shipment_status,payment_status from orders where user_id = ? limit ? offset ? ", userId, count, offset).Scan(&orderDetails)
	fmt.Println("order details is ", orderDetails)

	var fullOrderDetails []models.FullOrderDetails
	// for each order select all the associated products and their details
	for _, od := range orderDetails {

		var orderProductDetails []models.OrderProductDetails
		o.DB.Raw("select order_items.product_id,products.name as product_name,order_items.quantity,order_items.total_price from order_items inner join products on order_items.product_id = products.id where order_items.order_id = ?", od.OrderId).Scan(&orderProductDetails)
		fullOrderDetails = append(fullOrderDetails, models.FullOrderDetails{OrderDetails: od, OrderProductDetails: orderProductDetails})

	}

	return fullOrderDetails, nil

}
func(o *orderRepository) UserOrderRelationship(orderID string, userID int) (int, error) {

	var testUserID int
	err := o.DB.Raw("select user_id from orders where order_id = ?", orderID).Scan(&testUserID).Error
	if err != nil {
		return -1, err
	}
	return testUserID, nil
}
func(o *orderRepository) GetProductDetailsFromOrders(orderID string) ([]models.OrderProducts, error) {

	var orderProductDetails []models.OrderProducts
	if err := o.DB.Raw("select product_id,quantity from order_items where order_id = ?", orderID).Scan(&orderProductDetails).Error; err != nil {
		return []models.OrderProducts{}, err
	}

	return orderProductDetails, nil
}
func (o *orderRepository)GetShipmentStatus(orderID string) (string, error) {

	var shipmentStatus string
	err := o.DB.Raw("select shipment_status from orders where order_id = ?", orderID).Scan(&shipmentStatus).Error
	if err != nil {
		return "", err
	}

	return shipmentStatus, nil

}
func(o *orderRepository) CancelOrders(orderID string) error {
	shipmentStatus := "cancelled"
	err := o.DB.Exec("update orders set shipment_status = ? where order_id = ?", shipmentStatus, orderID).Error
	if err != nil {
		return err
	}
	var paymentMethod int
	err = o.DB.Raw("select payment_method_id from orders where order_id = ?", orderID).Scan(&paymentMethod).Error
	if err != nil {
		return err
	}
	if paymentMethod == 3 || paymentMethod == 2 {
		err = o.DB.Exec("update orders set payment_status = 'refunded'  where order_id = ?", orderID).Error
		if err != nil {
			return err
		}
	}
	return nil

}
func (o *orderRepository)UpdateQuantityOfProduct(orderProducts []models.OrderProducts) error {

	for _, od := range orderProducts {

		var quantity int
		if err := o.DB.Raw("select quantity from products where id = ?", od.ProductId).Scan(&quantity).Error; err != nil {
			return err
		}

		od.Quantity += quantity
		if err := o.DB.Exec("update products set quantity = ? where id = ?", od.Quantity, od.ProductId).Error; err != nil {
			return err
		}
	}
	return nil

}
