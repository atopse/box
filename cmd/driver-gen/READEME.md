# Driver-gen
用于生成各驱动的配置信息，不需要人为处理。通过识别类、方法的注释信息自动生成对应的配置。

## 驱动配置项
每个配置信息均已`// @`开头。

+ @title 驱动标题
+ @namespace 驱动命名空间
+ @tags 以“#”开头标识的标签内容，每个标签使用空格分割
+ @desc 驱动描述信息
+ @ignore  如果不需要注册解析该驱动则可填写该标记。

例子：
```go
// ExecDriver 命令执行驱动器
// @title Command
// @namespace exec.driver.atopse
// @tags #command #bash
// @desc 用于执行操作系统的各项命令
type ExecDriver struct {
	drivers.Driver
}
```

## 驱动的Action配置信息

+ @title Action的Title
+ @namespace Action的命名空间
+ @param 执行Action的入参信息，内容依次是：name type [required] [description] ，其中required和description不是必须项。
+ @ignore 表示忽略该Action
+ @return Action返还的结果数据类型

示例：
```go
// execute 执行OS的cmd命令
// @title execute
// @namespace execute
// @param command string required 执行命令内容
// @param args string 配置信息
// @return map[string]string  命令执行结果信息
func (d *ExecDriver) execute(action *drivers.Action) (output interface{}, err error) {
    //...
}
```

## 使用方法
在要执行的自动生成的Driver文件下添加 `go:generate` 说明。在运行代码前需自行执行go命令`generate`。

Driver文件添加内容
```go
//go:generate driver-gen
```
执行代码
```bash
go generate packagePath
// go generate github.com/atopse/box/drivers/exec
```
执行`generate`后将在Driver所在的目录下生成对应的go文件，格式为：`gen_%Driver类名称%.go`,如：`gen_ExecDriver.go`