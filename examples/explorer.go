package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/raff/goble"
)

var (
	debug   = flag.Bool("debug", false, "log debug messages")
	verbose = flag.Bool("verbose", false, "dump all events")
	dups    = flag.Bool("allow-duplicates", false, "allow duplicates when scanning")
)

func DebugPrint(params ...interface{}) {
	if *debug {
		log.Println(params...)
	}
}

type Result struct {
	count int
	data  string
}

func explore(ble *goble.BLE, peripheral *goble.Peripheral) {
	results := map[string]Result{}

	// connect
	ble.On("connect", func(ev goble.Event) (done bool) {
		DebugPrint("connected", ev)
		ble.DiscoverServices(ev.DeviceUUID, nil)

		go func() {
			time.Sleep(2 * time.Minute)
			ble.Disconnect(ev.DeviceUUID)
		}()

		return
	})

	// discover services
	ble.On("servicesDiscover", func(ev goble.Event) (done bool) {
		DebugPrint("serviceDiscovered", ev)
		for sid, service := range ev.Peripheral.Services {
			// this is a map that contains services UUIDs (string) and service startHandle (int)
			// for now we only process the "strings"
			if _, ok := sid.(string); ok {
				serviceInfo := service.Uuid

				if len(service.Name) > 0 {
					serviceInfo += " (" + service.Name + ")"
				}

				results[service.Uuid] = Result{data: serviceInfo}
				ble.DiscoverCharacterstics(ev.DeviceUUID, service.Uuid, nil)
			}
		}

		return
	})

	// discover characteristics
	ble.On("characteristicsDiscover", func(ev goble.Event) (done bool) {
		DebugPrint("characteristicsDiscovered", ev)
		serviceUuid := ev.ServiceUuid
		serviceResult := results[serviceUuid]

		for cid, characteristic := range ev.Peripheral.Services[serviceUuid].Characteristics {
			// this is a map that contains services UUIDs (string) and service startHandle (int)
			// for now we only process the "strings"
			if _, ok := cid.(string); ok {
				characteristicInfo := "  " + characteristic.Uuid

				if len(characteristic.Name) > 0 {
					characteristicInfo += " (" + characteristic.Name + ")"
				}

				characteristicInfo += "\n    properties  " + characteristic.Properties.String()
				serviceResult.data += characteristicInfo

				if characteristic.Properties.Readable() {
					serviceResult.count += 1
					ble.Read(ev.DeviceUUID, serviceUuid, characteristic.Uuid)
				}

				//ble.DiscoverDescriptors(ev.DeviceUUID, serviceUuid, characteristic.Uuid)
				results[serviceUuid] = serviceResult

				if *verbose {
					log.Println(results[serviceUuid])
				}
			}
		}

		return
	})

	// discover descriptors
	ble.On("descriptorsDiscover", func(ev goble.Event) (done bool) {
		DebugPrint("descriptorsDiscovered", ev)
		fmt.Println("    descriptors  ", ev.Peripheral.Services[ev.ServiceUuid].Characteristics[ev.CharacteristicUuid].Descriptors)
		return
	})

	// read
	ble.On("read", func(ev goble.Event) (done bool) {
		DebugPrint("read", ev)
		serviceUuid := ev.ServiceUuid
		serviceResult := results[serviceUuid]
		serviceResult.data += fmt.Sprintf("    value        %x | %q\n", ev.Data, ev.Data)
		serviceResult.count -= 1

		if serviceResult.count <= 0 {
			fmt.Println(serviceResult.data)
			return true
		} else {
			results[serviceUuid] = serviceResult
		}

		return
	})

	// disconnect
	ble.On("disconnect", func(ev goble.Event) (done bool) {
		DebugPrint("disconnected", ev)
		os.Exit(0)
		return true
	})

	fmt.Println("services and characteristics:")
	ble.Connect(peripheral.Uuid)

	/*
	           async.series([
	             function(callback) {
	               characteristic.discoverDescriptors(function(error, descriptors) {
	                 async.detect(
	                   descriptors,
	                   function(descriptor, callback) {
	                     return callback(descriptor.uuid === '2901")
	                   },
	                   function(userDescriptionDescriptor){
	                     if (userDescriptionDescriptor) {
	                       userDescriptionDescriptor.readValue(function(error, data) {
	                         characteristicInfo += ' (' + data.toString() + ')';
	                         callback();
	                       });
	                     } else {
	                       callback();
	                     }
	                   }
	                 );
	               });
	             },
	             function(callback) {
	                   characteristicInfo += '\n    properties  ' + characteristic.properties.join(', ")

	               if (characteristic.properties.indexOf('read') !== -1) {
	                 characteristic.read(function(error, data) {
	                   if (data) {
	                     var string = data.toString('ascii")

	                     characteristicInfo += '\n    value       ' + data.toString('hex') + ' | \'' + string + '\'';
	                   }
	                   callback();
	                 });
	               } else {
	                 callback();
	               }
	             },
	             function() {
	               console.log(characteristicInfo);
	               characteristicIndex++;
	               callback();
	             }
	           ]);
	         },
	         function(error) {
	           serviceIndex++;
	           callback();
	         }
	       );
	     });
	   },
	*/
}

func main() {
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Println("usage:", os.Args[0], "[options] peripheral-uuid")
		os.Exit(1)
	}

	peripheralUuid := flag.Args()[0]

	var done chan bool

	ble := goble.New()
	ble.SetVerbose(*verbose)

	if *verbose {
		ble.On(goble.ALL, func(ev goble.Event) (done bool) {
			log.Println("Event", ev)
			return
		})
	}

	ble.On("stateChange", func(ev goble.Event) (done bool) {
		DebugPrint("stateChanged", ev)
		if ev.State == "poweredOn" {
			ble.StartScanning(nil, *dups)
		} else {
			ble.StopScanning()
			done = true
		}

		return
	})

	ble.On("discover", func(ev goble.Event) (done bool) {
		DebugPrint("discovered", ev)
		if peripheralUuid == ev.DeviceUUID.String() {
			ble.StopScanning()

			fmt.Println()
			fmt.Println("peripheral with UUID", ev.DeviceUUID, "found")

			advertisement := ev.Peripheral.Advertisement

			DebugPrint("advertised", advertisement)

			localName := advertisement.LocalName
			txPowerLevel := advertisement.TxPowerLevel
			manufacturerData := advertisement.ManufacturerData
			serviceData := advertisement.ServiceData
			//serviceUuids := advertisement.ServiceUuids

			if ev.Peripheral.Connectable {
				fmt.Println("  Connectable")
			}
			if len(localName) > 0 {
				fmt.Println("  Local Name        =", localName)
			}

			if txPowerLevel != 0 {
				fmt.Println("  TX Power Level    =", txPowerLevel)
			}

			if len(manufacturerData) > 0 {
				fmt.Println("  Manufacturer Data =", manufacturerData)
			}

			if len(serviceData) > 0 {
				fmt.Println("  Service Data      =", serviceData)
			}

			fmt.Println()

			DebugPrint("explore", ev.Peripheral)
			explore(ble, &ev.Peripheral)
		}

		return
	})

	if *verbose {
		log.Println("Init...")
	}

	ble.Init()

	fmt.Println("waiting...")
	<-done

	fmt.Println("goodbye!")
	os.Exit(0)
}
