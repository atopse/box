package exec

import (
	"testing"

	"github.com/atopse/box/drivers"

	. "github.com/smartystreets/goconvey/convey"
)

func TestExec(t *testing.T) {
	driver := ExecDriver{}
	Convey("Exec Command", t, func() {
		action := drivers.Action{
			Input: drivers.Values{
				"command": "ls",
			},
		}
		output, err := driver.execute(&action)
		So(err, ShouldBeNil)
		So(output, ShouldHaveSameTypeAs, map[string]string{})
		data := output.(map[string]string)
		So(data, ShouldContainKey, "stdout")
		So(data, ShouldContainKey, "stderr")
		So(data["stdout"], ShouldContainSubstring, "execDriver_test.go")
		t.Logf("Exec Command Stdout:\n%s", data["stdout"])
	})
	Convey("Exec Command", t, func() {
		action := drivers.Action{
			Input: drivers.Values{
				"command": "miss",
			},
		}
		_, err := driver.execute(&action)
		So(err, ShouldNotBeNil)
	})
	Convey("Exec Command", t, func() {
		action := drivers.Action{
			Input: drivers.Values{
				"command": "ping www.baidu.com -c 1",
			},
		}
		output, err := driver.execute(&action)
		So(err, ShouldBeNil)
		data := output.(map[string]string)
		t.Logf("Exec Command Stdout:\n%s", data["stdout"])
	})
}
