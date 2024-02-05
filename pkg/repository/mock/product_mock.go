// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/repository/interface/product.go

// Package mock is a generated GoMock package.
package mock

import (
	models "clean/pkg/utils/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockProductRepository is a mock of ProductRepository interface.
type MockProductRepository struct {
	ctrl     *gomock.Controller
	recorder *MockProductRepositoryMockRecorder
}

// MockProductRepositoryMockRecorder is the mock recorder for MockProductRepository.
type MockProductRepositoryMockRecorder struct {
	mock *MockProductRepository
}

// NewMockProductRepository creates a new mock instance.
func NewMockProductRepository(ctrl *gomock.Controller) *MockProductRepository {
	mock := &MockProductRepository{ctrl: ctrl}
	mock.recorder = &MockProductRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProductRepository) EXPECT() *MockProductRepositoryMockRecorder {
	return m.recorder
}

// CheckValidityOfCategory mocks base method.
func (m *MockProductRepository) CheckValidityOfCategory(data map[string]int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckValidityOfCategory", data)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckValidityOfCategory indicates an expected call of CheckValidityOfCategory.
func (mr *MockProductRepositoryMockRecorder) CheckValidityOfCategory(data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckValidityOfCategory", reflect.TypeOf((*MockProductRepository)(nil).CheckValidityOfCategory), data)
}

// GetPriceOfProductFromID mocks base method.
func (m *MockProductRepository) GetPriceOfProductFromID(prodcut_id int) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPriceOfProductFromID", prodcut_id)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPriceOfProductFromID indicates an expected call of GetPriceOfProductFromID.
func (mr *MockProductRepositoryMockRecorder) GetPriceOfProductFromID(prodcut_id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPriceOfProductFromID", reflect.TypeOf((*MockProductRepository)(nil).GetPriceOfProductFromID), prodcut_id)
}

// GetProductFromCategory mocks base method.
func (m *MockProductRepository) GetProductFromCategory(id int) ([]models.ProductBrief, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductFromCategory", id)
	ret0, _ := ret[0].([]models.ProductBrief)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductFromCategory indicates an expected call of GetProductFromCategory.
func (mr *MockProductRepositoryMockRecorder) GetProductFromCategory(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductFromCategory", reflect.TypeOf((*MockProductRepository)(nil).GetProductFromCategory), id)
}

// GetQuantityFromProductID mocks base method.
func (m *MockProductRepository) GetQuantityFromProductID(id int) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetQuantityFromProductID", id)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetQuantityFromProductID indicates an expected call of GetQuantityFromProductID.
func (mr *MockProductRepositoryMockRecorder) GetQuantityFromProductID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetQuantityFromProductID", reflect.TypeOf((*MockProductRepository)(nil).GetQuantityFromProductID), id)
}

// ShowAllProducts mocks base method.
func (m *MockProductRepository) ShowAllProducts(page, count int) ([]models.ProductBrief, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShowAllProducts", page, count)
	ret0, _ := ret[0].([]models.ProductBrief)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ShowAllProducts indicates an expected call of ShowAllProducts.
func (mr *MockProductRepositoryMockRecorder) ShowAllProducts(page, count interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowAllProducts", reflect.TypeOf((*MockProductRepository)(nil).ShowAllProducts), page, count)
}

// ShowIndividualProducts mocks base method.
func (m *MockProductRepository) ShowIndividualProducts(id int) (*models.ProductBrief, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShowIndividualProducts", id)
	ret0, _ := ret[0].(*models.ProductBrief)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ShowIndividualProducts indicates an expected call of ShowIndividualProducts.
func (mr *MockProductRepositoryMockRecorder) ShowIndividualProducts(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowIndividualProducts", reflect.TypeOf((*MockProductRepository)(nil).ShowIndividualProducts), id)
}
