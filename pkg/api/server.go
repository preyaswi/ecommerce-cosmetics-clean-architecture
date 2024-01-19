package http

import (
	"clean/pkg/api/handler"
	"clean/pkg/routes"
	"log"

	"github.com/gin-gonic/gin"
)

type ServerHTTP struct {
	engine *gin.Engine
}

func NewServerHTTP(userHandler *handler.UserHandler) *ServerHTTP {
	engine := gin.New()
	engine.Use(gin.Logger())

	routes.UserRoutes(engine.Group("/users"), userHandler)

	return &ServerHTTP{engine: engine}
}
func (sh *ServerHTTP) Start() {
	err := sh.engine.Run(":3000")
	if err != nil {
		log.Fatal("gin engine couldn't start")
	}
}
