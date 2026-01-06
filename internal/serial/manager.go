package serial

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"time"

	"go.bug.st/serial"
)

type SerialPortManager struct {
	PortName    string
	BaudRate    int
	Parity      serial.Parity
	DataBits    int
	StopBits    serial.StopBits
	Timeout     time.Duration
	BytesToRead int
	port        serial.Port
}

func ListPorts() ([]string, error) {
	return serial.GetPortsList()
}

func (spm *SerialPortManager) Open() error {
	// 设置默认值
	if spm.BaudRate == 0 {
		spm.BaudRate = 115200
	}
	if spm.Parity == 0 {
		spm.Parity = serial.NoParity
	}
	if spm.DataBits == 0 {
		spm.DataBits = 8
	}
	if spm.StopBits == 0 {
		spm.StopBits = serial.OneStopBit
	}
	if spm.Timeout == 0 {
		spm.Timeout = 2 * time.Second
	}
	// 根据波特率动态设置默认BytesToRead（约1帧数据，常用为波特率/10，最小64，最大4096）
	if spm.BytesToRead == 0 {
		bytesToRead := spm.BaudRate / 10
		if bytesToRead < 64 {
			bytesToRead = 64
		} else if bytesToRead > 4096 {
			bytesToRead = 4096
		}
		spm.BytesToRead = bytesToRead
	}

	mode := &serial.Mode{
		BaudRate: spm.BaudRate,
		Parity:   spm.Parity,
		DataBits: spm.DataBits,
		StopBits: spm.StopBits,
	}

	port, err := serial.Open(spm.PortName, mode)
	if err != nil {
		return err
	}

	port.SetReadTimeout(spm.Timeout)

	spm.port = port
	return nil
}

func (spm *SerialPortManager) IsOpen() bool {
	if spm.port != nil {
		return true
	}

	return false
}

func (spm *SerialPortManager) Close() error {
	if spm.port != nil {
		return spm.port.Close()
	}

	return nil
}

func (spm *SerialPortManager) SetBaudRate(baudRate int) error {
	spm.BaudRate = baudRate

	if spm.IsOpen() {
		err := spm.port.SetMode(&serial.Mode{
			BaudRate: baudRate,
			Parity:   spm.Parity,
			DataBits: spm.DataBits,
			StopBits: spm.StopBits,
		})
		if err != nil {
			return errors.New("failed to set new baud rate on serial port")
		}
	}

	return nil
}

func (spm *SerialPortManager) SetParity(parity serial.Parity) error {
	spm.Parity = parity

	if spm.IsOpen() {
		err := spm.port.SetMode(&serial.Mode{
			BaudRate: spm.BaudRate,
			Parity:   parity,
			DataBits: spm.DataBits,
			StopBits: spm.StopBits,
		})
		if err != nil {
			return errors.New("failed to set new parity on serial port")
		}
	}

	return nil
}

func (spm *SerialPortManager) SetDataBits(dataBits int) error {
	spm.DataBits = dataBits

	if spm.IsOpen() {
		err := spm.port.SetMode(&serial.Mode{
			BaudRate: spm.BaudRate,
			Parity:   spm.Parity,
			DataBits: dataBits,
			StopBits: spm.StopBits,
		})
		if err != nil {
			return errors.New("failed to set new data bits on serial port")
		}
	}

	return nil
}

func (spm *SerialPortManager) SetStopBits(stopBits serial.StopBits) error {
	spm.StopBits = stopBits

	if spm.IsOpen() {
		err := spm.port.SetMode(&serial.Mode{
			BaudRate: spm.BaudRate,
			Parity:   spm.Parity,
			DataBits: spm.DataBits,
			StopBits: stopBits,
		})
		if err != nil {
			return errors.New("failed to set new stop bits on serial port")
		}
	}

	return nil
}

func (spm *SerialPortManager) SetReadTimeout(timeout time.Duration) error {
	spm.Timeout = timeout

	return spm.port.SetReadTimeout(timeout)
}

// Read 读取一次可用数据（最多 bytesToRead 字节）
func (spm *SerialPortManager) Read() ([]byte, error) {
	if spm.port == nil || !spm.IsOpen() {
		return nil, errors.New("port is not open")
	}

	buf := make([]byte, spm.BytesToRead)
	n, err := spm.port.Read(buf)
	if err != nil {
		return nil, err
	}

	return buf[:n], nil
}

