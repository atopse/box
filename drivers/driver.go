package drivers

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"

	"bytes"

	"reflect"

	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/atopse/comm/bee"
	uuid "github.com/nu7hatch/gouuid"
)

// DriverInterface 驱动器接口
type DriverInterface interface {

	// Title 此驱动名称
	Title() string
	// Namespace 驱动标识符
	Namespace() string

	// 驱动器描述
	Description() string

	// 获取驱动器配置选项信息
	Options() Options

	Logln(args ...interface{})
	Logf(format string, args ...interface{})
	LogErrorf(format string, args ...interface{})
	LogFatal(args ...interface{})

	// Actions 行为描述
	Actions() []ActionDescriptor

	NewAction(actionNamespace string, input Values) (*Action, error)

	// Action 执行行为命令
	ExecAction(action *Action) (interface{}, error)

	String() string
}

// DriverConfig 驱动器配置
type DriverConfig struct {
	Title       string  //驱动名称,简称
	Namespace   string  //驱动标识符，必须全局唯一。依次同其他驱动器进行区分
	Description string  //驱动描述，在使用驱动时，可以显示对驱动的描述，以方便了解驱动功能
	Options     Options //驱动器配置信息
}

// Driver 驱动
type Driver struct {
	config DriverConfig
}

var (
	drivers = make(map[string]DriverInterface)
)

// RegisterDriver 将驱动器注册到全局中, 注册时需保证Namespace唯一。
func RegisterDriver(driver DriverInterface) {
	if driver.Title() == "" {
		logs.Warn("拒绝注册Title为空的驱动器:%s", driver.String())
		return
	}
	namespace := driver.Namespace()
	if _, ok := drivers[namespace]; ok {
		logs.Warn("拒绝重复注册驱动器:%s", driver.String())
		return
	}
	logs.Info("受理驱动器注册:%s", driver.String())
	actions := driver.Actions()
	for _, a := range actions {
		logs.Info("驱动器:%s,Action: %+v", driver.Namespace(), a)
	}
	drivers[namespace] = driver
}

// GetDrivers 获取所有已注册的Drivers
func GetDrivers() []DriverInterface {
	items := []DriverInterface{}
	for _, v := range drivers {
		if v == nil {
			continue
		}
		items = append(items, v)
	}
	return items
}

// GetDriver 获取指定类型的驱动器
func GetDriver(namespace string) (DriverInterface, error) {
	d, ok := drivers[namespace]
	if !ok {
		return nil, errors.New("全局环境中无此驱动 " + namespace + " ,请确认在使用前已注册该驱动")
	}
	return d, nil
}

// NewDriver 新建驱动器
func NewDriver(cfg DriverConfig) Driver {
	if cfg.Title == "" {
		panic("驱动器名称不能为空")
	}
	if cfg.Namespace == "" {
		panic("驱动器Namespace不能为空")
	}

	return Driver{
		config: cfg,
	}
}

// Title 此驱动名称
func (d *Driver) Title() string {
	return d.config.Title
}

// Namespace 获取驱动器命名空间，以此注册驱动器
func (d *Driver) Namespace() string {
	return d.config.Namespace
}

// Description 获取驱动器描述
func (d *Driver) Description() string {
	return d.config.Description
}

// Options 获取驱动器配置选项信息
func (d *Driver) Options() Options {
	return d.config.Options
}

// Logln 写日志
func (d *Driver) Logln(args ...interface{}) {
	logs.Info("[", d.Title(), "]:", args)
}

// Logf logs a formatted string
func (d *Driver) Logf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	d.Logln(s)
}

// LogErrorf logs a formatted error string
func (d *Driver) LogErrorf(format string, args ...interface{}) {
	logs.Error("[", d.Title(), "]:", fmt.Sprintf(format, args...))
}

// LogFatal logs a fatal error
func (d *Driver) LogFatal(args ...interface{}) {
	logs.Error("[Fatal][", d.Title(), "]:", args)
}

// Tags 获取驱动器标签
// func (d *Driver) Tags() []string {
// 	return d.config.Tags
// }

// Actions 返回Action描述信息
// func (d *Driver) Actions() []ActionDescriptor {
// 	panic("")
// }

// NewAction 初始化一个Action
func (d *Driver) NewAction(actionNamespace string, values Values) (*Action, error) {
	dd, err := GetDriver(d.Namespace())
	if err != nil {
		return nil, err
	}
	actions := dd.Actions()
	if len(actions) == 0 {
		return nil, fmt.Errorf("NewAction:驱动%s(%s)下无Action", d.Title(), d.Namespace())
	}
	for _, a := range actions {
		fmt.Println(a.Namespace)
		if a.Namespace == actionNamespace {
			action := Action{
				ID:        UUID(),
				Driver:    d.Namespace(),
				Namespace: a.Namespace,
				Title:     a.Title,
				Input:     values,
			}
			// TODO: 验证输入值是否复核要求
			return &action, nil
		}
	}
	return nil, fmt.Errorf("NewAction:驱动%s(%s)下无%qAction", d.Title(), d.Namespace(), actionNamespace)
}

// ExecAction 执行Action，但基础Driver中不具体实现。
func (d *Driver) ExecAction(a *Action) (Output, error) {
	return Output{}, errors.New("driver:尚未尚未实现ExecAction")
}

// ActionHandler 执行Action
func (d *Driver) ActionHandler(ctx *bee.Context) {
	var action *Action
	err := json.Unmarshal(ctx.Input.RequestBody, action)
	if err != nil {
		ctx.JSON(err)
		return
	}
	if action.Namespace == "" {
		ctx.JSON(errors.New("Action的Name为空,无法进行Action匹配查询"))
		return
	}

	action, err = d.NewAction(action.Namespace, action.Input)
	if err != nil {
		ctx.JSON(err)
		return
	}
	output, err := d.ExecAction(action)
	if err != nil {
		ctx.JSON(err)
		return
	}
	data := Output{
		ID:       UUID(),
		ActionID: action.ID,
		Value:    output,
	}
	ctx.JSON(data)
}

func (d *Driver) String() string {
	if d.Description() == "" {
		return fmt.Sprintf("%s(%s)", d.Title(), d.Namespace())
	}
	return fmt.Sprintf("%s(%s)-%s", d.Title(), d.Namespace(), d.Description())
}

// UUID 生成一个唯一ID
func UUID() string {
	u, _ := uuid.NewV4()
	return u.String()
}

// ConvertInput 解析Input输入
func ConvertInput(input interface{}, data interface{}) (interface{}, error) {
	switch v := input.(type) {
	case string:
		if !strings.Contains(v, "{") {
			return v, nil
		}
		buf := bytes.NewBuffer(nil)
		err := template.Must(template.New(UUID()).Parse(v)).Execute(buf, data)
		if err != nil {
			return nil, err
		}

		return buf.String(), nil
	}
	return nil, fmt.Errorf("尚未实现对类型：%s的解析", reflect.TypeOf(input))
}
