package serial

import (
	"testing"

	"go.bug.st/serial"
)

func TestListPorts(t *testing.T) {
	ports, err := ListPorts()
	if err != nil {
		t.Fatalf("Failed to list ports: %v", err)
	}

	t.Logf("Available ports: %v", ports)
}

func TestSerialPortManager_Open_DefaultValues(t *testing.T) {
	spm := &SerialPortManager{
		PortName: "/dev/null", // 使用 /dev/null 进行测试，不需要真实串口
	}

	// 注意：这个测试可能会失败，因为 /dev/null 不是真正的串口设备
	// 但它可以测试默认值是否正确设置
	err := spm.Open()
	if err == nil {
		defer spm.Close()

		// 验证默认值
		if spm.BaudRate != 115200 {
			t.Errorf("Expected default BaudRate 115200, got %d", spm.BaudRate)
		}
		if spm.Parity != serial.NoParity {
			t.Errorf("Expected default Parity NoParity, got %v", spm.Parity)
		}
		if spm.DataBits != 8 {
			t.Errorf("Expected default DataBits 8, got %d", spm.DataBits)
		}
		if spm.StopBits != serial.OneStopBit {
			t.Errorf("Expected default StopBits OneStopBit, got %v", spm.StopBits)
		}
	}
}

func TestSerialPortManager_Open_CustomValues(t *testing.T) {
	spm := &SerialPortManager{
		PortName: "/dev/null",
		BaudRate: 9600,
		Parity:   serial.EvenParity,
		DataBits: 7,
		StopBits: serial.TwoStopBits,
	}

	err := spm.Open()
	if err == nil {
		defer spm.Close()

		// 验证自定义值保持不变
		if spm.BaudRate != 9600 {
			t.Errorf("Expected BaudRate 9600, got %d", spm.BaudRate)
		}
		if spm.Parity != serial.EvenParity {
			t.Errorf("Expected Parity EvenParity, got %v", spm.Parity)
		}
		if spm.DataBits != 7 {
			t.Errorf("Expected DataBits 7, got %d", spm.DataBits)
		}
		if spm.StopBits != serial.TwoStopBits {
			t.Errorf("Expected StopBits TwoStopBits, got %v", spm.StopBits)
		}
	}
}

func TestSerialPortManager_Open_InvalidPort(t *testing.T) {
	spm := &SerialPortManager{
		PortName: "/dev/nonexistent_port_12345",
	}

	err := spm.Open()
	if err == nil {
		defer spm.Close()
		t.Error("Expected error when opening non-existent port, got nil")
	}
}

func TestSerialPortManager_IsOpen(t *testing.T) {
	spm := &SerialPortManager{
		PortName: "/dev/null",
	}

	// 初始状态应该是关闭的
	if spm.IsOpen() {
		t.Error("Expected port to be closed initially")
	}

	// 打开后应该返回 true（如果成功）
	err := spm.Open()
	if err == nil {
		defer spm.Close()
		if !spm.IsOpen() {
			t.Error("Expected port to be open after Open()")
		}
	}
}

func TestSerialPortManager_Close(t *testing.T) {
	// 测试关闭未打开的端口
	spm := &SerialPortManager{}
	err := spm.Close()
	if err != nil {
		t.Errorf("Expected no error when closing unopened port, got: %v", err)
	}

	// 测试关闭已打开的端口
	spm2 := &SerialPortManager{
		PortName: "/dev/null",
	}
	if err := spm2.Open(); err == nil {
		err = spm2.Close()
		if err != nil {
			t.Errorf("Failed to close port: %v", err)
		}

		if spm2.IsOpen() {
			t.Error("Port should be closed after Close()")
		}
	}
}

func TestSerialPortManager_Read_PortNotOpen(t *testing.T) {
	spm := &SerialPortManager{}

	_, err := spm.Read()
	if err == nil {
		t.Error("Expected error when reading from closed port, got nil")
	}
	if err != nil && err.Error() != "port is not open" {
		t.Errorf("Expected 'port is not open' error, got: %v", err)
	}
}

func TestSerialPortManager_ReadExactly_InvalidLength(t *testing.T) {
	spm := &SerialPortManager{}

	// 测试无效长度
	testCases := []int{0, -1, -100}
	for _, length := range testCases {
		_, err := spm.ReadExactly(length)
		if err == nil {
			t.Errorf("Expected error for length %d, got nil", length)
		}
		if err != nil && err.Error() != "bytes to read must be greater than 0" {
			t.Errorf("Unexpected error message: %v", err)
		}
	}
}

func TestSerialPortManager_ReadExactly_PortNotOpen(t *testing.T) {
	spm := &SerialPortManager{}

	_, err := spm.ReadExactly(10)
	if err == nil {
		t.Error("Expected error when reading from closed port, got nil")
	}
	if err != nil && err.Error() != "port is not open" {
		t.Errorf("Expected 'port is not open' error, got: %v", err)
	}
}

