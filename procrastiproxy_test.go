package procrastiproxy_test

import (
	"errors"
	"net"
	"net/http"
	"testing"

	procrastiproxy "github.com/laychopy/procastiproxy"
)

// TODO: REFACTOR AND PROXY REQUEST
// TODO: ADD TESTS FOR PROXY SERVER
// Your tests go here!
func TestProxyServerRespondsWithStatusForbiddenRunningOnAddress(t *testing.T) {
	t.Parallel()
	server := procrastiproxy.NewServer()
	address, err := procrastiproxy.ListenTCP(":0")
	defer server.Close()
	// start proxy server
	go func(address net.Listener, server *http.Server) {
		err = procrastiproxy.Serve(server, address)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Errorf("Error starting proxy server: %v", err)
		}
	}(address, server)

	// make request to proxy server
	response, err := http.Get("http://" + address.Addr().String())

	// assert err is nil
	if err != nil {
		t.Fatalf("Error making request to proxy server: %v", err)
	}
	// assert response status code is 403
	if response.StatusCode != http.StatusForbidden {
		t.Fatalf("Expected status code 403, got %v", response.StatusCode)
	}
}

//func TestScript(t *testing.T) {
//	testscript.Run(t, testscript.Params{
//		Dir: "testdata/scripts",
//	})
//}

//func TestMain(m *testing.M) {
//	os.Exit(testscript.RunMain(m, map[string]func() int{
//		"procrastiproxy": procrastiproxy.Main,
//	}))
//}
