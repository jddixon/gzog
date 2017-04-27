package main

// gzog_go/cmd/gzog/gzog.go

import (
	"crypto/rand"
	"crypto/rsa"
	"flag"
	"fmt"
	"github.com/jddixon/gzog"
	xi "github.com/jddixon/xlNodeID_go"
	xn "github.com/jddixon/xlNode_go"
	xt "github.com/jddixon/xlTransport_go"
	xf "github.com/jddixon/xlUtil_go/lfs"
	"io/ioutil"
	"log"
	"os"
	"path"
)

func Usage() {
	fmt.Printf("Usage: %s [OPTIONS]\n", os.Args[0])
	fmt.Printf("where the options are:\n")
	flag.PrintDefaults()
}

const (
	DEFAULT_ADDR      = "127.0.0.1" // listen on all ports
	DEFAULT_NAME      = "gzog"
	DEFAULT_LFS       = "/var/app/gzog"
	TEST_DEFAULT_PORT = 44442 // for the server, not clients
	DEFAULT_PORT      = 55552 // for the server, not clients
)

var (
	// these need to be referenced as pointers
	address = flag.String("a", DEFAULT_ADDR,
		"server IP address (dotted quad)")
	justShow = flag.Bool("j", false,
		"display option settings and exit")
	lfs = flag.String("lfs", DEFAULT_LFS,
		"path to work directory")
	logFile = flag.String("l", "log",
		"path to log file")
	name = flag.String("n", DEFAULT_NAME,
		"server name")
	port = flag.Int("p", DEFAULT_PORT,
		"server listening port")
	testing = flag.Bool("T", false,
		"this is a test run")
	verbose = flag.Bool("v", false,
		"be talkative")
)

func init() {
	fmt.Println("init() invocation") // DEBUG
}

