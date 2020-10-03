package vd_test

import (
	gtest "github.com/og/x/test"
	vd "github.com/og/juice/validator"
	"testing"
)

type SpecStringMinLen struct {
	Name string
}
func (s SpecStringMinLen) VD(r *vd.Rule) {
	r.String(s.Name, vd.StringSpec{
		Name:              "姓名",
		MinRuneLen:        4,
	})
};
type SpecStringMinLenCustomMessage struct {
	Name string
}
func (s SpecStringMinLenCustomMessage) VD(r *vd.Rule) {
	r.String(s.Name, vd.StringSpec{
		Name:              "姓名",
		MinRuneLen:        4,
		MinRuneLenMessage: "姓名长度不能小于{{MinRuneLen}}位,你输入的是{{Value}}",
	})
}
func Test_SpecString_MinLen(t *testing.T) {
	c := vd.NewCN()
	as := gtest.NewAS(t)
	as.Equal(c.Check(SpecStringMinLen{Name:"ni"}), vd.Report{
		Fail:    true,
		Message: "姓名长度不能小于4",
	})
	as.Equal(c.Check(SpecStringMinLen{Name:"nim"}), vd.Report{
		Fail:    true,
		Message: "姓名长度不能小于4",
	})
	as.Equal(c.Check(SpecStringMinLen{Name:"nimo"}), vd.Report{
		Fail:    false,
		Message: "",
	})
	as.Equal(c.Check(SpecStringMinLen{Name:"nimoc"}), vd.Report{
		Fail:    false,
		Message: "",
	})
	as.Equal(c.Check(SpecStringMinLenCustomMessage{Name:"ni"}), vd.Report{
		Fail:    true,
		Message: "姓名长度不能小于4位,你输入的是ni",
	})
}

type SpecStringMaxLen struct {
	Name string 
}
func (s SpecStringMaxLen) VD(r *vd.Rule) {
	r.String(s.Name, vd.StringSpec{
		Name:              "姓名",
		MaxRuneLen:        4,
	})
};
type SpecStringMaxLenCustomMessage struct {
	Name string
}
func (s SpecStringMaxLenCustomMessage) VD(r *vd.Rule) {
	r.String(s.Name, vd.StringSpec{
		Name:              "姓名",
		MaxRuneLen:        4,
		MaxRuneLenMessage: "姓名长度不能大于{{MaxRuneLen}}位,你输入的是{{Value}}",
	})
}
func Test_SpecString_MaxLen(t *testing.T) {
	c := vd.NewCN()
	as := gtest.NewAS(t)
	as.Equal(c.Check(SpecStringMaxLen{Name:"nimoc"}), vd.Report{
		Fail:    true,
		Message: "姓名长度不能大于4",
	})
	as.Equal(c.Check(SpecStringMaxLen{Name:"nimo"}), vd.Report{
		Fail:    false,
		Message: "",
	})
	as.Equal(c.Check(SpecStringMaxLen{Name:"nim"}), vd.Report{
		Fail:    false,
		Message: "",
	})
	as.Equal(c.Check(SpecStringMaxLenCustomMessage{Name:"nimoc"}), vd.Report{
		Fail:    true,
		Message: "姓名长度不能大于4位,你输入的是nimoc",
	})
}
type SpecStringPattern struct {
	Name string
	Title string
	More string 
}
func (s SpecStringPattern) VD (r *vd.Rule){
	r.String(s.Name, vd.StringSpec{
		Name:              "姓名",
		Pattern:		   []string{"^nimo"},
	})
	r.String(s.Title, vd.StringSpec{
		Name: "标题",
		Pattern: []string{`abc$`},
		PatternMessage: "{{Name}}必须以abc为结尾",
	})
	r.String(s.More, vd.StringSpec{
		AllowEmpty: true,
		Name: "更多",
		Pattern:[]string{`^a`, `a$`},
		PatternMessage: "{{Name}}开始结尾必须是a",
	})
}
func TestSpecStringPattern(t *testing.T) {
	as := gtest.NewAS(t)
	c := vd.NewCN()
	{
		as.Equal(c.Check(SpecStringPattern{
			Name: "nimo",
			Title: "abc",
		}), vd.Report{
			Fail:    true,
			Message: "更多开始结尾必须是a",
		})
	}
	{
		as.Equal(c.Check(SpecStringPattern{
			Name: "xnimo",
			Title: "abc",
		}), vd.Report{
			Fail:    true,
			Message: "姓名格式错误",
		})
	}
	{
		as.Equal(c.Check(SpecStringPattern{
			Name: "nimo",
			Title: "abcd",
		}), vd.Report{
			Fail:    true,
			Message: "标题必须以abc为结尾",
		})
	}
	{
		as.Equal(c.Check(SpecStringPattern{
			Name: "nimo",
			Title: "abcd",
			More: "c",
		}), vd.Report{
			Fail:    true,
			Message: "标题必须以abc为结尾",
		})
	}
}

