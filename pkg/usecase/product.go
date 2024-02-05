package usecase

import (
	interfaces "clean/pkg/repository/interface"
	services "clean/pkg/usecase/interface"
	"clean/pkg/utils/models"
)

type productUseCase struct {
	productRepo interfaces.ProductRepository
}

func NewProductUseCase(repo interfaces.ProductRepository) services.ProductUseCase {
	return &productUseCase{
		productRepo: repo,
	}
}

func(pr *productUseCase) ShowAllProducts(page int, count int) ([]models.ProductBrief, error) {
	productsBrief, err := pr.productRepo.ShowAllProducts(page, count)
	if err != nil {
		return []models.ProductBrief{}, err
	}
	for i := range productsBrief {
		p := &productsBrief[i]
		if p.Quantity == 0 {
			p.ProductStatus = "out of stock"
		} else {
			p.ProductStatus = "in stock"
		}
	}
	return productsBrief, nil
}
func(pr *productUseCase) ShowIndividualProducts(id int) (*models.ProductBrief, error) {
	product, err := pr.productRepo.ShowIndividualProducts(id)
	if err != nil {
		return &models.ProductBrief{}, err
	}
	return product, nil

}
func (pr *productUseCase)FilterCategory(data map[string]int) ([]models.ProductBrief, error) {

	err := pr.productRepo.CheckValidityOfCategory(data)
	if err != nil {
		return []models.ProductBrief{}, err
	}

	var productFromCategory []models.ProductBrief
	for _, id := range data {

		product, err := pr.productRepo.GetProductFromCategory(id)
		if err != nil {
			return []models.ProductBrief{}, err
		}
		for _, product := range product {

			quantity, err := pr.productRepo.GetQuantityFromProductID(product.ID)
			if err != nil {
				return []models.ProductBrief{}, err
			}
			if quantity == 0 {
				product.ProductStatus = "out of stock"
			} else {
				product.ProductStatus = "in stock"
			}
			if product.ID != 0 {
				productFromCategory = append(productFromCategory, product)
			}
		}

		// if a product exist for that genre. Then only append it

	}
	return productFromCategory, nil

}