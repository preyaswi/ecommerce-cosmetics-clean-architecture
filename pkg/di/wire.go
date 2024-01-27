//go:build wireinject
// +build wireinject

package di

import (
	http "clean/pkg/api"
	handler "clean/pkg/api/handler"
	config "clean/pkg/config"
	db "clean/pkg/db"
	repository "clean/pkg/repository"
	usecase "clean/pkg/usecase"

	"github.com/google/wire"
)

func InitializeAPI(cfg config.Config) (*http.ServerHTTP, error) {
	wire.Build(db.ConnectDatabase,
		 repository.NewuserRepository, usecase.NewUserUseCase, handler.NewUserHandler, 
		 http.NewServerHTTP, 
		 repository.NewProductRepository, usecase.NewProductUseCase, handler.NewProductHandler,
		  handler.NewOtpHandler, usecase.NewOtpUseCase, repository.NewOtpRepository,
		   handler.NewCartHandler, usecase.NewCartUseCase, repository.NewCartRepository,
		   repository.NewOrderRepository,usecase.NewOrderUseCase, handler.NewOrderHandler,)
	return &http.ServerHTTP{}, nil
}
