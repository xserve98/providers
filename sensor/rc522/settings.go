package main

import (
	"encoding/hex"
	"fmt"

	"github.com/pkg/errors"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

// Settings has data required to start RFID reader.
type Settings struct {
	BusID       int               `yaml:"bus" default:"0"`
	DeviceID    int               `yaml:"device" default:"0"`
	ResetPin    int               `yaml:"reset" default:"13"`
	IRQPin      int               `yaml:"irq" default:"11"`
	AntennaGain int               `yaml:"antennaGain" validate:"gte=0,lte=7" default:"4"`
	Sector      int               `yaml:"sector" default:"1"`
	Block       int               `yaml:"block" default:"0"`
	Users       map[string]string `yaml:"users" validate:"required,min=1"`
	Key         string            `yaml:"key" validate:"required,len=12"`

	encKey [6]byte
	users  map[string][]byte

	rstGPIO gpio.PinIO
	irqGPIO gpio.PinIO
}

// Validate checks settings.
func (s *Settings) Validate() error {
	s.users = make(map[string][]byte)

	for k, v := range s.Users {
		d, err := hex.DecodeString(v)
		if err != nil || 16 != len(d) {
			continue
		}

		s.users[k] = d
	}

	if 0 == len(s.Users) {
		return errors.New("users are not defined")
	}

	d, err := hex.DecodeString(s.Key)
	if err != nil || 6 != len(d) {
		return errors.New("key must be 6 bytes long")
	}

	copy(s.encKey[:], d)

	s.rstGPIO = gpioreg.ByName(fmt.Sprintf("P1_%d", s.ResetPin))
	if nil == s.rstGPIO {
		return errors.New("incorrect reset pin")
	}

	s.irqGPIO = gpioreg.ByName(fmt.Sprintf("P1_%d", s.IRQPin))
	if nil == s.irqGPIO {
		return errors.New("incorrect irq pin")
	}

	return nil
}
