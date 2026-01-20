package server

import (
	"auth/internal/config"
	"auth/internal/errors"
	"auth/internal/handler"
	"auth/internal/routes"
	"auth/internal/service"
	"fmt"

	"github.com/gin-gonic/gin"
)

type Server struct {
	cfg    *config.Config
	router *gin.Engine
}

func NewServer(cfg *config.Config) (*Server, error) {
	if cfg == nil {
		return nil, fmt.Errorf("Кфг сервера не может быть nil")
	}

	service, err := service.NewService(cfg)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrServiceCreation, err)
	}

	handler := handler.NewHandler(service, cfg)
	if handler == nil {
		return nil, fmt.Errorf("Не удалось создать обработчик сервера")
	}
	fmt.Println("Обработчик сервера успешно создан")

	router := routes.SetupRouter(handler)

	return &Server{
		router: router,
		cfg:    cfg,
	}, nil
}

func (s *Server) Stop() error {
	fmt.Println("Server stopped")
	return nil
}

func (s *Server) Serve() error {
	address := fmt.Sprintf("%s:%s", s.cfg.Host, s.cfg.Port)

	fmt.Printf("Server is ready on %s...\n", address)
	return s.router.Run(address)
}
