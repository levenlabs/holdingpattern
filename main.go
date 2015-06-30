// Usage: ./holdingpattern [host] [addr]
// You can also pass command-line flags
package main

import (
	"github.com/mediocregopher/skyapi/client"
	"github.com/mediocregopher/lever"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	l := lever.New("holderingpattern", nil)
	l.Add(lever.Param{
		Name:        "--api",
		Description: "Address skyapi is listening on",
		Default:     "127.0.0.1:8053",
	})
	l.Add(lever.Param{
		Name:        "--hostname",
		Description: "Hostname to advertise",
		Default:     "",
	})
	l.Add(lever.Param{
		Name:        "--addr",
		Description: "Address and port to advertise",
		Default:     "",
	})
	l.Add(lever.Param{
		Name:        "--weight",
		Description: "Weight to advertise",
		Default:     "100",
	})
	l.Add(lever.Param{
		Name:        "--priority",
		Description: "Priority to advertise",
		Default:     "1",
	})
	l.Parse()

	apiAddr, _ := l.ParamStr("--api")
	hostname, _ := l.ParamStr("--hostname")
	addr, _ := l.ParamStr("--addr")
	weight, _ := l.ParamInt("--weight")
	priority, _ := l.ParamInt("--priority")

	argsFound := 0
	var argHost string
	var argAddr string
	for _, v := range os.Args[1:] {
		if strings.HasPrefix(v, "-") {
			continue
		}
		switch argsFound {
		case 0:
			argHost = v
		case 1:
			argAddr = v
		}
		argsFound++
	}

	if hostname == "" {
		if argHost != "" {
			hostname = argHost
			argHost = "" //reset so addr doesn't pick it up
		} else {
			log.Fatal("No hostname sent")
		}
	}

	if addr == "" {
		// allow them to run ./holdingpattern --host=test 127.0.0.1:8000
		if argHost != "" {
			addr = argHost
		} else if argAddr != "" {
			addr = argAddr
		} else {
			log.Fatal("No address sent")
		}
	}

	log.Printf("Advertising [%s]: %s on %s with priority %d weight %d", apiAddr, hostname, addr, priority, weight)

	log.Fatal(client.Provide(apiAddr, hostname, addr, weight, priority, -1, 5*time.Second))
}
