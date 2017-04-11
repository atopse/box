package main

const fileTpl = `//生成时间：{{.now}}

package {{.pkg}}

import (
	"github.com/atopse/box/drivers"
	"github.com/atopse/comm/kind"
) 
// Actions 驱动{{.structName}}所包含的Action信息
// @ignore
func (d *{{.structName}}) Actions() []drivers.ActionDescriptor {
	return []drivers.ActionDescriptor{	{{range .actions}}
			drivers.ActionDescriptor{ 
				Title:       "{{.title}}",
				Namespace:   "{{.namespace}}",
				Description: "{{.desc}}",
				Input: []drivers.InputDescriptor{ {{range .inputs}}
					{
						OptionDescriptor: drivers.OptionDescriptor{
							Name:        "{{.name}}",
							Description: "{{.desc}}",
							ValueType:   kind.{{.typ}},
						},
						Required: {{.required}},
					},{{end}}
				},
				Output: []drivers.OutputDescriptor{},
			},{{end}}
	}
} 
func init() { 
	drivers.RegisterDriver(&{{.structName}}{
		Driver: drivers.NewDriver(
			drivers.DriverConfig{
				Title:       "{{.title}}",
				Namespace:   "{{.namespace}}",
				Description: "{{.desc}}",
			},
		),
	} ) 
}
`
