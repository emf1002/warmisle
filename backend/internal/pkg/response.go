package pkg

import "github.com/gin-gonic/gin"

// Response is the standard API response format.
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Success sends a success response.
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Response{Code: 0, Message: "ok", Data: data})
}

// Error sends an error response.
func Error(c *gin.Context, httpCode, bizCode int, msg string) {
	c.JSON(httpCode, Response{Code: bizCode, Message: msg})
}
