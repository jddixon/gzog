package gzog

// gzog_go/zog_options.go

import (
	xt "github.com/jddixon/xlTransport_go"
	"log"
)

// Options normally set from the command line or derived from those.
// Not used in this package but used by cmd/gzog/gzog
type ZogOptions struct {
	Address   string       // port to listen on, default 127.0.0.1 = all
	EndPoint  xt.EndPointI // derived from Address, Port
	Ephemeral bool         // XXX probably don't need
	Lfs       string
	Logger    *log.Logger
	Name      string
	Port      string // gets restricted to being an uint16
	Testing   bool
	Verbose   bool
}
