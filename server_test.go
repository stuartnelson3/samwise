package main

import (
    "testing"
    "net/http/httptest"
    "net/http"
)

func setup(t *testing.T) *httptest.Server {
    t.Log("Setup testing environment")
    mux := http.NewServeMux()
    mux.HandleFunc("/", heartBeat)
    return httptest.NewServer(mux)
}

func heartBeat(w http.ResponseWriter, r *http.Request) {
    if r.FormValue("token") == "secret" {
        w.WriteHeader(200)
    } else {
        w.WriteHeader(404)
    }
}

func servers() map[string]bool {
    return map[string]bool{"http://localhost:8080": true}
}

func TestStreamEndpoint(t *testing.T) {
    server := setup(t)
    resp, err := http.Get(server.URL + "?token=secret")
    if err != nil {
        t.Fail()
    }
    if resp.StatusCode != 200 {
        t.Error("request unsuccessful")
    }
}

func TestStreamEndpointNoToken(t *testing.T) {
    server := setup(t)
    resp, err := http.Get(server.URL)
    if err != nil {
        t.Fail()
    }
    if resp.StatusCode != 404 {
        t.Error("made successful request without token")
    }
}
