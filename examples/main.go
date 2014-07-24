package main

import (
	"flag"
	"log"
	"time"

	"../../goble"
)

func main() {
	verbose := flag.Bool("verbose", false, "dump all events")
	advertise := flag.Int("advertise", 5, "Duration of advertising")
	scan := flag.Int("scan", 10, "Duration of scanning")
	uuid := flag.String("uuid", "", "device uuid")
	connect := flag.Bool("connect", false, "connect to device")
	disconnect := flag.Bool("disconnect", false, "disconnect from device")
	rssi := flag.Bool("rssi", false, "update rssi for device")

	flag.Parse()

	ble := goble.NewBLE()

	ble.SetVerbose(*verbose)

	log.Println("Init...")
	ble.Init()

	if *advertise > 0 {
		time.Sleep(1 * time.Second)
		log.Println("Start Advertising...")
		ble.StartAdvertising("gobble", []goble.UUID{})

		time.Sleep(time.Duration(*advertise) * time.Second)
		log.Println("Stop Advertising...")
		ble.StopAdvertising()
	}

	if *scan > 0 {
		time.Sleep(1 * time.Second)
		log.Println("Start Scanning...")
		ble.StartScanning([]goble.UUID{}, true)

		time.Sleep(time.Duration(*scan) * time.Second)
		log.Println("Stop Scanning...")
		ble.StopScanning()
	}

	if *connect {
		time.Sleep(1 * time.Second)
		uuid := goble.MakeUUID(*uuid)
		log.Println("Connect", uuid)
		ble.Connect(uuid)
		time.Sleep(5 * time.Second)
	}

	if *rssi {
		time.Sleep(1 * time.Second)
		uuid := goble.MakeUUID(*uuid)
		log.Println("UpdateRssi", uuid)
		ble.UpdateRssi(uuid)
		time.Sleep(5 * time.Second)
	}

	if *disconnect {
		time.Sleep(1 * time.Second)
		uuid := goble.MakeUUID(*uuid)
		log.Println("Disconnect", uuid)
		ble.Disconnect(uuid)
		time.Sleep(5 * time.Second)
	}

	time.Sleep(5 * time.Second)
	log.Println("Goodbye!")
}
