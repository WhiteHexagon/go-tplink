package main

import (
	"fmt"
	"time"

	"github.com/whitehexagon/go-tplink/tcp"
	"github.com/whitehexagon/go-tplink/udp"
)

//https://github.com/softScheck/tplink-smartplug

//reminder to self, allow firewall access for testing
func main() {
	devices := udp.Ping(time.Second * 2)
	for _, addr := range devices {
		info := tcp.FetchInfo(addr)
		a, p := tcp.ExtractSummaryFrom(info)
		fmt.Printf("%s->%s(%s)\n", addr, a, p)
	}
}
