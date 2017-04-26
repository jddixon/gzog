// gzog/const.go

package gzog

const (
	MAX_MSG     = 512   // maximum number of characters in a message
	PEER_PORT   = 55552 // expect connections from other daemons on this port
	CLIENT_PORT = 55551 // listen for client messages on this port
)

// Go doesn't like maps as constants.
var (
	RING_IP_ADDR = map[string]string{
		"losaltos": "192.168.152.253",
		"losgatos": "192.168.136.254",
		"mtview":   "192.168.152.242",
		"paloalto": "192.168.152.252",
		"saratoga": "192.168.152.11",
	}
)
