package api

import (
	"errors"
	"net/http"
	"time"

	"scaffold/model"

	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
)

func init() {
	RegisterController(new(userController))
}

type userController struct {
	X *xorm.Engine
}

type userContent struct {
	ID       int64
	Username string
	Password string
	Created  time.Time
	Updated  time.Time
}

func (c *userContent) Check(ctx *gin.Context, scfd *Scaffold, t ScaffoldRouteType) bool {
	println("sss:", t)
	if t == ScaffoldRouteTypeNew || t == ScaffoldRouteTypeUpdate {
		if len(c.Username) == 0 {
			ctx.AbortWithError(http.StatusBadRequest, errors.New("用户名不能为空"))
			return false
		}
	}
	return true
}

func (c *userController) Register(g *gin.RouterGroup) {
	route := g.Group("/user")
	scraffold := NewScaffold(c.X, new(model.User), new(userContent), ScaffoldRouteTypeALL)
	scraffold.HiddenField = []string{"Password"}
	scraffold.Register(route)

}
