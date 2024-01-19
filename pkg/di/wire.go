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
	wire.Build(db.ConnectDatabase, repository.NewuserRepository, usecase.NewUserUseCase, handler.NewUserHandler, http.NewServerHTTP)
	return &http.ServerHTTP{}, nil
}
