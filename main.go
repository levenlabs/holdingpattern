// Usage: ./holdingpattern [host] [addr]
// You can also pass command-line flags
package main

import (
	"strings"

	"github.com/levenlabs/go-llog"
	"github.com/mediocregopher/lever"
	"github.com/mediocregopher/skyapi/client"
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
	})
	l.Add(lever.Param{
		Name:        "--category",
		Description: "Category to advertise under. If unset the skyapi instance's global default will be used",
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
	category, _ := l.ParamStr("--category")
	addr, _ := l.ParamStr("--addr")
	weight, _ := l.ParamInt("--weight")
	priority, _ := l.ParamInt("--priority")

	argsFound := 0
	var argHost string
	var argAddr string
	for _, v := range l.ParamRest() {
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
			llog.Fatal("no hostname sent")
		}
	}

	if addr == "" {
		// allow them to run ./holdingpattern --host=test 127.0.0.1:8000
		if argHost != "" {
			addr = argHost
		} else if argAddr != "" {
			addr = argAddr
		} else {
			llog.Fatal("no address sent")
		}
	}

	hostcat := hostname
	if category != "" {
		hostcat = hostname + "." + category
	}
	llog.Info("advertising", llog.KV{
		"apiAddr":  apiAddr,
		"host":     hostcat,
		"thisAddr": addr,
		"priority": priority,
		"weight":   weight,
	})

	err := client.ProvideOpts(client.Opts{
		SkyAPIAddr:        apiAddr,
		Service:           hostname,
		ThisAddr:          addr,
		Category:          category,
		Priority:          priority,
		Weight:            weight,
		ReconnectAttempts: 3,
	})
	llog.Fatal("skyapi client failed", llog.KV{"err": err})
}
