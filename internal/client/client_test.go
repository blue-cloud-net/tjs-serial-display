package client

import (
	"testing"

	"github.com/blue-cloud-net/tjc-serial-display/internal/serial"
	"github.com/blue-cloud-net/tjc-serial-display/pkg/consts"
)

// TestParseResponse_Success 测试解析成功响应
func TestParseResponse_Success(t *testing.T) {
	// 成功响应: 0x01 0xFF 0xFF 0xFF
	data := []byte{0x01, 0xFF, 0xFF, 0xFF}
	resp, err := parseResponse(data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.Type != ResponseTypeSuccess {
		t.Errorf("Expected ResponseTypeSuccess, got %v", resp.Type)
	}
	if resp.Code != consts.CodeSuccess {
		t.Errorf("Expected code 0x01, got 0x%02X", resp.Code)
	}
}

// TestParseResponse_Error 测试解析错误响应
func TestParseResponse_Error(t *testing.T) {
	testCases := []struct {
		name string
		data []byte
		code byte
	}{
		{"InvalidInstruction", []byte{0x00, 0xFF, 0xFF, 0xFF}, consts.CodeInvalidInstruction},
		{"InvalidComponentID", []byte{0x02, 0xFF, 0xFF, 0xFF}, consts.CodeInvalidComponentID},
		{"InvalidPageID", []byte{0x03, 0xFF, 0xFF, 0xFF}, consts.CodeInvalidPageID},
		{"InvalidBaudrate", []byte{0x11, 0xFF, 0xFF, 0xFF}, consts.CodeInvalidBaudrate},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := parseResponse(tc.data)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if resp.Type != ResponseTypeError {
				t.Errorf("Expected ResponseTypeError, got %v", resp.Type)
			}
			if resp.Code != tc.code {
				t.Errorf("Expected code 0x%02X, got 0x%02X", tc.code, resp.Code)
			}
		})
	}
}

// TestParseResponse_Event 测试解析事件响应
func TestParseResponse_Event(t *testing.T) {
	testCases := []struct {
		name     string
		data     []byte
		code     byte
		dataLen  int
		expected []byte
	}{
		{"TouchEvent", []byte{0x65, 0x01, 0x02, 0xFF, 0xFF, 0xFF}, consts.CodeTouchEvent, 2, []byte{0x01, 0x02}},
		{"PageID", []byte{0x66, 0x05, 0xFF, 0xFF, 0xFF}, consts.CodePageID, 1, []byte{0x05}},
		{"StartupSuccess", []byte{0x88, 0xFF, 0xFF, 0xFF}, consts.CodeStartupSuccess, 0, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := parseResponse(tc.data)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if resp.Type != ResponseTypeEvent {
				t.Errorf("Expected ResponseTypeEvent, got %v", resp.Type)
			}
			if resp.Code != tc.code {
				t.Errorf("Expected code 0x%02X, got 0x%02X", tc.code, resp.Code)
			}
			if len(resp.Data) != tc.dataLen {
				t.Errorf("Expected data length %d, got %d", tc.dataLen, len(resp.Data))
			}
		})
	}
}

// TestParseResponse_Data 测试解析数据响应
func TestParseResponse_Data(t *testing.T) {
	testCases := []struct {
		name     string
		data     []byte
		code     byte
		expected []byte
	}{
		{"StringData", []byte{0x70, 'H', 'e', 'l', 'l', 'o', 0xFF, 0xFF, 0xFF}, consts.CodeStringData, []byte{'H', 'e', 'l', 'l', 'o'}},
		{"NumberData", []byte{0x71, 0x01, 0x02, 0x03, 0x04, 0xFF, 0xFF, 0xFF}, consts.CodeNumberData, []byte{0x01, 0x02, 0x03, 0x04}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := parseResponse(tc.data)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if resp.Type != ResponseTypeData {
				t.Errorf("Expected ResponseTypeData, got %v", resp.Type)
			}
			if resp.Code != tc.code {
				t.Errorf("Expected code 0x%02X, got 0x%02X", tc.code, resp.Code)
			}
			if len(resp.Data) != len(tc.expected) {
				t.Errorf("Expected data length %d, got %d", len(tc.expected), len(resp.Data))
			}
		})
	}
}

