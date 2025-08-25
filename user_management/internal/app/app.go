package app

import (
	"log"

	"enterdev.com.vn/user_management/internal/config"
	"enterdev.com.vn/user_management/internal/routes"
	"enterdev.com.vn/user_management/internal/validation"
	"github.com/gin-gonic/gin"
	"github.com/lpernett/godotenv"
)

type Module interface {
	Routes() routes.Route
}

type Application struct {
	config  *config.Config
	router  *gin.Engine
	modules []Module
}

func NewApplication(cfg *config.Config) *Application {

	if err := validation.InitValidator(); err != nil {
		log.Fatalf("validator init failed %v", err)
	}

	loadEnv()

	server := gin.Default()

	modules := []Module{
		NewUserModule(),
	}

	routes.RegisterRoutes(server, getModuleRoutes(modules)...)

	return &Application{
		config:  cfg,
		router:  server,
		modules: modules,
	}
}

func (a *Application) Run() error {
	return a.router.Run(a.config.ServerAddress)
}

func getModuleRoutes(modules []Module) []routes.Route {
	routeList := make([]routes.Route, len(modules))
	for i, module := range modules {
		routeList[i] = module.Routes()
	}
	return routeList
}

func loadEnv() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println(".env file not found")
	}
}
