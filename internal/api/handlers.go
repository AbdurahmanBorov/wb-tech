package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"wb-tech/internal/models"
)

type services interface {
	GetCache(uuid string) (models.Order, error)
}

type Api struct {
	services services
}

func New(serv services) *gin.Engine {
	h := Api{services: serv}

	r := gin.New()

	r.Static("/css", "templates/css/")
	r.LoadHTMLGlob("templates/index.html")

	r.GET("/order/:orderId", h.GetOrder)
	r.GET("/order", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	return r
}
