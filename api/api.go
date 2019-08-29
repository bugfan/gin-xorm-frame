package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
)

type APIBackend struct {
	G *gin.Engine
	X *xorm.Engine
}

func NewAPIBackend(x *xorm.Engine) (*APIBackend, error) {
	// if setting.IsProduction() {
	// 	gin.SetMode(gin.ReleaseMode)
	// }
	b := &APIBackend{
		G: gin.Default(),
		X: x,
	}
	b.G.Use(gin.Recovery())
	b.G.Use(gin.ErrorLogger())

	// static dir
	// b.G.StaticFile("/", "./panel/dist/index.html")
	// b.G.Static("/static", "./panel/dist/static")
	// b.G.StaticFile("/favicon.ico", "./panel/dist/favicon.ico")

	// api
	api := b.G.Group("/api")
	// api.Use(setToken)
	Sign(api)
	api.Use(AuthMiddleware)
	// routes
	b.initRoute(api)
	return b, nil
}

func (b *APIBackend) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b.G.ServeHTTP(w, r)
}

func setToken(ctx *gin.Context) {
	if h := ctx.Request.Header.Get("Authorization"); h == "" {
		if token := ctx.Param("token"); token != "" {
			ctx.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			return
		}
		if token, err := ctx.Cookie("Admin-Token"); err == nil && token != "" {
			ctx.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			return
		}
	}
}
