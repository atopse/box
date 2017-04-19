//go:generate driver-gen

package stringDriver

import (
	"bytes"
	"errors"
	"html/template"
	"strings"

	"fmt"

	"html"

	"github.com/atopse/box/drivers"
)

// Driver 字符串处理驱动
// @title String
// @namespace string.driver.atopse
// @tags #string
// @desc 用于对字符串进行加工处理
type Driver struct {
	drivers.Driver
}

// ExecAction 执行Action
func (d *Driver) ExecAction(action *drivers.Action) (output interface{}, err error) {
	switch action.Namespace {
	case "format":
		return d.execFormatAction(action)
	case "":
		return nil, errors.New("Action的Namespace不能为空")
	default:
		return nil, errors.New("未知的Action:" + action.Namespace)
	}
}

// execFormatAction 格式化字符串
// @title 格式化
// @namespace format
// @param format string required 待格式化字符串文本
// @param fmtway string          格式化方式，所支持的方式有：fmt.Sprintf(缺省方式),html.Template
// @param data    any             带格式数据
// @return  string  已格式化内容
func (d *Driver) execFormatAction(action *drivers.Action) (output string, err error) {
	var (
		format string
		fmtway string
		data   interface{}
	)
	if err = action.Input.Bind("format", &format); err != nil {
		return "", err
	} else if format == "" {
		return "", errors.New("format参数不能为空")
	}

	if err = action.Input.Bind("fmtway", &fmtway); err != nil {
		if _, ok := err.(*drivers.NotFoundItemError); !ok {
			return "", err
		}
	}
	data = action.Input["data"]

	fmtway = strings.ToLower(strings.TrimSpace(fmtway))
	if fmtway == "" {
		fmtway = "fmt.sprintf"
	}

	switch fmtway {
	case "fmt.sprintf":
		return fmt.Sprintf(format, data), nil
	case "html.template":
		buf := bytes.NewBuffer(nil)
		//TODO:需要支持更多的模板函数
		err := template.Must(template.New(drivers.UUID()).Parse(format)).Execute(buf, data)
		if err != nil {
			return "", err
		}
		return html.UnescapeString(buf.String()), nil
	default:
		return "", errors.New("尚未实现格式化方式:" + fmtway)
	}

}
