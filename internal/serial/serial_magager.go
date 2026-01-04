package serial

import (
	"errors"

	"go.bug.st/serial"
)

type SerialPortManager struct {
	PortName string
	BaudRate int
	Parity   serial.Parity
	DataBits int
	StopBits serial.StopBits
	port     serial.Port
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

	port.SetReadTimeout(10 * 1000) // 设置读取超时为10秒

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

func (spm *SerialPortManager) Read(len int) ([]byte, error) {
	if len <= 0 {
		return nil, errors.New("Invalid length for read. Must be greater than 0.")
	}

	if spm.port != nil {
		buf := make([]byte, len)
		n, err := spm.port.Read(buf)
		if err != nil {
			return nil, err
		}
		return buf[:n], nil
	}

	return nil, errors.New("Port is not open")
}

func (spm *SerialPortManager) Write(p []byte) error {
	if len(p) == 0 {
		return errors.New("Invalid data length for write. Must be greater than 0.")
	}

	if spm.port != nil {
		n, err := spm.port.Write(p)
		if err != nil {
			return err
		}

		if n != len(p) {
			return errors.New("Failed to write all data to serial port.")
		}

		return nil
	}

	return errors.New("Port is not open.")
}

func (spm *SerialPortManager) Flush() error {
	if spm.port != nil {
		return spm.port.Drain()
	}

	return errors.New("Port is not open.")
}
