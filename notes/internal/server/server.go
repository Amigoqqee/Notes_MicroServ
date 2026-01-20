package server

import (
	"fmt"

	"notes/internal/config"
	"notes/internal/handler"
	"notes/internal/routes"
	"notes/internal/service"

	"github.com/gin-gonic/gin"
)

type Server struct {
	cfg    *config.Config
	router *gin.Engine
}

func NewServer(cfg *config.Config) (*Server, error) {
	if cfg == nil {
		return nil, fmt.Errorf("Кфг не может быть nil")
	}

	service, err := service.NewService(cfg)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать сервис: %w", err)
	}

	handler := handler.NewHandler(cfg, service)

	if handler == nil {
		return nil, fmt.Errorf("не удалось создать обработчик сервера")
	}

	fmt.Printf("Обработчик сервера успешно создан")

	router := routes.SetupRouter(handler)

	return &Server{
		router: router,
		cfg:    cfg,
	}, nil
}

func (s *Server) Start() error {
	fmt.Printf("Сервер запускается на %s:%s\n", s.cfg.Host, s.cfg.Port)
	return nil
}

func (s *Server) Stop() error {
	fmt.Println("Сервер остановлен")
	return nil
}

func (s *Server) Serve() error {
	if err := s.Start(); err != nil {
		return err
	}

	address := fmt.Sprintf("%s:%s", s.cfg.Host, s.cfg.Port)
	fmt.Printf("Сервер готов к обработке запросов на %s...\n", address)
	return s.router.Run(address)
}
