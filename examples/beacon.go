package main

import (
	"flag"
	"log"
	"time"

	"github.com/raff/goble"
	"github.com/raff/goble/xpc"
)

func main() {
	uuid := flag.String("uuid", "1BEAC099-BEAC-BEAC-BEAC-BEAC09BEAC09", "iBeacon UUID")
	major := flag.Int("major", 0, "iBeacon major value (uint16)")
	minor := flag.Int("minor", 0, "iBeacon minor value (uint16)")
	power := flag.Int("power", -57, "iBeacon measured power (int8)")
	d := flag.Duration("duration", 1*time.Minute, "advertising duration")
	verbose := flag.Bool("verbose", false, "dump all events")
	flag.Parse()

	ble := goble.New()
	ble.SetVerbose(*verbose)
	ble.Init()

	var utsname xpc.Utsname
	xpc.Uname(&utsname)
	log.Println(utsname.Release)

	time.Sleep(1 * time.Second)

	log.Println("Start Advertising", *uuid, *major, *minor, *power)
	ble.StartAdvertisingIBeacon(xpc.MustUUID(*uuid), uint16(*major), uint16(*minor), int8(*power))

	time.Sleep(*d)

	log.Println("Stop Advertising")
	ble.StopAdvertising()
}
