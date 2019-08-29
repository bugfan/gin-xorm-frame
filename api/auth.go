package api

import (
	"errors"
	"gin-xorm-frame/models"
	"gin-xorm-frame/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	AUTH_USERNAME = "username"
)

func AuthMiddleware(ctx *gin.Context) {
	_, err := utils.GetJWTDataFromCookie(ctx.Request)
	if err != nil {
		ctx.Writer.WriteHeader(http.StatusUnauthorized)
		ctx.Abort()
	}
	ctx.Next()
}

func Sign(g *gin.RouterGroup) {
	g.POST("/sign_in", signin)
	g.POST("/sign_out", signout)
	g.GET("/user_info", userinfo)
}

func signin(ctx *gin.Context) {
	type admin struct {
		Username, Password string
	}

	body := &admin{}
	err := ctx.ShouldBind(body)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
	}
	u := models.FindAdmin(body.Username)
	if u == nil {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("帐号密码错误"))
	}
	err = u.CheckPassword(body.Password)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("帐号密码错误"))
	}
	// todo set cookie
	if err := utils.SetJWTDataToCookie(ctx.Writer, map[string]string{AUTH_USERNAME: u.Username}); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
	ctx.JSON(http.StatusOK, u)
}

func signout(ctx *gin.Context) {
	// todo delete cookie
	username := CurrentName(ctx.Request)
	utils.DeleteJWTCookie(ctx.Writer)
	ctx.JSON(http.StatusOK, username)
}

func userinfo(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "")
}

func CurrentName(req *http.Request) string {
	cookieData, err := utils.GetJWTDataFromCookie(req)
	if err != nil {
		return ""
	}
	if username, ok := cookieData[AUTH_USERNAME]; ok {
		return username
	}
	return ""
}
