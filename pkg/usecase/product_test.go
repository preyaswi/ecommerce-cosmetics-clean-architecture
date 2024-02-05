package usecase

import (
	"clean/pkg/repository/mock"
	"clean/pkg/utils/models"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_ShowIndividualProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock Product Repository
	productRepo := mock.NewMockProductRepository(ctrl)

	// Create Product UseCase with the mocked repository
	productUseCase := NewProductUseCase(productRepo)

	// Test Case 1: Successful retrieval of product
	expectedProduct := &models.ProductBrief{
		ID:            1,
		Name:          "TestProduct",
		SKU:           "TEST123",
		CategoryName:  "TestCategory",
		BrandID:       1,
		Quantity:      10,
		Price:         99.99,
		ProductStatus: "Active",
	}

	productRepo.EXPECT().
		ShowIndividualProducts(gomock.Any()).
		Return(expectedProduct, nil)

	resultProduct, err := productUseCase.ShowIndividualProducts(1)
	assert.Nil(t, err)
	assert.Equal(t, expectedProduct, resultProduct)
}
