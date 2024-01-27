package repository

import (
	interfaces "clean/pkg/repository/interface"
	"clean/pkg/utils/models"
	"errors"

	"gorm.io/gorm"
)

type productDatabase struct {
	DB *gorm.DB
}

func NewProductRepository(DB *gorm.DB) interfaces.ProductRepository {
	return &productDatabase{DB}
}

func(p *productDatabase) ShowAllProducts(page int, count int) ([]models.ProductBrief, error) {
	if page <= 0 {
		page = 1
	}

	if count <= 0 {
		count = 5
	}

	offset := (page - 1) * count
	var productsBrief []models.ProductBrief

	err := p.DB.Raw(`
		SELECT * FROM products limit ? offset ?
	`, count, offset).Scan(&productsBrief).Error

	if err != nil {
		return nil, err
	}

	return productsBrief, nil

}
func (p *productDatabase)ShowIndividualProducts(id int) (*models.ProductBrief, error) {
	var product models.ProductBrief
	result := p.DB.Raw(`
SELECT 
       p.id,
	   p.name,
	   p.sku,
	   c.category_name,
	   p.brand_id,
	   p.quantity,
	   p.price,
	   p.product_status
FROM
	   products p
JOIN
	   categories c ON p.category_id=c.id
WHERE
	   p.id=?`, id).Scan(&product)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}
	return &product, nil
}
func(p *productDatabase) CheckValidityOfCategory(data map[string]int) error {

	for _, id := range data {
		var count int
		err := p.DB.Raw("select count(*) from categories where id = ?", id).Scan(&count).Error
		if err != nil {
			return err
		}

		if count < 1 {
			return errors.New("genre does not exist")
		}
	}
	return nil
}
func (p *productDatabase)GetProductFromCategory(id int) ([]models.ProductBrief, error) {

	var product []models.ProductBrief
	err := p.DB.Raw(`
		SELECT *
		FROM products
		JOIN categories ON products.category_id = categories.id
		 where categories.id = ?
	`, id).Scan(&product).Error

	if err != nil {
		return []models.ProductBrief{}, err
	}
	return product, nil
}
func(p *productDatabase) GetQuantityFromProductID(id int) (int, error) {

	var quantity int
	err := p.DB.Raw("select quantity from products where id = ?", id).Scan(&quantity).Error
	if err != nil {
		return 0.0, err
	}

	return quantity, nil

}
func(p *productDatabase) GetPriceOfProductFromID(prodcut_id int) (float64, error) {
	var productPrice float64

	if err := p.DB.Raw("select price from products where id = ?", prodcut_id).Scan(&productPrice).Error; err != nil {
		return 0.0, err
	}
	return productPrice, nil
}