// TestParseResponse_InvalidLength 测试无效长度的响应
func TestParseResponse_InvalidLength(t *testing.T) {
	testCases := [][]byte{
		{},
		{0x01},
		{0x01, 0xFF},
		{0x01, 0xFF, 0xFF},
	}

	for i, data := range testCases {
		t.Run(string(rune(i)), func(t *testing.T) {
			_, err := parseResponse(data)
			if err == nil {
				t.Error("Expected error for invalid length, got nil")
			}
		})
	}
}

// TestParseResponse_InvalidEndBytes 测试无效的结束符
func TestParseResponse_InvalidEndBytes(t *testing.T) {
	testCases := [][]byte{
		{0x01, 0xFF, 0xFF, 0xFE},
		{0x01, 0xFF, 0xFE, 0xFF},
		{0x01, 0xFE, 0xFF, 0xFF},
		{0x01, 0x00, 0x00, 0x00},
	}

	for i, data := range testCases {
		t.Run(string(rune(i)), func(t *testing.T) {
			_, err := parseResponse(data)
			if err == nil {
				t.Error("Expected error for invalid end bytes, got nil")
			}
		})
	}
}

// TestResponse_ToError 测试响应转换为错误
func TestResponse_ToError(t *testing.T) {
	// 测试错误响应
	errorResp := &Response{
		Type: ResponseTypeError,
		Code: consts.CodeInvalidInstruction,
	}
	err := errorResp.toError()
	if err == nil {
		t.Error("Expected error, got nil")
	}
	tjcErr, ok := err.(*TjcError)
	if !ok {
		t.Error("Expected TjcError type")
	}
	if tjcErr.Code != consts.CodeInvalidInstruction {
		t.Errorf("Expected code 0x%02X, got 0x%02X", consts.CodeInvalidInstruction, tjcErr.Code)
	}

	// 测试成功响应
	successResp := &Response{
		Type: ResponseTypeSuccess,
		Code: consts.CodeSuccess,
	}
	err = successResp.toError()
	if err != nil {
		t.Errorf("Expected no error for success response, got %v", err)
	}
}

// TestParseDeviceInfo_Valid 测试解析有效的设备信息
func TestParseDeviceInfo_Valid(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		expect struct {
			Type                  int
			Address               string
			Model                 string
			FirmwareVersion       int
			MainControlChipNumber int
			Number                string
			FlashSize             int
		}
	}{
		{
			name:  "WithComokPrefix",
			input: "comok 1,101-0,TJC4024T032_011R,52,61488,D264B8204F0E1828,16777216",
			expect: struct {
				Type                  int
				Address               string
				Model                 string
				FirmwareVersion       int
				MainControlChipNumber int
				Number                string
				FlashSize             int
			}{
				Type:                  1,
				Address:               "101-0",
				Model:                 "TJC4024T032_011R",
				FirmwareVersion:       52,
				MainControlChipNumber: 61488,
				Number:                "D264B8204F0E1828",
				FlashSize:             16777216,
			},
		},
		{
			name:  "WithoutSpace",
			input: "comok1,101-0,TJC4024T032_011R,52,61488,D264B8204F0E1828,16777216",
			expect: struct {
				Type                  int
				Address               string
				Model                 string
				FirmwareVersion       int
				MainControlChipNumber int
				Number                string
				FlashSize             int
			}{
				Type:                  1,
				Address:               "101-0",
				Model:                 "TJC4024T032_011R",
				FirmwareVersion:       52,
				MainControlChipNumber: 61488,
				Number:                "D264B8204F0E1828",
				FlashSize:             16777216,
			},
		},
		{
			name:  "NoPrefix",
			input: "0,0-0,TestModel,10,12345,ABCDEF123456,8388608",
			expect: struct {
				Type                  int
				Address               string
				Model                 string
				FirmwareVersion       int
				MainControlChipNumber int
				Number                string
				FlashSize             int
			}{
				Type:                  0,
				Address:               "0-0",
				Model:                 "TestModel",
				FirmwareVersion:       10,
				MainControlChipNumber: 12345,
				Number:                "ABCDEF123456",
				FlashSize:             8388608,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			info, err := parseDeviceInfo(tc.input)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if info.Type != tc.expect.Type {
				t.Errorf("Expected Type %d, got %d", tc.expect.Type, info.Type)
			}
			if info.Address != tc.expect.Address {
				t.Errorf("Expected Address %s, got %s", tc.expect.Address, info.Address)
			}
			if info.Model != tc.expect.Model {
				t.Errorf("Expected Model %s, got %s", tc.expect.Model, info.Model)
			}
			if info.FirmwareVersion != tc.expect.FirmwareVersion {
				t.Errorf("Expected FirmwareVersion %d, got %d", tc.expect.FirmwareVersion, info.FirmwareVersion)
			}
			if info.MainControlChipNumber != tc.expect.MainControlChipNumber {
				t.Errorf("Expected MainControlChipNumber %d, got %d", tc.expect.MainControlChipNumber, info.MainControlChipNumber)
			}
			if info.Number != tc.expect.Number {
				t.Errorf("Expected Number %s, got %s", tc.expect.Number, info.Number)
			}
			if info.FlashSize != tc.expect.FlashSize {
				t.Errorf("Expected FlashSize %d, got %d", tc.expect.FlashSize, info.FlashSize)
			}
		})
	}
}