type SpecStringBanPattern struct {
	Name string
	Title string
	More string
}
func (s SpecStringBanPattern) VD (r *vd.Rule){
	r.String(s.Name, vd.StringSpec{
		Name:              "姓名",
		BanPattern:		   []string{"fuck"},
		PatternMessage: "{{Name}}不允许出现敏感词",
	})
	r.String(s.Title, vd.StringSpec{
		Name: "标题",
		BanPattern: []string{`fuck`},
		PatternMessage: "{{Name}}不允许出现敏感词",
	})
	r.String(s.More, vd.StringSpec{
		AllowEmpty: true,
		Name: "更多",
		BanPattern: []string{`fuck`, `dick`},
		PatternMessage: "{{Name}}不允许出现敏感词:{{BanPattern}}",
	})
}
func TestSpecStringBanPattern(t *testing.T) {
	as := gtest.NewAS(t)
	c := vd.NewCN()
	{
		as.Equal(c.Check(SpecStringBanPattern{
			Name: "nimo",
			Title: "nimo",
			More: "nimo",
		}), vd.Report{
			Fail:    false,
			Message: "",
		})
	}
	{
		as.Equal(c.Check(SpecStringBanPattern{
			Name: "fuck",
			Title: "nimo",
			More: "nimo",
		}), vd.Report{
			Fail:    true,
			Message: "姓名不允许出现敏感词",
		})
	}
	{
		as.Equal(c.Check(SpecStringBanPattern{
			Name: "nimo",
			Title: "fuck",
			More: "nimo",
		}), vd.Report{
			Fail:    true,
			Message: "标题不允许出现敏感词",
		})
	}
	{
		as.Equal(c.Check(SpecStringBanPattern{
			Name: "nimo",
			Title: "nimo",
			More: "fuck",
		}), vd.Report{
			Fail:    true,
			Message: "更多不允许出现敏感词:[fuck dick]",
		})
	}
	{
		as.Equal(c.Check(SpecStringBanPattern{
			Name: "nimo",
			Title: "nimo",
			More: "dick",
		}), vd.Report{
			Fail:    true,
			Message: "更多不允许出现敏感词:[fuck dick]",
		})
	}
}
type SpecStringEnum struct {
	Type string
}
func (SpecStringEnum) Dict() (dict struct{
	Type struct {
		Normal string
		Danger string
	}
}) {
	dict.Type.Normal = "normal"
	dict.Type.Danger = "danger"
	return
}
func (s SpecStringEnum) VD(r *vd.Rule) {
	r.String(s.Type, vd.StringSpec{
		Name: "类型",
		Enum: vd.EnumValues(s.Dict().Type),
	})
}
func TestStringSpec_CheckEnum (t *testing.T) {
	as := gtest.NewAS(t)
	c := vd.NewCN()
	as.Equal(c.Check(SpecStringEnum{
		Type: "normal1",
	}), vd.Report{
		Fail:    true,
		Message: "类型参数错误，只允许(normal danger)",
	})
}
type SpecStringMinMax struct {
	Name string
}
func (v SpecStringMinMax) VD(r *vd.Rule) {
	r.String(v.Name, vd.StringSpec{
		Name:              "姓名",
		AllowEmpty: 	   true,
		MinRuneLen:        2,
		MaxRuneLen:        4,
	})
}
func TestSpectStringMinMax(t *testing.T) {
	as := gtest.NewAS(t)
	c := vd.NewCN()
	as.Equal(c.Check(SpecStringMinMax{
		Name: "",
	}), vd.Report{
		Fail:    true,
		Message: "姓名长度不能小于2",
	})
	as.Equal(c.Check(SpecStringMinMax{
		Name: "1",
	}), vd.Report{
		Fail:    true,
		Message: "姓名长度不能小于2",
	})
	as.Equal(c.Check(SpecStringMinMax{
		Name: "12",
	}), vd.Report{
		Fail:    false,
		Message: "",
	})
	as.Equal(c.Check(SpecStringMinMax{
		Name: "123",
	}), vd.Report{
		Fail:    false,
		Message: "",
	})
	as.Equal(c.Check(SpecStringMinMax{
		Name: "1234",
	}), vd.Report{
		Fail:    false,
		Message: "",
	})
	as.Equal(c.Check(SpecStringMinMax{
		Name: "12345",
	}), vd.Report{
		Fail:    true,
		Message: "姓名长度不能大于4",
	})

}

type SpecStringEmail struct {
	Email string
	OtherEmail string
}
func (v SpecStringEmail) VD(r *vd.Rule) {
	r.String(v.Email, vd.StringSpec{
		Name: "邮箱",
		Ext:  []vd.StringSpec{
			vd.Email(),
		},
	})
	r.String(v.OtherEmail, vd.Email().NameIs("附属邮箱"))
}
func TestStringEmail(t *testing.T) {
	as := gtest.NewAS(t)
	_=as
	c := vd.NewCN()
	as.Equal(c.Check(SpecStringEmail{
		Email: "12345",
		OtherEmail: "mail@github.com",
	}), vd.Report{
		Fail:    true,
		Message: "邮箱格式错误",
	})
	as.Equal(c.Check(SpecStringEmail{
		Email: "12345@qq.com",
		OtherEmail: "mailithub.com",
	}), vd.Report{
		Fail:    true,
		Message: "附属邮箱格式错误",
	})
}
