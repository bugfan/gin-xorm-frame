package api

import (
	"errors"
	"net/http"
	"time"

	"gin-xorm-frame/models"

	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
)

func init() {
	RegisterController(new(adminController))
}

type adminController struct {
	X *xorm.Engine
}

type adminContent struct {
	ID       int64
	Username string
	Password string
	Created  time.Time
	Updated  time.Time
}

func (c *adminContent) Check(ctx *gin.Context, scfd *Scaffold, t ScaffoldRouteType) bool {
	if len(c.Username) == 0 {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("用户名不能为空"))
		return false
	}
	if t == ScaffoldRouteTypeNew {
		if has, _ := scfd.engine.Exist(&models.Admin{
			Username: c.Username,
		}); has {
			ctx.AbortWithError(http.StatusBadRequest, errors.New("用户名已存在"))
			return false
		}
	} else {
		if has, _ := scfd.engine.Where("ID != ?", c.ID).Exist(&models.Admin{
			Username: c.Username,
		}); has {
			ctx.AbortWithError(http.StatusBadRequest, errors.New("用户名已存在"))
			return false
		}
	}
	return true
}

func (c *adminController) Register(g *gin.RouterGroup) {
	route := g.Group("/admin")
	scraffold := NewScaffold(c.X, new(models.Admin), new(adminContent), ScaffoldRouteTypeALL)
	scraffold.HiddenField = []string{"Password"}
	scraffold.Register(route)

}
