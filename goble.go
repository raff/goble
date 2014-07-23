package goble

/*
#include <xpc/xpc.h>
#include "xpc_wrapper.h"
*/
import "C"

import (
	"fmt"
	"log"
	"time"
)

//
// BLE support
//

var STATES = []string{"unknown", "resetting", "unsupported", "unauthorized", "poweredOff", "poweredOn"}

type Advertisement struct {
	localName        string
	txPowerLevel     int
	manufacturerData []byte
	serviceData      []byte
	serviceUuids     []UUID
}

type Peripheral struct {
	uuid          UUID
	advertisement Advertisement
	rssi          int
}

type BLE struct {
	conn C.xpc_connection_t
}

func NewBLE() *BLE {
	ble := BLE{}
	ble.conn = XpcConnect(ble.eventHandler)
	return &ble
}

// process BLE events and asynchronous errors
func (ble *BLE) eventHandler(event dict, err error) {
	id := event["kCBMsgId"].(int64)
	args := event["kCBMsgArgs"].(dict) // what happens if there are no args ?

	log.Printf("event: %v %v\n", id, args)

	switch id {
	case 6: // state change
		state := args["kCBMsgArgState"].(int64)
		log.Printf("event: stateChange %v\n", STATES[state])

	case 16: // advertising start
		result := args["kCBMsgArgResult"].(int64)
		if result != 0 {
			log.Printf("event: error in advertisingStart %v\n", result)
		} else {
			log.Println("event: advertisingStart")
		}

	case 17: // advertising stop
		result := args["kCBMsgArgResult"].(int64)
		if result != 0 {
			log.Printf("event: error in advertisingStop %v\n", result)
		} else {
			log.Println("event: advertisingStop")
		}
	}
}

// send a message to Blued
func (ble *BLE) sendCBMsg(id int, args dict) {
	C.XpcSendMessage(ble.conn, goToXpc(dict{"kCBMsgId": id, "kCBMsgArgs": args}), true)
}

// initialize BLE
func (ble *BLE) Init() {
	ble.sendCBMsg(1, dict{"kCBMsgArgName": fmt.Sprintf("node-%v", time.Now().Unix()),
		"kCBMsgArgOptions": dict{"kCBInitOptionShowPowerAlert": 0}, "kCBMsgArgType": 0})
}

// start advertising
func (ble *BLE) StartAdvertising(name string, serviceUuids []UUID) {
	uuids := []string{}

	for uuid := range serviceUuids {
		uuids = append(uuids, fmt.Sprintf("%x", uuid))
	}

	ble.sendCBMsg(8, dict{"kCBAdvDataLocalName": name, "kCBAdvDataServiceUUIDs": uuids})
}

// start advertising as IBeacon
func (ble *BLE) StartAdvertisingIBeacon(name string, data []byte) {
	ble.sendCBMsg(8, dict{"kCBAdvDataAppleBeaconKey": data})
}

// stop advertising
func (ble *BLE) StopAdvertising() {
	ble.sendCBMsg(9, nil)
}

// start scanning
func (ble *BLE) StartScanning(serviceUuids []UUID, allowDuplicates bool) {
	uuids := []string{}

	for uuid := range serviceUuids {
		uuids = append(uuids, fmt.Sprintf("%x", uuid))
	}

	args := dict{"kCBMsgArgUUIDs": uuids}
	if allowDuplicates {
		args["kCBMsgArgOptions"] = dict{"kCBScanOptionAllowDuplicates": 1}
	} else {
		args["kCBMsgArgOptions"] = dict{}
	}

	ble.sendCBMsg(29, args)
}

// stop scanning
func (ble *BLE) StopScanning() {
	ble.sendCBMsg(30, nil)
}
