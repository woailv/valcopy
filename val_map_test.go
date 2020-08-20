package algoutil

import (
	"encoding/json"
	"testing"
)

type A struct {
	Age  int
	Name string
	C    *CA
	CS   []CA
}

type CA struct {
	Sex int
	TT  *T
}

type B struct {
	Age  int
	Name string
	C    *CB
	CS   []CB
}

type CB struct {
	TT  *T
	Sex int
}

type T struct {
	H string
}

func TestValMap(t *testing.T) {
	b := &B{
		//CA: &CB{},
	}
	ValMap(
		&A{
			Age:  1,
			Name: "abc",
			C:    &CA{Sex: 123, TT: &T{"hello1"}},
			CS: []CA{
				{Sex: 1, TT: &T{"hello1"}},
				{Sex: 2, TT: &T{"hello2"}},
				{Sex: 3, TT: &T{"hello3"}},
			},
		},
		b, nil,
	)
	bs, _ := json.MarshalIndent(b, "", "  ")

	var a []B
	var c []A = []A{{
		Age:  1,
		Name: "abc",
		C:    &CA{Sex: 123, TT: &T{"hello1"}},
		CS: []CA{
			{Sex: 1, TT: &T{"hello1"}},
			{Sex: 2, TT: &T{"hello2"}},
			{Sex: 3, TT: &T{"hello3"}},
		},
	},{
		Age:  111111111,
		Name: "abc111111111111",
		C:    &CA{Sex: 123, TT: &T{"hello1"}},
		CS: []CA{
	{Sex: 1, TT: &T{"hello1"}},
	{Sex: 2, TT: &T{"hello2"}},
	{Sex: 3, TT: &T{"hello3"}},
	},
	}}
	ValMap(c, &a, nil)
	bs, _ = json.MarshalIndent(a, "", "  ")
	t.Log(string(bs))
}

var mapFunc = map[string]func(interface{}) interface{}{
	/*".Age": func(age interface{}) interface{} {
		return age.(int) + 1
	},
	".Name": func(i interface{}) interface{} {
		return i
	},
	".CS.Sex": func(srcVal interface{}) interface{} {
		return 100
	},
	".C.Sex": func(srcVal interface{}) interface{} {
		return 100
	},
	".C.Sex.T": func(srcVal interface{}) interface{} {
		return 100
	},
	".C.TT.H": func(srcVal interface{}) interface{} {
		return "100haha"
	},
	".CS.TT.H": func(srcVal interface{}) interface{} {
		return "yayayayayay"
	},*/
}
