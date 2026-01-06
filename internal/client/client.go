package client

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/blue-cloud-net/tjc-serial-display/internal/serial"
	"github.com/blue-cloud-net/tjc-serial-display/pkg/consts"
	"github.com/blue-cloud-net/tjc-serial-display/pkg/models"
)

// 控制命令的起始和结束符
const (
	EndStr = "\xFF\xFF\xFF"
)

var (
	EndSymbol = []byte(EndStr)
)

type TjcDisplayClient struct {
	PortName      string
	BaudRate      int
	Timeout       time.Duration
	serialManager *serial.SerialPortManager
	optLock       sync.Mutex
}

func (c *TjcDisplayClient) connect() error {
	ports, err := serial.ListPorts()
	if err != nil {
		return err
	}

	// 检查是否存在指定的串口
	found := slices.Contains(ports, c.PortName)
	if !found {
		return fmt.Errorf("serial port %s not found", c.PortName)
	}

	// 检查波特率是否支持
	baudrateSupported := slices.Contains(consts.SupportedBaudrate, c.BaudRate)
	if !baudrateSupported {
		return fmt.Errorf("baud rate %d is not supported", c.BaudRate)
	}

	if c.serialManager == nil {
		manager := &serial.SerialPortManager{
			PortName: c.PortName,
			BaudRate: c.BaudRate,
			Timeout:  c.Timeout,
		}

		err := manager.Open()

		if err != nil {
			return err
		}

		c.serialManager = manager
	}

	if !c.serialManager.IsOpen() {
		return c.serialManager.Open()
	}

	// 退出主动解析模式
	_ = c.sendCommand("DRAKJHSUYDGBNCJHGJKSHBDN", false)

	return nil
}

// GetDeviceInfo 获取设备信息（示例实现，实际需根据协议解析串口返回数据）
func (c *TjcDisplayClient) GetDeviceInfo() (*models.DeviceInfo, error) {
	err := c.connect()
	if err != nil {
		return nil, err
	}

	// 发送connect
	result, err := c.sendCommandAndWaitResult("connect", false)
	if err != nil {
		return nil, err
	}

	return parseDeviceInfo(string(result))
}

// parseDeviceInfo 解析设备信息字符串
func parseDeviceInfo(data string) (*models.DeviceInfo, error) {
	// 以TJC4024T032_011R设备为例，设备返回如下8组数据(每组数据逗号隔开):
	// comok 1,101-0,TJC4024T032_011R,52,61488,D264B8204F0E1828,16777216
	// 格式说明:
	// comok [屏幕类型],[设备地址],[设备型号],[固件版本号],[主控芯片编号],[设备唯一编号],[Flash存储大小]

	// 去掉开头的 "comok " 或其他前缀
	data = strings.TrimSpace(data)
	if strings.HasPrefix(data, "comok ") {
		data = strings.TrimPrefix(data, "comok ")
	} else if strings.HasPrefix(data, "comok") {
		data = strings.TrimPrefix(data, "comok")
	}

	data = strings.TrimSpace(data)

	// 按逗号分割
	parts := strings.Split(data, ",")
	if len(parts) < 7 {
		return nil, fmt.Errorf("invalid device info format, expected 7 fields, got %d: %s", len(parts), data)
	}

	deviceInfo := &models.DeviceInfo{}

	// 解析屏幕类型
	if screenType, err := strconv.Atoi(strings.TrimSpace(parts[0])); err == nil {
		deviceInfo.Type = screenType
	}

	// 设备地址
	deviceInfo.Address = strings.TrimSpace(parts[1])

	// 设备型号
	deviceInfo.Model = strings.TrimSpace(parts[2])

	// 固件版本号
	if version, err := strconv.Atoi(strings.TrimSpace(parts[3])); err == nil {
		deviceInfo.FirmwareVersion = version
	}

	// 主控芯片编号
	if chipNum, err := strconv.Atoi(strings.TrimSpace(parts[4])); err == nil {
		deviceInfo.MainControlChipNumber = chipNum
	}

	// 设备唯一编号
	deviceInfo.Number = strings.TrimSpace(parts[5])

	// Flash 存储大小
	if flashSize, err := strconv.Atoi(strings.TrimSpace(parts[6])); err == nil {
		deviceInfo.FlashSize = flashSize
	}

	return deviceInfo, nil
}

