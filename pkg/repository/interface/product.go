package interfaces

import "clean/pkg/utils/models"

type ProductRepository interface {
	ShowAllProducts(page int, count int) ([]models.ProductBrief, error)
	ShowIndividualProducts(id int) (*models.ProductBrief, error)
	CheckValidityOfCategory(data map[string]int) error
	GetProductFromCategory(id int) ([]models.ProductBrief, error)
	GetQuantityFromProductID(id int) (int, error)
	GetPriceOfProductFromID(prodcut_id int) (float64, error)
}
