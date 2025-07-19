package rest

import "github.com/gin-gonic/gin"

// Response represents a standardized JSON response structure.
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// RespondOK sends a successful JSON response.
func RespondOK(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Status: "success",
		Data:   data,
	})
}

// RespondCreated sends a successful creation JSON response.
func RespondCreated(c *gin.Context, data interface{}) {
	c.JSON(201, Response{
		Status: "success",
		Data:   data,
	})
}

// RespondNoContent sends a 204 No Content response.
func RespondNoContent(c *gin.Context) {
	c.Status(204)
}

// RespondError sends an error JSON response.
func RespondError(c *gin.Context, httpStatus int, message string, err error) {
	resp := Response{
		Status:  "error",
		Message: message,
	}
	if err != nil {
		resp.Error = err.Error()
	}
	c.JSON(httpStatus, resp)
}
