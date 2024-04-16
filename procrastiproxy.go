package procrastiproxy

import (
	"fmt"
	"io"
	"net"
	"net/http"
)

type ProxyServer struct {
	Addr     string
	server   *http.Server
	listener net.Listener
	blocked  map[string]bool
}

func NewServer(addr string) *ProxyServer {
	p := &ProxyServer{
		Addr:    addr,
		blocked: map[string]bool{},
	}
	p.server = &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if p.Deny(r.Host) {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			resp, err := http.Get(r.RequestURI)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, err)
			}
			defer resp.Body.Close()
			w.WriteHeader(resp.StatusCode)
			io.Copy(w, resp.Body)
		}),
	}
	return p
}

func (p *ProxyServer) Block(link string) {
	p.blocked[link] = true
}

func (p *ProxyServer) Close() error {
	return p.server.Close()
}

func (p *ProxyServer) Deny(link string) bool {
	return p.blocked[link]
}

func (p *ProxyServer) ListenAndServe() error {
	return p.server.ListenAndServe()
}

func Main() int {
	err := NewServer(":0").ListenAndServe()
	if err != nil {
		fmt.Printf("Error starting proxy server: %v", err)
		return 1
	}
	return 0
}
