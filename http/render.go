package jhttp

import (
	ogjson "github.com/og/json"
)

type Helper struct {

}
func (Helper) JSON(v interface{}) string {
	s, err := ogjson.StringWithErr(v)
	if err != nil {
		return "Error: render helper JSON Fail"
	}
	return s
}
