package main

import (
        "log"
        "time"

	"../../goble"
)

func main() {
	ble := goble.NewBLE()

	/*
		ble.sendCBMsg(0, dict{
			"kCBMsgArgAlert": 1,
			"kCBMsgArgName":  "node",
		})
	*/


        log.Println("Init...")
	ble.Init()

        time.Sleep(1 * time.Second)
        log.Println("Start Advertising...")
        ble.StartAdvertising("gobble", []goble.UUID{})

        time.Sleep(5 * time.Second)
        log.Println("Stop Advertising...")
        ble.StopAdvertising()

        time.Sleep(1 * time.Second)
        log.Println("Start Scanning...")
        ble.StartScanning([]goble.UUID{}, true)

        time.Sleep(5 * time.Second)
        log.Println("Stop Scanning...")
        ble.StopScanning()

        time.Sleep(5 * time.Second)
        log.Println("Goodbye!")
}