func main() {
	var err error

	flag.Usage = Usage
	flag.Parse()

	// FIXUPS ///////////////////////////////////////////////////////

	if err != nil {
		fmt.Println("error processing NodeID: %s\n", err.Error())
		os.Exit(-1)
	}
	if *testing {
		if *name == DEFAULT_NAME || *name == "" {
			*name = "testReg"
		}
		if *lfs == DEFAULT_LFS || *lfs == "" {
			*lfs = "./myApp/gzog"
		} else {
			*lfs = path.Join("tmp", *lfs)
		}
		if *address == DEFAULT_ADDR {
			*address = "127.0.0.1"
		}
		if *port == DEFAULT_PORT || *port == 0 {
			*port = TEST_DEFAULT_PORT
		}
	}
	addrAndPort := fmt.Sprintf("%s:%d", *address, *port)
	endPoint, err := xt.NewTcpEndPoint(addrAndPort)
	if err != nil {
		fmt.Printf("not a valid endPoint: %s\n", addrAndPort)
		Usage()
		os.Exit(-1)
	}

	// SANITY CHECKS ////////////////////////////////////////////////
	if err == nil {
		err = xf.CheckLFS(*lfs, 0700) // tries to create if it doesn't exist
		if err == nil {
			if *logFile != "" {
				*logFile = path.Join(*lfs, *logFile)
			}
		}
	}
	// DISPLAY STUFF ////////////////////////////////////////////////
	if *verbose || *justShow {
		fmt.Printf("endPoint         = %v\n", endPoint)
		fmt.Printf("justShow         = %v\n", *justShow)
		fmt.Printf("lfs              = %s\n", *lfs)
		fmt.Printf("logFile          = %s\n", *logFile)
		fmt.Printf("name             = %s\n", *name)
		fmt.Printf("server address   = %v\n", *address)
		fmt.Printf("server port      = %d\n", *port)
		fmt.Printf("testing          = %v\n", *testing)
		fmt.Printf("verbose          = %v\n", *verbose)
	}
	if *justShow {
		return
	}
	// SET UP OPTIONS ///////////////////////////////////////////////
	var (
		f      *os.File
		logger *log.Logger
		opt    gzog.ZogOptions
		// rs     *gzog.Gzog
	)
	if *logFile != "" {
		f, err = os.OpenFile(*logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
		if err == nil {
			logger = log.New(f, "", log.Ldate|log.Ltime)
		}
	}
	if f != nil {
		defer f.Close()
	}
	if err == nil {
		opt.Address = *address
		opt.Lfs = *lfs
		opt.Logger = logger
		opt.Lfs = *lfs
		opt.Port = fmt.Sprintf("%d", *port)
		opt.Testing = *testing
		opt.Verbose = *verbose

		rs, err = setup(&opt)
		if err == nil {
			err = serve(rs)
		}
	}
	_ = logger // NOT YET
	_ = err
}

// DEBUG //////////////////////////
func nothing(foo int) int {
	return foo
}

// END DEBUG //////////////////////

func setup(opt *gzog.ZogOptions) (rs *gzog.Gzog, err error) {
	// If LFS/.xlattice/gzog.config exists, we load that.  Otherwise we
	// create a node.  In either case we force the node to listen on
	// the designated port

	var (
		e                []xt.EndPointI
		node             *xn.Node
		pathToConfigFile string
		rn               *gzog.RegNode
		ckPriv, skPriv   *rsa.PrivateKey
	)
	logger := opt.Logger
	verbose := opt.Verbose

	greetings := fmt.Sprintf("xlReg v%s %s start run\n",
		gzog.VERSION, gzog.VERSION_DATE)
	if verbose {
		fmt.Print(greetings)
	}
	logger.Print(greetings)

	pathToConfigFile = path.Join(path.Join(opt.Lfs, ".xlattice"), "gzog.config")
	found, err := xf.PathExists(pathToConfigFile)
	if err == nil {
		if found {
			logger.Printf("Loading existing reg config from %s\n",
				pathToConfigFile)
			// The registry node already exists.  Parse it and we are done.
			var data []byte
			data, err = ioutil.ReadFile(pathToConfigFile)
			if err == nil {
				rn, _, err = gzog.ParseRegNode(string(data))
			}
		} else {
			logger.Println("No config file found, creating new registry.")
			// We need to create a registry node from scratch.
			nodeID, _ := xi.New(nil)
			ep, err = xt.NewTcpEndPoint(opt.Address + ":" + opt.Port)
			if err == nil {
				e = []xt.EndPointI{ep}
				ckPriv, err = rsa.GenerateKey(rand.Reader, 2048)
				if err == nil {
					skPriv, err = rsa.GenerateKey(rand.Reader, 2048)
				}
				if err == nil {
					node, err = xn.New("xlReg", nodeID, opt.Lfs, ckPriv, skPriv,
						nil, e, nil)
					if err == nil {
						node.OpenAcc() // XXX needs a complementary close
						if err == nil {
							// DEBUG
							fmt.Printf("XLattice node successfully created\n")
							fmt.Printf("  listening on %s\n", ep.String())
							// END
							rn, err = gzog.NewRegNode(node, ckPriv, skPriv)
							if err == nil {
								// DEBUG
								fmt.Printf("regNode successfully created\n")
								// END
								err = xf.MkdirsToFile(pathToConfigFile, 0700)
								if err == nil {
									err = ioutil.WriteFile(pathToConfigFile,
										[]byte(rn.String()), 0400)
									// DEBUG
								} else {
									fmt.Printf("error writing config file: %v\n",
										err.Error())
								}
								// END --------------

								// DEBUG
							} else {
								fmt.Printf("error creating regNode: %v\n",
									err.Error())
								// END
							}
						}
					}
				}
			}
		}
	}
	// TEMPORARILY COMMENTED OUT
	//if err == nil {
	//	var r *gzog.Registry
	//	r, err = gzog.NewRegistry(nil, // nil = clusters so far
	//		rn, opt) // regNode, options
	//	if err == nil {
	//		logger.Printf("Registry name: %s\n", rn.GetName())
	//		logger.Printf("         ID:   %s\n", rn.GetNodeID().String())
	//	}
	//	if err == nil {
	//		var verbosity int
	//		if opt.Verbose {
	//			verbosity++
	//		}
	//		rs, err = gzog.NewGzog(r, opt.Testing, verbosity)
	//	}
	//}
	if err != nil {
		logger.Printf("ERROR: %s\n", err.Error())
	}
	return
}
func serve(zog *gzog.Gzog) (err error) {

	err = zog.Start() // non-blocking
	if err == nil {
		<-zog.DoneCh
	}

	// XXX STUB XXX

	// ORDERLY SHUTDOWN /////////////////////////////////////////////

	return
}
