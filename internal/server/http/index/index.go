package index

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) HandlerIndex(c *gin.Context) {
	c.JSON(http.StatusOK, "Service is running :))")
}
