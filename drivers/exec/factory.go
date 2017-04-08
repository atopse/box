package exec

import (
	"github.com/atopse/box/drivers"
	"github.com/atopse/comm"
)

// Actions 返回Exec驱动所执行的行为描述
func (d *ExecDriver) Actions() []drivers.ActionDescriptor {
	return []drivers.ActionDescriptor{
		{
			Title:       "execute",
			Namespace:   "execute",
			Description: "在目标服务器执行cmd命令",
			Input: []drivers.InputDescriptor{
				{
					OptionDescriptor: drivers.OptionDescriptor{
						Name:        "command",
						Description: "命令内容",
						ValueType:   comm.VTSingleString,
					},
					Mandatory: true,
				},
			},
			Output: []drivers.OutputDescriptor{
				{
					OptionDescriptor: drivers.OptionDescriptor{
						Name:        "stdout",
						Description: "cmd输出结果",
						ValueType:   comm.VTSingleString,
					},
				},
				{
					OptionDescriptor: drivers.OptionDescriptor{
						Name:        "stderr",
						Description: "cmd输出结果",
						ValueType:   comm.VTSingleString,
					},
				},
			},
		},
	}
}

// NewDriver 新实例化
func NewDriver() drivers.DriverInterface {
	d := ExecDriver{
		Driver: drivers.NewDriver(
			drivers.DriverConfig{
				Title:       "cmd",
				Namespace:   "atopse.driver.cmd",
				Description: "用于执行操作系统的各项命令",
				Options: []drivers.Option{
					{
						Name: "tags",
						Desc: drivers.OptionDescriptor{
							Description: "标签信息",
						},
						Value: []string{"command"},
					},
				},
			},
		),
	}
	return &d
}

func init() {
	drivers.RegisterDriver(NewDriver())
}