// GetPage 获取当前页面
func (c *TjcDisplayClient) GetPage() (int, error) {
	err := c.connect()
	if err != nil {
		return 0, err
	}

	result, err := c.sendCommandAndWaitResult("sendme", false)
	if err != nil {
		return 0, err
	}

	if len(result) != 1 {
		return 0, fmt.Errorf("invalid response length for GetPage: %d", len(result))
	}

	page := int(result[0])

	return page, nil
}

// JumpPage 跳转到指定页面
func (c *TjcDisplayClient) JumpPage(page int) error {
	err := c.connect()
	if err != nil {
		return err
	}

	return c.sendCommand(fmt.Sprintf("page %d", page), false)
}

// Prints 打印目标的值或者输入内容
func (c *TjcDisplayClient) Prints(target string) (string, error) {
	err := c.connect()
	if err != nil {
		return "", err
	}

	result, err := c.sendCommandAndWaitResult(fmt.Sprintf("print %s", target), true)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

// ClickUp 模拟弹起目标按钮
func (c *TjcDisplayClient) ClickUp(target string) error {
	err := c.connect()
	if err != nil {
		return nil
	}

	return c.sendCommand(fmt.Sprintf("click %s,0", target), false)
}

// ClickDown 模拟按下目标按钮
func (c *TjcDisplayClient) ClickDown(target string) error {
	err := c.connect()
	if err != nil {
		return err
	}

	return c.sendCommand(fmt.Sprintf("click %s,1", target), false)
}

// Hide 隐藏指定目标
func (c *TjcDisplayClient) Hide(target string) error {
	err := c.connect()
	if err != nil {
		return err
	}

	return c.sendCommand(fmt.Sprintf("vis %s,0", target), false)
}

// Show 显示指定目标
func (c *TjcDisplayClient) Show(target string) error {
	err := c.connect()
	if err != nil {
		return err
	}

	return c.sendCommand(fmt.Sprintf("vis %s,1", target), false)
}

// ExecuteCommand 执行原始 TJC 命令
func (c *TjcDisplayClient) ExecuteCommand(cmd string) ([]byte, error) {
	err := c.connect()
	if err != nil {
		return nil, err
	}

	return c.sendCommandAndWaitRawResult(cmd, false)
}

// Upgrade 升级面板程序
func (c *TjcDisplayClient) Upgrade(programPath string, baudRate int, progressCallback models.UpgradeProgressCallback) error {
	if baudRate == 0 {
		baudRate = 921600
	}

	err := c.connect()
	if err != nil {
		return err
	}

	f, err := os.Open(programPath)
	if err != nil {
		return fmt.Errorf("failed to open program file: %w", err)
	}
	defer f.Close()

	fileInfo, err := f.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	fileSize := fileInfo.Size()

	// 发送 whmi-wri 命令（使用当前连接的波特率）
	cmd := []byte(fmt.Sprintf("whmi-wri %d,%d,0", fileSize, baudRate))
	cmd = append(cmd, EndSymbol...)
	err = c.serialManager.Write(cmd)
	if err != nil {
		return fmt.Errorf("failed to initiate upgrade: %w", err)
	}

	// 等待350ms，确保设备已准备好
	time.Sleep(350 * time.Millisecond)

	// 切换到下载波特率
	err = c.serialManager.SetBaudRate(baudRate)
	if err != nil {
		return fmt.Errorf("failed to set new baud rate: %w", err)
	}

	// 读取响应，等待设备准备好接收数据（0x05 = ENQ）
	resp, err := c.serialManager.ReadExactly(1)
	if err != nil {
		return fmt.Errorf("failed to read upgrade response: %w", err)
	}

	if len(resp) != 1 || resp[0] != 0x05 {
		return fmt.Errorf("unexpected upgrade response: got 0x%02X, expected 0x05", resp[0])
	}

	// 初始化进度信息
	var totalSent int64 = 0
	startTime := time.Now()

	// 初始进度回调
	if progressCallback != nil {
		progressCallback(&models.UpgradeProgress{
			Current:    0,
			Total:      fileSize,
			Percentage: 0.0,
			Speed:      0,
			Elapsed:    0,
			Remaining:  0,
		})
	}

	buf := make([]byte, 4096)
	for {
		n, err := f.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return fmt.Errorf("failed to read program file: %w", err)
			}
		}

		if n > 0 {
			// 写入数据块
			err = c.serialManager.Write(buf[:n])
			if err != nil {
				return fmt.Errorf("failed to write program data at byte %d: %w", totalSent, err)
			}

			// 更新已发送字节数
			totalSent += int64(n)

			// 计算进度信息并调用回调
			if progressCallback != nil {
				elapsed := time.Since(startTime)
				percentage := float64(totalSent) / float64(fileSize) * 100

				var speed int64
				var remaining time.Duration
				if elapsed.Seconds() > 0 {
					speed = int64(float64(totalSent) / elapsed.Seconds())
					if speed > 0 {
						remainingBytes := fileSize - totalSent
						remaining = time.Duration(float64(remainingBytes)/float64(speed)) * time.Second
					}
				}

				progressCallback(&models.UpgradeProgress{
					Current:    totalSent,
					Total:      fileSize,
					Percentage: percentage,
					Speed:      speed,
					Elapsed:    elapsed,
					Remaining:  remaining,
				})
			}

			// 等待设备响应准备好接收下一块数据（0x05 = ENQ）
			resp, err := c.serialManager.ReadExactly(1)
			if err != nil {
				return fmt.Errorf("failed to read upgrade response at byte %d: %w", totalSent, err)
			}

			if len(resp) != 1 || resp[0] != 0x05 {
				return fmt.Errorf("unexpected upgrade response: got 0x%02X, expected 0x05", resp[0])
			}
		}
	}

	// 完成回调
	if progressCallback != nil {
		elapsed := time.Since(startTime)
		var speed int64
		if elapsed.Seconds() > 0 {
			speed = int64(float64(fileSize) / elapsed.Seconds())
		}

		progressCallback(&models.UpgradeProgress{
			Current:    fileSize,
			Total:      fileSize,
			Percentage: 100.0,
			Speed:      speed,
			Elapsed:    elapsed,
			Remaining:  0,
		})
	}

	// 重连串口
	return c.serialManager.SetBaudRate(c.BaudRate)
}

