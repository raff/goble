package main

import (
	"flag"
	"log"
	"time"

	"../../goble"
)

func main() {
	advertise := flag.Int("advertise", 5, "Duration of advertising")
	scan := flag.Int("scan", 10, "Duration of scanning")

	flag.Parse()

	ble := goble.NewBLE()

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

	time.Sleep(2 * time.Second)
	log.Println("Goodbye!")
}
