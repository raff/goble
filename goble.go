package goble

/*
#include <xpc/xpc.h>
#include "xpc_wrapper.h"
*/
import "C"

import (
	"fmt"
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
	ble.conn = C.XpcConnectBlued()
	return &ble
}

func (ble *BLE) sendCBMsg(id int, args dict) {
	C.XpcSendMessage(ble.conn, goToXpc(dict{"kCBMsgId": 1, "kCBMsgArgs": args}), true)
}

func (ble *BLE) Init() {
	ble.sendCBMsg(1, dict{"kCBMsgArgName": fmt.Sprintf("node-%v", time.Now().Unix()),
		"kCBMsgArgOptions": dict{"kCBInitOptionShowPowerAlert": 0}, "kCBMsgArgType": 0})
}

func (ble *BLE) StartAdvertising(name string, serviceUuids []UUID) {
	uuids := []string{}

	for uuid := range serviceUuids {
		uuids = append(uuids, fmt.Sprintf("%x", uuid))
	}

	ble.sendCBMsg(8, dict{"kCBAdvDataLocalName": name, "kCBAdvDataServiceUUIDs": uuids})
}
