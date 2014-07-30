package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"../../goble"
)

func main() {
	verbose := flag.Bool("verbose", false, "dump all events")
	advertise := flag.Int("advertise", 5, "Duration of advertising")
	scan := flag.Int("scan", 10, "Duration of scanning")
	uuid := flag.String("uuid", "", "device uuid (for ibeacon uuid,major,minor,power)")
	connect := flag.Bool("connect", false, "connect to device")
	disconnect := flag.Bool("disconnect", false, "disconnect from device")
	rssi := flag.Bool("rssi", false, "update rssi for device")
	ibeacon := flag.Int("ibeacon", 0, "Duration of IBeacon advertising")
	remove := flag.Bool("remove", false, "Remove all services")

	flag.Parse()

	ble := goble.New()

	ble.SetVerbose(*verbose)

	log.Println("Init...")
	ble.Init()

	if *advertise > 0 {
		uuids := []goble.UUID{}

		if len(*uuid) > 0 {
			uuids = append(uuids, goble.MakeUUID(*uuid))
		}

		time.Sleep(1 * time.Second)
		log.Println("Start Advertising...")
		ble.StartAdvertising("gobble", uuids)

		time.Sleep(time.Duration(*advertise) * time.Second)
		log.Println("Stop Advertising...")
		ble.StopAdvertising()
	}

	if *ibeacon > 0 {
		parts := strings.Split(*uuid, ",")
		id := parts[0]

		var major, minor uint16
		var measuredPower int8

		if len(parts) > 1 {
			fmt.Sscanf(parts[1], "%d", &major)
		}
		if len(parts) > 2 {
			fmt.Sscanf(parts[2], "%d", &minor)
		}
		if len(parts) > 2 {
			fmt.Sscanf(parts[3], "%d", &measuredPower)
		}

		time.Sleep(1 * time.Second)
		log.Println("Start Advertising IBeacon...")
		ble.StartAdvertisingIBeacon(goble.MakeUUID(id), major, minor, measuredPower)

		time.Sleep(time.Duration(*ibeacon) * time.Second)
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

	if *remove {
		time.Sleep(1 * time.Second)
		log.Println("Remove all services")
		ble.RemoveServices()
		time.Sleep(5 * time.Second)
	}

	time.Sleep(5 * time.Second)
	log.Println("Goodbye!")
}
