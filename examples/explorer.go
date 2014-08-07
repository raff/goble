package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"../../goble"
)

func explore(ble *goble.BLE, peripheral *goble.Peripheral) {

	ble.On("connect", func(ev goble.Event) (done bool) {
		log.Println("connected", ev)
		ble.DiscoverServices(ev.DeviceUUID, nil)
		return
	})

	ble.On("disconnect", func(ev goble.Event) (done bool) {
		log.Println("disconnected", ev)
		os.Exit(0)
		return true
	})

	ble.On("servicesDiscover", func(ev goble.Event) (done bool) {
		log.Println("services", ev)
		return
	})

	fmt.Println("services and characteristics:")
	ble.Connect(peripheral.Uuid)

	/*

	     peripheral.on('disconnect', function() {
	       process.exit(0);
	     });

	     peripheral.connect(function(error) {
	       peripheral.discoverServices([], function(error, services) {
	         var serviceIndex = 0;

	         async.whilst(
	           function () {
	             return (serviceIndex < services.length);
	           },
	           function(callback) {
	             var service = services[serviceIndex];
	             var serviceInfo = service.uuid;

	             if (service.name) {
	               serviceInfo += ' (' + service.name + ')';
	             }
	             console.log(serviceInfo);

	             service.discoverCharacteristics([], function(error, characteristics) {
	               var characteristicIndex = 0;

	               async.whilst(
	                 function () {
	                   return (characteristicIndex < characteristics.length);
	                 },
	                 function(callback) {
	                   var characteristic = characteristics[characteristicIndex];
	                   var characteristicInfo = '  ' + characteristic.uuid;

	                   if (characteristic.name) {
	                     characteristicInfo += ' (' + characteristic.name + ')';
	                   }

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
	           function (err) {
	             peripheral.disconnect();
	           }
	         );
	       });
	     });
	   }
	*/
}

func main() {
	verbose := flag.Bool("verbose", false, "dump all events")
	dups := flag.Bool("allow-duplicates", false, "allow duplicates when scanning")
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
		if ev.State == "poweredOn" {
			ble.StartScanning(nil, *dups)
		} else {
			ble.StopScanning()
			done = true
		}

		return
	})

	ble.On("discover", func(ev goble.Event) (done bool) {
		if peripheralUuid == ev.DeviceUUID.String() {
			ble.StopScanning()

			fmt.Println()
			fmt.Println("peripheral with UUID", ev.DeviceUUID, "found")

			advertisement := ev.Peripheral.Advertisement

			localName := advertisement.LocalName
			txPowerLevel := advertisement.TxPowerLevel
			manufacturerData := advertisement.ManufacturerData
			serviceData := advertisement.ServiceData
			//serviceUuids := advertisement.ServiceUuids

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
