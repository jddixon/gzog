package gzog

// gzog_go/sprog_options.go

import (
	xt "github.com/jddixon/xlTransport_go"
	"log"
)

// Options normally set from the command line or derived from those.
// Not used in this package but used by cmd/sprog/sprog
type SprogOptions struct {
	EC2Host       bool
	EndPoint      xt.EndPointI // derived from Address, Port
	Interval      uint         // XXX WRONG TYPE -- seconds between messages
	JustShow      bool         // display options and exit
	Lfs           string
	Logger        *log.Logger
	MsgText       string // message text
	N             uint   // number of messages to send
	Name          string
	ServerPort    string // constrained to be a uint16
	ServerAddress string // dotted quad, defaults to 127.0.0.1
	Testing       bool
	Verbose       bool
}
