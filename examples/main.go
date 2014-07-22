package main

import (
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

	ble.Init()

	waiter := make(chan int)
	<-waiter
}
