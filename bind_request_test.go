package juice

import (
	"encoding/base64"
	gtest "github.com/og/x/test"
	"net/http/httptest"
	"testing"
)

type URLBase64 struct {
	Value string
}
// 实现 juice.QueryValuer
func (b *URLBase64) UnmarshalQuery(queryValue string) error {
	valueBytes, err := base64.URLEncoding.DecodeString(queryValue)
	if err != nil {return err}
	b.Value = string(valueBytes)
	return nil
}
func TestBindRequest(t *testing.T) {
	as := gtest.NewAS(t)
	type UserID string
	type School struct {
		School string `query:"school"`
	}
	type Job struct {
		Title string `query:"jobTitle"`
	}
	type Req struct {
		Name string `query:"name"`
		Age uint `query:"age"`
		Elevation int `query:"elevation"`
		Happy bool `query:"happy"`
		UserID UserID `query:"userID"`
		School
		Job Job
		Website URLBase64 `query:"website"`
	}
	r := httptest.NewRequest(
		"GET",
		"http://github.com/og/juice?"+
			"name=nimoc&"+
			"age=27&"+
			"elevation=-100&"+
			"happy=true&"+
			"userID=a&"+
			"school=xjtu&"+
			"jobTitle=Programmer&"+
			"website=aHR0cHM6Ly9naXRodWIuY29tL25pbW9j",
			nil,
	)

	r.URL.Query()
	req := Req{}
	as.NoError(BindRequest(&req, r))
	as.Equal(req, Req{
		Name: "nimoc",
		Age: 27,
		Elevation: -100,
		Happy: true,
		UserID: UserID("a"),
		School: School{School: "xjtu"},
		Job: Job{Title: "Programmer"},
		Website: URLBase64{Value: "https://github.com/nimoc"},
	})
}
