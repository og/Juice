package juice

import (
	ogjson "github.com/og/json"
	gconv "github.com/og/x/conv"
	greflect "github.com/og/x/reflect"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)



type bindRequestEachCounter struct {
	QueryCount uint
}

func BindRequest(ptr interface{}, r *http.Request) error {
	contentType := r.Header.Get("Content-Type")
	query := r.URL.Query()
	queryCount := len(query)
	var formCount int
	// 后续代码会重新定义 formGet
	var formGet = func(key string)  string {return ""}
	switch {
	case strings.Contains(contentType, "application/x-www-form-urlencoded"):
		err := r.ParseForm()
		if err != nil { return err }
		formCount = len(r.PostForm)
		formGet = func(key string) string {
			return r.PostForm.Get(key)
		}
	case strings.Contains(contentType, "multipart/form-data"):
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {return err}
		formCount = len(r.MultipartForm.Value)
		formGet = func(key string) string {
			return r.FormValue(key)
		}
	case strings.Contains(contentType, "application/json"):
		jsonb , err := ioutil.ReadAll(r.Body)
		if err != nil {
			return err
		}
		err = ogjson.ParseBytesWithErr(jsonb, ptr)
		if err != nil {
			return err
		}
	default:
	}
	if formCount == 0 && queryCount == 0 {
		return nil
	}
	return greflect.DeepEach(ptr, func(rValue reflect.Value, rType reflect.Type, field reflect.StructField) (op greflect.EachOperator) {
		if formCount == 0 && queryCount == 0 {
			return op.Break()
		}
		/* parse query */ {
			err := parserField(&queryCount, field.Tag.Get(queryTag), query.Get, rValue, rType)
			if err != nil {
				return op.Error(err)
			}
		}
		/* parse form */ {
			err := parserField(&formCount, field.Tag.Get(formTag), formGet, rValue, rType)
			if err != nil {
				return op.Error(err)
			}
		}
		return
	})
}

func parserField(unresolvedCount *int, key string, get func(key string) string, rValue reflect.Value, rType reflect.Type)  error {
		if *unresolvedCount == 0 {
			return nil
		}
		if key == "" {return nil}
		value := get(key)
		if value == "" { return nil }
		/* 转换赋值 */ {
			if reflect.PtrTo(rType).Implements(requestMarshalerType) {
				err := rValue.Addr().Interface().(RequestMarshaler).MarshalRequest(value)
				if err != nil { return err }
				*unresolvedCount--
			} else {
				err := gconv.StringReflect(value, rValue)
				if err != nil { return err }
				*unresolvedCount--
			}
		}
		return nil
}