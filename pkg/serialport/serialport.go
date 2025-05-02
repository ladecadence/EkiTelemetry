package serialport

import (
	"bytes"
	"errors"

	"go.bug.st/serial"
)

var (
	startPattern = []byte{0xaa, 0x55, 0xaa, 0x55}
	endPattern   = []byte{0x33, 0xcc, 0x33, 0xcc}
)

type Serial struct {
	Port    serial.Port
	Packets int
}

func GetSerialPorts() []string {
	ports, err := serial.GetPortsList()
	if err != nil {
		return nil
	}
	return ports
}

func NewSerial(port string, speed int) (*Serial, error) {
	s := Serial{Packets: 0}

	// prepare port
	mode := &serial.Mode{
		BaudRate: speed,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	// open port
	var err error
	s.Port, err = serial.Open(port, mode)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (s *Serial) Close() error {
	if s.Port != nil {
		err := s.Port.Close()
		return err
	} else {
		return errors.New("Serial port not opened")
	}
}

func (s *Serial) ListenAndDecode(telemChan chan string, ssdvChan chan []byte) error {
	var err error
	go func(err error) {
		serialBuf := make([]byte, 256)
		var dataBuf []byte

		for {
			// read data
			num, err := s.Port.Read(serialBuf)
			if err != nil {
				break
			}
			// if data received has a complete packet
			if bytes.Contains(serialBuf, startPattern) && bytes.Contains(serialBuf, endPattern) && (bytes.Index(dataBuf, startPattern) < bytes.Index(dataBuf, endPattern)) {
				// ok, cut the data
				_, dataBuf, _ = bytes.Cut(serialBuf, startPattern) // keep after startPattern
				dataBuf, _, _ = bytes.Cut(dataBuf, endPattern)     // keep before end pattern
				// send it
				if dataBuf[0] == '$' && dataBuf[1] == '$' {
					telemChan <- string(dataBuf)
					s.Packets++
				} else if dataBuf[0] == 0x66 && len(dataBuf) == 255 {
					ssdvChan <- dataBuf
				}
				// clear
				dataBuf = []byte{}
			} else {
				// if not keep adding
				dataBuf = append(dataBuf, serialBuf[0:num]...)
				// until we have a complete packet
				if bytes.Contains(dataBuf, startPattern) && bytes.Contains(dataBuf, endPattern) && (bytes.Index(dataBuf, startPattern) < bytes.Index(dataBuf, endPattern)) {
					// ok, cut the data
					_, dataBuf, _ = bytes.Cut(dataBuf, startPattern) // keep after startPattern
					dataBuf, _, _ = bytes.Cut(dataBuf, endPattern)   // keep before end pattern
					// send it
					if dataBuf[0] == '$' && dataBuf[1] == '$' {
						s.Packets++
						telemChan <- string(dataBuf)
					} else if dataBuf[0] == 0x66 && len(dataBuf) == 255 {
						ssdvChan <- dataBuf
					}
					// clear
					dataBuf = []byte{}
				}
			}

		}
	}(err)

	return err
}
