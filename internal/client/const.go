package client

// TJC串口屏支持波特率
var SupportedBaudrate = []int{2400, 4800, 9600, 19200, 38400, 57600, 115200, 230400, 256000, 512000, 921600}

// 返回错误码定义（bkcmd非0时）
const (
	CodeInvalidInstruction   = 0x00 // 无效指令
	CodeSuccess              = 0x01 // 指令成功执行
	CodeInvalidComponentID   = 0x02 // 控件ID无效
	CodeInvalidPageID        = 0x03 // 页面ID无效
	CodeInvalidPictureID     = 0x04 // 图片ID无效
	CodeInvalidFontID        = 0x05 // 字库ID无效
	CodeFileOperationFailed  = 0x06 // 文件操作失败
	CodeCRCCheckFailed       = 0x09 // CRC校验失败
	CodeInvalidBaudrate      = 0x11 // 波特率设置无效
	CodeInvalidCurveID       = 0x12 // 曲线控件ID号或通道号无效
	CodeInvalidVariableName  = 0x1A // 变量名称无效
	CodeInvalidVariableOp    = 0x1B // 变量运算无效
	CodeAssignmentFailed     = 0x1C // 赋值操作失败
	CodeEEPROMFailed         = 0x1D // 掉电存储空间操作失败
	CodeInvalidParamCount    = 0x1E // 参数数量无效
	CodeIOFailed             = 0x1F // IO操作失败
	CodeEscapeCharError      = 0x20 // 转义字符使用错误
	CodeVariableNameTooLong  = 0x23 // 变量名称太长
	CodeSerialBufferOverflow = 0x24 // 串口缓冲区溢出
)

// 事件返回码定义（不受bkcmd影响）
const (
	CodeTouchEvent       = 0x65 // 控件点击事件返回
	CodePageID           = 0x66 // 当前页面的ID号返回
	CodeTouchCoordinate  = 0x67 // 触摸坐标数据返回
	CodeSleepTouch       = 0x68 // 睡眠模式触摸事件
	CodeStringData       = 0x70 // 字符串变量数据返回
	CodeNumberData       = 0x71 // 数值变量数据返回
	CodeAutoSleep        = 0x86 // 设备自动进入睡眠模式
	CodeAutoWake         = 0x87 // 设备自动唤醒
	CodeStartupSuccess   = 0x88 // 系统启动成功
	CodeStartSDUpgrade   = 0x89 // 开始SD卡升级
	CodeTransparentReady = 0xFE // 数据透传就绪
	CodeTransparentDone  = 0xFD // 透传数据完成
)
