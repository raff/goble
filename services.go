package goble

// A dictionary of known service names and type (keyed by service uuid)
var knownServices = map[string]struct{ Name, Type string }{
	"1800": {Name: "Generic Access", Type: "org.bluetooth.service.generic_access"},
	"1801": {Name: "Generic Attribute", Type: "org.bluetooth.service.generic_attribute"},
	"1802": {Name: "Immediate Alert", Type: "org.bluetooth.service.immediate_alert"},
	"1803": {Name: "Link Loss", Type: "org.bluetooth.service.link_loss"},
	"1804": {Name: "Tx Power", Type: "org.bluetooth.service.tx_power"},
	"1805": {Name: "Current Time Service", Type: "org.bluetooth.service.current_time"},
	"1806": {Name: "Reference Time Update Service", Type: "org.bluetooth.service.reference_time_update"},
	"1807": {Name: "Next DST Change Service", Type: "org.bluetooth.service.next_dst_change"},
	"1808": {Name: "Glucose", Type: "org.bluetooth.service.glucose"},
	"1809": {Name: "Health Thermometer", Type: "org.bluetooth.service.health_thermometer"},
	"180a": {Name: "Device Information", Type: "org.bluetooth.service.device_information"},
	"180d": {Name: "Heart Rate", Type: "org.bluetooth.service.heart_rate"},
	"180e": {Name: "Phone Alert Status Service", Type: "org.bluetooth.service.phone_alert_service"},
	"180f": {Name: "Battery Service", Type: "org.bluetooth.service.battery_service"},
	"1810": {Name: "Blood Pressure", Type: "org.bluetooth.service.blood_pressuer"},
	"1811": {Name: "Alert Notification Service", Type: "org.bluetooth.service.alert_notification"},
	"1812": {Name: "Human Interface Device", Type: "org.bluetooth.service.human_interface_device"},
	"1813": {Name: "Scan Parameters", Type: "org.bluetooth.service.scan_parameters"},
	"1814": {Name: "Running Speed and Cadence", Type: "org.bluetooth.service.running_speed_and_cadence"},
	"1815": {Name: "Cycling Speed and Cadence", Type: "org.bluetooth.service.cycling_speed_and_cadence"},
}
