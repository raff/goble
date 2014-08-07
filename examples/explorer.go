package main

import (
	"flag"
	"fmt"
	"log"

	"../../goble"
)

func main() {
	verbose := flag.Bool("verbose", false, "dump all events")
	dups := flag.Bool("allow-duplicates", false, "allow duplicates when scanning")
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Println("usage: explorer [options] peripheral-uuid")
		return
	}

	peripheralUuid := flag.Args()[0]

	var done chan bool

	ble := goble.New()
	ble.SetVerbose(*verbose)

	if *verbose {
		ble.On(goble.ALL, func(ev goble.Event) {
			log.Println("Event", ev)
		})
	}

	ble.On("stateChange", func(ev goble.Event) {
		if ev.State == "poweredOn" {
			ble.StartScanning(nil, *dups)
		} else {
			ble.StopScanning()
		}
	})

	ble.On("discover", func(ev goble.Event) {
		if peripheralUuid == ev.DeviceUUID.String() {
			ble.StopScanning()

			fmt.Println()
			fmt.Println("peripheral with UUID", ev.DeviceUUID, "found")

			advertisement := ev.Peripheral.Advertisement

			localName := advertisement.LocalName
			txPowerLevel := advertisement.TxPowerLevel
			manufacturerData := advertisement.ManufacturerData
			serviceData := advertisement.ServiceData
			//serviceUuids := advertisement.ServiceUuids

			if len(localName) > 0 {
				fmt.Println("  Local Name        =", localName)
			}

			if txPowerLevel != 0 {
				fmt.Println("  TX Power Level    =", txPowerLevel)
			}

			if len(manufacturerData) > 0 {
				fmt.Println("  Manufacturer Data =", manufacturerData)
			}

			if len(serviceData) > 0 {
				fmt.Println("  Service Data      =", serviceData)
			}

			fmt.Println()
			//explore(peripheral)

			done <- true
		}
	})

	if *verbose {
		log.Println("Init...")
	}

	ble.Init()
	_ = <-done
}