// Open 开启串口连接
func (c *TjcDisplayClient) Open() error {
	return c.connect()
}

// Close 关闭串口连接
func (c *TjcDisplayClient) Close() error {
	if c.serialManager != nil && c.serialManager.IsOpen() {
		return c.serialManager.Close()
	}
	return nil
}

func (c *TjcDisplayClient) sendCommand(cmd string, appendReturnEndBytes bool) error {
	c.optLock.Lock()
	defer c.optLock.Unlock()

	cmdBytes := append([]byte(cmd), EndSymbol...)
	if appendReturnEndBytes {
		cmdBytes = append(append([]byte("printh "), consts.CodeStringData), EndSymbol...)
	}

	err := c.serialManager.Write(cmdBytes)
	if err != nil {
		return err
	}

	err = c.serialManager.Flush()
	if err != nil {
		return err
	}

	// 读取响应
	resData, err := c.serialManager.ReadUntil(EndSymbol)
	if err != nil {
		return err
	}

	// 清理多余数据
	_, _ = c.serialManager.ReadWithTimeout(50 * time.Millisecond)

	// 解析响应
	resp, err := parseResponse(resData)
	if err != nil {
		return fmt.Errorf("parse response failed: %w", err)
	}

	// 检查是否有错误
	if respErr := resp.toError(); respErr != nil {
		return respErr
	}

	return nil
}

