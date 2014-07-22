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

type UUID [16]byte

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

	switch id {
	case 6: // state change
		state := args["kCBMsgArgState"].(int64)
		states := []string{"unknown", "resetting", "unsupported", "unauthorized", "poweredOff", "poweredOn"}
		log.Printf("event: adapterState %v\n", states[state])

	default:
		log.Printf("event: %#v\n", event)
	}
}

// send a message to Blued
func (ble *BLE) sendCBMsg(id int, args dict) {
	C.XpcSendMessage(ble.conn, goToXpc(dict{"kCBMsgId": 1, "kCBMsgArgs": args}), true)
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
