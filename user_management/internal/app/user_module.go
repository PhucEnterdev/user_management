package app

import (
	"enterdev.com.vn/user_management/internal/handler"
	"enterdev.com.vn/user_management/internal/repository"
	"enterdev.com.vn/user_management/internal/routes"
	"enterdev.com.vn/user_management/internal/services"
)

type UserModule struct {
	routes routes.Route
}

func NewUserModule() *UserModule {
	// init the layers
	userRepo := repository.NewInMemoryUserRepositoryImpe()
	userService := services.NewUserServiceImpl(userRepo)
	userHandler := handler.NewUserHandler(userService)
	userRoute := routes.NewUserRoute(userHandler)
	return &UserModule{routes: userRoute}
}

func (m *UserModule) Routes() routes.Route {
	return m.routes
}
