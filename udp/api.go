package udp

import (
	"fmt"
	"net"
	"time"
)

//see https://www.softscheck.com/en/reverse-engineering-tp-link-hs110/

func checkError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		panic("failed.")
	}
}

func connect() *net.UDPConn {
	localAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:4242")
	checkError(err)
	conn, err := net.ListenUDP("udp", localAddr)
	checkError(err)
	return conn
}

func sendBroadcastPing(conn *net.UDPConn) {
	bcastAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:9999")
	checkError(err)
	_, err = conn.WriteToUDP([]byte("\xd0\xf2\x81\xf8\x8b\xff\x9a\xf7\xd5\xef\x94\xb6\xd1\xb4\xc0\x9f\xec\x95\xe6\x8f\xe1\x87\xe8\xca\xf0\x8b\xf6\x8b\xf6"), bcastAddr)
	checkError(err)
}

func waitForResponses(conn *net.UDPConn, maxWait time.Duration, handler func(string)) error {
	buffer := make([]byte, 1024) //expecting ~600bytes per device
	err := conn.SetReadDeadline(time.Now().Add(maxWait))
	checkError(err)
	for {
		readLen, fromAddr, err := conn.ReadFromUDP(buffer)
		//check for timeout
		if err != nil {
			e, isa := err.(net.Error)
			if isa && e.Timeout() {
				//we expect this error
				break
			} else {
				return err
			}
		}
		if readLen > 0 {
			fmt.Printf("found [%s]\n", fromAddr)
			handler(fromAddr.IP.String())
		} //else ignore the response
	}
	return nil
}

//Ping - wait maxTime for any responses, returns a list of discovered TCP IP addresses, ignores errors
func Ping(maxWait time.Duration) []string {
	result := []string{}
	err := PingAndDo(maxWait, func(addr string) {
		result = append(result, addr)
	})
	if err != nil {
		fmt.Println("Ping::", err)
	}
	return result
}

//PingAndDo - wait maxTime for any responses, for each response the handler function is called with the address as a param
func PingAndDo(maxWait time.Duration, handler func(string)) error {
	conn := connect()
	defer conn.Close()
	sendBroadcastPing(conn)
	return waitForResponses(conn, maxWait, handler)
}
