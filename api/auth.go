package api

import (
	"errors"
	"time"

	jwt "github.com/appleboy/gin-jwt"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {

	authMiddleware := &jwt.GinJWTMiddleware{
		Realm:      "code_name",        //setting.Get("code_name"),
		Key:        []byte("test_key"), //[]byte(setting.SecretKey()),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour * 24,
		// Authenticator: authenticator,
		// Authorizator: authorizator,
		// IdentityHandler: identityHandler,
		Unauthorized: func(c *gin.Context, code int, message string) {
			_ = c.AbortWithError(code, errors.New(message))
		},
		// PayloadFunc: loginPayload,
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		TokenLookup: "header:Authorization",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

	}
	RegisterGlobal("POST", "/api/sign_in", authMiddleware.LoginHandler)
	RegisterGlobal("POST", "/api/sign_out", func(*gin.Context) {})
	RegisterGlobal("GET", "/api/refresh_token", authMiddleware.MiddlewareFunc(), authMiddleware.RefreshHandler)

	return authMiddleware.MiddlewareFunc()
}

type AuthBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// func authenticator(username string, password string, c *gin.Context) (string, bool) {
// func authenticator(c *gin.Context) (interface{}, error) {
// 	var body AuthBody
// 	if err := c.ShouldBind(&body); err != nil {
// 		return nil, jwt.ErrMissingLoginValues
// 	}

// 	u := model.FindAdmin(body.Username)
// 	if u == nil {
// 		return nil, jwt.ErrFailedAuthentication
// 	}
// 	err := u.CheckPassword(body.Password)
// 	if err != nil {
// 		return nil, jwt.ErrFailedAuthentication
// 	}

// 	return u, nil
// }

// func loginPayload(data interface{}) jwt.MapClaims {
// 	if v, ok := data.(*model.Admin); ok {
// 		return jwt.MapClaims{
// 			"id": v.ID,
// 		}
// 	}
// 	return jwt.MapClaims{}
// }

// func identityHandler(claims jwtgo.MapClaims) interface{} {
// 	if id, ok := claims["id"]; ok {
// 		i64 := int64(id.(float64))
// 		user := new(model.Admin)
// 		has, err := model.Get(i64, user)
// 		if err != nil || !has {
// 			return nil
// 		}
// 		return user
// 	}
// 	return nil
// }

// func authorizator(data interface{}, c *gin.Context) bool {
// 	if data == nil {
// 		return false
// 	}
// 	if _, ok := data.(*model.Admin); ok {
// 		return true
// 	}

// 	return false
// }

// func currentUser(c *gin.Context) *model.Admin {
// 	claims := jwt.ExtractClaims(c)
// 	res := identityHandler(claims)
// 	if u, ok := res.(*model.Admin); ok {
// 		return u
// 	}
// 	return nil
// }

// func currentUserName(c *gin.Context) string {
// 	username := "Unknown"
// 	if user := currentUser(c); user != nil {
// 		username = user.Username
// 	}
// 	return username
// }
