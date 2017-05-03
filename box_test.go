package box

import (
	"testing"

	_ "github.com/atopse/box/drivers/exec"

	"os"

	"github.com/atopse/box/drivers"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFaildBox(t *testing.T) {

	Convey("数据非法测试", t, func() {
		Convey("Title不能为空", func() {
			_, err := New("", "")
			So(err, ShouldNotBeNil)
		})
		Convey("无Action可执行", func() {
			box, err := New("title", "one.box.atopse")
			So(err, ShouldBeNil)

			_, err = box.Exec()
			So(err, ShouldNotBeNil)
		})
		Convey("Driver无效", func() {
			box, err := New("title", "one.box.atopse")
			So(err, ShouldBeNil)
			box.Actions = []ActionOption{
				{Driver: "notfound", Action: "execute"},
			}
			_, err = box.Exec()
			So(err, ShouldNotBeNil)
		})
		Convey("Action无效", func() {
			box, err := New("title", "one.box.atopse")
			So(err, ShouldBeNil)
			box.Actions = []ActionOption{
				{Driver: "exec.driver.atopse", Action: "notfound"},
			}
			_, err = box.Exec()
			So(err, ShouldNotBeNil)
		})

	})
}

func TestBoxExec(t *testing.T) {
	Convey("单个Action", t, func() {
		box, err := New("title", "one.box.atopse")
		So(err, ShouldBeNil)
		box.Actions = []ActionOption{
			{Driver: "exec.driver.atopse", Action: "execute", Input: drivers.Values{"command": "ls"}},
		}
		output, err := box.Exec()
		So(err, ShouldBeNil)
		So(output, ShouldNotBeNil)

		result := output.(string)
		So(result, ShouldContainSubstring, "box_test.go")
	})
	Convey("多个Action关联执行", t, func() {
		box, err := New("title", "one.box.atopse")
		So(err, ShouldBeNil)
		box.Actions = []ActionOption{
			{Driver: "exec.driver.atopse", Action: "execute", Input: drivers.Values{"command": "go list"}, OutputVar: "pkg"},
			{Driver: "exec.driver.atopse", Action: "execute", Input: drivers.Values{"command": "ls {{.gopath}}/src/{{.pkg}}", "gopath": os.Getenv("GOPATH")}},
		}
		output, err := box.Exec()
		So(err, ShouldBeNil)
		result := output.(string)
		t.Logf("%+v", result)
		So(result, ShouldContainSubstring, "box_test.go")
	})
}
