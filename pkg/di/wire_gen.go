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
	userRepository := repository.NewuserRepository(gormDB)
	userUseCase := usecase.NewUserUseCase(userRepository, cfg)
	userHandler := handler.NewUserHandler(userUseCase)

	serverHTTP := http.NewServerHTTP(userHandler)

	return serverHTTP, nil

}
