package main

import (
	"encoding/xml"
	// "fmt"
	"io/ioutil"
)

type Root struct {
	Device Device `xml:"device"`
}

type Device struct {
	FriendlyName string      `xml:"friendlyName"`
	ServiceList  ServiceList `xml:"serviceList"`
}

type ServiceList struct {
	Services []Service `xml:"service"`
}

type Service struct {
	ServiceType string `xml:"serviceType"`
	ServiceId   string `xml:"serviceId"`
	ControlURL  string `xml:"controlURL"`
}

func Parse(data []byte) Device {
	bytes, err := ioutil.ReadFile("tst.xml")
	if err != nil {
		panic(err)
	}

	var q Root
	err = xml.Unmarshal(bytes, &q)
	if err != nil {
		panic(err)
	}

	// // fmt.Println(q.Services)
	// for _, service := range q.Device.ServiceList.Services {
	// 	fmt.Printf("\t%s\n", service)
	// }

	return q.Device
}
