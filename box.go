package box

import (
	"errors"

	"strings"

	"github.com/atopse/box/drivers"
)

// ActionOption 信息
type ActionOption struct {
	Driver    string // driver namespace
	Action    string // action namespace
	OutputVar string // 定义输出结果存储的变量名
	Input     drivers.Values
}

// Box Driver的变种
type Box struct {
	Title       string          //魔盒名称,简称
	Namespace   string          //魔盒标识符，必须全局唯一。依次同其他魔盒器进行区分
	Description string          //魔盒描述，在使用魔盒时，可以显示对魔盒的描述，以方便了解魔盒功能
	Options     drivers.Options //魔盒器配置信息
	Actions     []ActionOption
	Input       drivers.Values
}

// New 新建驱动器
func New(title, namespace string, description ...string) (*Box, error) {
	if namespace == "" {
		return nil, errors.New("namespace不能为空")
	}
	if title == "" {
		return nil, errors.New("title不能为空")
	}
	return &Box{
		Title:       title,
		Namespace:   namespace,
		Description: strings.Join(description, ","),
		Actions:     []ActionOption{},
		Input:       make(drivers.Values),
		Options:     drivers.Options{},
	}, nil
}
func mapJoin(maps ...drivers.Values) drivers.Values {
	m := make(drivers.Values)
	for _, item := range maps {
		for k, v := range item {
			m[k] = v
		}
	}
	return m
}

// Exec 执行Action
func (b *Box) Exec() (output interface{}, err error) {
	if len(b.Actions) == 0 {
		return nil, errors.New("无Action可执行")
	}
	outputValues := make(drivers.Values)
	var lastOutput interface{}
	for _, a := range b.Actions {
		driver, err := drivers.GetDriver(a.Driver)
		if err != nil {
			return nil, err
		}
		input := mapJoin(outputValues, a.Input, b.Input)
		action, err := driver.NewAction(a.Action, input)
		if err != nil {
			return nil, err
		}
		lastOutput, err = driver.ExecAction(action)
		if err != nil {
			return nil, err
		}
		va := a.OutputVar
		if a.OutputVar == "" {
			va = a.Driver + "." + a.Action
		}
		outputValues[va] = lastOutput
	}
	return lastOutput, nil
}
