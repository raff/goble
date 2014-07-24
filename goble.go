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

type ServiceData struct {
	uuid string
	data []byte
}

type Advertisement struct {
	localName        string
	txPowerLevel     int64
	manufacturerData []byte
	serviceData      []ServiceData
	serviceUuids     []string
}

type Peripheral struct {
	uuid          string
	advertisement Advertisement
	rssi          int64
}

type BLE struct {
	conn        C.xpc_connection_t
	peripherals map[string]Peripheral
}

func NewBLE() *BLE {
	ble := &BLE{peripherals: map[string]Peripheral{}}
	ble.conn = XpcConnect("com.apple.blued", ble)
	return ble
}

// process BLE events and asynchronous errors
// (implements XpcEventHandler)
func (ble *BLE) HandleXpcEvent(event dict, err error) {
	id := event["kCBMsgId"].(int64)
	args := event["kCBMsgArgs"].(dict) // what happens if there are no args ?

	log.Printf("event: %v %#v\n", id, args)

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

	case 37:
		advdata := args["kCBMsgArgAdvertisementData"].(dict)
		if len(advdata) == 0 {
			//log.Println("event: discover with no advertisment data")
			break
		}

		devuuid := args["kCBMsgArgDeviceUUID"].(UUID).String()

		advertisement := Advertisement{
			localName:        advdata.GetString("kCBAdvDataLocalName", args.GetString("kCBMsgArgName", "")),
			txPowerLevel:     advdata.GetInt("kCBAdvDataTxPowerLevel", 0),
			manufacturerData: advdata.GetBytes("kCBAdvDataManufacturerData", nil),
			serviceData:      []ServiceData{},
			serviceUuids:     []string{},
		}

		rssi := args.GetInt("kCBMsgArgRssi", 0)

		if uuids, ok := advdata["kCBAdvDataServiceUUIDs"]; ok {
			for _, uuid := range uuids.(array) {
				advertisement.serviceUuids = append(advertisement.serviceUuids, GetUUID(uuid).String())
			}
		}

		if sdata, ok := advdata["kCBAdvDataServiceData"]; ok {
			for _, data := range sdata.(array) {
				bytes := data.([]byte)
				sd := ServiceData{
					uuid: fmt.Sprintf("%x", bytes[0]),
					data: bytes[1:],
				}

				advertisement.serviceData = append(advertisement.serviceData, sd)
			}
		}

		ble.peripherals[devuuid] = Peripheral{
			uuid:          devuuid,
			advertisement: advertisement,
			rssi:          rssi,
		}

		log.Println("event: discover", devuuid, advertisement, rssi)
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

	for _, uuid := range serviceUuids {
		uuids = append(uuids, uuid.String())
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

	for _, uuid := range serviceUuids {
		uuids = append(uuids, uuid.String())
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
