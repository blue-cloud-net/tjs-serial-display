package models

// TJC串口显示屏设备信息
type DeviceInfo struct {
	Type                  int    // 屏幕类型（0:非触摸屏；1:电阻屏；2:电容屏）
	Address               string // 设备地址，唯一标识设备的通信地址
	Model                 string // 设备型号，如产品型号标识
	FirmwareVersion       int    // 固件版本号，表示设备固件的软件版本
	MainControlChipNumber int    // 主控芯片编号，主控 MCU 的唯一编号
	Number                string // 设备唯一编号，设备序列号
	FlashSize             int    // Flash 存储大小（单位：字节），设备内置存储容量
}