func (c *TjcDisplayClient) sendCommandAndWaitResult(cmd string, startSymbol bool) ([]byte, error) {
	c.optLock.Lock()
	defer c.optLock.Unlock()

	cmdBytes := append([]byte(cmd), EndSymbol...)
	if startSymbol {
		cmdBytes = append(append([]byte("printh "), consts.CodeStringData), EndSymbol...)
	}

	err := c.serialManager.Write(cmdBytes)
	if err != nil {
		return nil, err
	}

	err = c.serialManager.Flush()
	if err != nil {
		return nil, err
	}

	// 读取响应
	resData, err := c.serialManager.ReadUntil(EndSymbol)
	if err != nil {
		return nil, err
	}

	// 清理多余数据
	_, _ = c.serialManager.ReadWithTimeout(50 * time.Millisecond)

	// 解析响应
	resp, err := parseResponse(resData)
	if err != nil {
		return nil, fmt.Errorf("parse response failed: %w", err)
	}

	// 检查是否有错误
	if respErr := resp.toError(); respErr != nil {
		return nil, respErr
	}

	return resp.Data, nil
}

func (c *TjcDisplayClient) sendCommandAndWaitRawResult(cmd string, startSymbol bool) ([]byte, error) {
	c.optLock.Lock()
	defer c.optLock.Unlock()

	cmdBytes := append([]byte(cmd), EndSymbol...)
	if startSymbol {
		cmdBytes = append(append([]byte("printh "), consts.CodeStringData), EndSymbol...)
	}

	err := c.serialManager.Write(cmdBytes)
	if err != nil {
		return nil, err
	}

	err = c.serialManager.Flush()
	if err != nil {
		return nil, err
	}

	// 读取响应
	resData, err := c.serialManager.Read()
	if err != nil {
		return nil, err
	}

	return resData, nil
}

// parseResponse 解析串口屏返回数据
func parseResponse(data []byte) (*Response, error) {
	if len(data) < 1 {
		return nil, fmt.Errorf("Invalid response length: %d", len(data))
	}

	resp := &Response{
		RawData: data,
		Code:    data[0],
	}

	// 根据第一个字节判断响应类型
	switch data[0] {
	case consts.CodeSuccess:
		resp.Type = ResponseTypeSuccess
	case consts.CodeInvalidInstruction,
		consts.CodeInvalidComponentID,
		consts.CodeInvalidPageID,
		consts.CodeInvalidPictureID,
		consts.CodeInvalidFontID,
		consts.CodeFileOperationFailed,
		consts.CodeCRCCheckFailed,
		consts.CodeInvalidBaudrate,
		consts.CodeInvalidCurveID,
		consts.CodeInvalidVariableName,
		consts.CodeInvalidVariableOp,
		consts.CodeAssignmentFailed,
		consts.CodeEEPROMFailed,
		consts.CodeInvalidParamCount,
		consts.CodeIOFailed,
		consts.CodeEscapeCharError,
		consts.CodeVariableNameTooLong,
		consts.CodeSerialBufferOverflow:
		resp.Type = ResponseTypeError
	case consts.CodeTouchEvent,
		consts.CodePageID,
		consts.CodeTouchCoordinate,
		consts.CodeSleepTouch,
		consts.CodeAutoSleep,
		consts.CodeAutoWake,
		consts.CodeStartupSuccess,
		consts.CodeStartSDUpgrade,
		consts.CodeTransparentReady,
		consts.CodeTransparentDone:
		resp.Type = ResponseTypeEvent
		// 提取事件数据（去掉第一个字节）
		if len(data) > 1 {
			resp.Data = data[1:]
		}
	case consts.CodeStringData,
		consts.CodeNumberData:
		resp.Type = ResponseTypeData
		// 提取数据（去掉第一个字节和结束符）
		if len(data) > 1 {
			resp.Data = data[1:]
		}
	default:
		// 未知响应码，可能是数据返回
		resp.Type = ResponseTypeData
		resp.Data = data
	}

	return resp, nil
}
