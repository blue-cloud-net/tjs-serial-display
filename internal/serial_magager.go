package serialmanager

import (
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
	mode := &serial.Mode{
		BaudRate: spm.BaudRate,
	}

	port, err := serial.Open(spm.PortName, mode)
	if err != nil {
		return err
	}

	spm.port = port
	return nil
}

func (spm *SerialPortManager) IsOpen() (bool, error) {
	if spm.port != nil {
		return true, nil
	}

	return false, nil
}

func (spm *SerialPortManager) Close() error {
	if spm.port != nil {
		return spm.port.Close()
	}

	return nil
}
