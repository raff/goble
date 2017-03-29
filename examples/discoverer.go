package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/raff/goble"
)

func main() {
	verbose := flag.Bool("verbose", false, "dump all events")
	dups := flag.Bool("allow-duplicates", false, "allow duplicates when scanning")
	flag.Parse()

	ble := goble.New()
	ble.SetVerbose(*verbose)

	if *verbose {
		ble.On(goble.ALL, func(ev goble.Event) (done bool) {
			log.Println("Event", ev)
			return
		})
	}

	ble.On("stateChange", func(ev goble.Event) (done bool) {
		if ev.State == "poweredOn" {
			ble.StartScanning(nil, *dups)
		} else {
			ble.StopScanning()
			done = true
		}

		return
	})

	ble.On("discover", func(ev goble.Event) (done bool) {
		fmt.Println()
		fmt.Println("peripheral discovered (", ev.DeviceUUID, "):")
		fmt.Println("\thello my local name is:")
		fmt.Println("\t\t", ev.Peripheral.Advertisement.LocalName)
		fmt.Println("\tcan I interest you in any of the following advertised services:")
		fmt.Println("\t\t", ev.Peripheral.Advertisement.ServiceUuids)

		serviceData := ev.Peripheral.Advertisement.ServiceData
		if len(serviceData) > 0 {
			fmt.Println("\there is my service data:")
			for _, d := range serviceData {
				fmt.Println("\t\t", d.Uuid, ":", d.Data)
			}
		}

		if len(ev.Peripheral.Advertisement.ManufacturerData) > 0 {
			fmt.Println("\there is my manufacturer data:")
			fmt.Println("\t\t", ev.Peripheral.Advertisement.ManufacturerData)
		}

		if ev.Peripheral.Advertisement.TxPowerLevel != 0 {
			fmt.Println("\tmy TX power level is:")
			fmt.Println("\t\t", ev.Peripheral.Advertisement.TxPowerLevel)
		}

		return
	})

	if *verbose {
		log.Println("Init...")
	}

	ble.Init()

	var done chan bool
	<-done
}
