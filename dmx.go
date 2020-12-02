package dmx

import (
	"io"
	"log"
	"time"

	"github.com/distributed/sers"
)

const (
	defaultDevice = "/dev/ttyUSB"
	frameSize     = 512
)

// DMX is the model for serial DMX connection
type DMX struct {
	dev    string
	frame  [frameSize]byte
	serial io.ReadWriteCloser
	port   sers.SerialPort
}

// NewDMXConnection creates a new DMX connection using an specific serial device
func NewDMXConnection(device string) (dmx *DMX, err error) {
	dmx = &DMX{}
	dmx.dev = device

	if len(dmx.dev) == 0 {
		dmx.dev = defaultDevice
	}

	c, err := sers.Open(dmx.dev)
	if err != nil {
		return nil, err
	}

	if c.SetMode(250000, 8, 0, 2, 0); err != nil {
		return nil, err
	}

	dmx.serial = c
	dmx.port = c
	dmx.frame[0] = 0
	log.Printf("Success opening port %s.", dmx.dev)
	return dmx, nil
}

// SetAddress -
func (dmx *DMX) SetAddress(address int, value byte) error {
	checkAdress(address)
	dmx.frame[address] = value
	return nil
}

// Render sends frame to serial device
func (dmx *DMX) Render() error {
	f := dmx.frame

	dmx.port.SetBreak(true)
	time.Sleep(500 * time.Millisecond)
	dmx.port.SetBreak(false)
	time.Sleep(500 * time.Millisecond)

	if _, err := dmx.serial.Write(f[:]); err != nil {
		return err
	}

	return nil
}

// Close serial port
func (dmx *DMX) Close() error {
	return dmx.serial.Close()
}

func checkAdress(id int) {
	if (id > 512) || (id < 1) {
		log.Fatalln("Channel format error: ", id)
	}
}
