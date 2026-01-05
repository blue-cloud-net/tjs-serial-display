package client

import "fmt"

// TjcError TJC串口屏错误
type TjcError struct {
	Code    byte   // 错误码
	Message string // 错误消息
}

// ResponseType 响应类型
type ResponseType int

const (
	ResponseTypeError   ResponseType = iota // 错误响应
	ResponseTypeSuccess                     // 成功响应
	ResponseTypeEvent                       // 事件响应
	ResponseTypeData                        // 数据响应
)

// Response TJC响应
type Response struct {
	Type    ResponseType // 响应类型
	Code    byte         // 响应码
	Data    []byte       // 响应数据
	RawData []byte       // 原始响应数据
}

func (e *TjcError) Error() string {
	return fmt.Sprintf("TJC Error 0x%02X: %s", e.Code, e.Message)
}

// toError 将响应转换为错误
func (r *Response) toError() error {
	if r.Type == ResponseTypeError {
		msg, ok := errorMessages[r.Code]
		if !ok {
			msg = "未知错误"
		}
		return &TjcError{
			Code:    r.Code,
			Message: msg,
		}
	}
	return nil
}

// 错误码到错误消息的映射
var errorMessages = map[byte]string{
	0x00: "无效指令",
	0x01: "指令成功执行",
	0x02: "控件ID无效",
	0x03: "页面ID无效",
	0x04: "图片ID无效",
	0x05: "字库ID无效",
	0x06: "文件操作失败",
	0x09: "CRC校验失败",
	0x11: "波特率设置无效",
	0x12: "曲线控件ID号或通道号无效",
	0x1A: "变量名称无效",
	0x1B: "变量运算无效",
	0x1C: "赋值操作失败",
	0x1D: "掉电存储空间操作失败",
	0x1E: "参数数量无效",
	0x1F: "IO操作失败",
	0x20: "转义字符使用错误",
	0x23: "变量名称太长",
	0x24: "串口缓冲区溢出",
}
