package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func setup(t *testing.T) *httptest.Server {
	t.Log("Setup testing environment")
	os.Setenv("TOKEN", "secret")
	c.Add("example")
	mux := http.NewServeMux()
	mux.HandleFunc("/", heartBeat)
	mux.HandleFunc("/stream", stream)
	return httptest.NewServer(mux)
}

func endpoints() map[string]bool {
	return map[string]bool{"http://localhost:8080": true}
}

func TestHeartbeatEndpoint(t *testing.T) {
	server := setup(t)
	resp, err := http.Get(server.URL + "?token=secret")
	if err != nil {
		t.Fail()
	}
	if resp.StatusCode != 200 {
		t.Error("request unsuccessful")
	}
}

func TestHeartbeatEndpointNoToken(t *testing.T) {
	server := setup(t)
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fail()
	}
	if resp.StatusCode != 404 {
		t.Error("made successful request without token")
	}
}

func TestStreamEndpointNoToken(t *testing.T) {
	server := setup(t)
	url := server.URL + "/stream"
	resp, err := http.Get(url)
	if err != nil {
		t.Fail()
	}
	if resp.StatusCode != 404 {
		t.Error("stream made successful request without token")
	}
}

func TestStreamEndpointNoStream(t *testing.T) {
	server := setup(t)
	url := server.URL + "/stream"
	resp, err := http.Get(url)
	if err != nil {
		t.Fail()
	}
	if resp.StatusCode != 404 {
		t.Error("stream made successful request without stream")
	}
}

func TestStreamEndpoint(t *testing.T) {
	server := setup(t)
	url := server.URL + "/stream?token=secret&stream=stream1"
	resp, err := http.Get(url)
	if err != nil {
		t.Fail()
	}
	if resp.StatusCode != 404 {
		// got redirected and then 404'd
		t.Error("did not get redirected")
		t.Errorf("status code: %d", resp.StatusCode)
	}
}

func TestMonitorServersError(t *testing.T) {
	c.Add("garbageUrl")
	servers := map[string]bool{"garbageUrl": true}
	monitor(servers)
	if servers["garbageUrl"] {
		t.Error("garbage url not marked inactive")
		t.Fail()
	}

	if inHash("garbageUrl") {
		t.Error("url not removed from consistent hashing")
		t.Fail()
	}
}

func TestMonitorServersSuccess(t *testing.T) {
	server := setup(t)
	c.Add(server.URL)
	servers := map[string]bool{server.URL: false}
	monitor(servers)
	if !servers[server.URL] {
		t.Error("correct url not marked active")
		t.Fail()
	}

	if !inHash(server.URL) {
		t.Error("correct url not marked active")
		t.Fail()
	}
}
