package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sinasadeghi83/aut-grader/internal/api/server/user"
	"github.com/sinasadeghi83/aut-grader/pkg/config"
	"gorm.io/gorm"
)

type Server struct {
	Engine *gin.Engine
	Addr   string
	DB     *gorm.DB
	Config *config.AppConfig
}

func NewServer(addr string, db *gorm.DB, cfg *config.AppConfig) *Server {
	return &Server{
		Engine: gin.Default(),
		Addr:   addr,
		DB:     db,
		Config: cfg,
	}
}

func (s *Server) SetupRoutes() {
	api := s.Engine.Group("/api")
	//Register modules
	authHandler := user.RegisterHandler(s.DB, s.Config)
	authHandler.RegisterRoutes(api.Group("/auth"))

	//Health Check
	s.Engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
		})
	})
}

func (s *Server) Start() error {
	log.Printf("Server listening on %s", s.Addr)
	return s.Engine.Run(s.Addr)
}
