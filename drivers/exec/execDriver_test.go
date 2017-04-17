package exec

import (
	"testing"

	"github.com/atopse/box/drivers"

	. "github.com/smartystreets/goconvey/convey"
)

func TestExec(t *testing.T) {
	driver := Driver{}
	Convey("Exec Command", t, func() {
		action := drivers.Action{
			Input: drivers.Values{
				"command": "ls",
			},
		}
		output, err := driver.execute(&action)
		So(err, ShouldBeNil)
		So(output, ShouldContainSubstring, "execDriver_test.go")
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
	Convey("Exec Command with var", t, func() {
		action := drivers.Action{
			Input: drivers.Values{
				"command": "go build {{.pkg}}",
				"pkg":     "github.com/atopse/box/drivers/exec_notfound",
			},
		}
		output, err := driver.execute(&action)
		So(output, ShouldBeEmpty)
		So(err.Error(), ShouldStartWith, "can't load package: package github.com/atopse/box/drivers/exec_notfound")
	})
}
