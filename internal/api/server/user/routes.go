package user

import (
	"github.com/gin-gonic/gin"
	"github.com/sinasadeghi83/aut-grader/internal/api/config"
	"github.com/sinasadeghi83/aut-grader/internal/api/user"
	"gorm.io/gorm"
)

func RegisterHandler(db *gorm.DB, cfg *config.AppConfig) *Handler {
	repo := user.NewUserRepo(db)
	svc := user.NewUserService(*repo, cfg)
	return NewHandler(*svc)
}

func (handler *Handler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/login", handler.Login)
}
