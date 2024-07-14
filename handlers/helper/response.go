package helper

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func base(c *gin.Context, statusCode int, body interface{}) {
	if body == nil {
		c.JSON(statusCode, nil)
	} else {
		c.JSON(statusCode, body)
	}
}

func CriticalError(c *gin.Context, err error) {
	base(c, http.StatusInternalServerError, err.Error())
}

func Error(c *gin.Context, statusCode int, err error) {
	base(c, statusCode, err)
}

func Ok(c *gin.Context, v interface{}) {
	base(c, http.StatusOK, v)
}
