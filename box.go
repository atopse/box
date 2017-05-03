package box

import (
	"errors"

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
	drivers.Driver
	Actions []ActionOption
	Input   drivers.Values
}

// Config Box配置
type Config drivers.DriverConfig

// New 新建驱动器
func New(cfg Config) Box {
	if cfg.Title == "" {
		panic("Box名称不能为空")
	}
	if cfg.Namespace == "" {
		cfg.Namespace = "one.box.atopse"
	}
	d := drivers.NewDriver(drivers.DriverConfig(cfg))
	return Box{
		Driver: d,
	}
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
