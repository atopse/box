package drivers

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/atopse/comm/kind"
)

// NotFoundItemError 不存在该项Error
type NotFoundItemError struct {
	Item string
}

func (n *NotFoundItemError) Error() string {
	if n.Item == "" {
		return "未找到项"
	}
	return fmt.Sprintf("未找到项%q", n.Item)
}

// OutputType 数据输出类型
type OutputType kind.Kind

// Output 输出信息
type Output struct {
	ID       string      `json:"id"`
	ActionID string      `json:"action_id"`
	Type     OutputType  `json:"typ"`
	Value    interface{} `json:"value"`
}

// Values 驱动器选项输入值
type Values map[string]interface{}

// Bind 将Map的值转换为指定类型的数据，如果未找到Key所对应的项或转换失败则返回错误
func (v Values) Bind(key string, dst interface{}) error {
	value, ok := v[key]
	if !ok {
		return &NotFoundItemError{key}
	}
	return ConvertValue(value, dst)
}

// ConvertValue 尝试将 v 转换为 dst
func ConvertValue(v interface{}, dst interface{}) error {
	switch d := dst.(type) {
	case *string:
		switch vt := v.(type) {
		case string:
			*d = vt
		case []string:
			*d = strings.Join(vt, ",")
		case bool:
			*d = strconv.FormatBool(vt)
		case int64:
			*d = strconv.FormatInt(vt, 10)
		case float64:
			*d = strconv.FormatFloat(vt, 'f', -1, 64)
		case int:
			*d = strconv.FormatInt(int64(vt), 10)
		default:
			return fmt.Errorf("无法将类型 %+v 转换为string", vt)
		}

	case *[]string:
		switch vt := v.(type) {
		case []string:
			*d = vt
		case string:
			*d = strings.Split(vt, ",")
		default:
			return fmt.Errorf("无法将类型 %+v 转换为[]string", vt)
		}

	case *bool:
		switch vt := v.(type) {
		case bool:
			*d = vt
		case string:
			vt = strings.ToLower(vt)
			if vt == "true" || vt == "on" || vt == "yes" || vt == "1" || vt == "t" {
				*d = true
			}
		case int64:
			*d = vt > 0
		case int:
			*d = vt > 0
		case uint64:
			*d = vt > 0
		case uint:
			*d = vt > 0
		case float64:
			*d = vt > 0
		default:
			return fmt.Errorf("无法将类型 %+v 转换为bool", vt)
		}

	case *float64:
		switch vt := v.(type) {
		case int64:
			*d = float64(vt)
		case int32:
			*d = float64(vt)
		case int16:
			*d = float64(vt)
		case int8:
			*d = float64(vt)
		case int:
			*d = float64(vt)
		case uint64:
			*d = float64(vt)
		case uint32:
			*d = float64(vt)
		case uint16:
			*d = float64(vt)
		case uint8:
			*d = float64(vt)
		case uint:
			*d = float64(vt)
		case float64:
			*d = vt
		case float32:
			*d = float64(vt)
		case string:
			x, _ := strconv.Atoi(vt)
			*d = float64(x)
		default:
			return fmt.Errorf("无法将类型 %+v 转换为float64", vt)
		}

	case *int:
		switch vt := v.(type) {
		case int64:
			*d = int(vt)
		case int32:
			*d = int(vt)
		case int16:
			*d = int(vt)
		case int8:
			*d = int(vt)
		case int:
			*d = vt
		case uint64:
			*d = int(vt)
		case uint32:
			*d = int(vt)
		case uint16:
			*d = int(vt)
		case uint8:
			*d = int(vt)
		case uint:
			*d = int(vt)
		case float64:
			*d = int(vt)
		case float32:
			*d = int(vt)
		case string:
			*d, _ = strconv.Atoi(vt)
		default:
			return fmt.Errorf("无法将类型 %+v 转换为int", vt)
		}

	case *url.Values:
		switch vt := v.(type) {
		case string:
			*d, _ = url.ParseQuery(vt)
		default:
			return fmt.Errorf("无法将类型 %+v 转换为 url.Values", vt)
		}
	case interface{}:
		dst = v
	default:
		return fmt.Errorf("无法处理的dst数据类型 %+v", dst)
	}

	return nil
}
