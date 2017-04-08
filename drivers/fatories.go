package drivers

// DriverFactoryInterface 用于创建驱动器.
type DriverFactoryInterface interface {
	// 驱动ID
	ID() string
	// 驱动名称
	Name() string
	// 驱动描述
	Description() string

	// 创建驱动所需的选项配置信息
	Options() []OptionDescriptor

	// 驱动所还有的行为描述
	Actions() []ActionDescriptor

	// New 新建驱动实例
	New(name, description string, options Options) DriverInterface
}

// DriverFactory 驱动Factory
type DriverFactory struct {
}
