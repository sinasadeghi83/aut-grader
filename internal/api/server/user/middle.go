package user

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sinasadeghi83/aut-grader/internal/api/platform/rest"
)

func (h *Handler) CheckAuth(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		rest.RespondError(c, http.StatusUnauthorized, "Authorization header is missing", nil)
		c.Abort()
		return
	}

	authToken := strings.Split(authHeader, " ")
	if len(authToken) != 2 || authToken[0] != "Bearer" {
		rest.RespondError(c, http.StatusUnauthorized, "Invalid token format", nil)
		c.Abort()
		return
	}

	tokenString := authToken[1]

	user, err := h.Service.CheckToken(tokenString)
	if err != nil {
		rest.RespondError(c, http.StatusUnauthorized, "Invalid token", nil)
		c.Abort()
		return
	}

	c.Set("curUser", user)
	c.Next()
}
