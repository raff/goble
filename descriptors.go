package goble

// A dictionary of known descriptor names and type (keyed by descriptor uuid)
var knownDescriptors = map[string]struct{ Name, Type string }{
	"2900": {Name: "Characteristic Extended Properties", Type: "org.bluetooth.descriptor.gatt.characteristic_extended_properties"},
	"2901": {Name: "Characteristic User Description", Type: "org.bluetooth.descriptor.gatt.characteristic_user_description"},
	"2902": {Name: "Client Characteristic Configuration", Type: "org.bluetooth.descriptor.gatt.client_characteristic_configuration"},
	"2903": {Name: "Server Characteristic Configuration", Type: "org.bluetooth.descriptor.gatt.server_characteristic_configuration"},
	"2904": {Name: "Characteristic Presentation Format", Type: "org.bluetooth.descriptor.gatt.characteristic_presentation_format"},
	"2905": {Name: "Characteristic Aggregate Format", Type: "org.bluetooth.descriptor.gatt.characteristic_aggregate_format"},
	"2906": {Name: "Valid Range", Type: "org.bluetooth.descriptor.valid_range"},
	"2907": {Name: "External Report Reference", Type: "org.bluetooth.descriptor.external_report_reference"},
	"2908": {Name: "Report Reference", Type: "org.bluetooth.descriptor.report_reference"},
}
