package juice

import (
	gconv "github.com/og/x/conv"
	greflect "github.com/og/x/reflect"
	"net/http"
	"reflect"
)
type QueryValuer interface {
	UnmarshalQuery(queryValue string) error
}
var queryValueType = reflect.TypeOf((*QueryValuer)(nil)).Elem()


var queryTag = "query"
func BindRequest(ptr interface{}, r *http.Request) error {
	query := r.URL.Query()
	queryLen := len(query)
	queryBindCount := 0
	err := greflect.DeepEach(ptr, func(rValue reflect.Value, rType reflect.Type, field reflect.StructField) (op greflect.EachOperator) {
		queryKey := field.Tag.Get(queryTag)
		if queryKey == "" {
			return
		}
		queryValue := query.Get(queryKey)
		if queryValue == "" {
			return
		}
		/* 转换赋值 */{
			if reflect.PtrTo(rType).Implements(queryValueType) {
				err := rValue.Addr().Interface().(QueryValuer).UnmarshalQuery(queryValue)
				if err != nil {
					return op.Error(err)
				}
			} else {
				err := gconv.StringReflect(queryValue, rValue)
				if err != nil {
					return op.Error(err)
				}
			}
			queryBindCount++
			if queryBindCount == queryLen {
				return op.Break()
			}
		}
		return
	})
	if err != nil {
		return err
	}
	return nil
}