// TestParseDeviceInfo_Invalid 测试解析无效的设备信息
func TestParseDeviceInfo_Invalid(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"Empty", ""},
		{"TooFewFields", "comok 1,101-0,TJC4024T032_011R"},
		{"OnlyPrefix", "comok "},
		{"LessThan7Fields", "1,2,3,4,5,6"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parseDeviceInfo(tc.input)
			if err == nil {
				t.Error("Expected error for invalid format, got nil")
			}
		})
	}
}

// TestTjcError_Error 测试错误消息格式化
func TestTjcError_Error(t *testing.T) {
	err := &TjcError{
		Code:    0x02,
		Message: "控件ID无效",
	}

	expected := "TJC Error 0x02: 控件ID无效"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}

// TestTjcDisplayClient_Connect 测试连接功能
func TestTjcDisplayClient_Connect(t *testing.T) {
	// 测试使用无效端口（应该失败）
	client := &TjcDisplayClient{
		PortName: "/dev/nonexistent_port_99999",
		BaudRate: 921600,
	}

	err := client.connect()
	if err == nil {
		t.Error("Expected error when connecting to non-existent port, got nil")
	}
}

// TestTjcDisplayClient_Close 测试关闭未连接的客户端
func TestTjcDisplayClient_Close(t *testing.T) {
	client := &TjcDisplayClient{}

	err := client.Close()
	if err != nil {
		t.Errorf("Expected no error when closing unconnected client, got %v", err)
	}
}

// TestTjcDisplayClient_GetDeviceInfo_RealDevice 测试获取真实设备信息
// 此测试需要连接到真实的串口设备
func TestTjcDisplayClient_GetDeviceInfo_RealDevice(t *testing.T) {
	// 获取可用端口
	ports, err := serial.ListPorts()
	if err != nil || len(ports) == 0 {
		t.Skip("Skipping real device test: no serial ports available")
	}

	client := &TjcDisplayClient{
		PortName: ports[0],
		BaudRate: 921600,
	}
	defer client.Close()

	info, err := client.GetDeviceInfo()
	if err != nil {
		t.Skipf("Cannot get device info from %s: %v (device may not be a TJC display)", ports[0], err)
	}

	t.Logf("Device Info: %+v", info)

	// 验证基本信息
	if info.Model == "" {
		t.Error("Expected non-empty Model")
	}
}

// TestTjcDisplayClient_BaseWorkflow 测试基础工作流程
// 此测试需要连接到真实的TJC串口屏
func TestTjcDisplayClient_BaseWorkflow(t *testing.T) {
	// 获取可用端口
	ports, err := serial.ListPorts()
	if err != nil || len(ports) == 0 {
		t.Skip("Skipping full workflow test: no serial ports available")
	}

	t.Logf("Testing with port: %s", ports[0])

	client := &TjcDisplayClient{
		PortName: ports[0],
		BaudRate: 921600,
	}
	defer client.Close()

	// 测试获取设备信息
	info, err := client.GetDeviceInfo()
	if err != nil {
		t.Skipf("Cannot connect to TJC display on %s: %v", ports[0], err)
	}
	t.Logf("Device Info: %+v", info)

	// 测试获取当前页面
	page, err := client.GetPage()
	if err != nil {
		t.Errorf("GetPage test: %v (may be expected if device doesn't support this command)", err)
	} else {
		t.Logf("Current page: %d", page)
	}
}
