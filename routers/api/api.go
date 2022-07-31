package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func DefaultOK(c *gin.Context) { c.String(http.StatusOK, "default api content") }
