package drivers

import "fmt"
import "github.com/atopse/comm"

// Action 驱动动作
type Action struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Driver    string `json:"driver"` //对应Driver的Namespace
	Namespace string `json:"namespace"`
	Input     Values `json:"inputValues"`
}

// NewAction 初始化一个Action
func NewAction(driverNamespace, actionNamespace string, values Values) (*Action, error) {
	d, err := GetDriver(driverNamespace)
	if err != nil {
		return nil, err
	}
	actions := d.Actions()
	for _, a := range actions {
		if a.Namespace == actionNamespace {
			action := Action{
				ID:     UUID(),
				Driver: d.Namespace(),
				Title:  a.Title,
				Input:  values,
			}
			// TODO: 验证输入值是否复核要求
			return &action, nil
		}
	}
	return nil, fmt.Errorf("NewAction:驱动%s(%s)下无%qAction", d.Title(), d.Namespace(), actionNamespace)
}

// ActionDescriptor Action描述信息
type ActionDescriptor struct {
	Title       string
	Namespace   string
	Description string
	Input       []InputDescriptor
	Output      []OutputDescriptor
}

// OptionDescriptor Option描述
type OptionDescriptor struct {
	Name        string         //参数名称
	Description string         //参数描述
	ValueType   comm.ValueType //参数值数据类型
}

// InputDescriptor 输入信息描述
type InputDescriptor struct {
	OptionDescriptor
	Mandatory bool //是否是必要参数
	// ValeRules   []m.Rule       //参数值验证信息
}

// OutputDescriptor 结果描述
type OutputDescriptor struct {
	OptionDescriptor
}
