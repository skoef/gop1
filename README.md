[![Go Report Card](https://goreportcard.com/badge/github.com/skoef/gop1)](https://goreportcard.com/report/github.com/skoef/gop1) [![Documentation](https://godoc.org/github.com/skoef/gop1?status.svg)](http://godoc.org/github.com/skoef/gop1)

# Golang P1 protocol library

This is a golang library to read P1 data from a so called *smart* energy meter, used primarily in The Netherlands. P1 is the protocol Dutch power grid companies designed together and is described on [netbeheernederland.nl](https://www.netbeheernederland.nl/_upload/Files/Slimme_meter_15_a727fce1f1.pdf). The smart meters which are being deployed in Belgium implement the same protocol, but some additional data types were defined by the power grid companies. These types are defined in the [e-MUCS H](https://www.fluvius.be/sites/fluvius/files/2019-12/e-mucs_h_ed_1_3.pdf) specification.

To read P1 data, you'll need something like a P1-to-USB cable. The P1 port is essentially a serial port where data (a so called P1 telegram) is dumped every second.

## Example usage:
```golang
package main

import (
	"fmt"

	"github.com/skoef/gop1"
)

func main() {
	// open a new reader to given USB serial device
	p1, err := gop1.New(gop1.P1Config{
		USBDevice: "/dev/ttyUSB0",
	})
	if err != nil {
		panic(err)
	}

	// start reading data
	// this will send new telegrams to the channel p1.Incoming
	p1.Start()

	for telegram := range p1.Incoming {
		// loop over the objects in the telegram to find types we're interested in
		for _, obj := range telegram.Objects {
			switch obj.Type {
			case gop1.OBISTypeInstantaneousPowerDeliveredL1:
				fmt.Printf("actual power usage: %s %s\n", obj.Values[0].Value, obj.Values[0].Unit)
			}
		}
	}
}
```

In the [example/](https://github.com/skoef/gop1/tree/master/example) folder is an example application that collects relevant metrics and offers them over a prometheus-compatible HTTP endpoint for scraping.

## Acknowledgements
The [smartmeter](https://github.com/marceldegraaf/smartmeter) project from Marcel de Graaf inspired me to write something like this. I like his work, but was looking for a more pluggable library rather than an actual application. Also there is a whole lot python projects out there with P1 support that gave some insight.
