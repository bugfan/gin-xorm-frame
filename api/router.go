package api

import (
	"reflect"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	Register(g *gin.RouterGroup)
}

var controllers = make([]Controller, 0)

func RegisterController(c Controller) {
	controllers = append(controllers, c)
}

type route struct {
	httpMethod   string
	relativePath string
	handlers     []gin.HandlerFunc
}

var globalRoutes = make([]*route, 0)
var globalMiddleware = make([]gin.HandlerFunc, 0)

func RegisterGlobal(method, path string, handlers ...gin.HandlerFunc) {
	for _, r := range globalRoutes {
		if r.httpMethod == method && r.relativePath == path {
			return
		}
	}

	r := &route{
		httpMethod:   method,
		relativePath: path,
		handlers:     handlers,
	}
	globalRoutes = append(globalRoutes, r)
}

func RegisterMiddleware(handler gin.HandlerFunc) {
	if handler != nil {
		globalMiddleware = append(globalMiddleware, handler)
	}
}

func (b *APIBackend) initRoute(r *gin.RouterGroup) {
	for _, m := range globalMiddleware {
		r.Use(m)
	}
	// register routes
	for _, c := range controllers {
		t := reflect.TypeOf(c)
		ot := t.Elem()
		instValue := reflect.New(ot)
		ifce := instValue.Interface()
		b.setExports(ifce)
		ctl, _ := ifce.(Controller)
		ctl.Register(r)
	}

	// global routes
	for _, r := range globalRoutes {
		b.G.RouterGroup.Handle(r.httpMethod, r.relativePath, r.handlers...)
	}
}

func (b *APIBackend) exports() map[string]interface{} {
	fieldVal := make(map[string]interface{})
	ptrObjVal := reflect.ValueOf(b)
	objVal := ptrObjVal.Elem()
	objType := objVal.Type()
	fieldNum := objType.NumField()
	for i := 0; i < fieldNum; i++ {
		sf := objType.Field(i)
		valField := objVal.Field(i)
		if valField.CanInterface() {
			fieldVal[sf.Name] = valField.Interface()
		}
	}
	return fieldVal
}

func (b *APIBackend) setExports(obj interface{}) {
	exports := b.exports()
	for name, val := range exports {
		valVal := reflect.ValueOf(val)
		ptrObjVal := reflect.ValueOf(obj)
		objVal := ptrObjVal.Elem()
		fieldVal := objVal.FieldByName(name)
		if fieldVal.IsValid() && fieldVal.Type() == valVal.Type() {
			fieldVal.Set(valVal)
		}
	}
}
