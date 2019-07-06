package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	RegisterController(new(systemController))
}

type systemController struct {
}

func (c *systemController) Register(g *gin.RouterGroup) {
	route := g.Group("/system")
	route.GET("/foo", c.Foo)
}
func (c *systemController) Foo(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "foo")
}
