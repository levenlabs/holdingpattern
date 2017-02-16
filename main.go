// Usage: ./holdingpattern [host] [addr]
// You can also pass command-line flags
package main

import (
	"strings"
	"time"

	"github.com/levenlabs/go-llog"
	"github.com/levenlabs/go-srvclient"
	"github.com/mediocregopher/lever"
	"github.com/mediocregopher/skyapi/client"
	"github.com/miekg/dns"
)

// prefixResolver implements the Resolver interface
type prefixResolver struct {
	prefix string
	srv    *srvclient.SRVClient
}

func (r *prefixResolver) preprocess(m *dns.Msg) {
	if r.prefix == "" {
		return
	}
	for i := range m.Answer {
		if ansSRV, ok := m.Answer[i].(*dns.SRV); ok {
			tar := ansSRV.Target
			if strings.HasPrefix(tar, r.prefix+"-") {
				if ansSRV.Priority < 2 {
					ansSRV.Priority = uint16(0)
				} else {
					ansSRV.Priority = ansSRV.Priority - 1
				}
			}
		}
	}
}

func (r *prefixResolver) Resolve(h string) (string, error) {
	if r.srv == nil {
		r.srv = new(srvclient.SRVClient)
		r.srv.EnableCacheLast()
		r.srv.Preprocess = r.preprocess
	}
	return r.srv.SRV(h)
}

func main() {
	l := lever.New("holdingpattern", nil)
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
	l.Add(lever.Param{
		Name:        "--prefix",
		Description: "Prefix to pass to skyapi. This will be prefixed to the unique id that is actually stored by skyapi",
	})
	l.Add(lever.Param{
		Name:        "--prefer-prefixed-skyapi",
		Description: "Prefer skyapi hostnames with the same --prefix",
		Flag:        true,
	})
	l.Parse()

	apiAddr, _ := l.ParamStr("--api")
	hostname, _ := l.ParamStr("--hostname")
	category, _ := l.ParamStr("--category")
	addr, _ := l.ParamStr("--addr")
	weight, _ := l.ParamInt("--weight")
	priority, _ := l.ParamInt("--priority")
	prefix, _ := l.ParamStr("--prefix")

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

	opts := client.Opts{
		SkyAPIAddr:        apiAddr,
		Service:           hostname,
		ThisAddr:          addr,
		Category:          category,
		Priority:          priority,
		Weight:            weight,
		Prefix:            prefix,
		ReconnectAttempts: 0, // do not attempt to reconnect, we'll do that here
	}

	if l.ParamFlag("--prefer-prefixed-skyapi") {
		opts.Resolver = &prefixResolver{
			prefix: prefix,
		}
	}

	kv := llog.KV{
		"apiAddr":  apiAddr,
		"host":     hostcat,
		"thisAddr": addr,
		"priority": priority,
		"weight":   weight,
		"prefix":   prefix,
	}
	for {
		llog.Info("advertising", kv)
		err := client.ProvideOpts(opts)
		if err != nil {
			llog.Warn("skyapi error", kv, llog.ErrKV(err))
		}
		time.Sleep(1 * time.Second)
	}
}