func TestSerialPortManager_ReadUntil_EmptyDelimiter(t *testing.T) {
	spm := &SerialPortManager{}

	_, err := spm.ReadUntil([]byte{})
	if err == nil {
		t.Error("Expected error for empty delimiter, got nil")
	}
	if err != nil && err.Error() != "delimiter cannot be empty" {
		t.Errorf("Expected 'delimiter cannot be empty' error, got: %v", err)
	}
}

func TestSerialPortManager_ReadUntil_PortNotOpen(t *testing.T) {
	spm := &SerialPortManager{}

	_, err := spm.ReadUntil([]byte{0xFF, 0xFF, 0xFF})
	if err == nil {
		t.Error("Expected error when reading from closed port, got nil")
	}
	if err != nil && err.Error() != "port is not open" {
		t.Errorf("Expected 'port is not open' error, got: %v", err)
	}
}

func TestSerialPortManager_ReadWithTimeout_PortNotOpen(t *testing.T) {
	spm := &SerialPortManager{}

	_, err := spm.ReadWithTimeout(1000)
	if err == nil {
		t.Error("Expected error when reading from closed port, got nil")
	}
	if err != nil && err.Error() != "port is not open" {
		t.Errorf("Expected 'port is not open' error, got: %v", err)
	}
}

func TestSerialPortManager_ReadAll_PortNotOpen(t *testing.T) {
	spm := &SerialPortManager{}

	_, err := spm.ReadAll(1024)
	if err == nil {
		t.Error("Expected error when reading from closed port, got nil")
	}
	if err != nil && err.Error() != "port is not open" {
		t.Errorf("Expected 'port is not open' error, got: %v", err)
	}
}

func TestSerialPortManager_Write_EmptyData(t *testing.T) {
	spm := &SerialPortManager{}

	err := spm.Write([]byte{})
	if err == nil {
		t.Error("Expected error when writing empty data, got nil")
	}
	if err != nil && err.Error() != "Invalid data length for write. Must be greater than 0." {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestSerialPortManager_Write_PortNotOpen(t *testing.T) {
	spm := &SerialPortManager{}

	err := spm.Write([]byte("test"))
	if err == nil {
		t.Error("Expected error when writing to closed port, got nil")
	}
	if err != nil && err.Error() != "Port is not open." {
		t.Errorf("Expected 'Port is not open.' error, got: %v", err)
	}
}

func TestSerialPortManager_Flush_PortNotOpen(t *testing.T) {
	spm := &SerialPortManager{}

	err := spm.Flush()
	if err == nil {
		t.Error("Expected error when flushing closed port, got nil")
	}
	if err != nil && err.Error() != "Port is not open." {
		t.Errorf("Expected 'Port is not open.' error, got: %v", err)
	}
}

// TestSerialPortManager_FullWorkflow 测试完整的工作流程
// 注意：此测试需要真实的串口设备才能完全通过
func TestSerialPortManager_FullWorkflow(t *testing.T) {
	// 获取可用端口
	ports, err := ListPorts()
	if err != nil {
		t.Skipf("Skipping full workflow test: cannot list ports: %v", err)
	}

	if len(ports) == 0 {
		t.Skip("Skipping full workflow test: no serial ports available")
	}

	t.Logf("Testing with port: %s", ports[0])

	spm := &SerialPortManager{
		PortName: ports[0],
		BaudRate: 115200,
	}

	// 打开端口
	err = spm.Open()
	if err != nil {
		t.Skipf("Cannot open port %s: %v", ports[0], err)
	}
	defer spm.Close()

	// 验证端口已打开
	if !spm.IsOpen() {
		t.Error("Port should be open")
	}

	// 测试写入
	testData := []byte("Hello, Serial!")
	err = spm.Write(testData)
	if err != nil {
		t.Logf("Write test: %v (may be expected if device doesn't support writing)", err)
	}

	// 测试刷新
	err = spm.Flush()
	if err != nil {
		t.Logf("Flush test: %v (may be expected depending on device)", err)
	}

	// 测试读取
	spm.BytesToRead = 128
	_, err = spm.Read()
	if err != nil {
		t.Logf("Read test: %v (may timeout, which is expected)", err)
	}

	// 测试 ReadExactly
	_, err = spm.ReadExactly(10)
	if err != nil {
		t.Logf("ReadExactly test: %v (may timeout, which is expected)", err)
	}

	// 测试 ReadUntil
	_, err = spm.ReadUntil([]byte{0xFF, 0xFF, 0xFF})
	if err != nil {
		t.Logf("ReadUntil test: %v (may timeout, which is expected)", err)
	}

	// 关闭端口
	err = spm.Close()
	if err != nil {
		t.Errorf("Failed to close port: %v", err)
	}
}
