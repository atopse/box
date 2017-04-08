package drivers

import "errors"

// Options 一组选项信息
type Options []Option

// Option 选项信息
type Option struct {
	Name  string           //参数名称
	Desc  OptionDescriptor //参数描述
	Value interface{}
}

// Value 从选项配置中获取对应项的值
func (opts Options) Value(name string) interface{} {
	for _, opt := range opts {
		if opt.Name == name {
			return opt.Value
		}
	}
	return nil
}

// SetValue 给一组参数名为 name 的设置值
func (opts Options) SetValue(name string, value interface{}) error {
	find := false
	for _, opt := range opts {
		if opt.Name == name {
			//TODO: 根据要求的参数数据类型进行数据判断
			opt.Value = value
			find = true
		}
	}

	if !find {
		return errors.New("选项中不不存在此项" + name)
	}
	return nil
}

// Bind 给选项绑定一值
func (opts Options) Bind(name string, dst interface{}) error {
	v := opts.Value(name)
	if v == nil {
		return errors.New("选项 " + name + " 不存在")
	}

	return ConvertValue(v, dst)
}
