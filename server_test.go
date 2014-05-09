package main

import (
    "testing"
    "net/http/httptest"
    "net/http"
)

func setup(t *testing.T) *httptest.Server {
    t.Log("Setup testing environment")
    return httptest.NewServer(http.HandlerFunc(streamFunc))
}

func streamFunc(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(200)
}

func servers() map[string]bool {
    return map[string]bool{"http://localhost:8080": true}
}

func TestStreamEndpoint(t *testing.T) {
    server := setup(t)
    resp, err := http.Get(server.URL)
    if err != nil {
        t.Fail()
    }
    if resp.StatusCode != 200 {
        t.Error("status code wasn't 200")
    }
}
