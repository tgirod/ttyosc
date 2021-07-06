package main

import (
	"flag"
	"log"
	"net"

	"github.com/Lobaro/slip"
	"go.bug.st/serial"
)

var tty = flag.String("t", "/dev/ttyUSB0", "serial device")
var remote = flag.String("r", ":6010", "remote address: osc from serial is sent here")
var local = flag.String("l", ":6020", "local address: send your osc here")
var baudRate = flag.Int("baud", 115200, "serial connection's baudrate")
var verbose = flag.Bool("v", false, "verbose")

func main() {
	flag.Parse()

	// Serial connection

	mode := &serial.Mode{
		BaudRate: *baudRate,
	}

	port, err := serial.Open(*tty, mode)
	if err != nil {
		log.Fatal(err)
	}

	var done = make(chan error)

	// udp -> serial communication

	localAddress, err := net.ResolveUDPAddr("udp", *local)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("local address:", localAddress.String())

	go func() {
		writer := slip.NewWriter(port)
		// open UDP server
		conn, err := net.ListenUDP("udp", localAddress)
		if err != nil {
			done <- err
			return
		}
		defer conn.Close()
		// listen to incoming data from UDP and send to serial
		buf := make([]byte, 1024)
		for {
			n, _, err := conn.ReadFromUDP(buf)
			writer.WritePacket(buf[:n])
			if *verbose {
				log.Println("udp -> serial")
			}
			if err != nil {
				done <- err
				return
			}
		}
	}()

	// serial -> udp communication

	remoteAddress, err := net.ResolveUDPAddr("udp", *remote)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("remote address:", remoteAddress.String())

	go func() {
		reader := slip.NewReader(port)
		// open UDP client
		conn, err := net.DialUDP("udp", nil, remoteAddress)
		if err != nil {
			done <- err
			return
		}
		defer conn.Close()
		// listen to incoming data from serial and send to UDP
		for {
			p, _, err := reader.ReadPacket()
			if err != nil {
				done <- err
				return
			}
			_, err = conn.Write(p)
			if err != nil {
				done <- err
				return
			}
			if *verbose {
				log.Println("serial -> udp")
			}
		}
	}()

	log.Fatal(<-done)
}
