package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/raff/goble"
	"github.com/raff/goble/xpc"
)

func main() {
	verbose := flag.Bool("verbose", false, "dump all events")
	advertise := flag.Duration("advertise", 0, "Duration of advertising - 0: does not advertise")
	dups := flag.Bool("allow-duplicates", false, "allow duplicates when scanning")
	ibeacon := flag.Duration("ibeacon", 0, "Duration of IBeacon advertising - 0: does not advertise")
	scan := flag.Duration("scan", 10, "Duration of scanning - 0: does not scan")
	uuid := flag.String("uuid", "", "device uuid (for ibeacon uuid,major,minor,power)")
	connect := flag.Bool("connect", false, "connect to device")
	disconnect := flag.Bool("disconnect", false, "disconnect from device")
	rssi := flag.Bool("rssi", false, "update rssi for device")
	remove := flag.Bool("remove", false, "Remove all services")
	discover := flag.Bool("discover", false, "Discover services")

	flag.Parse()

	ble := goble.New()

	ble.SetVerbose(*verbose)

	log.Println("Init...")
	ble.Init()

	if *advertise > 0 {
		uuids := []xpc.UUID{}

		if len(*uuid) > 0 {
			uuids = append(uuids, xpc.MakeUUID(*uuid))
		}

		time.Sleep(1 * time.Second)
		log.Println("Start Advertising...")
		ble.StartAdvertising("gobble", uuids)

		time.Sleep(*advertise)
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
		ble.StartAdvertisingIBeacon(xpc.MakeUUID(id), major, minor, measuredPower)

		time.Sleep(*ibeacon)
		log.Println("Stop Advertising...")
		ble.StopAdvertising()
	}

	if *scan > 0 {
		time.Sleep(1 * time.Second)
		log.Println("Start Scanning...")
		ble.StartScanning(nil, *dups)

		time.Sleep(*scan)
		log.Println("Stop Scanning...")
		ble.StopScanning()
	}

	if *connect {
		time.Sleep(1 * time.Second)
		uuid := xpc.MakeUUID(*uuid)
		log.Println("Connect", uuid)
		ble.Connect(uuid)
		time.Sleep(5 * time.Second)
	}

	if *rssi {
		time.Sleep(1 * time.Second)
		uuid := xpc.MakeUUID(*uuid)
		log.Println("UpdateRssi", uuid)
		ble.UpdateRssi(uuid)
		time.Sleep(5 * time.Second)
	}

	if *discover {
		time.Sleep(1 * time.Second)
		uuid := xpc.MakeUUID(*uuid)
		log.Println("DiscoverServices", uuid)
		ble.DiscoverServices(uuid, nil)
		time.Sleep(5 * time.Second)
	}

	if *disconnect {
		time.Sleep(1 * time.Second)
		uuid := xpc.MakeUUID(*uuid)
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
