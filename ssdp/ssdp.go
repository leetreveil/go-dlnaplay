package ssdp

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
	SSDP_PORT           = 1900
	SSDP_ALIVE          = "ssdp:alive"
	SSDP_BYEBYE         = "ssdp:byebye"
	SSDP_UPDATE         = "ssdp:update"
	SSDP_ALL            = "ssdp:all"
	UDP_MAX_PACKET_SIZE = 65536
)

var (
	BROADCAST_ADDR = fmt.Sprintf("239.255.255.250:%d", SSDP_PORT)
)

func makeMSearchString(searchType string) string {
	return fmt.Sprintf(
		"M-SEARCH * HTTP/1.1\r\nHost:%s\r\nST:%s\r\nMan:\"ssdp:discover\"\r\nMX:3\r\n\r\n",
		BROADCAST_ADDR, searchType)
}

type controlPoint struct {
	callback func(map[string][]string)
}

func readData(conn *net.UDPConn, callback func(*bufio.Reader)) {
	msg := make([]byte, UDP_MAX_PACKET_SIZE)
	for {
		if n, err := conn.Read(msg); nil != err {
			panic(err)
		} else {
			callback(bufio.NewReaderSize(bytes.NewBuffer(msg), n))
		}
	}
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
	addr, err := net.ResolveUDPAddr(BROADCAST_VERSION, "0.0.0.0:0")
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP(BROADCAST_VERSION, addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	waddr, err := net.ResolveUDPAddr(BROADCAST_VERSION, BROADCAST_ADDR)
	if err != nil {
		panic(err)
	}

	_, err = conn.WriteTo([]byte(makeMSearchString(searchType)), waddr)
	if err != nil {
		panic(err)
	}

	readData(conn, func(rdr *bufio.Reader) {
		res, err := http.ReadResponse(rdr, nil)
		// TODO: malformed request - do we really want to kill the program??
		if err != nil {
			panic(err)
		}
		c.callback(res.Header)
	})
}

func (c *controlPoint) listen(iface *net.Interface) {
	addr, err := net.ResolveUDPAddr(BROADCAST_VERSION, BROADCAST_ADDR)
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenMulticastUDP(BROADCAST_VERSION, iface, addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	p := ipv4.NewPacketConn(conn)

	err = p.SetMulticastTTL(4)
	if err != nil {
		panic(err)
	}

	err = p.SetMulticastLoopback(true)
	if err != nil {
		panic(err)
	}

	readData(conn, func(rdr *bufio.Reader) {
		req, err := http.ReadRequest(rdr)
		// TODO: malformed request - do we really want to kill the program??
		if err != nil {
			panic(err)
		}
		c.callback(req.Header)
	})
}

func MakeControlPoint() *controlPoint {
	return &controlPoint{}
}

func TestMe() {
	cp := MakeControlPoint()
	cp.callback = func(headers map[string][]string) {
		// fmt.Println(headers)
		fmt.Println(headers["Location"])
	}

	for _, iface := range findCompatibleInterfaces() {
		go cp.listen(&iface)
	}

	cp.search("urn:schemas-upnp-org:service:AVTransport:1")
}
