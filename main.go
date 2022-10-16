package main

import (
	"time"

	"tinygo.org/x/bluetooth"

	"github.com/go-vgo/robotgo"
)

var bleMac = ""

var bleRssi int16 = -50

var scanTimeout = 30

var isLock = false

var adapter = bluetooth.DefaultAdapter

func main() {
	must("enable BLE stack", adapter.Enable())

	// Start scanning.
	go timeoutStopScan()
	doLock := true
	err := adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		if bleMac == "scan" {
			println(device.Address.String(), device.RSSI, device.LocalName())
		}
		if device.Address.String() == bleMac {
			adapter.StopScan()
			if device.RSSI < bleRssi {
				println("Found RSSI low", device.Address.String(), device.RSSI, device.LocalName())
			} else {
				doLock = false
				println("Found OK", device.Address.String(), device.RSSI, device.LocalName())
			}
		}
	})
	must("start scan", err)
	if err == nil && doLock {
		keyLock()
	}
}

func timeoutStopScan() {
	time.Sleep(time.Duration(scanTimeout) * time.Second)
	adapter.StopScan()
}

func keyLock() {
	isLock = true
	robotgo.KeyTap("l", "command")
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
