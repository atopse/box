package drivers

import "github.com/atopse/comm/kind"
import "fmt"

// Action 驱动动作
type Action struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Driver    string `json:"driver"` //对应Driver的Namespace
	Namespace string `json:"namespace"`
	Input     Values `json:"inputValues"`
}

func (a *Action) String() string {
	return fmt.Sprintf("%s.%s.%s", a.Title, a.ID, a.Namespace)
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
	Name        string    //参数名称
	Description string    //参数描述
	ValueType   kind.Kind //参数值数据类型
}

// InputDescriptor 输入信息描述
type InputDescriptor struct {
	OptionDescriptor
	Required bool //是否是必要参数
	// ValeRules   []m.Rule       //参数值验证信息
}

// OutputDescriptor 结果描述
type OutputDescriptor struct {
	OptionDescriptor
}
