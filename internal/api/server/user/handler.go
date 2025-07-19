package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sinasadeghi83/aut-grader/internal/api/platform/rest"
	"github.com/sinasadeghi83/aut-grader/internal/api/user"
)

type Handler struct {
	Service user.UserService
}

func NewHandler(svc user.UserService) *Handler {
	return &Handler{
		Service: svc,
	}
}

func (h *Handler) Login(c *gin.Context) {
	var lgIn LoginInput
	if err := c.ShouldBindJSON(&lgIn); err != nil {
		rest.RespondError(c, http.StatusBadRequest, "Inavlid input", err)
		return
	}
	user, token, err := h.Service.Login(lgIn.Username, lgIn.Password)
	if err != nil {
		rest.RespondError(c, http.StatusUnauthorized, "Login failed", err)
		return
	}

	rest.RespondOK(c, gin.H{
		"user":  user,
		"token": token,
	})
}
