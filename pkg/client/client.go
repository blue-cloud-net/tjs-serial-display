package client

import (
	"github.com/blue-cloud-net/tjc-serial-display/internal/client"
	"github.com/blue-cloud-net/tjc-serial-display/pkg/models"
)

// 显示屏的客户端接口，定义了设备操作相关方法。
type DisplayClient interface {
	// 获取设备信息
	GetDeviceInfo() (*models.DeviceInfo, error)
	// 执行原始 TJC 命令
	ExecuteCommand(cmd string) ([]byte, error)
	// 升级面板程序
	Upgrade(programPath string, baudRate int, progressCallback models.UpgradeProgressCallback) error

	// 获取当前页面
	GetPage() (int, error)
	// 跳转到指定页面
	JumpPage(page int) error
	// 打印目标的值或者输入内容
	Prints(target string) (string, error)
	// 模拟弹起目标按钮
	ClickUp(target string) error
	// 模拟按下目标按钮
	ClickDown(target string) error
	// 隐藏指定目标
	Hide(target string) error
	// 显示指定目标
	Show(target string) error
}

func CreateClient(portName string, baudRate int) DisplayClient {
	return &client.TjcDisplayClient{
		PortName: portName,
		BaudRate: baudRate,
	}
}
