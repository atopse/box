package stringDriver

import (
	"testing"

	"github.com/atopse/box/drivers"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFormat(t *testing.T) {
	driver := Driver{}
	Convey("数据非法测试", t, func() {
		Convey("format不允许为空", func() {
			action := drivers.Action{
				Namespace: "format",
				Input:     drivers.Values{"format": ""},
			}
			_, err := driver.ExecAction(&action)
			So(err, ShouldNotBeNil)

			action = drivers.Action{Namespace: "format", Input: drivers.Values{}}
			_, err = driver.ExecAction(&action)
			So(err, ShouldNotBeNil)
		})
		Convey("默认使用fmt.Sprintf方式", func() {
			action := drivers.Action{
				Namespace: "format",
				Input:     drivers.Values{"format": "data:%v", "data": 12},
			}
			output, err := driver.ExecAction(&action)
			So(err, ShouldBeNil)
			So(output, ShouldEqual, "data:12")
		})
	})
	Convey("Format之fmt.Sprintf", t, func() {

		cases := []struct {
			format string
			data   interface{}
			want   string
		}{
			{format: "%v", data: nil, want: "<nil>"},
			{format: "%s", data: "", want: ""},
			{format: "%s", data: "1", want: "1"},
			{format: "%d", data: 1, want: "1"},
			{format: "%v", data: true, want: "true"},
			{format: "%+v", data: user{"ysqi", 27}, want: "{Name:ysqi Age:27}"},
		}
		for _, c := range cases {
			action := drivers.Action{
				Namespace: "format",
				Input:     drivers.Values{"format": c.format, "fmtway": "fmt.sprintf", "data": c.data},
			}
			output, err := driver.ExecAction(&action)
			So(err, ShouldBeNil)
			So(output, ShouldEqual, c.want)
		}
	})

}

func TestFormat_HTML(t *testing.T) {
	driver := Driver{}

	Convey("Format之html.Template", t, func() {
		cases := []struct {
			format string
			data   interface{}
			want   string
		}{
			{format: "{{.}}", data: nil, want: ""},
			{format: "{{.}}", data: "", want: ""},
			{format: "{{.}}", data: "1", want: "1"},
			{format: "{{.}}", data: 1, want: "1"},
			{format: "{{.}}", data: true, want: "true"},
			{format: `"<{{printf "%.2f" .}}>"`, data: 1.0, want: `"<1.00>"`},
			{format: "{{.Name}}-{{.Age}}", data: user{"ysqi", 27}, want: "ysqi-27"},
			{format: "{{.Name}}-{{.Age}}", data: map[string]interface{}{"Name": "ysqi", "Age": 27}, want: "ysqi-27"},
			{format: "{{.name}}-{{.age}}", data: map[string]interface{}{"name": "ysqi", "age": 27}, want: "ysqi-27"},
		}
		for _, c := range cases {
			action := drivers.Action{
				Namespace: "format",
				Input:     drivers.Values{"format": c.format, "fmtway": "html.template", "data": c.data},
			}
			output, err := driver.ExecAction(&action)
			So(err, ShouldBeNil)
			So(output, ShouldEqual, c.want)
		}

	})
}

type user struct {
	Name string
	Age  int
}
