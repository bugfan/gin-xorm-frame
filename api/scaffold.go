package api

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
)

type Scaffold struct {
	model       interface{}
	modelType   reflect.Type
	content     APIContent
	contentType reflect.Type
	engine      *xorm.Engine
	routes      ScaffoldRouteType
	HiddenField []string
	NotCopy     []string
}

type APIContent interface {
	Check(*gin.Context, *Scaffold, ScaffoldRouteType) bool
}

type ModelNew interface {
	New(*xorm.Engine) (interface{}, error)
}

type ModelUpdate interface {
	Update(*xorm.Engine, interface{}) (interface{}, error)
}

type ScaffoldRouteType int

const (
	ScaffoldRouteTypeNew    = 1 << iota
	ScaffoldRouteTypeList   // get query all
	ScaffoldRouteTypeGet    // get query one
	ScaffoldRouteTypeUpdate // put
	ScaffoldRouteTypePatch  // patch
	ScaffoldRouteTypeDelete // delete
	// disable path by default
	ScaffoldRouteTypeALL = ScaffoldRouteTypeNew | ScaffoldRouteTypeList | ScaffoldRouteTypeGet | ScaffoldRouteTypeUpdate | ScaffoldRouteTypeDelete
	// ScaffoldRouteTypeALL = ScaffoldRouteTypeNew | ScaffoldRouteTypeList | ScaffoldRouteTypeGet | ScaffoldRouteTypeUpdate | ScaffoldRouteTypePatch | ScaffoldRouteTypeDelete
)

func NewScaffold(engine *xorm.Engine, model interface{}, content APIContent, routes ...ScaffoldRouteType) *Scaffold {
	modelT := reflect.TypeOf(model)
	if modelT.Kind() == reflect.Ptr {
		modelT = modelT.Elem()
	}
	contentT := reflect.TypeOf(content)
	if contentT.Kind() == reflect.Ptr {
		contentT = contentT.Elem()
	}

	var route ScaffoldRouteType
	for _, r := range routes {
		route |= r
	}
	return &Scaffold{
		model:       model,
		content:     content,
		engine:      engine,
		modelType:   modelT,
		contentType: contentT,
		routes:      route,
		NotCopy:     []string{"ID", "Created", "Updated"},
	}
}

func (b *Scaffold) Register(g *gin.RouterGroup) {
	route := b.routes
	if (route & ScaffoldRouteTypeNew) != 0 {
		g.POST("", b.New)
	}

	if (route & ScaffoldRouteTypeList) != 0 {
		g.GET("", b.List)
	}
	if (route & ScaffoldRouteTypeGet) != 0 {
		g.GET("/:id", b.Get)
	}
	if (route & ScaffoldRouteTypeUpdate) != 0 {
		g.PUT("/:id", b.Update)
	}
	if (route & ScaffoldRouteTypePatch) != 0 {
		g.PATCH("/:id", b.Patch)
	}
	if (route & ScaffoldRouteTypeDelete) != 0 {
		g.DELETE("/:id", b.Delete)
	}
}

func (b *Scaffold) New(c *gin.Context) {
	contentValue := reflect.New(b.contentType)
	content := contentValue.Interface().(APIContent)
	err := c.BindJSON(content)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if ok := content.Check(c, b, ScaffoldRouteTypeNew); !ok {
		return
	}

	var model interface{}
	if getter, ok := content.(ModelNew); ok {
		model, err = getter.New(b.engine)
		if err != nil {
			c.AbortWithError(http.StatusUnprocessableEntity, err)
		}
	} else {
		model = reflect.New(b.modelType).Interface()
		err = copyField(model, content, b.NotCopy)
		if err != nil {
			c.AbortWithError(http.StatusUnprocessableEntity, err)
			return
		}
	}
	_, err = b.engine.Insert(model)
	if err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}

	err = copyField(content, model, b.HiddenField)
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}

	c.JSON(http.StatusCreated, content)
}

func (b *Scaffold) List(c *gin.Context) {
	slice := reflect.MakeSlice(reflect.SliceOf(reflect.PtrTo(b.modelType)), 0, 0)
	slicePtr := reflect.New(slice.Type())
	sliceVal := slicePtr.Elem()
	err := b.engine.Find(slicePtr.Interface())
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
	}
	// contentSlice := reflect.MakeSlice(reflect.SliceOf(), 0, 0).Interface().([]APIContent)
	contentSlice := make([]APIContent, 0, sliceVal.Len())
	for i := 0; i < sliceVal.Len(); i++ {
		content := reflect.New(b.contentType).Interface().(APIContent)
		err = copyField(content, sliceVal.Index(i).Interface(), b.HiddenField)
		if err != nil {
			c.AbortWithError(http.StatusUnprocessableEntity, err)
			return
		}

		contentSlice = append(contentSlice, content)
	}
	c.JSON(http.StatusOK, contentSlice)
}

func (b *Scaffold) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	if id == 0 {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("not found %s", c.Param("id")))
		return
	}
	inst := reflect.New(b.modelType).Interface()
	has, err := b.engine.ID(id).Get(inst)
	if !has {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("not found %s", c.Param("id")))
		return
	}
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	content := reflect.New(b.contentType).Interface().(APIContent)
	err = copyField(content, inst, b.HiddenField)
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	c.JSON(200, content)
}

