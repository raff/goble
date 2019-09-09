package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/raff/goble"
)

func main() {
	verbose := flag.Bool("verbose", false, "dump all events")
	compact := flag.Bool("compact", true, "compact messages")
	dups := flag.Bool("allow-duplicates", false, "allow duplicates when scanning")
	flag.Parse()

	var quit chan bool

	ble := goble.New()
	ble.SetVerbose(*verbose)

	if *verbose {
		ble.On(goble.ALL, func(ev goble.Event) (done bool) {
			log.Println("Event", ev)
			return
		})
	}

	ble.On("stateChange", func(ev goble.Event) (done bool) {
                if *verbose {
                        fmt.Println("stateChange", ev.State)
                }
		if ev.State == "poweredOn" {
			ble.StartScanning(nil, *dups)
		} else {
			ble.StopScanning()
			done = true
			quit <- true
		}

		return
	})

	ble.On("discover", func(ev goble.Event) (done bool) {
                if *verbose {
                        fmt.Println("discover", ev.State)
                }
		if *compact {
			fmt.Println("peripheral:", ev.DeviceUUID)
			if ev.Peripheral.Advertisement.LocalName != "" {
				fmt.Println("  name:", ev.Peripheral.Advertisement.LocalName)
			}
			if len(ev.Peripheral.Advertisement.ServiceUuids) > 0 {
				fmt.Println("  services:", ev.Peripheral.Advertisement.ServiceUuids)
			}
		} else {
			fmt.Println()
			fmt.Println("peripheral discovered (", ev.DeviceUUID, "):")
			fmt.Println("\thello my local name is:")
			fmt.Println("\t\t", ev.Peripheral.Advertisement.LocalName)
			fmt.Println("\tcan I interest you in any of the following advertised services:")
			fmt.Println("\t\t", ev.Peripheral.Advertisement.ServiceUuids)
		}

		serviceData := ev.Peripheral.Advertisement.ServiceData
		if len(serviceData) > 0 {
			prefix := "\t\t"

			if *compact {
				prefix = "    "
				fmt.Println("  service data:")
			} else {
				fmt.Println("\there is my service data:")
			}

			for _, d := range serviceData {
				fmt.Println(prefix, d.Uuid, ":", d.Data)
			}
		}

		if len(ev.Peripheral.Advertisement.ManufacturerData) > 0 {
			if *compact {
				fmt.Printf("  manufacturer data: %x\n", ev.Peripheral.Advertisement.ManufacturerData)
			} else {
				fmt.Println("\there is my manufacturer data:")
				fmt.Println("\t\t", ev.Peripheral.Advertisement.ManufacturerData)
			}
		}

		if ev.Peripheral.Advertisement.TxPowerLevel != 0 {
			if *compact {
				fmt.Println("  TX power level:", ev.Peripheral.Advertisement.TxPowerLevel)
			} else {
				fmt.Println("\tmy TX power level is:")
				fmt.Println("\t\t", ev.Peripheral.Advertisement.TxPowerLevel)
			}
		}

		if *compact {
			fmt.Println()
		}

		return
	})

	if *verbose {
		log.Println("Init...")
	}

	ble.Init()

	<-quit
}
