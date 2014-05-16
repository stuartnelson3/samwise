package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
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
	mux.HandleFunc("/add_server", addServer)
	mux.HandleFunc("/update_stream", updateStream)
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
		t.Errorf("Response: %d, expected 404", resp.StatusCode)
	}
}

func TestHeartbeatEndpointNoToken(t *testing.T) {
	server := setup(t)
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fail()
	}
	if resp.StatusCode != 404 {
		t.Errorf("Response: %d, expected 404", resp.StatusCode)
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
		t.Errorf("Response: %d, expected 404", resp.StatusCode)
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
		// aka FIXME
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
		t.Error("url not marked active")
		t.Fail()
	}

	if !inHash(server.URL) {
		t.Error("url not added to consistent hash")
		t.Fail()
	}
}

func TestAddServers(t *testing.T) {
	server := setup(t)
	data := url.Values{}
	url := server.URL
	data.Set("server", url)
	servers = map[string]bool{url: false}
	res, err := http.PostForm(url+"/add_server?token=secret", data)
	if err != nil {
		t.Error("Error posting to add_server")
	}
	if res.StatusCode != 200 {
		t.Errorf("Response: %d, expected 200", res.StatusCode)
	}
	if !inHash(url) {
		t.Error("Server not added to hash")
	}
	if !servers[url] {
		t.Error("Server status incorrect")
	}
}

func TestUpdateStream(t *testing.T) {
	server := setup(t)
	res, err := http.PostForm(server.URL+"/update_stream?token=secret&stream=stream123", url.Values{})
	if err != nil {
		t.Error("Error posting to add_server")
	}
	if res.StatusCode != 200 {
		t.Errorf("Response: %d, expected 200", res.StatusCode)
	}
}
