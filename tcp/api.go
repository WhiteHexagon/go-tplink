package tcp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
)

//see https://www.softscheck.com/en/reverse-engineering-tp-link-hs110/

//CheckError ...
func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func send(conn net.Conn, cmd string) {
	mess := encrypt(cmd)
	sent, err := fmt.Fprintf(conn, mess)
	CheckError(err)
	if len(mess) != sent {
		CheckError(errors.New("partial send"))
	}
}

// SendCommand - returns result of call
func SendCommand(addr string, cmd string) string {
	//fmt.Printf("connect to addr: [%s]\n", addr)
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:9999", addr))
	CheckError(err)
	defer conn.Close()

	send(conn, cmd)

	reader := bufio.NewReader(conn)
	buffer := make([]byte, 1024)
	got, err := reader.Read(buffer)
	if err == io.EOF {
		fmt.Println("end")
		return "failed"
	}
	result := decrypt(string(buffer[4:got]))
	return result
}

//FetchInfo - ...
func FetchInfo(addr string) string {
	return SendCommand(addr, `{"system":{"get_sysinfo":{}}}`)
}

//Switch - ...
func Switch(addr string, state bool) string {
	onOff := "0"
	if state {
		onOff = "1"
	}
	return SendCommand(addr, fmt.Sprintf(`{"system":{"set_relay_state":{"state":%s}}}`, onOff))
}

//ExtractSummaryFrom - extract name and signal from info
func ExtractSummaryFrom(info string) (string, string) {
	pos := strings.Index(info, "alias\":\"")
	if pos != -1 {
		rest := info[pos+8:]
		pos = strings.Index(rest, "\"")
		rest = rest[:pos]
		pos = strings.Index(info, "rssi\":")
		if pos != -1 {
			rssi := info[pos+6:]
			pos = strings.Index(rssi, ",")
			rssi = rssi[:pos]
			return rest, rssi
		}
		return rest, "unknown"
	}
	return "unknown", "unknown"
}