func (b *Scaffold) Update(c *gin.Context) {
	// get content
	contentValue := reflect.New(b.contentType)
	content := contentValue.Interface().(APIContent)
	err := c.BindJSON(content)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if ok := content.Check(c, b, ScaffoldRouteTypeUpdate); !ok {
		return
	}
	// get model
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	if id == 0 {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("not found %s", c.Param("id")))
		return
	}
	inst := reflect.New(b.modelType).Interface()
	has, err := b.engine.ID(id).Get(inst)
	if !has {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("not found %s", c.Param("id")))
		return
	}
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	if getter, ok := content.(ModelUpdate); ok {
		inst, err = getter.Update(b.engine, inst)
		if err != nil {
			c.AbortWithError(http.StatusUnprocessableEntity, err)
		}
	} else {
		err = copyField(inst, content, b.NotCopy)
		if err != nil {
			c.AbortWithError(http.StatusUnprocessableEntity, err)
			return
		}
	}
	_, err = b.engine.ID(id).AllCols().Update(inst)
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	err = copyField(content, inst, b.HiddenField)
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	c.JSON(http.StatusOK, content)
}

func (b *Scaffold) Patch(c *gin.Context) {
	// get content
	contentValue := reflect.New(b.contentType)
	content := contentValue.Interface().(APIContent)
	err := c.BindJSON(content)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if ok := content.Check(c, b, ScaffoldRouteTypePatch); !ok {
		return
	}
	// get model
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	if id == 0 {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("not found %s", c.Param("id")))
		return
	}
	inst := reflect.New(b.modelType).Interface()
	has, err := b.engine.ID(id).Get(inst)
	if !has {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("not found %s", c.Param("id")))
		return
	}
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	err = copyField(inst, content, b.NotCopy)
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	_, err = b.engine.ID(id).Update(inst)
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	err = copyField(content, inst, b.HiddenField)
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	c.JSON(http.StatusOK, content)
}

func (b *Scaffold) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	if id == 0 {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("not found %s", c.Param("id")))
		return
	}
	inst := reflect.New(b.modelType).Interface()
	has, err := b.engine.ID(id).Get(inst)
	if !has {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("not found %s", c.Param("id")))
		return
	}
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	_, err = b.engine.ID(id).Delete(inst)
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func copyField(to interface{}, from interface{}, excepts []string) error {
	toVal := reflect.ValueOf(to)
	if toVal.Kind() == reflect.Ptr {
		toVal = toVal.Elem()
	}
	fromVal := reflect.ValueOf(from)
	if fromVal.Kind() == reflect.Ptr {
		fromVal = fromVal.Elem()
	}
	// to fileld
	toType := toVal.Type()
	fieldNum := toType.NumField()
	for i := 0; i < fieldNum; i++ {
		toField := toType.Field(i)
		if excepts != nil && stringInSlice(toField.Name, excepts) {
			continue
		}
		toValField := toVal.Field(i)
		if !toValField.CanSet() {
			continue
		}
		if fromValField := fromVal.FieldByName(toField.Name); fromValField.IsValid() && fromValField.Type() == toValField.Type() {
			toValField.Set(fromValField)
			continue
		}
		if fromFunc := fromVal.Addr().MethodByName(toField.Name); fromFunc.IsValid() &&
			fromFunc.Type().NumOut() >= 1 &&
			fromFunc.Type().Out(0) == toValField.Type() &&
			fromFunc.Type().NumIn() == 0 {
			res := fromFunc.Call(make([]reflect.Value, 0))
			if len(res) > 1 {
				last := res[len(res)-1]
				if last.CanInterface() && !last.IsNil() {
					if err, ok := last.Interface().(error); ok {
						return err
					}
				}

			}
			toValField.Set(res[0])
			continue
		}
	}
	// to func

	toVal = toVal.Addr()
	toType = toVal.Type()
	funcNum := toType.NumMethod()
	for i := 0; i < funcNum; i++ {
		// method from type
		toMethod := toType.Method(i)
		if !strings.HasPrefix(toMethod.Name, "Set") {
			// only SetXXX methods
			continue
		}

		name := strings.TrimPrefix(toMethod.Name, "Set")
		// skip excepts
		if excepts != nil && stringInSlice(name, excepts) {
			continue
		}

		// func from value
		toFunc := toVal.MethodByName(toMethod.Name)
		argType := toFunc.Type().In(0)

		// from field
		if fromValField := fromVal.FieldByName(name); fromValField.IsValid() && fromValField.Type() == argType {
			res := toFunc.Call([]reflect.Value{fromValField})
			if len(res) > 0 {
				last := res[len(res)-1]
				if last.CanInterface() && !last.IsNil() {
					if err, ok := last.Interface().(error); ok {
						return err
					}
				}

			}
			continue
		}
		// from func

		if fromFunc := fromVal.Addr().MethodByName(name); fromFunc.IsValid() &&
			fromFunc.Type().NumOut() >= 1 &&
			fromFunc.Type().Out(0) == argType &&
			fromFunc.Type().NumIn() == 0 {
			res := fromFunc.Call(make([]reflect.Value, 0))
			if len(res) > 1 {
				last := res[len(res)-1]

				if last.CanInterface() && !last.IsNil() {
					if err, ok := last.Interface().(error); ok {
						return err
					}
				}

			}

			res = toFunc.Call([]reflect.Value{res[0]})
			if len(res) > 0 {
				last := res[len(res)-1]
				if last.CanInterface() && !last.IsNil() {
					if err, ok := last.Interface().(error); ok {
						return err
					}
				}

			}
			continue
		}

	}
	return nil
}

// func copyStaticField(to interface{}, from interface{}, excepts []string) {
// 	exports := exports(from)
// 	ptrObjVal := reflect.ValueOf(to)
// 	objVal := ptrObjVal.Elem()
// 	for name, val := range exports {
// 		if excepts != nil && stringInSlice(name, excepts) {
// 			continue
// 		}
// 		valVal := reflect.ValueOf(val)
// 		fieldVal := objVal.FieldByName(name)
// 		if fieldVal.IsValid() && fieldVal.Type() == valVal.Type() {
// 			fieldVal.Set(valVal)
// 		}
// 	}
// }

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
