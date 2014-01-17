package main

import (
	"io/ioutil"
	"testing"
	// "fmt"
)

func TestParse(t *testing.T) {
	bytes, err := ioutil.ReadFile("tst.xml")
	if err != nil {
		panic(err)
	}

	device := Parse(bytes)

	expected := "[TV][LG]47LN578V-ZE"
	if device.FriendlyName != expected {
		t.Errorf("actual = %v, expected = %v", device.FriendlyName, expected)
	}

	expected = "/AVTransport/09eb9c7c-410a-29c6-bd82-593e90df1082/control.xml"
	result := device.ServiceList.Services[2].ControlURL
	if result != expected {
		t.Errorf("actual = %v, expected = %v", result, expected)
	}
}
