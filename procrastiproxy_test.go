package procrastiproxy_test

import (
	"bytes"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/laychopy/procrastiproxy"
)

func TestProxyServerAllowsAllRequestsByDefault(t *testing.T) {
	t.Parallel()
	ts := newTestServer(t)
	_, client := newTestProxyAndClient(t)
	requestIsAllowed(t, client, ts.URL)
}

func TestProxyServerDeniesRequestsToBlockedURLs(t *testing.T) {
	t.Parallel()
	ts := newTestServer(t)
	proxy, client := newTestProxyAndClient(t)
	proxy.Block(strings.TrimPrefix(ts.URL, "http://"))
	requestIsDenied(t, client, ts.URL)
}

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	want := []byte("Hello, world")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(want)
	}))
	t.Cleanup(ts.Close)
	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatal(resp.StatusCode)
	}
	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("want %q, got %q", want, got)
	}
	return ts
}

func requestIsAllowed(t *testing.T, client *http.Client, link string) {
	t.Helper()
	resp, err := client.Get(link)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatal(resp.StatusCode)
	}
	want := []byte("Hello, world")
	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("want %q, got %q", want, got)
	}
}

func requestIsDenied(t *testing.T, client *http.Client, link string) {
	t.Helper()
	resp, err := client.Get(link)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusForbidden {
		t.Fatal(resp.StatusCode)
	}
}

func newTestProxyAndClient(t *testing.T) (*procrastiproxy.ProxyServer, *http.Client) {
	proxy := procrastiproxy.NewServer(randomFreePortAddr(t))
	t.Cleanup(func() { proxy.Close() })
	go func() {
		err := proxy.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Error(err)
		}
	}()
	return proxy, &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(&url.URL{
				Scheme: "http",
				Host:   proxy.Addr,
			}),
		},
	}
}

func randomFreePortAddr(t *testing.T) string {
	t.Helper()
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	return l.Addr().String()
}