// ReadExactly 读取精确的字节数，会循环读取直到读满或出错
func (spm *SerialPortManager) ReadExactly(n int) ([]byte, error) {
	if n <= 0 {
		return nil, errors.New("bytes to read must be greater than 0")
	}

	if spm.port == nil || !spm.IsOpen() {
		return nil, errors.New("port is not open")
	}

	buf := make([]byte, n)
	totalRead := 0

	for totalRead < n {
		bytesRead, err := spm.port.Read(buf[totalRead:])
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		if bytesRead == 0 {
			// 超时或没有数据
			return buf[:totalRead], errors.New("read timeout or no data")
		}

		totalRead += bytesRead
	}

	if totalRead < n {
		return buf[:totalRead], fmt.Errorf("incomplete read: expected %d bytes, got %d", n, totalRead)
	}

	return buf, nil
}

// ReadUntil 读取数据直到遇到指定的分隔符（如 []byte{0xFF, 0xFF, 0xFF}）
func (spm *SerialPortManager) ReadUntil(delimiter []byte) ([]byte, error) {
	if len(delimiter) == 0 {
		return nil, errors.New("delimiter cannot be empty")
	}

	if spm.port == nil || !spm.IsOpen() {
		return nil, errors.New("port is not open")
	}

	var result bytes.Buffer
	buf := make([]byte, 1)

	for {
		n, err := spm.port.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		if n == 0 {
			// 超时或没有数据
			return nil, errors.New("read timeout or no data.")
		}

		if n > 0 {
			result.WriteByte(buf[0])

			// 检查是否以分隔符结尾
			if result.Len() >= len(delimiter) {
				data := result.Bytes()
				if bytes.HasSuffix(data, delimiter) {
					// 移除分隔符后返回
					return data[:len(data)-len(delimiter)], nil
				}
			}
		}
	}

	return result.Bytes(), nil
}

// ReadWithTimeout 在指定时间内读取数据
func (spm *SerialPortManager) ReadWithTimeout(timeout time.Duration) ([]byte, error) {
	if spm.port == nil {
		return nil, errors.New("port is not open")
	}

	// 保存原始超时设置
	originalTimeout := 3 * time.Second
	defer spm.port.SetReadTimeout(originalTimeout)

	// 设置新的超时
	spm.port.SetReadTimeout(timeout)
	var result bytes.Buffer
	buf := make([]byte, 4096)

	for {
		n, err := spm.port.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}
		if n > 0 {
			result.Write(buf[:n])
		} else {
			// 没有更多数据，可能是超时
			break
		}
	}

	return result.Bytes(), nil
}

// ReadAll 持续读取直到没有更多数据（会等待超时）
func (spm *SerialPortManager) ReadAll(maxSize int) ([]byte, error) {
	if spm.port == nil {
		return nil, errors.New("port is not open")
	}

	if maxSize <= 0 {
		maxSize = 65536 // 默认最大 64KB
	}

	var result bytes.Buffer
	buf := make([]byte, spm.BytesToRead)

	for result.Len() < maxSize {
		n, err := spm.port.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		if n > 0 {
			result.Write(buf[:n])
		} else {
			// 没有更多数据
			break
		}
	}

	return result.Bytes(), nil
}

func (spm *SerialPortManager) Write(p []byte) error {
	if len(p) == 0 {
		return errors.New("invalid data length for write")
	}

	if spm.port == nil {
		return errors.New("port is not open")
	}

	// 循环写入，确保所有数据都写入成功（处理部分写入情况）
	totalWritten := 0
	startTime := time.Now()
	writeTimeout := 5 * time.Second

	for totalWritten < len(p) {
		// 检查是否超时
		if time.Since(startTime) > writeTimeout {
			return errors.New("write timeout: unable to write all data")
		}

		n, err := spm.port.Write(p[totalWritten:])
		if err != nil {
			return err
		}

		if n == 0 {
			// 无法写入数据，短暂延迟后重试
			time.Sleep(10 * time.Millisecond)
			continue
		}

		totalWritten += n

		if totalWritten < len(p) {
			err = spm.port.Drain()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (spm *SerialPortManager) Flush() error {
	if spm.port != nil {
		return spm.port.Drain()
	}

	return errors.New("port is not open")
}
