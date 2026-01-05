package serial

import (
	"bytes"
	"errors"
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

	// 设置读取超时为3秒
	port.SetReadTimeout(3 * time.Second)

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
		return buf[:totalRead], errors.New("incomplete read: expected " + string(rune(n)) + " bytes, got " + string(rune(totalRead)))
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

	if spm.port != nil {
		n, err := spm.port.Write(p)
		if err != nil {
			return err
		}

		if n != len(p) {
			return errors.New("failed to write all data to serial port")
		}

		return nil
	}

	return errors.New("port is not open")
}

func (spm *SerialPortManager) Flush() error {
	if spm.port != nil {
		return spm.port.Drain()
	}

	return errors.New("port is not open")
}
