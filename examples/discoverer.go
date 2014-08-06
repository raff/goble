package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"../../goble"
)

func main() {
	verbose := flag.Bool("verbose", false, "dump all events")
	dups := flag.Bool("allow-duplicates", false, "allow duplicates when scanning")
	flag.Parse()

	ble := goble.New()
	ble.SetVerbose(*verbose)

	if *verbose {
		ble.On(goble.ALL, func(ev goble.Event) {
			log.Println("event", ev)
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
		fmt.Println("peripheral discovered (", ev.DeviceUUID, "):")
		fmt.Println("\thello my local name is:")
		fmt.Println("\t\t", ev.Peripheral.Advertisement.LocalName)
		fmt.Println("\tcan I interest you in any of the following advertised services:")
		fmt.Println("\t\t", ev.Peripheral.Services)
	})

	/*
	   var serviceData = peripheral.advertisement.serviceData;
	   if (serviceData && serviceData.length) {
	     console.log('\there is my service data:');
	     for (var i in serviceData) {
	       console.log('\t\t' + JSON.stringify(serviceData[i].uuid) + ': ' + JSON.stringify(serviceData[i].data.toString('hex')));
	     }
	   }
	   if (peripheral.advertisement.manufacturerData) {
	     console.log('\there is my manufacturer data:');
	     console.log('\t\t' + JSON.stringify(peripheral.advertisement.manufacturerData.toString('hex')));
	   }
	   if (peripheral.advertisement.txPowerLevel !== undefined) {
	     console.log('\tmy TX power level is:');
	     console.log('\t\t' + peripheral.advertisement.txPowerLevel);
	   }

	   console.log();
	*/

	if *verbose {
		log.Println("Init...")
	}

	ble.Init()

	time.Sleep(60 * time.Second)
	log.Println("Goodbye!")
}
