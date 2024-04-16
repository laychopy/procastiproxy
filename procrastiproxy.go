package procrastiproxy

import (
	"fmt"
	"net"
	"net/http"
)

type ProxyServer struct {
	server   *http.Server
	listener net.Listener
}

func NewServer() *http.Server {
	return &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
		}),
	}
}

func ListenTCP(address string) (net.Listener, error) {
	return net.Listen("tcp", address)
}

func Serve(server *http.Server, listener net.Listener) error {
	println("Starting server on port", listener.Addr().String())
	return server.Serve(listener)
}

func ListenAndServe(address string) error {
	listener, err := ListenTCP(address)
	if err != nil {
		return err
	}

	return Serve(NewServer(), listener)
}

func Main() int {
	err := ListenAndServe(":0")
	if err != nil {
		fmt.Printf("Error starting proxy server: %v", err)
		return 1
	}
	return 0
}
