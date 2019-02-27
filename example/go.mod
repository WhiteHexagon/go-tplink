module github.com/whitehexagon/go-tplink

require (
	github.com/whitehexagon/go-tplink/tcp v0.0.0
	github.com/whitehexagon/go-tplink/udp v0.0.0
)

replace github.com/whitehexagon/go-tplink/udp v0.0.0 => ../udp

replace github.com/whitehexagon/go-tplink/tcp v0.0.0 => ../tcp
