package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	// "log"
	"os"
)

func setup(t *testing.T) *httptest.Server {
	t.Log("Setup testing environment")
	os.Setenv("TOKEN", "secret")
	c.Add("example")
	mux := http.NewServeMux()
	mux.HandleFunc("/", heartBeat)
	mux.HandleFunc("/stream", Stream)
	return httptest.NewServer(mux)
}

func heartBeat(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("token") == "secret" {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(404)
	}
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
