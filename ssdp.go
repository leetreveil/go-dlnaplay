package main

import (
	"bufio"
	"bytes"
	"code.google.com/p/go.net/ipv4"
	"fmt"
	"net"
	"net/http"
)

const (
	BROADCAST_VERSION   = "udp4"
	SSDP_IP             = "239.255.255.250"
	SSDP_PORT           = 1900
	SSDP_ALIVE          = "ssdp:alive"
	SSDP_BYEBYE         = "ssdp:byebye"
	SSDP_UPDATE         = "ssdp:update"
	SSDP_ALL            = "ssdp:all"
	UDP_MAX_PACKET_SIZE = 65536
)

var (
	SSDP_ADDR = net.UDPAddr{IP: net.ParseIP(SSDP_IP), Port: SSDP_PORT}
)

func makeMSearchString(searchType string) string {
	return fmt.Sprintf(
		"M-SEARCH * HTTP/1.1\r\nHost:%s\r\nST:%s\r\nMan:\"ssdp:discover\"\r\nMX:3\r\n\r\n",
		SSDP_ADDR.String(), searchType)
}

type controlPoint struct {
	callback func(map[string][]string)
}

func readData(conn *net.UDPConn) <-chan *bufio.Reader {
	ch := make(chan *bufio.Reader)
	go func() {
		msg := make([]byte, UDP_MAX_PACKET_SIZE)
		for {
			if n, err := conn.Read(msg); err != nil {
				panic(err)
			} else {
				ch <- bufio.NewReaderSize(bytes.NewBuffer(msg), n)
			}
		}
		close(ch)
	}()
	return ch
}

func findCompatibleInterfaces() []net.Interface {
	interfaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	out := make([]net.Interface, 0, len(interfaces))

	for _, iface := range interfaces {
		if (iface.Flags & net.FlagLoopback) != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			panic(err)
		}

		if len(addrs) == 0 {
			continue
		}

		out = append(out, iface)
	}

	return out
}

func (c *controlPoint) search(searchType string) {
	conn, err := net.ListenUDP(BROADCAST_VERSION, &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	if _, err := conn.WriteTo([]byte(makeMSearchString(searchType)), &SSDP_ADDR); err != nil {
		panic(err)
	}

	for rdr := range readData(conn) {
		res, err := http.ReadResponse(rdr, nil)
		// TODO: malformed request - do we really want to kill the program??
		if err != nil {
			panic(err)
		}
		c.callback(res.Header)
	}
}

func (c *controlPoint) listen(iface *net.Interface) {
	conn, err := net.ListenMulticastUDP(BROADCAST_VERSION, iface, &SSDP_ADDR)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	p := ipv4.NewPacketConn(conn)

	if err := p.SetMulticastTTL(4); err != nil {
		panic(err)
	}

	if err := p.SetMulticastLoopback(true); err != nil {
		panic(err)
	}

	for rdr := range readData(conn) {
		req, err := http.ReadRequest(rdr)
		// TODO: malformed request - do we really want to kill the program??
		if err != nil {
			panic(err)
		}
		c.callback(req.Header)
	}
}

func MakeControlPoint() *controlPoint {
	return &controlPoint{}
}

func TestMe() {
	cp := MakeControlPoint()
	cp.callback = func(headers map[string][]string) {
		fmt.Println(headers)
		// fmt.Println(headers["Location"])
	}

	for _, iface := range findCompatibleInterfaces() {
		go cp.listen(&iface)
	}

	cp.search("urn:schemas-upnp-org:service:AVTransport:1")
}
