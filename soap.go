package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	soapTemplate = `<?xml version="1.0"?>
<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" 
          s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
  <s:Body>%s</s:Body>
</s:Envelope>`
)

func Play() {
	createSoapRequest(
		"http://192.168.0.19:1365/AVTransport/09eb9c7c-410a-29c6-bd82-593e90df1082/control.xml",
		"urn:upnp-org:serviceId:AVTransport:1#Play",
		`<u:Play xmlns:u="urn:schemas-upnp-org:service:AVTransport:1">
    <InstanceID>0</InstanceID>
    <Speed>1</Speed>
</u:Play>`)
}

func SetAVTransportURI() {
	createSoapRequest(
		"http://192.168.0.19:1365/AVTransport/09eb9c7c-410a-29c6-bd82-593e90df1082/control.xml",
		"urn:upnp-org:serviceId:AVTransport:1#SetAVTransportURI",
		`<u:SetAVTransportURI xmlns:u="urn:schemas-upnp-org:service:AVTransport:1">
    <InstanceID>0</InstanceID>
    <CurrentURI>http://192.168.0.17:50008/Kinsky?file=%2FUsers%2Fleetreveil%2FDownloads%2FEtherwood%20-%20Etherwood%20(2013)%20MP3-MEDIC36-CD-LP%2F13-etherwood-the_time_is_here_at_last_(feat_hybrid_minds).mp3</CurrentURI>
    <CurrentURIMetaData></CurrentURIMetaData>
</u:SetAVTransportURI>`)
}

func GetProtocolInfo() {
	createSoapRequest(
		"http://192.168.0.19:1365/ConnectionManager/09eb9c7c-410a-29c6-bd82-593e90df1082/control.xml",
		"urn:schemas-upnp-org:service:ConnectionManager:1#GetProtocolInfo",
		`<u:GetProtocolInfo xmlns:u="urn:schemas-upnp-org:service:ConnectionManager:1">
    </u:GetProtocolInfo>`)
}

func createSoapRequest(url string, soapAction string, soapBody string) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(fmt.Sprintf(soapTemplate, soapBody)))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "text/xml; charset=\"utf-8\"")
	req.Header.Add("SOAPAction", soapAction)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
	fmt.Println(string(body))
}
