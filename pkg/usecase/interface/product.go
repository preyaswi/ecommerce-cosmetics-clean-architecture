package interfaces

import "clean/pkg/utils/models"

type ProductUseCase interface {
	ShowAllProducts(page int, count int) ([]models.ProductBrief, error)
	ShowIndividualProducts(id int) (*models.ProductBrief, error)
	FilterCategory(data map[string]int) ([]models.ProductBrief, error)
}
