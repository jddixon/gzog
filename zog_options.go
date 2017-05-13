package gzog

// gzog_go/zog_options.go

import (
	"fmt"
	xc "github.com/jddixon/xlCrypto_go"
	xt "github.com/jddixon/xlTransport_go"
	"log"
	"strconv"
	"strings"
)

var _ = fmt.Printf

// Options normally set from the command line or derived from those.
// Not used in this package but used by cmd/gzog/gzog
type ZogOptions struct {
	Address   string       // AddressI to listen on, default 127.0.0.1 = all
	EndPoint  xt.EndPointI // derived from Address, Port
	Ephemeral bool         // XXX probably don't need
	Lfs       string       // a path as in path/filepath
	Logger    *log.Logger
	Name      string // a valid server name containing no delimiters
	Port      uint16 // interpret 0 as DEFAULT_PORT or TEST_DEFAULT_PORT
	Testing   bool
	Verbose   bool
}

// SERIALIZATION

func (o *ZogOptions) String() string {
	ss := []string{
		"ZogOptions {",
		fmt.Sprintf("\tAddress:  %s", o.Address),
		// EndPoint is determined by Address and Port
		fmt.Sprintf("\tEphemeral: %s", o.Ephemeral),
		fmt.Sprintf("\tLfs:      %s", o.Lfs),
		// Logger is ephemeral, not serialized
		fmt.Sprintf("\tName:     %s", o.Name),
		fmt.Sprintf("\tPort:     %d", o.Port),
		fmt.Sprintf("\tTesting:  %s", o.Testing),
		fmt.Sprintf("\tverbose:  %s", o.Verbose),
		"}",
	}
	return strings.Join(ss, "\n") + "\n"
}

// DESERIALIZATION

// Parse a serialized ZogOptions
func ParseZogOptions(s string) (*ZogOptions, error) {
	ss := strings.Split(s, "\n") // yields slice of strings
	return ParseZogOptionsFromStrings(ss)
}

// Parse a string array containing a serialized ZogOptions, ignoring any
// leading or trailing whitespace.  xlNode.ParseBaseNode is the model.
func ParseZogOptionsFromStrings(ss []string) (*ZogOptions, error) {
	var (
		err     error
		address string // xt.AddressI
		// endPoint xt.EndPointI
		lfs     string
		name    string
		port    uint16
		testing bool
		verbose bool
	)
	foundBrace := false

	s, err := xc.NextNBLine(&ss)
	if err == nil {
		if s != "ZogOptions {" {
			err = MissingZogOptionsOpen
		}
	}
	if err == nil {
		for {
			s, err = xc.NextNBLine(&ss)
			if s == "}" {
				foundBrace = true
				break
			}
			parts := strings.SplitAfterN(s, ":", 1)
			switch parts[0] {
			case "Address":
				address = parts[1]
			case "Lfs":
				lfs = parts[1]
			case "Name":
				name = parts[1]
			case "Port":
				var port64 uint64
				port64, err = strconv.ParseUint(parts[1], 10, 16)
				if err == nil {
					port = uint16(port64)
				} else {
					break
				}
			case "Testing":
				val := parts[1]
				if val == "true" {
					testing = true
				} else if val == "false" {
					testing = false
				} else {
					err = NotAValidBoolean
					break
				}
			case "Verbose":
				val := parts[1]
				if val == "true" {
					verbose = true
				} else if val == "false" {
					verbose = false
				} else {
					err = NotAValidBoolean
					break
				}
			default:
				err = NotAnOption
				break
			}
		}
	}
	if err == nil {
		if foundBrace == false {
			err = MissingZogOptionsClose
		}
	}

	if err == nil {
		zo := ZogOptions(address, endPoint, // ephemeral,
			lfs, nil,
			name, port, testing, verbose)
		return &zo, nil
	} else {
		return nil, err
	}
}
