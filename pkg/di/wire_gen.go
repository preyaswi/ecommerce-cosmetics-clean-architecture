package di

import (
	http "clean/pkg/api"
	handler "clean/pkg/api/handler"
	"clean/pkg/config"
	db "clean/pkg/db"
	"clean/pkg/repository"
	usecase "clean/pkg/usecase"
)

func InitializeAPI(cfg config.Config) (*http.ServerHTTP, error) {
	gormDB, err := db.ConnectDatabase(cfg)
	if err != nil {
		return nil, err
	}

	otpRepository:=repository.NewOtprepository(gormDB)
	otpUseCase:=usecase.NewOtpUseCase(cfg,otpRepository)
	otpHandler := handler.NewOtpHandler(otpUseCase)
	productRepository := repository.NewProductRepository(gormDB)
	productUseCase := usecase.NewProductUseCase(productRepository, cfg)
	productHandler := handler.NewProductHandler(productUseCase)
	cartRepository := repository.NewCartRepository(gormDB)
	cartUseCase := usecase.NewCartUseCase(cartRepository, productRepository)
	cartHandler := handler.NewCartHandler(cartUseCase)
	userRepository := repository.NewuserRepository(gormDB)
	userUseCase := usecase.NewUserUseCase(userRepository,cartRepository, cfg)
	userHandler := handler.NewUserHandler(userUseCase)
	orderRepository := repository.NewOrderRepository(gormDB)
	orderUseCase := usecase.NewOrderUseCase(orderRepository, cartRepository, userRepository)
	orderHandler := handler.NewOrderHandler(orderUseCase)
	serverHTTP := http.NewServerHTTP(userHandler,otpHandler,productHandler,cartHandler,orderHandler)

	return serverHTTP, nil

}
