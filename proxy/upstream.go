package proxy

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
)

type Upstream struct {
	Name   string
	handle *http.Client
}

func NewSocket(socket string) Upstream {
	stream := Upstream{Name: socket}
	stream.handle = &http.Client{
		Transport: &http.Transport{
			Dial: func(proto, addr string) (net.Conn, error) {
				conn, err := net.Dial("unix", socket)
				return conn, err
			},
		},
	}
	return stream
}

func (r Upstream) Pass() func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			http.Error(res, "400 Bad request; only GET allowed.", 400)
			return
		}
		param := ""
		if len(req.URL.RawQuery) > 0 {
			param = "?" + req.URL.RawQuery
		}
		body, _ := r.Get("http://docker"+req.URL.Path+param, res)
		fmt.Fprintf(res, "%s", body)
	}
}

func (r Upstream) Get(url string, res http.ResponseWriter) ([]byte, error) {
	req, err := r.handle.Get(url)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	contentType := req.Header.Get("Content-type")
	if contentType != "" {
		res.Header().Set("Content-Type", contentType)
	}
	return io.ReadAll(req.Body)
}

func (r Upstream) PassStream() func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			http.Error(res, "400 Bad request; only GET allowed.", 400)
			return
		}
		param := ""
		if len(req.URL.RawQuery) > 0 {
			param = "?" + req.URL.RawQuery
		}

		resp, err := r.handle.Get("http://docker" + req.URL.Path + param)
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}
		defer resp.Body.Close()
		contentType := resp.Header.Get("Content-Type")

		if contentType != "" {
			res.Header().Set("Content-Type", contentType)
		}

		reader := bufio.NewReader(resp.Body)

		for {
			line, err := reader.ReadBytes('\n')

			if err != nil {
				break
			}

			fmt.Fprintf(res, "%s", line)
		}
	}
}
