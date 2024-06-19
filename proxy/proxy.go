package proxy

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/seabreeze-project/seabreeze/util"
)

func SupportedEndpoints() []string {
	return []string{"containers", "images", "volumes", "networks", "services", "tasks", "events", "version", "info", "ping"}
}

type Proxy struct {
	opt      ProxyOptions
	mux      *http.ServeMux
	listener net.Listener
}

type ProxyOptions struct {
	AllowedEndpoints   []string
	UpstreamSocketPath string
	SocketPath         string
	SocketMode         int32
	SocketOwner        string
}

func New(rootMux *http.ServeMux, opt ProxyOptions) (*Proxy, error) {
	supportedEndpoints := SupportedEndpoints()
	var allowedMap map[string]bool = make(map[string]bool)
	for _, s := range opt.AllowedEndpoints {
		if !util.StringInSlice(s, supportedEndpoints) {
			return nil, fmt.Errorf("unsupported endpoint: %s", s)
		}
		allowedMap[s] = true
	}

	var routers [2]*mux.Router
	routers[0] = mux.NewRouter()
	routers[1] = routers[0].PathPrefix("/{version:[v][0-9]+[.][0-9]+}").Subrouter()

	upstream := NewSocket(opt.UpstreamSocketPath)

	if allowedMap["containers"] {
		for _, m := range routers {
			containers := m.PathPrefix("/containers").Subrouter()
			containers.HandleFunc("/json", upstream.Pass())
			containers.HandleFunc("/{name}/json", upstream.Pass())
		}
	}

	if allowedMap["images"] {
		for _, m := range routers {
			containers := m.PathPrefix("/images").Subrouter()
			containers.HandleFunc("/json", upstream.Pass())
			containers.HandleFunc("/{name}/json", upstream.Pass())
			containers.HandleFunc("/{name}/history", upstream.Pass())
		}
	}

	if allowedMap["volumes"] {
		for _, m := range routers {
			m.HandleFunc("/volumes", upstream.Pass())
			m.HandleFunc("/volumes/{name}", upstream.Pass())
		}
	}

	if allowedMap["networks"] {
		for _, m := range routers {
			m.HandleFunc("/networks", upstream.Pass())
			m.HandleFunc("/networks/{name}", upstream.Pass())
		}
	}

	if allowedMap["services"] {
		for _, m := range routers {
			m.HandleFunc("/services", upstream.Pass())
			m.HandleFunc("/services/{name}", upstream.Pass())
		}
	}

	if allowedMap["tasks"] {
		for _, m := range routers {
			m.HandleFunc("/tasks", upstream.Pass())
			m.HandleFunc("/tasks/{name}", upstream.Pass())
		}
	}

	if allowedMap["events"] {
		for _, m := range routers {
			m.HandleFunc("/events", upstream.PassStream())
		}
	}

	if allowedMap["version"] {
		for _, m := range routers {
			m.HandleFunc("/version", upstream.Pass())
		}
	}

	if allowedMap["info"] {
		for _, m := range routers {
			m.HandleFunc("/info", upstream.Pass())
		}
	}

	if allowedMap["ping"] {
		for _, m := range routers {
			m.HandleFunc("/_ping", upstream.Pass())
		}
	}

	rootMux.Handle("/", routers[0])

	err := setupSocket(opt)
	if err != nil {
		return nil, err
	}

	return &Proxy{mux: rootMux, opt: opt}, nil
}

func (p *Proxy) Run() error {
	l, err := net.Listen("unix", p.opt.SocketPath)
	if err != nil {
		return err
	}

	p.listener = l

	err = http.Serve(l, p.mux)
	if err != nil {
		return err
	}

	return nil
}

func (p *Proxy) Close() error {
	if p.listener == nil {
		return errors.New("proxy is not running")
	}
	return p.listener.Close()
}

func setupSocket(opt ProxyOptions) error {
	// create the socket file's parent directories if they don't exist
	os.MkdirAll(filepath.Dir(opt.SocketPath), 0777)

	os.Chmod(opt.SocketPath, os.FileMode(opt.SocketMode))
	u, err := user.Lookup(opt.SocketOwner)
	if err != nil {
		return err
	}

	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return err
	}
	gid, err := strconv.Atoi(u.Gid)
	if err != nil {
		return err
	}

	err = os.Chown(opt.SocketPath, uid, gid)
	if err != nil {
		return err
	}

	return nil
}
