/*

MIT License

Copyright (c) 2017 Peter Bjorklund

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/

package server

import (
	"encoding/hex"
	"fmt"

	"net"
	"time"
)

type Server struct {
	connection          *net.UDPConn
	destinationEndpoint *net.UDPAddr
	mainClientEndpoint  *net.UDPAddr
	verbose             bool
}

// New : Creates a new server
func New(listenPort int, destination string, verbose bool) (*Server, error) {
	udpAddr, udpAddrErr := net.ResolveUDPAddr("udp", destination)
	if udpAddrErr != nil {
		return nil, udpAddrErr
	}

	s := &Server{destinationEndpoint: udpAddr, verbose: verbose}
	serverConnection, serverConnectionErr := s.listen(listenPort)
	if serverConnectionErr != nil {
		return nil, serverConnectionErr
	}
	s.connection = serverConnection

	return s, nil
}

func (server *Server) listen(port int) (*net.UDPConn, error) {
	portString := fmt.Sprintf(":%d", port)
	listenAddr, err := net.ResolveUDPAddr("udp", portString)
	if err != nil {
		return nil, fmt.Errorf("Error:%v ", err)
	}
	serverConnection, err := net.ListenUDP("udp", listenAddr)
	if err != nil {
		return nil, fmt.Errorf("Error: %v", err)
	}

	fmt.Printf("Listening to %s\n", portString)
	return serverConnection, nil
}

func (server *Server) fromDestination(fromAddr *net.UDPAddr) bool {
	return (fromAddr.IP.Equal(server.destinationEndpoint.IP) && fromAddr.Port == server.destinationEndpoint.Port)
}

func (server *Server) handlePacket(payload []byte, fromAddr *net.UDPAddr) error {
	toAddr := server.destinationEndpoint
	if server.fromDestination(fromAddr) {
		if server.mainClientEndpoint == nil {
			return nil
		}
		toAddr = server.mainClientEndpoint
	} else {
		server.mainClientEndpoint = fromAddr
	}
	if server.verbose {
		hexPayload := hex.Dump(payload)
		fmt.Println("Sending ", hexPayload, " to ", toAddr)
	}
	server.connection.WriteToUDP(payload, toAddr)
	return nil
}

func (server *Server) handleIncomingUDP() {
	for {
		buf := make([]byte, 1800)
		n, addr, err := server.connection.ReadFromUDP(buf)
		packet := buf[0:n]
		if server.verbose {
			hexPayload := hex.Dump(packet)
			fmt.Println("Received ", hexPayload, " from ", addr)
		}
		if err != nil {
			fmt.Println("Error: ", err)
		}
		packetErr := server.handlePacket(packet, addr)
		if packetErr != nil {
			fmt.Printf("Problem with packet:%s\n", packetErr)
		}
	}
}

func (server *Server) tick() error {
	return nil
}

func (server *Server) start(ticker *time.Ticker) {
	go func() {
		for range ticker.C {
			err := server.tick()
			if err != nil {
				fmt.Printf("Start err %s \n", err)
			}
		}
	}()
}

func (server *Server) Forever() error {
	go server.handleIncomingUDP()
	//defer serverConnection.Close()
	ticker := time.NewTicker(time.Millisecond * 100)
	server.start(ticker)
	select {}
}
